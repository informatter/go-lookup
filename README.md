
# Go Lookup

**Go Lookup** is a custom hash table implementation in Go, and I am currently developing it for recreational purposes. I have recorded myself doing the development process and can be seen in my YouTube channel:

- [Part 1](https://youtu.be/QQLxtAdRw4w)
- [Part 2](https://youtu.be/q2v7FqwIwHI)
- [Part 3](https://youtu.be/Cyf0M_f3FjE)
- [Part 4](https://youtu.be/YH2FcnruuIE)

---

## Features

- **Custom hash table** implementation.
- Uses **FNV-1a** (Fowler–Noll–Vo) non-cryptographic hash function for fast lookups.
- Users **Double Hashing** as a collision handling technique
- Uses Prime numbers for hash table sizing to reduce collisions.

---

## Quick Start

**Run all tests:**
```bash
go test .
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

**Format code:**
```bash
go fmt .
```

## Benchmarks

### System
- **CPU:** Intel® Core™ i7-8750H @ 2.20GHz
- **Arch:** amd64
- **OS:** Windows

### FNV-1a Custom Implementation

| **Keys**    | **Runs** | **ns/op**    | **B/op** | **Allocs/op** |
|-------------|----------|--------------|----------|---------------|
| 1,000       | 100      | 11,288       | 0        | 0             |
| 50,000      | 100      | 639,210      | 0        | 0             |
| 1,000,000   | 100      | 12,491,896   | 0        | 0             |

### Go Stdlib `hash/fnv`

| **Keys**    | **Runs** | **ns/op**    | **B/op** | **Allocs/op** |
|-------------|----------|--------------|----------|---------------|
| 1,000       | 100      | 30,330       | 0        | 0             |
| 50,000      | 100      | 1,452,023    | 0        | 0             |
| 1,000,000   | 100      | 29,547,213   | 0        | 0             |

---

## References

- **Go Language Reference:** [Go Spec][1]
- **Go Testing Package:** [testing docs][2], [testing flags][3]
- **FNV Hash Function:** [Wikipedia][4]
- **Prime Numbers:** [Wikipedia][5]
- **Double Hashing:** [Wikipedia][6]

[1]: https://go.dev/ref/spec
[2]: https://pkg.go.dev/testing
[3]: https://golang.org/cmd/go/#hdr-Testing_flags.
[4]: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
[5]: https://en.wikipedia.org/wiki/Prime_number
[6]: https://en.wikipedia.org/wiki/Double_hashing
