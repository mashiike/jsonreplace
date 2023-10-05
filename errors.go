package jsonreplace

import "errors"

var (
	ErrNilEncoder    = errors.New("jsonreplace: nil encoder")
	ErrInvalidSchema = errors.New("jsonreplace: invalid schema")
	ErrNilReplacer   = errors.New("jsonreplace: nil replacer")
	ErrAbortReplacer = errors.New("jsonreplace: abort replacer")
)
