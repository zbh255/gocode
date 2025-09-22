package atomicx

import "sync/atomic"

type Uintptr struct {
	atomic.Uintptr
}

func (x *Uintptr) SwapIfGt(v uintptr) (old uintptr, swapped bool) {
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

func (x *Uintptr) SwapIfGte(v uintptr) (old uintptr, swapped bool) {
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

func (x *Uintptr) SwapIfGl(v uintptr) (old uintptr, swapped bool) {
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

func (x *Uintptr) SwapIfGle(v uintptr) (old uintptr, swapped bool) {
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
