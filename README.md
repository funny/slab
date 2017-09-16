Introduction
============

[![Go Report Card](https://goreportcard.com/badge/github.com/funny/slab)](https://goreportcard.com/report/github.com/funny/slab)
[![Build Status](https://travis-ci.org/funny/slab.svg?branch=master)](https://travis-ci.org/funny/slab)
[![codecov](https://codecov.io/gh/funny/slab/branch/master/graph/badge.svg)](https://codecov.io/gh/funny/slab)
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg)](https://godoc.org/github.com/funny/slab)

Slab allocation memory pools for Go.

Usage
=====
Slab pool supports sync pool, chan pool, unlock pool and lock pool types.

const (
	// 10: sync pool, 20: chan pool, 30: atom pool; 40: lock pool
	TYPE__SLAB_POOL__SYNC = 10
	TYPE__SLAB_POOL__CHAN = 20
	TYPE__SLAB_POOL__ATOM = 30 // no lock
	TYPE__SLAB_POOL__LOCK = 40 // lock
)

Use lock-free memory pool:

```go
pool := slab.NewSlabPool(
    TYPE__SLAB_POOL__ATOM, // unlock pool type
	64,          // The smallest chunk size is 64B.
	64 * 1024,   // The largest chunk size is 64KB.
	2,           // Power of 2 growth in chunk size.
)

buf1 := pool.Alloc(64)

    ... use the buf ...
	
pool.Free(buf)
```


Performance
===========

Result of `GOMAXPROCS=4 go test -v -bench=. -benchmem`:

```
Benchmark_AtomPool_AllocAndFree_128-4                   10000000               151 ns/op               0 B/op          0 allocs/op
Benchmark_AtomPool_AllocAndFree_256-4                   10000000               145 ns/op               0 B/op          0 allocs/op
Benchmark_AtomPool_AllocAndFree_512-4                   10000000               175 ns/op               0 B/op          0 allocs/op
Benchmark_ChanPool_AllocAndFree_128-4                    3000000               448 ns/op               0 B/op          0 allocs/op
Benchmark_ChanPool_AllocAndFree_256-4                    5000000               459 ns/op               0 B/op          0 allocs/op
Benchmark_ChanPool_AllocAndFree_512-4                    3000000               405 ns/op               0 B/op          0 allocs/op
Benchmark_LockPool_AllocAndFree_128-4                    5000000               315 ns/op               0 B/op          0 allocs/op
Benchmark_LockPool_AllocAndFree_256-4                    5000000               291 ns/op               0 B/op          0 allocs/op
Benchmark_LockPool_AllocAndFree_512-4                    5000000               259 ns/op               0 B/op          0 allocs/op
Benchmark_SyncPool_AllocAndFree_128-4                   50000000                24.7 ns/op            32 B/op          1 allocs/op
Benchmark_SyncPool_AllocAndFree_256-4                   100000000               24.7 ns/op            32 B/op          1 allocs/op
Benchmark_SyncPool_AllocAndFree_512-4                   50000000                26.1 ns/op            32 B/op          1 allocs/op
Benchmark_SyncPool_AllocAndFree_NonIntFactor-4          50000000                27.2 ns/op            32 B/op          1 allocs/op
Benchmark_SyncPool_CacheMiss_128-4                       3000000               387 ns/op             160 B/op          2 allocs/op
Benchmark_SyncPool_CacheMiss_256-4                       5000000               346 ns/op             288 B/op          2 allocs/op
Benchmark_SyncPool_CacheMiss_512-4                       5000000               390 ns/op             544 B/op          2 allocs/op
Benchmark_Make_128-4                                    50000000                29.4 ns/op           128 B/op          1 allocs/op
Benchmark_Make_256-4                                    30000000                53.3 ns/op           256 B/op          1 allocs/op
Benchmark_Make_512-4                                    20000000               105 ns/op             512 B/op          1 allocs/op
```
