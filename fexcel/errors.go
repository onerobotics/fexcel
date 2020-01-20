package fexcel

import (
	"fmt"
	"sync"
)

type errorList struct {
	Errors []error
	mux    sync.Mutex
}

func (e *errorList) Add(err error) int {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.Errors = append(e.Errors, err)
	return len(e.Errors)
}

func (e errorList) Error() string {
	switch len(e.Errors) {
	case 0:
		return "no errors"
	case 1:
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", e.Errors[0], len(e.Errors)-1)
}

func (e errorList) Err() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}
