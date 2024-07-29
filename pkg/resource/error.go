package resource

import "fmt"

// TODO: phase this out. Doesn't work with errors.As
type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

type NotFoundByIDError struct {
	Type string
	ID   string
}

func (e *NotFoundByIDError) Error() string {
	return fmt.Sprintf("%s %q not found", e.Type, e.ID)
}

type NotFoundByNameError struct {
	Type string
	Name string
}

func (e *NotFoundByNameError) Error() string {
	return fmt.Sprintf("%s with name %q not found", e.Type, e.Name)
}

// TODO: This error could potentially replace NotFoundByIDError and NotFoundByNameError
type NotFoundByError struct {
	Type  string
	Name  string
	Value string
}

func (e *NotFoundByError) Error() string {
	return fmt.Sprintf("%s with %s %q not found", e.Type, e.Name, e.Value)
}
