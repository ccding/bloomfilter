package bloomfilter

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(10000, 5)
	d1, d2 := []byte("Hello"), []byte("World")
	bf.Add(d1)
	if !bf.Check(d1) {
		t.Errorf("%s should be in the BloomFilter", d1)
	}
	if bf.Check(d2) {
		t.Errorf("%s shouldn't be in the BloomFilter", d2)
	}
}

func BenchmarkAdd(b *testing.B) {
	bf := NewBloomFilter(1000000, 20)
	d := []byte("Hello")
	for i := 0; i < b.N; i++ {
		bf.Add(d)
	}
}

func BenchmarkCheck(b *testing.B) {
	bf := NewBloomFilter(1000000, 20)
	d1, d2 := []byte("Hello"), []byte("World")
	bf.Add(d1)
	for i := 0; i < b.N; i++ {
		bf.Check(d2)
	}
}
