package raft

import (
	"time"
)

// if the peer has not acked in this duration, it's considered inactive.
const activeWindowWidth = 2 * baseElectionTimeout * time.Millisecond

type PeerTracker struct {
	nextIndex  uint64
	matchIndex uint64

	lastAck time.Time
}

func (rf *Raft) resetTrackedIndexes() {
	for i := range rf.peerTrackers {
		rf.peerTrackers[i].nextIndex = rf.log.lastIndex() + 1
		// warning: cannot set the initial match index to the snapshot index since there might be new peers or way too lag-behind peers.
		rf.peerTrackers[i].matchIndex = 0
	}
}

func (rf *Raft) quorumActive() bool {
	activePeers := 1
	for i, tracker := range rf.peerTrackers {
		if i != rf.me && time.Since(tracker.lastAck) <= activeWindowWidth {
			activePeers++
		}
	}
	return 2*activePeers > len(rf.peers)
}
