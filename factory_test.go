package slab

import (
	"testing"

	"github.com/1046102779/utest"
)

func Test_Factory_Slab_NewInstance(t *testing.T) {
	abpool, _ := NewSlabPool(TYPE__SLAB_POOL__SYNC, 64, 1024, 2, 1024*1024)
	abpool.Free(abpool.Alloc(512))
	utest.NotNilNow(t, abpool)

	abpool2, _ := NewSlabPool(TYPE__SLAB_POOL__CHAN, 64, 1024, 2, 1024*1024)
	abpool2.Free(abpool2.Alloc(512))
	utest.NotNilNow(t, abpool2)

	abpool3, _ := NewSlabPool(TYPE__SLAB_POOL__ATOM, 64, 1024, 2, 1024*1024)
	abpool3.Free(abpool3.Alloc(512))
	utest.NotNilNow(t, abpool3)

	abpool4, _ := NewSlabPool(TYPE__SLAB_POOL__LOCK, 64, 1024, 2, 1024*1024)
	utest.NotNilNow(t, abpool4)

	abpool5, _ := NewSlabPool(5 /**unsupported slab pool type**/, 64, 1024, 2, 1024*1024)
	utest.IsNilNow(t, abpool5)
}
