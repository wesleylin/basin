# Run with 1, 4, and 8 cores

go test -bench=. -cpu=1,4,8 ./concurrentmap
