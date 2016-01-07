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