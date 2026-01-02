## data structures

OrderedMap note this is insertion order

OrderedSet

Naively benchmarked can be 5x faster, but uses twice the memory for keys. will test more

UnorderedMap not included yet as one can directly use built-in map

Heap (priorityQueue) is not stable

## api

| **Data Stucture** | Operation | Core Method    | Returns     | Fluent Method | returns       |
| ----------------- | --------- | -------------- | ----------- | ------------- | ------------- |
| **Set**           | Insertion | `Insert(v)`    | `bool`      | `Add(v)`      | `*Set[T]`     |
| -                 | Removal   | `Delete(v)`    | `bool`      | `Remove(v)`   | `*Set[T]`     |
| -                 | Query     | `Has(v)`       | `bool`      | —             | —             |
| **OrderedMap**    | Insertion | `Put(k, v)`    | `bool`      | `Set(k, v)`   | `*Map[K, V]`  |
| -                 | Removal   | `Delete(k)`    | `bool`      | `Remove(k)`   | `*Map[K, V]`  |
| -                 | Access    | `Get(k)`       | `(V, bool)` | -             | -             |
| **Heap**          | Insertion | `Insert(v, p)` | `void`      | `Push(v, p)`  | `*Heap[V, P]` |
| -                 | Removal   | `Pop()`        | `(V, bool)` | `Drop()`      | `*Heap[V, P]` |

## development

go test ./map

go test -bench=. -benchmem
