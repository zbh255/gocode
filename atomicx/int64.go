package atomicx

import "sync/atomic"

type Int64 struct {
	atomic.Int64
}

func (x *Int64) SwapIfGt(v int64) (old int64, swapped bool) {
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

func (x *Int64) SwapIfGte(v int64) (old int64, swapped bool) {
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

func (x *Int64) SwapIfGl(v int64) (old int64, swapped bool) {
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

func (x *Int64) SwapIfGle(v int64) (old int64, swapped bool) {
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
