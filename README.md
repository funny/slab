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

Use `chan` based memory pool:

```go
pool := slab.NewChanPool(
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
Benchmark_AtomPool_AllocAndFree_128-8   	20000000	        99.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_AtomPool_AllocAndFree_256-8   	20000000	       106 ns/op	       0 B/op	       0 allocs/op
Benchmark_AtomPool_AllocAndFree_512-8   	20000000	       109 ns/op	       0 B/op	       0 allocs/op

Benchmark_ChanPool_AllocAndFree_128-8   	 5000000	       290 ns/op	       0 B/op	       0 allocs/op
Benchmark_ChanPool_AllocAndFree_256-8   	 5000000	       308 ns/op	       0 B/op	       0 allocs/op
Benchmark_ChanPool_AllocAndFree_512-8   	 5000000	       304 ns/op	       0 B/op	       0 allocs/op

Benchmark_SyncPool_AllocAndFree_128-8   	100000000	        22.1 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_AllocAndFree_256-8   	100000000	        22.0 ns/op	      32 B/op	       1 allocs/op
Benchmark_SyncPool_AllocAndFree_512-8   	100000000	        23.0 ns/op	      32 B/op	       1 allocs/op

Benchmark_SyncPool_CacheMiss_128-8      	10000000	       208 ns/op	     160 B/op	       2 allocs/op
Benchmark_SyncPool_CacheMiss_256-8      	10000000	       238 ns/op	     288 B/op	       2 allocs/op
Benchmark_SyncPool_CacheMiss_512-8      	 5000000	       256 ns/op	     544 B/op	       2 allocs/op

Benchmark_Make_128-8                    	100000000	        17.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Make_256-8                    	50000000	        31.3 ns/op	     256 B/op	       1 allocs/op
Benchmark_Make_512-8                    	20000000	        62.0 ns/op	     512 B/op	       1 allocs/op
```
