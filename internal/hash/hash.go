package hash

import (
	"unsafe"
)

// Maphash returns a high-speed hash for any comparable key.
// It uses the same internal algorithm as Go's native maps.
func Maphash[K comparable](key K) uint64 {
	// We get the internal 'runtime._type' of the key by
	// putting it in an 'any' (interface) and inspecting it.
	var i any = key

	// An 'any' is a 2-word header: [type_pointer, data_pointer]
	header := (*[2]uintptr)(unsafe.Pointer(&i))
	typePtr := header[0] // Points to the internal runtime._type
	dataPtr := header[1] // Points to the actual value

	// We use the typePtr to ensure we are hashing the value
	// correctly according to its specific type.
	return uint64(runtime_typehash(typePtr, unsafe.Pointer(dataPtr), 0))
}

// These link to the private functions inside the Go runtime.
// The go:linkname directive tells the compiler to use the
// actual runtime implementation of these functions.

//go:noescape
//go:linkname runtime_memhash runtime.memhash
func runtime_memhash(p unsafe.Pointer, h, s uintptr) uintptr

// runtime_typehash is the internal function Go uses to hash map keys.
// It takes the type information into account to handle strings,
// structs, and interfaces properly.
//
//go:noescape
//go:linkname runtime_typehash runtime.typehash
func runtime_typehash(t uintptr, p unsafe.Pointer, seed uintptr) uintptr
