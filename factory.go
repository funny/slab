package slab

import "github.com/pkg/errors"

const (
	// 10: sync pool, 20: chan pool, 30: atom pool
	TYPE__SLAB_POOL__SYNC = 10
	TYPE__SLAB_POOL__CHAN = 20
	TYPE__SLAB_POOL__ATOM = 30
)

type Slab interface {
	Alloc(size int) []byte
	Free(mem []byte)
	ErrChan() <-chan error
}

func NewSlabPool(typ int16, minSize int, maxSize int, factor int) (Slab, error) {
	var (
		abpool Slab
		err    error
	)
	switch typ {
	case TYPE__SLAB_POOL__SYNC:
		abpool = newSyncPool(minSize, maxSize, factor)
	case TYPE__SLAB_POOL__ATOM:
		abpool = newAtomPool(minSize, maxSize, factor)
	case TYPE__SLAB_POOL__CHAN:
		abpool = newChanPool(minSize, maxSize, factor)
	default:
		err = errors.Errorf("unsupport type: [%d] slab pool!!", typ)
	}
	return abpool, err
}
