# python3 dstest --iter 100 --workers 10 --timeout 60 --output out.log TestBasic3A TestSpeed3A TestConcurrent3A TestUnreliable3A TestUnreliableOneKey3A TestOnePartition3A TestManyPartitionsOneClient3A TestManyPartitionsManyClients3A TestPersistOneClient3A TestPersistConcurrent3A TestPersistConcurrentUnreliable3A TestPersistPartition3A TestPersistPartitionUnreliable3A TestPersistPartitionUnreliableLinearizable3A
# python3 dstest --iter 100 --workers 10 --timeout 60 --output out.log TestPersistPartitionUnreliableLinearizable3A
# python3 dstest --iter 100 --workers 10 --timeout 60 --output out.log TestSnapshotRPC3B TestSnapshotSize3B TestSpeed3B TestSnapshotRecover3B TestSnapshotRecoverManyClients3B TestSnapshotUnreliable3B TestSnapshotUnreliableRecover3B TestSnapshotUnreliableRecoverConcurrentPartition3B TestSnapshotUnreliableRecoverConcurrentPartitionLinearizable3B
python3 dstest --iter 100 --workers 10 --timeout 60 --output out.log TestSnapshotRPC3B TestSnapshotSize3B TestSpeed3B
# python3 dstest --iter 100 --workers 10 --timeout 60 --output out.log TestSnapshotRecover3B TestSnapshotRecoverManyClients3B TestSnapshotUnreliable3B TestSnapshotUnreliableRecover3B TestSnapshotUnreliableRecoverConcurrentPartition3B TestSnapshotUnreliableRecoverConcurrentPartitionLinearizable3B
# python3 dstest --iter 100 --workers 10 --timeout 40 --output out.log TestSnapshotRecover3B