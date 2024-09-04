package server

// statusStoreCtxKey is used for storing our StatusStore as a Context value
type statusStoreCtxKey struct{}

// StatusStore is a basic store for a status value over multiple RPC calls
type StatusStore struct {
	status int
}

// GetStatus returns the current value of the status store
func (s *StatusStore) GetStatus() int {
	return s.status
}

// SetStatus sets the value of the status store, and returns this new value
func (s *StatusStore) SetStatus(value int) int {
	s.status = value

	return s.status
}
