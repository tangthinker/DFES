package id_generator

import (
	"fmt"
	"sync/atomic"
)

type SequenceIdGenerator struct {
	seq    atomic.Uint64
	Prefix string
}

func NewSequenceIdGenerator(prefix string) *SequenceIdGenerator {
	return &SequenceIdGenerator{
		Prefix: prefix,
	}
}

func (s *SequenceIdGenerator) Next() string {
	var ok = false
	var seq uint64
	for !ok {
		seq = s.seq.Load()
		ok = s.seq.CompareAndSwap(seq, seq+1)
	}
	result := fmt.Sprintf("%s%020d", s.Prefix, seq)
	return result
}
