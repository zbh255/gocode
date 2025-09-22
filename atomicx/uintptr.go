package atomicx

import "sync/atomic"

type Uintptr struct {
	atomic.Uintptr
}

func (x *Uintptr) GtAndSwap(v uintptr) {
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

func (x *Uintptr) GteAndSwap(v uintptr) {
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

func (x *Uintptr) GlAndSwap(v uintptr) {
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

func (x *Uintptr) GleAndSwap(v uintptr) {
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
