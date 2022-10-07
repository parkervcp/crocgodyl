package crocgodyl

import "fmt"

type Error struct {
	Code   string      `json:"code"`
	Status string      `json:"status"`
	Detail string      `json:"detail"`
	Meta   interface{} `json:"meta,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Status, e.Code, e.Detail)
}

type ApiError struct {
	Errors []*Error `json:"errors"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d unexpected error(s)", len(e.Errors))
}
