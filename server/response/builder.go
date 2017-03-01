// Package builder provides interfaces and mechanisms to make client pulls very easy.
package builder

// Serializable provides functions for serializing into json
type Serializable interface {
	// Data returns a representation of the data that can be converted into json
	Data() interface{}
}
