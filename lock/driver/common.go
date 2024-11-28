package driver

type LockDriver interface {
	Acquire(wait bool) bool
	Release() error
}
