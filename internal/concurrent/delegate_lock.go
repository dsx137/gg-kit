package concurrent

type DelegateLock interface {
	Lock()
	UnLock()
	Compute(func())
}

type DelegateRWLock interface {
	DelegateLock
	RLock()
	RUnLock()
	ComputeR(func())
}
