package s3_bucket

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, _ := newOptions()
	return NewResourceWithOptions(options)
}

func NewResourceWithOptions(options *BucketOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "s3-bucket",
			Description: "Amazon S3 Bucket",
			Aliases:     []string{"s3-buckets", "s3", "buckets", "bucket"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
