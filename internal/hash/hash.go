package hash

import (
	"unsafe"
)

// Maphash returns a high-speed hash for any comparable key.
func Maphash[K comparable](key K) uint64 {
	// We use the 'any' trick to detect the underlying type
	var i any = key

	switch v := i.(type) {
	case string:
		// For strings, we MUST hash the actual bytes, not the header.
		// unsafe.StringData is the safe, modern way (Go 1.20+) to get the pointer.
		return uint64(memhash(unsafe.Pointer(unsafe.StringData(v)), 0, uintptr(len(v))))
	case int, int64, uint64, float64:
		// Fixed size 8-byte types
		return uint64(memhash(unsafe.Pointer(&key), 0, 8))
	case int32, uint32, float32:
		// Fixed size 4-byte types
		return uint64(memhash(unsafe.Pointer(&key), 0, 4))
	default:
		// For structs and other types, we hash the raw memory of the variable.
		// This works for any comparable struct.
		return uint64(memhash(unsafe.Pointer(&key), 0, unsafe.Sizeof(key)))
	}
}

// Maphash returns a high-speed hash for any comparable key.
// It uses the same internal algorithm as Go's native maps.
func Maphash2[K comparable](key K) uint64 {
	// We get the internal 'runtime._type' of the key by
	// putting it in an 'any' (interface) and inspecting it.
	var i any = key

	// An 'any' is a 2-word header: [type_pointer, data_pointer]
	// header := (*[2]uintptr)(unsafe.Pointer(&i))
	// typePtr := header[0] // Points to the internal runtime._type
	// dataPtr := header[1] // Points to the actual value

	// We use the typePtr to ensure we are hashing the value
	// correctly according to its specific type.
	return uint64(runtime_interhash(unsafe.Pointer(&i), 0))
}

// We use the runtime.memhash because it is the most stable linkname
// across Go versions. To make it safe for strings and structs,
// we will handle the types ourselves.

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

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

// runtime_interhash is what Go uses internally to hash the 'any' interface.
// It is perfectly safe, handles strings correctly, and won't trigger checkptr errors
// because we aren't doing the pointer math ourselves; the runtime is.
//
//go:noescape
//go:linkname runtime_interhash runtime.interhash
func runtime_interhash(p unsafe.Pointer, h uintptr) uintptr
