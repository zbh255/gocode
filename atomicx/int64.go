package atomicx

import "sync/atomic"

type Int64 struct {
	atomic.Int64
}

func (x *Int64) GtAndSwap(v int64) {
	old2 := x.Load()
	for old2 < v {
		rt_procPin()
		x.Store(v)
		old2 = x.Load()
		if old2 == v {
			rt_procUnpin()
			break
		}
		// pause()
		rt_procUnpin()
	}
}

func (x *Int64) GteAndSwap(v int64) {
	old2 := x.Load()
	for old2 <= v {
		rt_procPin()
		x.Store(v)
		old2 = x.Load()
		if old2 == v {
			rt_procUnpin()
			break
		}
		// pause()
		rt_procUnpin()
	}
}

func (x *Int64) GlAndSwap(v int64) {
	old2 := x.Load()
	for old2 > v {
		rt_procPin()
		x.Store(v)
		old2 = x.Load()
		if old2 == v {
			rt_procUnpin()
			break
		}
		// pause()
		rt_procUnpin()
	}
}

func (x *Int64) GleAndSwap(v int64) {
	old2 := x.Load()
	for old2 >= v {
		rt_procPin()
		x.Store(v)
		old2 = x.Load()
		if old2 == v {
			rt_procUnpin()
			break
		}
		// pause()
		rt_procUnpin()
	}
}
