package donkeydb

const (
	kvHeaderBytes     = 3
	kvHeaderTypeBits  = 4
	kvHeaderKeyBits   = 8
	kvHeaderValueBits = kvHeaderBytes*8 - kvHeaderTypeBits - kvHeaderKeyBits

	kvTypeSentinel = 0
	kvTypeString   = 1
)

type kvHeader int

func (h kvHeader) Type() int {
	return int(h) >> (kvHeaderKeyBits + kvHeaderValueBits)
}

func (h kvHeader) Key() int {
	return (int(h) >> kvHeaderValueBits) & ((1 << kvHeaderKeyBits) - 1)
}

func (h kvHeader) Value() int {
	return int(h) & ((1 << kvHeaderValueBits) - 1)
}

func kvh(t, k, v int) kvHeader {
	value := (((t << kvHeaderKeyBits) | k) << kvHeaderValueBits) | v
	return kvHeader(value)
}
