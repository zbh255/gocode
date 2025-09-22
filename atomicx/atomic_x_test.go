package atomicx

import (
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
)

func TestUint64(t *testing.T) {
	t.Run("Serial", func(t *testing.T) {
		var u Uint64
		u.CompareAndSwap(0, 1024)
		for i := 1023; i >= 512; i-- {
			u.SwapIfGl(uint64(i))
		}
		require.Equal(t, u.Load(), uint64(512))
	})
}

type mutexUint64 struct {
	sync.RWMutex
	u uint64
}

func (m *mutexUint64) Load() (v uint64) {
	m.RLock()
	v = m.u
	m.RUnlock()
	return
}

func (m *mutexUint64) SwapIfGt(v uint64) (old uint64, swapped bool) {
	m.Lock()
	if m.u < v {
		old = m.u
		swapped = true
		m.u = v
	}
	m.Unlock()
	return
}

func BenchmarkUint64(b *testing.B) {
	b.StopTimer()
	dataSheet := make([]uint64, 1024*1024*16)
	for i := 0; i < len(dataSheet); i++ {
		dataSheet[i] += 12
	}
	b.StartTimer()
	b.Run("SpinSwap", func(b *testing.B) {
		var (
			u  Uint64
			cu Uint64
		)
		b.SetParallelism(runtime.NumCPU())
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				idx := cu.Add(1)
				val := dataSheet[idx%uint64(len(dataSheet))]
				u.SwapIfGt(val)
			}
		})
		b.Log(cu.Load())
	})
	b.Run("MutexSwap", func(b *testing.B) {
		var (
			u  mutexUint64
			cu Uint64
		)
		b.SetParallelism(runtime.NumCPU())
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				idx := cu.Add(1)
				val := dataSheet[idx%uint64(len(dataSheet))]
				u.SwapIfGt(val)
			}
		})
		b.Log(cu.Load())
	})
}
