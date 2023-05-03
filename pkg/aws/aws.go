package aws

import (
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
)

const maxPages = 5

// Return cleaner error message for service API errors
func FormatError(err error) error {
	var ae smithy.APIError
	if err != nil && errors.As(err, &ae) {
		return ae
	}
	return err
}

// Return cleaner error message for service API errors
func FormatErrorAsMessageOnly(err error) error {
	var ae smithy.APIError
	if err != nil && errors.As(err, &ae) {
		return fmt.Errorf(ae.ErrorMessage())
	}
	return err
}
