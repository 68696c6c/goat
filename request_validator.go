package main

import "gopkg.in/go-playground/validator.v8"

type RequestValidator struct {
	*validator.Validate
}

func (v *RequestValidator) ValidateStruct(i interface{}) error {
	return v.Struct(i)
}
