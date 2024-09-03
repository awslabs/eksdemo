package aws

import (
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
)

const maxPages = 30

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
		return fmt.Errorf("%s", ae.ErrorMessage())
	}
	return err
}
