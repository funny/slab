[![Build Status](https://travis-ci.org/funny/slab.svg)](https://travis-ci.org/funny/slab)
[![Coverage Status](https://coveralls.io/repos/funny/slab/badge.svg?branch=master&service=github)](https://coveralls.io/github/funny/slab?branch=master)

Introduction
============

A lock-free slab allocation memory pool for Go.

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

Compare with `sync.Pool` and `make([]byte, n)` when `GOMAXPROCS=16`:

```
Benchmark_Slab_AllocAndFree_128-16 	10000000	       182 ns/op	       0 B/op	       0 allocs/op
Benchmark_Slab_AllocAndFree_256-16 	10000000	       183 ns/op	       0 B/op	       0 allocs/op
Benchmark_Slab_AllocAndFree_512-16 	10000000	       183 ns/op	       0 B/op	       0 allocs/op

Benchmark_SyncPool_GetAndPut_128-16	20000000	        60.9 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_GetAndPut_256-16	20000000	        62.1 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_GetAndPut_512-16	20000000	        61.7 ns/op	      32 B/op	       1 allocs/op

Benchmark_Make_128-16              	30000000	        53.5 ns/op	     128 B/op	       1 allocs/op
Benchmark_Make_256-16              	20000000	        82.6 ns/op	     256 B/op	       1 allocs/op
Benchmark_Make_512-16              	10000000	       139 ns/op	     512 B/op	       1 allocs/op
```
