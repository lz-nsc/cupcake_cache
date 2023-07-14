package cupcakecache

type ByteView struct {
	bytes []byte
}

func (bv ByteView) Len() int {
	return len(bv.bytes)
}

func (bv ByteView) ByteSlice() []byte {
	res := make([]byte, bv.Len())
	copy(res, bv.bytes)
	return res
}
func (bv ByteView) String() string {
	return string(bv.bytes)
}
