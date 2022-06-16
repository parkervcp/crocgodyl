package crocgodyl

import "fmt"

// TODO: finish this
type ApiError struct {
	Code   string
	Status string
	Detail string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s (code: %s): %s", e.Status, e.Code, e.Detail)
}
