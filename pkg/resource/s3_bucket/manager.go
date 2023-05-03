package s3_bucket

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun   bool
	s3Client *aws.S3Client
	s3Getter *Getter
}

func (m *Manager) Init() {
	if m.s3Client == nil {
		m.s3Client = aws.NewS3Client()
	}
	m.s3Getter = NewGetter(m.s3Client)
}

func (m *Manager) Create(options resource.Options) error {
	bucketOptions, ok := options.(*BucketOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to BucketOptions")
	}

	_, err := m.s3Getter.GetBucketByName(bucketOptions.BucketName)

	if err == nil {
		fmt.Printf("Bucket %q already exists\n", bucketOptions.BucketName)
		return nil
	} else {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if m.DryRun {
		return m.dryRun(bucketOptions)
	}

	fmt.Printf("Creating Bucket: %s...", bucketOptions.BucketName)

	err = m.s3Client.CreateBucket(bucketOptions.BucketName, options.Common().Region)
	if err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	bucketOptions, ok := options.(*BucketOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to BucketOptions")
	}

	fmt.Printf("Deletion of Bucket %q not supported. Please delete manually.\n", bucketOptions.BucketName)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *BucketOptions) error {
	fmt.Printf("\nS3 Bucket Manager Dry Run:\n")
	fmt.Printf("S3 API Call %q with request parameters:\n", "CreateBucket")
	fmt.Printf("Bucket: %q\n", options.BucketName)
	fmt.Printf("CreateBucketConfiguration.LocationConstraint: %q\n", options.Region)
	return nil
}
