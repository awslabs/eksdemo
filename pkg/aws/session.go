package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var awsConfig *aws.Config
var profile string
var region string
var debug, responseBodyDebug bool

func Init(awsProfile, awsRegion string, awsDebug, awsResponseBodyDebug bool) {
	profile = awsProfile
	region = awsRegion
	debug = awsDebug
	responseBodyDebug = awsResponseBodyDebug
}

func GetConfig() aws.Config {
	if awsConfig != nil {
		return *awsConfig
	}

	cfgOptions := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	}

	if debug || responseBodyDebug {
		logMode := aws.LogRetries | aws.LogRequestWithBody | aws.LogResponse

		if responseBodyDebug {
			logMode |= aws.LogResponseWithBody
		}

		cfgOptions = append(cfgOptions,
			config.WithClientLogMode(logMode),
		)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), cfgOptions...)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create AWS config: %w", err))
	}
	region = cfg.Region
	awsConfig = &cfg

	return cfg
}

func Region() string {
	if region == "" {
		GetConfig()
	}
	return region
}
