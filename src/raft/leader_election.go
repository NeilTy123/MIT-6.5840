package raft

import (
	"math/rand"
	"time"
)

// let the base election timeout be T.
// the election timeout is in the range [T, 2T).
const baseElectionTimeout = 300

func (rf *Raft) pastElectionTimeout() bool {
	return time.Since(rf.lastElection) > rf.electionTimeout
}

func (rf *Raft) resetElectionTimer() {
	electionTimeout := baseElectionTimeout + (rand.Int63() % baseElectionTimeout)
	rf.electionTimeout = time.Duration(electionTimeout) * time.Millisecond
	rf.lastElection = time.Now()
}

func (rf *Raft) becomeFollower(term uint64) bool {
	termChanged := false
	oldTerm := rf.term

	if term > rf.term {
		rf.term = term
		rf.votedTo = None
		termChanged = true
	}

	if rf.state != Follower {
		rf.state = Follower
		rf.logger.stateToFollower(oldTerm)
	}

	return termChanged
}

func (rf *Raft) becomeCandidate() {
	defer rf.persist()
	rf.term++
	rf.votedMe = make([]bool, len(rf.peers))
	rf.votedTo = rf.me
	rf.logger.stateToCandidate()
	rf.state = Candidate
	rf.resetElectionTimer()
}

func (rf *Raft) becomeLeader() {
	rf.resetTrackedIndexes()
	rf.logger.stateToLeader()
	rf.state = Leader
}

func (rf *Raft) makeRequestVoteArgs(to int) *RequestVoteArgs {
	args := new(RequestVoteArgs)
	lastLogIndex := rf.log.lastIndex()
	lastLogTerm, _ := rf.log.term(lastLogIndex)
	*args = RequestVoteArgs{From: rf.me, To: to, Term: rf.term, LastLogIndex: lastLogIndex, LastLogTerm: lastLogTerm}
	return args
}

func (rf *Raft) sendRequestVote(args *RequestVoteArgs) {
	reply := RequestVoteReply{}
	if ok := rf.peers[args.To].Call("Raft.RequestVote", args, &reply); ok {
		rf.handleRequestVoteReply(args, &reply)
	}
}

func (rf *Raft) broadcastRequestVote() {
	rf.logger.bcastRVOT()
	for i := range rf.peers {
		if i != rf.me {
			args := rf.makeRequestVoteArgs(i)
			go rf.sendRequestVote(args)
		}
	}
}

func (rf *Raft) eligibleToGrantVote(candidateLastLogIndex, candidateLastLogTerm uint64) bool {
	lastLogIndex := rf.log.lastIndex()
	lastLogTerm, _ := rf.log.term(lastLogIndex)
	return candidateLastLogTerm > lastLogTerm || (candidateLastLogTerm == lastLogTerm && candidateLastLogIndex >= lastLogIndex)
}

// RequestVote request handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.logger.recvRVOT(args)

	reply.From = rf.me
	reply.To = args.From
	reply.Term = rf.term
	reply.Voted = false

	m := Message{Type: Vote, From: args.From, Term: args.Term}
	ok, termChanged := rf.checkMessage(m)
	if termChanged {
		reply.Term = rf.term
		defer rf.persist()
	}
	if !ok {
		return
	}

	// FIXME: does a pending snapshot affect granting vote?
	if (rf.votedTo == None || rf.votedTo == args.From) && rf.eligibleToGrantVote(args.LastLogIndex, args.LastLogTerm) {
		rf.votedTo = args.From
		rf.resetElectionTimer()
		reply.Voted = true

		rf.logger.voteTo(args.From)

	} else {
		lastLogIndex := rf.log.lastIndex()
		lastLogTerm, _ := rf.log.term(lastLogIndex)
		rf.logger.rejectVoteTo(args.From, args.LastLogIndex, args.LastLogTerm, lastLogIndex, lastLogTerm)
	}
}

func (rf *Raft) receivedMajorityVotes() bool {
	votes := 1
	for i, votedMe := range rf.votedMe {
		if i != rf.me && votedMe {
			votes++
		}
	}
	if 2*votes > len(rf.peers) {
		rf.logger.recvVoteQuorum()
	}
	return 2*votes > len(rf.peers)
}

func (rf *Raft) handleRequestVoteReply(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.logger.recvRVOTRes(reply)

	m := Message{Type: VoteReply, From: reply.From, Term: reply.Term, ArgsTerm: args.Term}
	ok, termChanged := rf.checkMessage(m)
	if termChanged {
		defer rf.persist()
	}
	if !ok {
		return
	}

	if reply.Voted {
		rf.votedMe[reply.From] = true
		if rf.receivedMajorityVotes() {
			rf.becomeLeader()
		}
	}
}
