
# Go Lookup

**Go Lookup** is a custom hash table implementation in Go, and I am currently developing it for recreational purposes. I have recorded myself during the development process, feel free to check it out! 😎:

- [Part 1](https://youtu.be/QQLxtAdRw4w)
- [Part 2](https://youtu.be/q2v7FqwIwHI)
- [Part 3](https://youtu.be/Cyf0M_f3FjE)
- [Part 4](https://youtu.be/YH2FcnruuIE)


## Features

- **Custom hash table** implementation.
- Uses **FNV-1a** (Fowler–Noll–Vo) non-cryptographic hash function for fast lookups.
- Uses **Open Addressing** as a collision resolution technique
- Uses **Double Hashing** as a probing technique
- Uses Prime numbers for hash table sizing to reduce collisions.

---
## Open-Addressing with Tetrahedral Double Hashing

The implementation uses an open-addressed hash table using a double hashing strategy that incorporates tetrahedral numbers.

In an **open-addressed hash table**, all key-value pairs are stored directly in a single array. I  refer to each location in this array as a **slot**. When inserting or looking up a key, a hash function is used to determine its home slot.

For each key, we a hash is computed:

```go
hash(key) % tableSize
```

This gives the initial slot where the key should be placed. If that slot is empty, the key is inserted there. However, if the slot is already occupied ( a collision occurred) then a new slot need to be searched using a **probe sequence**.

### Probing with Double Hashing

**Double hashing** is used as a probing stretegy. Rather than stepping through the array linearly or quadratically, a second hash to compute a more varied step size is used to help avoid clustering.

The first and second hash functions are calculated as follows.

```go
hash1 = hash(key) % tableSize
hash2 = 1 + (hash(key) % (tableSize - 1))
```

This second hash ensures that the offset is always non-zero and coprime with the table size (which is always a prime number), allowing to eventually probe every slot in the table.

### Tetrahedral Offset

I implemented a tetrahedral offset to the double hashing function, but have not yet done rigorous benchmarks to prove its efficiency, to be honest, I just decided to add it as it provides a more "spread" out slots on a collision chain.

The tetrahedral number is defined as:

```go
T(n) = (n³ - n) / 6
```

This number grows slowly at first, then more rapidly as the number of collisions increases, and is used to further spread out probe indices in a non-linear way.

The final probe location is calculated as:

```go
index = (hash1 + c * hash2 + T(c)) % tableSize
```

Where `c` is the number of collisions encountered so far while probing.

This approach improves dispersion, especially for keys that would otherwise follow similar probe patterns, further reducing the chance of clustering.

### Example

Suppose we have a hash table with 17 slots and want to insert a key. Let’s say `hash(key) = 12345678`.

We calculate:

```
hash1 = 12345678 % 17 = 4
hash2 = 1 + (12345678 % 16) = 15
```

We then try to insert:

- **Attempt 0**: (4 + 0×15 + T(0)) % 17 = 4
- **Attempt 1**: (4 + 1×15 + T(1)) % 17 = (4 + 15 + 0) % 17 = 19 % 17 = 2
- **Attempt 2**: (4 + 2×15 + T(2)) % 17 = (4 + 30 + 1) % 17 = 35 % 17 = 1
- **Attempt 3**: (4 + 3×15 + T(3)) % 17 = (4 + 45 + 4) % 17 = 53 % 17 = 2 → already tried
- Continue probing...

If a slot is taken, the collision count is incremented and a new index is calculated until an empty slot is found.

### Search/Deletion

Searching and deletion follows the exact same probe sequence used during insertion. If a  key is being searched:
- Start at the computed `hash1`.
- If the slot doesn't match the key, we go to the next probe location using the same formula.
- Continue until the key is found or hit an empty slot (which means the key is not present).

### Resizing and Load Factors

This hash table keeps track of how full it is. As more elements are inserted, the average number of collisions and probe length increases. To maintain performance, it resizes the backing array when the **load factor** reaches 60% (this is currently an arbitrary chocie)

When resizing:
- The new size is the next appropriate **prime number** larger than twice the current size.
- Every key in the old table is reinserted using the new hash parameters.
- If the new size exceeds Go's `uint64` panic is raised.

---

## Tests

**Run all tests:**
```bash
go test .
```

**Format code:**
```bash
go fmt .
```

**Run all tests and benchmarks:**
```bash
go test -bench . -benchmem
```

**Benchmark a specific test:**
```bash
go test -bench=<test> -benchmem
```

**Benchmark multiple tests (regex):**
```bash
go test -bench='<test>|<test>' -benchmem
```

*Useful flags:*
- Set iterations/time: `-benchtime 2s` or `-benchtime 100x`
- Run specific test only: `go test -run `

## Benchmarks

### System

- **CPU:** Intel® Core™ i7-8750H @ 2.20GHz
- **Cores:** 12
- **Arch:** amd64
- **OS:** Windows

## FNV-1a Custom Implementation

| **Keys**   | **Runs** | **ns/op**   | **ns/op per key** | **B/op** | **Allocs/op** |
|------------|----------|-------------|-------------------|----------|---------------|
| 1,000      | 100      | 11,288      | **11.29**         | 0        | 0             |
| 50,000     | 100      | 639,210     | **12.78**         | 0        | 0             |
| 1,000,000  | 100      | 12,491,896  | **12.49**         | 0        | 0             |

## Go Stdlib `hash/fnv`

| **Keys**   | **Runs** | **ns/op**   | **ns/op per key** | **B/op** | **Allocs/op** |
|------------|----------|-------------|-------------------|----------|---------------|
| 1,000      | 100      | 30,330      | **30.33**         | 0        | 0             |
| 50,000     | 100      | 1,452,023   | **29.04**         | 0        | 0             |
| 1,000,000  | 100      | 29,547,213  | **29.55**         | 0        | 0             |

---


## Insert No Resize

### Go map

| **Benchmark**       | **Runs** | **Elements** |   **ns/op**  | **ns/op per key** | **B/op** | **Allocs/op** |
|---------------------|--------:|------------:|------------:|------------------:|--------:|-------------:|
| Insert No Resize    |     500 |   1,000,000 | 124,873,508 |            124.87 |   400,048 | 14,000 |

### Custom hash table

| **Benchmark**       | **Runs** | **Elements** |   **ns/op**  | **ns/op per key** | **B/op** | **Allocs/op** |
|---------------------|--------:|------------:|------------:|------------------:|--------:|-------------:|
| Insert No Resize    |     500 |   1,000,000 | 174,580,206 |            174.58 |   400,045 | 14,000 |



## Search existing key

### Go map

| **Benchmark**          | **Runs** | **ns/op** | **B/op** | **Allocs/op** |
|------------------------|--------:|---------:|--------:|-------------:|
| Search existing key    |     500 |    11.40 |       0 |            0 |

### Custom hash table

| **Benchmark**          | **Runs** | **ns/op** | **B/op** | **Allocs/op** |
|------------------------|--------:|---------:|--------:|-------------:|
| Search existing key    |     500 |    63.80 |      16 |            1 |



## Search non-existing key

### Go map

| **Benchmark**             | **Runs** | **ns/op** | **B/op** | **Allocs/op** |
|---------------------------|--------:|---------:|--------:|-------------:|
| Search non-existing key   |     500 |    12.00 |       0 |            0 |

### Custom hash table

| **Benchmark**             | **Runs** | **ns/op** | **B/op** | **Allocs/op** |
|---------------------------|--------:|---------:|--------:|-------------:|
| Search non-existing key   |     500 |    75.20 |      16 |            1 |

**Legend**

- **ns/op**: average number of nanoseconds per benchmark operation. For insert benchmarks the operation is the whole benchmark run (inserting all elements); for search benchmarks the operation is a single lookup.
- **ns/op per key/lookup**: average number of nanoseconds per individual element operation — for inserts this is `ns/op / Elements`, for searches this equals `ns/op` because each benchmark operation performs one lookup.
- **B/op**: average bytes allocated per benchmark operation.
- **Allocs/op**: average number of memory allocations performed per benchmark operation.



## References

- **Go Language Reference:** [Go Spec][1]
- **Go Testing Package:** [testing docs][2], [testing flags][3]
- **FNV Hash Function:** [Wikipedia][4]
- **Prime Numbers:** [Wikipedia][5]
- **Open Addressing:** [Wikipedia][6]
- **Double Hashing:** [Wikipedia][7]

[1]: https://go.dev/ref/spec
[2]: https://pkg.go.dev/testing
[3]: https://golang.org/cmd/go/#hdr-Testing_flags.
[4]: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
[5]: https://en.wikipedia.org/wiki/Prime_number
[6]: https://en.wikipedia.org/wiki/Open_addressing
[7]: https://en.wikipedia.org/wiki/Double_hashing
