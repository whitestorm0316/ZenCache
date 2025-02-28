package cache

// 只读结构 防止内存修改
type ByteView struct {
	bytes []byte
}

func (bv ByteView) Len() int {
	return len(bv.bytes)
}
func NewByteView(b []byte) ByteView {
	return ByteView{cloneBytes(b)}
}
func (bv ByteView) ByteSlices() []byte {
	return cloneBytes(bv.bytes)
}
func (bv ByteView) String() string {
	return string(bv.bytes)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
