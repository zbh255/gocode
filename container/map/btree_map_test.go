package _map

import (
	"strconv"
	"testing"
)

func TestBtree(t *testing.T) {
	bt := NewBtreeMap[int, string](16)
	for i := 0; i < 1024*1024; i++ {
		bt.put(i, strconv.Itoa(i))
	}
	bt.del(4097)
	bt.del(4119)
	bt.Range(4096, func(key int, val string) bool {
		if key > 4120 {
			return false
		} else {
			t.Log("Range : ", key, val)
			return true
		}
	})
	t.Log("Len : ", bt.Len())
	t.Log("MaxKey: ", bt.MaxKey())
	t.Log(bt.LoadOk(1024 * 512))
}
