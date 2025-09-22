package atomicx

import _ "unsafe"

//go:linkname rt_procUnpin runtime.procUnpin
func rt_procUnpin()

//go:linkname rt_procPin runtime.procPin
func rt_procPin()
