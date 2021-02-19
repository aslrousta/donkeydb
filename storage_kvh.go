package donkeydb

const (
	kvHeaderBytes       = 3
	kvHeaderTypeBits    = 4
	kvHeaderKeyBits     = 8
	kvHeaderValueBits   = kvHeaderBytes*8 - kvHeaderTypeBits - kvHeaderKeyBits
	kvHeaderMaxKeyLen   = 1 << kvHeaderKeyBits
	kvHeaderMaxValueLen = (1 << kvHeaderValueBits) - 1

	kvTypeSentinel = 0
	kvTypeString   = 1
)

type kvHeader int

func (h kvHeader) Type() int {
	return int(h) >> (kvHeaderKeyBits + kvHeaderValueBits)
}

func (h kvHeader) Key() int {
	return ((int(h) >> kvHeaderValueBits) & ((1 << kvHeaderKeyBits) - 1)) + 1
}

func (h kvHeader) Value() int {
	return int(h) & ((1 << kvHeaderValueBits) - 1)
}

func kvh(t, k, v int) kvHeader {
	value := (((t << kvHeaderKeyBits) | (k - 1)) << kvHeaderValueBits) | v
	return kvHeader(value)
}
