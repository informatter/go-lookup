
# Go Lookup 

A custom hash table implementation written in Golang.

## Resources

**Go Language Reference**
- [see](https://go.dev/ref/spec)

**Go Testing Package**
- [see](https://pkg.go.dev/testing)
- [flags](https://golang.org/cmd/go/#hdr-Testing_flags.)

**FNV Hash Function**
The base hash function being used is called "Fowler–Noll–Vo" (FNV). Its a non-cryptographic function which is much faster than cryptographic ones to compute.

- [fnv](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function)

** Prime Numbers**
- [see](https://en.wikipedia.org/wiki/Prime_number)

## Unit Tests

Run all tests:
`go test .`

Run all tests and all benchmarks:

`go test -bench . -benchmem`

To benchmark a specific test:

`go test  -bench=<benchMarkName> -benchmem`

Optional:
- Add `-benchtime t` for iterations or specified time.
- for example, `-benchtime 1h30s` or `-benchtime 5s`
- `-benchtime 100x` to run 100 times

Run a specific test:

`go test --run <test-name>`


## Formatting

`go fmt .`