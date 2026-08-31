package handler

import "encoding/json"

// Patch[T] is a three-state PATCH field: absent (not in request), null (explicit clear), or a value.
// Use in request structs paired with validate tags — the validator sees the inner value via RegisterCustomTypeFunc.
type Patch[T any] struct {
	Set   bool
	Value *T
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	p.Set = true
	if string(data) == "null" {
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	p.Value = &v
	return nil
}
