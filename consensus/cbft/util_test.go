package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/deckarep/golang-set"
	"math/big"
	"testing"
)

func TestRandomOffset(t *testing.T) {
	var expect = []struct{
		end 	 int
		wanted   bool
	}{
		{end: 10, wanted: true},
		{end: 0,  wanted: true},
	}
	for _, data := range expect {
		res := randomOffset(data.end)
		if !data.wanted && (0 > res || res > data.end) {
			t.Errorf("randomOffset has incorrect value. result:{%v}", res)
		}
	}
}

func TestRandomOffset_Collision(t *testing.T) {
	vals := make(map[int]struct{})
	for i := 0; i < 100; i++ {
		offset := randomOffset(2 << 30)
		if _, ok := vals[offset]; ok {
			t.Fatalf("got collision")
		}
		vals[offset] = struct{}{}
	}
}

func TestRandomOffset_Zero(t *testing.T) {
	offset := randomOffset(0)
	if offset != 0 {
		t.Fatalf("bad offset")
	}
}

func TestProduceHash(t *testing.T) {
	var hashes mapset.Set
	hashes = mapset.NewSet()
	hashes.Add(common.BigToHash(big.NewInt(10)))
	hashes.Add(common.BigToHash(big.NewInt(11)))
	hashes.Add(common.BigToHash(big.NewInt(12)))

	var con interface{} = common.BigToHash(big.NewInt(10))
	if hashes.Contains(con) {
		t.Log("exists")
	} else {
		t.Error("not exists")
	}
}

func TestUint64ToBytes(t *testing.T) {
	var wants = []struct{
		src uint64
		want string
	}{
		{
			src: 1558679713,
			want: "5fcb2251f5b31c73534c57718f0d60b23bc99898a0c4c4e69ae97b4a09f17205",
		},
		{
			src: 1558679714,
			want: "2b83fe25cd31f504192d7e9fa725f8d4d334d724feaf35cd26d225050b825683",
		},
		{
			src: 1558679715,
			want: "8ce02e5594c6da16f9c6d3958119c7ed6f0d25d3b45aeec10bcfdfe258aaf83f",
		},
	}
	for _, v := range wants {
		target := sha3.Sum256(uint64ToBytes(v.src))
		t_hex := common.Bytes2Hex(target[:])
		if t_hex != v.want {
			t.Error("error")
		}
	}
}