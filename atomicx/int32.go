package atomicx

import "sync/atomic"

type Int32 struct {
	atomic.Int32
}

func (x *Int32) SwapIfGt(v int32) (old int32, swapped bool) {
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

func (x *Int32) SwapIfGte(v int32) (old int32, swapped bool) {
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

func (x *Int32) SwapIfGl(v int32) (old int32, swapped bool) {
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

func (x *Int32) SwapIfGle(v int32) (old int32, swapped bool) {
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
