package render

// Bufferer should be implemented by all renderable components.
type Bufferer interface {
	Buffer() Buffer
}
