package application

import (
	"fmt"
)

func (a *Application) CreateDependencies() error {
	if len(a.Dependencies) > 0 {
		fmt.Printf("Creating %d dependencies for %s\n", len(a.Dependencies), a.Name)
	}

	for _, res := range a.Dependencies {
		fmt.Printf("Creating dependency: %s\n", res.Common().Name)

		a.AssignCommonResourceOptions(res)

		if err := res.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) DeleteDependencies() error {
	if len(a.Dependencies) > 0 {
		fmt.Printf("Deleting %d dependencies for %s\n", len(a.Dependencies), a.Name)
	}

	for _, res := range a.Dependencies {
		fmt.Printf("Deleting dependency: %s\n", res.Common().Name)

		a.AssignCommonResourceOptions(res)

		if err := res.Delete(); err != nil {
			return err
		}
	}
	return nil
}
