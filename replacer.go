package jsonreplace

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

// A Replacer replaces a JSON object with another JSON object.
type Replacer interface {
	ReplaceJSON(json.RawMessage) (json.RawMessage, error)
}

// ReplacerFunc is an adapter to allow the use of ordinary functions as Replacers.
type ReplacerFunc func(json.RawMessage) (json.RawMessage, error)

func (f ReplacerFunc) ReplaceJSON(bs json.RawMessage) (json.RawMessage, error) {
	return f(bs)
}

type muxEntry struct {
	schema string
	loader gojsonschema.JSONLoader
	r      Replacer
}

// A ReplaceMux is a multiplexer for Replacers.
type ReplaceMux struct {
	mu          sync.RWMutex
	es          []muxEntry
	noRecursive bool
}

// NewReplaceMux returns a new ReplaceMux.
func NewReplaceMux() *ReplaceMux {
	return &ReplaceMux{}
}

// VoidReplacer is a Replacer that does nothing.
var VoidReplacer Replacer = ReplacerFunc(func(raw json.RawMessage) (json.RawMessage, error) {
	return raw, nil
})

// DefaultReplaceMux is the default ReplaceMux.
var DefaultReplaceMux = &defaultReplaceMux

var defaultReplaceMux ReplaceMux

// Clone returns a clone of the ReplaceMux.
func (mux *ReplaceMux) Clone() *ReplaceMux {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	clone := &ReplaceMux{
		noRecursive: mux.noRecursive,
	}
	if mux.es != nil {
		clone.es = make([]muxEntry, len(mux.es))
		copy(clone.es, mux.es)
	}
	return clone
}

// Replace adds a Replacer to the DefaultReplaceMux based on the given JSON schema.
func Replace(schema string, r Replacer) error {
	return DefaultReplaceMux.Replace(schema, r)
}

// ReplaceFunc adds a Replacer function to the DefaultReplaceMux based on the given JSON schema.
func ReplaceFunc(schema string, f func(json.RawMessage) (json.RawMessage, error)) error {
	return DefaultReplaceMux.ReplaceFunc(schema, f)
}

// Replace adds a Replacer to the ReplaceMux based on the given JSON schema.
func (mux *ReplaceMux) Replace(schema string, r Replacer) error {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if !json.Valid([]byte(schema)) {
		return ErrInvalidSchema
	}
	if r == nil {
		return ErrNilReplacer
	}
	if mux.es == nil {
		mux.es = make([]muxEntry, 0, 1)
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)
	mux.es = append(mux.es, muxEntry{
		schema: schema,
		loader: schemaLoader,
		r:      r,
	})
	return nil
}

// ReplaceFunc adds a Replacer function to the ReplaceMux based on the given JSON schema.
func (mux *ReplaceMux) ReplaceFunc(schema string, f func(json.RawMessage) (json.RawMessage, error)) error {
	return mux.Replace(schema, ReplacerFunc(f))
}

// ReplaceJSON replaces a JSON object with another JSON object based on the JSON schema.
// It executes the ReplaceJSON function of the Replacer that matches the given JSON schema.
// It checks the registered JSON schemas in the order they were added.
// If no JSON schema matches, it recursively checks substructures such as objects and arrays.
// To turn off recursive checking, call mux.NoRecursive(true).
func (mux *ReplaceMux) ReplaceJSON(raw json.RawMessage) (json.RawMessage, error) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	return mux.replaceJSON(raw)
}

// NoRecursive sets whether the ReplaceMux should recursively check substructures such as objects and arrays.
func (mux *ReplaceMux) NoRecursive(flag bool) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.noRecursive = flag
}

func (mux *ReplaceMux) replaceJSON(raw json.RawMessage) (json.RawMessage, error) {
	target := gojsonschema.NewBytesLoader(raw)
	if len(raw) == 0 {
		return raw, nil
	}
	for _, e := range mux.es {
		result, err := gojsonschema.Validate(e.loader, target)
		if err != nil {
			return nil, err
		}
		if !result.Valid() {
			continue
		}
		replaced, err := e.r.ReplaceJSON(raw)
		if err != nil {
			if errors.Is(err, ErrAbortReplacer) {
				return replaced, nil
			}
			return nil, err
		}
		raw = replaced
	}
	if mux.noRecursive {
		return raw, nil
	}
	switch raw[0] {
	case '{':
		var obj map[string]json.RawMessage
		err := json.Unmarshal(raw, &obj)
		if err != nil {
			return nil, err
		}
		for k, v := range obj {
			v, err = mux.replaceJSON(v)
			if err != nil {
				return nil, err
			}
			obj[k] = v
		}
		replaced, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		return replaced, nil
	case '[':
		var arr []json.RawMessage
		err := json.Unmarshal(raw, &arr)
		if err != nil {
			return nil, err
		}
		for i, v := range arr {
			v, err := mux.replaceJSON(v)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		replaced, err := json.Marshal(arr)
		if err != nil {
			return nil, err
		}
		return replaced, nil
	}
	return raw, nil
}
