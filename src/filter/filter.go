package filter

type Filter int32

func NewBloomFilter(bitsPerKey int) Filter {
	return Filter(bitsPerKey)
}