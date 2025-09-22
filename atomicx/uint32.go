package atomicx

import "sync/atomic"

type Uint32 struct {
	atomic.Uint32
}

func (x *Uint32) GtAndSwap(v uint32) {
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

func (x *Uint32) GteAndSwap(v uint32) {
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

func (x *Uint32) GlAndSwap(v uint32) {
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

func (x *Uint32) GleAndSwap(v uint32) {
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
