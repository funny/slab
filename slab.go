package slab

type Pool interface {
	Alloc(int) []byte
	Free([]byte)
}

var _ = Pool((*ChanPool)(nil))
var _ = Pool((*SyncPool)(nil))
