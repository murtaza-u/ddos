package manifest

import "encoding/json"

// Marshal marshals an object in to bytes abstracting away the encoding.
func Marshal(in any) ([]byte, error) {
	return json.Marshal(in)
}

// Unmarshal unmarshals the raw data back to an object abstracting away
// the encoding.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
