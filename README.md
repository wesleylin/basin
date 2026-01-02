## data structures

OrderedMap with not only generics but also new iterator based looping.

For example:

```
package main

import "github.com/yourusername/basin"

type Animal struct {
	Name string
	Type string
}

// In Go 1.24, we use [Animal any] for the alias if we want it to be generic,
// but here we are pinning it to the specific 'Animal' struct.
type ZooMap = basin.OrderedMap[string, Animal]

func main() {
	// Initialize the map
	zoo := basin.NewOrderedMap[string, Animal]()

	// Using the Fluent API we designed:
	// 1. We use 'Set' for fluent Map insertion.
	// 2. Struct fields need strings in quotes.
	// 3. Keys (strings) must be passed separately from the Value (Animal).
	zoo.Set("kyle", Animal{"Kyle", "Kangaroo"}).
	    Set("sam", Animal{"Sam", "Tiger"})
}
```

Older libraries while pretty robust are usually still using slices requiring you to nest the function with

```
zoo.Range(func(k, v)...)
```

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

Above are the general methods. There are two variants the regular one that returns a bool on success and the fluent variant that can be chained as in `om := orderedMap.Set(1, "one").Set(2, "two").Remove(1)`

| **Data Stucture** | Operation    | Core Method       | Returns    | Fluent Method |
| ----------------- | ------------ | ----------------- | ---------- | ------------- |
| **Set**           | `All()`      | `iter.Seq[T]`     | `v T`      | Lazy          |
| -                 | `Query()`    | `Stream[T]`       | —          | Lazy (Fluent) |
| **OrderedMap**    | `All()`      | `iter.Seq2[K, V]` | `k K, v V` | Lazy          |
| -                 | `Keys()`     | `iter.Seq[K]`     | `k K`      | Lazy          |
| -                 | `Values()`   | `iter.Seq[V]`     | `v V`      | Lazy          |
| -                 | `Query()`    | `Stream[V]`       | —          | Lazy (Fluent) |
| **Heap**          | `Drain()`    | `iter.Seq[V]`     | `v V`      | **Consuming** |
| -                 | `Query()`    | `Stream[V]`       | —          | Lazy (Fluent) |
| **Stream**        | `Filter(fn)` | `Stream[T]`       | —          | Lazy          |
| -                 | `Take(n)`    | `Stream[T]`       | —          | Lazy          |
| -                 | `Collect()`  | `[]T`             | —          | **Terminal**  |

You can choose either way to retrieve values from the maps and sets. Use the normal .Values() .All() if you want the normal golang 1.23+ iterators. Use the .Stream() if you use the wrapped Stream().

The main advantage of the Stream() is that it is easier to chain, but it is slightly slower.

```
// Option A: Standard Go
for v := range m.Values() { ... }

// Option B: Basin Fluent
m.Stream().Filter(fn).Collect()
```

## streams

You can use other existing stream libraries with

```
import "github.com/samber/lo"

zoo := basin.NewOrderedMap[string, Animal]()

...


result := lo.Filter(slices.Collect(zoo.Values()), func(v Animal, _ int) bool {
    return v.Type == "Tiger"
})
```

## development

go test ./map

go test -bench=. -benchmem

## todo items

- add more iterator convenience methods
- wrap existing concurrent unordered hashmap
