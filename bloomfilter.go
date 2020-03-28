package bloomfilter

import (
	"hash"
	"hash/fnv"

	ba "github.com/Workiva/go-datastructures/bitarray"
)

type BloomFilter struct {
	m        uint64 // Size of the bloom filter
	k        uint64 // Number of hash functions
	bitarray ba.BitArray
	hashfn1  hash.Hash64 // The first hash function
	hashfn2  hash.Hash64 // The second hash function
}

func NewBloomFilter(bfSize, numHashFuncs uint64) *BloomFilter {
	bf := new(BloomFilter)
	bf.m = bfSize
	bf.k = numHashFuncs
	bf.bitarray = ba.NewBitArray(bfSize)
	bf.hashfn1 = fnv.New64()
	bf.hashfn2 = fnv.New64a()
	return bf
}

func NewBloomFilterFromBytes(bs []byte, numHashFuncs uint64) (*BloomFilter, error) {
	ba, err := ba.Unmarshal(bs)
	if err != nil {
		return nil, err
	}
	bf := new(BloomFilter)
	bf.m = ba.Capacity()
	bf.k = numHashFuncs
	bf.bitarray = ba
	bf.hashfn1 = fnv.New64()
	bf.hashfn2 = fnv.New64a()
	return bf, nil
}

func (bf *BloomFilter) getHash(b []byte) (uint64, uint64) {
	bf.hashfn1.Reset()
	bf.hashfn1.Write(b)
	h1 := bf.hashfn1.Sum64()

	bf.hashfn2.Reset()
	bf.hashfn2.Write(b)
	h2 := bf.hashfn2.Sum64()

	return h1, h2
}

func (bf *BloomFilter) Add(b []byte) {
	h1, h2 := bf.getHash(b)
	for i := uint64(0); i < bf.k; i++ {
		ind := (h1 + i*h2) % bf.m
		bf.bitarray.SetBit(ind)
	}
}

func (bf *BloomFilter) Check(b []byte) bool {
	h1, h2 := bf.getHash(b)
	res := true
	for i := uint64(0); i < bf.k; i++ {
		ind := (h1 + i*h2) % bf.m
		r, _ := bf.bitarray.GetBit(ind) // ignore the error
		res = res && r
	}
	return res
}

func (bf *BloomFilter) Dump() ([]byte, error) {
	return ba.Marshal(bf.bitarray)
}
