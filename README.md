## data structures

OrderedMap note this is insertion order

OrderedSet

Naively benchmarked can be 5x faster, but uses twice the memory for keys. will test more

UnorderedMap not included yet as one can directly use built-in map

Heap (priorityQueue) is not stable

## development

go test ./map

go test -bench=. -benchmem
