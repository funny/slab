[![Build Status](https://travis-ci.org/funny/slab.svg)](https://travis-ci.org/funny/slab)
[![Coverage Status](https://coveralls.io/repos/funny/slab/badge.svg?branch=master&service=github)](https://coveralls.io/github/funny/slab?branch=master)

Introduction
============

A lock-free slab allocation memory pool for Go.

Usage
=====

```go
// Create a lock-free pool
pool := slab.NewLockFreePool(
	64,          // The smallest chunk size is 128B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
	1024 * 1024, // Each slab will be 1MB in size.
)

buf1 := pool.Alloc(64)

    ... use the buf ...
	
pool.Free(buf)
```

```go
// Create a sync.Pool based memory pool
pool := slab.NewSyncPool(
	64,          // The smallest chunk size is 128B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
)

buf := pool.Alloc(64)

    ... use the buf ...
	
pool.Free(buf)
```

Performance
===========

When `GOMAXPROCS=16`:

```
Benchmark_LockFree_AllocAndFree_128-8	10000000	       185 ns/op	       0 B/op	       0 allocs/op
Benchmark_LockFree_AllocAndFree_256-8	10000000	       185 ns/op	       0 B/op	       0 allocs/op
Benchmark_LockFree_AllocAndFree_512-8	10000000	       186 ns/op	       0 B/op	       0 allocs/op

Benchmark_Sync_AllocAndFree_128-8    	20000000	        64.9 ns/op	      32 B/op	       1 allocs/op
Benchmark_Sync_AllocAndFree_256-8    	20000000	        66.3 ns/op	      32 B/op	       1 allocs/op
Benchmark_Sync_AllocAndFree_512-8    	20000000	        68.9 ns/op	      32 B/op	       1 allocs/op

Benchmark_Make_128-8                 	30000000	        47.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Make_256-8                 	20000000	        71.7 ns/op	     256 B/op	       1 allocs/op
Benchmark_Make_512-8                 	10000000	       119 ns/op	     512 B/op	       1 allocs/op
```
