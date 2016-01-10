package slab

import (
	"testing"

	"github.com/funny/utest"
)

func Test_Sync_AllocSmall(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(64)
	utest.EqualNow(t, len(mem), 64)
	utest.EqualNow(t, cap(mem), 128)
	pool.Free(mem)
}

func Test_Sync_AllocLarge(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(2048)
	utest.EqualNow(t, len(mem), 2048)
	utest.EqualNow(t, cap(mem), 2048)
	pool.Free(mem)
}
