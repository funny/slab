Introduction
============

A lock-free slab allocator for Go.

[![Build Status](https://travis-ci.org/funny/slab.svg)](https://travis-ci.org/funny/slab)
[![Coverage Status](https://coveralls.io/repos/funny/slab/badge.svg?branch=master&service=github)](https://coveralls.io/github/funny/slab?branch=master)

Usage
=====

```go
pool := slab.NewPool(
	64,          // The smallest chunk size is 128B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
	1024 * 1024, // Each slab will be 1MB in size.
)

buf := pool.Alloc(64)
    ... use the buf ...
pool.Free(buf)
```

Performance
===========

```
Benchmark_Alloc128_And_Free-4	20000000	        62.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Alloc256_And_Free-4	20000000	        59.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Alloc512_And_Free-4	20000000	        57.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Make128-4          	20000000	        73.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Make256-4          	20000000	        94.7 ns/op	     256 B/op	       1 allocs/op
Benchmark_Make512-4          	10000000	       146 ns/op	     512 B/op	       1 allocs/op
```
