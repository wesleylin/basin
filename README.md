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
    // instantiate orderedmap
    zoo := orderedmap.New[string, Animal]()

	// chain set calls
	zoo = zoo.Set("kyle", Animal{"Kyle", "Kangaroo"}).
		Set("sam", Animal{"Sam", "Tiger"}).Set("leo", Animal{"Leo", "Tiger"})


    // convert to a stream
    // alternatively can directly fetch as iter.Seq2
    // with zoo.All()
	zooStream := zoo.Stream2()

    // (delayed) filter and then return same stream
	zooStream = zooStream.Filter(func(k string, a Animal) bool {
		return a.Type == "Tiger" || a.Type == "Lion"
	})

    // (delayed) convert animal to string
	zooStream = stream.Map2(zooStream, func(k string, a Animal) (string, Animal) {
		a.Name = a.Name + " the Great"
		return k, a
	})

    // collect results
    // we can also check if there were any errors
	results, err := zooStream.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

    // []string{"Sam the Great", "Leo the Great"}
    fmt.Println(results)
}
```

We can easily create an orderedmap. And then in order loop through the results by converting to Stream which is a lightweight wrapper over iter.Seq2. Main benefit of the wrapper is we can check if any errors happened once Collect() is attempted.

```

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

Above are the general methods. There are two variants the regular one that returns a bool on success and the fluent variant that can be chained as in `om := orderedMap.Set(1, "one").Set(2, "two").Remove(1)` or on multiple lines as

```
om := orderedMap.Set(1, "one")
om = om.Set(2, "two")
om = om.Remove(1)
```

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

The map, heap, and set datastructures can then be iterated over with 2 ways. Use the normal .Values() .All() if you want the normal golang 1.23+ iterators. Use the Stream variants if you want to use the wrapped Stream().

The main advantage of the Stream() is that it is easier to chain, and keeps the error until the last Collect.

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

| benchmark                 | iterations    | time per op        | memory        | allocs |
| ------------------------- | ------------- | ------------------ | ------------- | ------ |
| **Put (Basin)** - 100     | $214,558,254$ | $5.891\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Put (StdMap)** - 100    | $217,225,171$ | $5.704\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Put (Basin)** - 1,000   | $202,188,237$ | $5.927\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Put (StdMap)** - 1,000  | $207,942,988$ | $6.135\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Put (Basin)** - 10,000  | $159,800,976$ | $7.499\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Put (StdMap)** - 10,000 | $166,129,954$ | $7.204\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Get (Basin)**           | $227,376,798$ | $5.114\\text{ ns}$ | $0\\text{ B}$ | $0$    |
| **Get (StdMap)**          | $241,516,156$ | $5.088\\text{ ns}$ | $0\\text{ B}$ | $0$    |

The put and get operations were around the same time.

| benchmark              | iterations   | time per op         | type                            |
| ---------------------- | ------------ | ------------------- | ------------------------------- |
| **Iterate (Basin)**    | $400,747$    | $3,377\\text{ ns}$  | **$\\approx 20\\times$ Faster** |
| **Iterate (StdMap)**   | $19,369$     | $67,479\\text{ ns}$ | Baseline                        |
| **Before Compact**     | $2,373,374$  | $524.0\\text{ ns}$  | Baseline                        |
| **After Compact**      | $3,816,990$  | $313.9\\text{ ns}$  | **$+40\\%$ Efficiency**         |
| **Basin (No Compact)** | $49,238,680$ | $25.08\\text{ ns}$  | Baseline                        |
| **StdMap (Memory)**    | $60,622,773$ | $19.64\\text{ ns}$  | $\\approx 21\\%$ faster         |

## todo items

    [x] add more iterator convenience methods
    [ ] wrap existing concurrent unordered hashmap

    [ ]compaction optimization segmented shards (probably 7 day effort)
