package atomicx

import "sync/atomic"

type Uint64 struct {
	atomic.Uint64
}

func (x *Uint64) SwapIfGt(v uint64) (old uint64, swapped bool) {
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

func (x *Uint64) SwapIfGte(v uint64) (old uint64, swapped bool) {
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

func (x *Uint64) SwapIfGl(v uint64) (old uint64, swapped bool) {
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

func (x *Uint64) SwapIfGle(v uint64) (old uint64, swapped bool) {
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
