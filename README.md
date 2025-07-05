

## Unit Tests

Run all tests:
`go test .`

Run all tests and all benchmarks:

`go test -bench . -benchmem`

To benchmark a specific test:

`go test  -bench=<benchMarkName> -benchmem`

Run a specific test:

`go test --run <test-name>`


## Formatting

`go fmt .`