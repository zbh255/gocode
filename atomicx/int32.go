package atomicx

import "sync/atomic"

type Int32 struct {
	atomic.Int32
}

func (x *Int32) GtAndSwap(v int32) {
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

func (x *Int32) GteAndSwap(v int32) {
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

func (x *Int32) GlAndSwap(v int32) {
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

func (x *Int32) GleAndSwap(v int32) {
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
