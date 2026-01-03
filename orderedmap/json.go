package orderedmap

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements json.Marshaler.
// It ensures the keys appear in the JSON object in the order they were inserted.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	first := true
	// We use the internal slice to maintain the order
	for _, entry := range m.slots {
		if entry.deleted {
			continue
		}

		if !first {
			buf.WriteByte(',')
		}
		first = false

		// Marshal Key
		keyJSON, err := json.Marshal(entry.key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')

		// Marshal Value
		valJSON, err := json.Marshal(entry.value)
		if err != nil {
			return nil, err
		}
		buf.Write(valJSON)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Note: JSON object key order is not guaranteed by all parsers,
// but the standard library's Decoder processes them sequentially.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	// We use a decoder to process the object key by key to preserve order
	dec := json.NewDecoder(bytes.NewReader(data))

	// Expect start of object '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{', got %v", t)
	}

	for dec.More() {
		// Read Key
		t, err := dec.Token()
		if err != nil {
			return err
		}

		// Convert token to string (JSON keys are always syntactically strings)
		t, ok := t.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", t)
		}
		keyStr := fmt.Sprintf("%v", t)

		var key K
		var kAny any = &key

		// Use a type switch on the pointer to the key
		switch k := kAny.(type) {
		case *string:
			*k = keyStr
		case *int:
			fmt.Sscanf(keyStr, "%d", k)
		case *int64:
			fmt.Sscanf(keyStr, "%d", k)
		case *float64:
			fmt.Sscanf(keyStr, "%f", k)
		default:
			// Fallback for custom types or other primitives
			if err := json.Unmarshal([]byte(fmt.Sprintf("%q", keyStr)), &key); err != nil {
				return err
			}
		}

		// Read Value
		var val V
		if err := dec.Decode(&val); err != nil {
			return err
		}

		m.Put(key, val)
	}

	// Expect end of object '}'
	_, err = dec.Token()
	return err
}
