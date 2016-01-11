[![Build Status](https://travis-ci.org/funny/slab.svg)](https://travis-ci.org/funny/slab)
[![Coverage Status](https://coveralls.io/repos/funny/slab/badge.svg?branch=master&service=github)](https://coveralls.io/github/funny/slab?branch=master)

Introduction
============

Slab allocation memory pools for Go.

Usage
=====

Use lock-free memory pool:
```go
pool := slab.NewAtomPool(
	64,          // The smallest chunk size is 64B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
	1024 * 1024, // Each slab will be 1MB in size.
)

buf1 := pool.Alloc(64)

    ... use the buf ...
	
pool.Free(buf)
```

Use `sync.Pool` based memory pool:
```go
pool := slab.NewSyncPool(
	64,          // The smallest chunk size is 64B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
)

buf := pool.Alloc(64)

    ... use the buf ...
	
pool.Free(buf)
```

Performance
===========

Result of `GOMAXPROCS=16 go test -v -bench=. -benchmem`:

```
Benchmark_AtomPool_AllocAndFree_128-16	10000000	       185 ns/op	       0 B/op	       0 allocs/op
Benchmark_AtomPool_AllocAndFree_256-16	10000000	       181 ns/op	       0 B/op	       0 allocs/op
Benchmark_AtomPool_AllocAndFree_512-16	10000000	       176 ns/op	       0 B/op	       0 allocs/op

Benchmark_SyncPool_AllocAndFree_128-16	20000000	        92.6 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_AllocAndFree_256-16	20000000	        76.6 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_AllocAndFree_512-16	20000000	        81.3 ns/op	      32 B/op	       1 allocs/op

Benchmark_Make_128-16                 	20000000	        52.9 ns/op	     128 B/op	       1 allocs/op
Benchmark_Make_256-16                 	20000000	        97.1 ns/op	     256 B/op	       1 allocs/op
Benchmark_Make_512-16                 	10000000	       128 ns/op	     512 B/op	       1 allocs/op
```
