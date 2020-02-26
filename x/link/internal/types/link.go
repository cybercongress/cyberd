package types

import (
	"encoding/binary"
	. "github.com/cybercongress/go-cyber/types"
)

type Link struct {
	From Cid `json:"from"`
	To   Cid `json:"to"`
}

type CompactLink struct {
	from CidNumber
	to   CidNumber
	acc  AccNumber
}

func NewLink(from CidNumber, to CidNumber, acc AccNumber) CompactLink {
	return CompactLink{from: from, to: to, acc: acc}
}

func (l CompactLink) From() CidNumber { return l.from }
func (l CompactLink) To() CidNumber   { return l.to }
func (l CompactLink) Acc() AccNumber  { return l.acc }

func UnmarshalBinaryLink(b []byte) CompactLink {
	return NewLink(
		CidNumber(binary.LittleEndian.Uint64(b[0:8])),
		CidNumber(binary.LittleEndian.Uint64(b[8:16])),
		AccNumber(binary.LittleEndian.Uint64(b[16:24])),
	)
}

func (l CompactLink) MarshalBinary() []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:8], uint64(l.From()))
	binary.LittleEndian.PutUint64(b[8:16], uint64(l.To()))
	binary.LittleEndian.PutUint64(b[16:24], uint64(l.Acc()))
	return b
}
