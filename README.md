# basin

Basin is a high-performance data structure library for Go, engineered for the post-iterator era. By leveraging Generics (1.18+) and Native Iterators (1.23+), it eliminates the "Interface Tax"—the runtime cost of type assertions and heap escapes.

Building on this foundation, Basin provides a fluent API that is strictly Type-Safe and Zero-Allocation—delivering lazy evaluation without the overhead of interface{} or extra heap allocations, maintaining the raw speed of Go.

| Feature     | Other Go Libs            | Basin                 |
| ----------- | ------------------------ | --------------------- |
| **Types**   | `interface{}` (Boxing)   | **Generics (Native)** |
| **Memory**  | Eager (Allocates Slices) | **Lazy (Zero-Alloc)** |
| **Syntax**  | Procedural / Nested      | **Fluent / Chained**  |
| **Runtime** | Reflect/Assertions       | **Compile-time Safe** |

## data structures

OrderedMap with not only generics but also new iterator based looping.

For example:

```
package main

type Animal struct {
    Name string
    Species string
}

func main() {
    zoo := basin.NewOrderedMap[string, Animal]()

    zoo.Put("A01", Animal{"Tony", "Tiger"})
       .Put("A02", Animal{"Leo", "Lion"})
       .Put("A03", Animal{"Shere Khan", "Tiger"})
       .Put("A04", Animal{"Baloo", "Bear"})

    // Goal: Find the names of the first 2 Tigers
    tigerNames := zoo.Stream().
        Filter(func(a Animal) bool {
            return a.Species == "Tiger"
        }).
        Take(2).
        Collect()
}
```

Older libraries while pretty robust are usually still using slices requiring you to nest the function with

```
zoo.Range(func(k, v)...)
```

OrderedMap is stable and updates will keep the key in place. For example if it has

```
{
    "dog": 7,
    "cat": 8
    "zebra": 5
}
```

and we update the cat to value 90, the place in the hashmap will not be moved to the end.

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

## iter with 3rd party streams

You can use other existing stream libraries for example

```
import "github.com/samber/lo"

zoo := basin.NewOrderedMap[string, Animal]()

...

// Old way convert iter back to list and then filter
result := lo.Filter(slices.Collect(zoo.Values()), func(v Animal, _ int) bool {
    return v.Type == "Tiger"
})


// This is the Go 1.23+ way with lo zoo.Values() returns iter.Seq[Animal]
// lo.Filter returns a slice, but it consumes the iterator lazily
result := lo.Filter(lo.FromPairs(zoo.All()), func(item lo.Entry[string, Animal], _ int) bool {
    return item.Value.Type == "Tiger"
})
```

## development

go test ./...

go test ./orderedmap

go test -bench=. -benchmem

go test -bench=. -benchmem > bench_results.txt

## todo items

    [x] add more iterator convenience methods
    [ ] wrap existing concurrent unordered hashmap

    [ ]compaction optimization segmented shards (probably 7 day effort)
