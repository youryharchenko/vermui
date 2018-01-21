package render

import "sync"

// Bufferer should be implemented by all renderable components.
type Bufferer interface {
	sync.Locker
	Buffer() Buffer
}
