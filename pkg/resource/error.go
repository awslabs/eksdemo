package resource

import "fmt"

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
