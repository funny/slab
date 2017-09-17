package slab

import "testing"

func Test_NoSlabPool_Alloc(t *testing.T) {
	pool := new(NoPool)
	pool.Alloc(64)
	pool.Free(make([]byte, 64))
}
