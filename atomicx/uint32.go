package atomicx

import "sync/atomic"

type Uint32 struct {
	atomic.Uint32
}

func (x *Uint32) SwapIfGt(v uint32) (old uint32, swapped bool) {
	for {
		old = x.Load()
		if old >= v {
			break
		}
		if !x.CompareAndSwap(old, v) {
			pause()
			continue
		} else {
			swapped = true
			break
		}
	}
	return
}

func (x *Uint32) SwapIfGte(v uint32) (old uint32, swapped bool) {
	for {
		old = x.Load()
		if old > v {
			break
		}
		if !x.CompareAndSwap(old, v) {
			pause()
			continue
		} else {
			swapped = true
			break
		}
	}
	return
}

func (x *Uint32) SwapIfGl(v uint32) (old uint32, swapped bool) {
	for {
		old = x.Load()
		if old <= v {
			break
		}
		if !x.CompareAndSwap(old, v) {
			pause()
			continue
		} else {
			swapped = true
			break
		}
	}
	return
}

func (x *Uint32) SwapIfGle(v uint32) (old uint32, swapped bool) {
	for {
		old = x.Load()
		if old < v {
			break
		}
		if !x.CompareAndSwap(old, v) {
			pause()
			continue
		} else {
			swapped = true
			break
		}
	}
	return
}
