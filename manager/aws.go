package manager

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func getAwsMetadata(
	client *imds.Client,
	ctx context.Context,
	path string,
) (string, error) {

	// Set up the GetMetadata input
	input := &imds.GetMetadataInput{
		Path: path,
	}

	// Attempt to call GetMetadata
	result, err := client.GetMetadata(ctx, input)
	if err != nil {
		return "", err
	}

	defer func() {
		// Perform close operation
		err = result.Content.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Read all the contents from the result
	value, err := io.ReadAll(result.Content)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func getAwsTags(timeout time.Duration) (Tags, error) {

	// Load a default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Don't wait too long for API call
	ctx, cancel := context.WithTimeout(
		context.Background(), timeout,
	)
	defer cancel()

	// TODO: Implement version 2 of IMDS

	// Create new IMDS client from config
	imdsClient := imds.NewFromConfig(cfg)

	// Check if the tags can be retrieved using instance metadata
	tags, err := getAwsMetadata(imdsClient, ctx, "tags/instance")
	if err == nil {

		// Convert the tags into a slice
		keys := strings.Split(tags, "\n")

		output := make(Tags)
		// Loop through tag keys
		for _, key := range keys {

			// Retrieve value for the key
			value, err := getAwsMetadata(imdsClient, ctx, fmt.Sprintf("tags/instance/%s", key))
			if err != nil {
				return nil, err
			}

			output[key] = value
		}

		return output, nil

	} else {

		// Attempt to get own instance ID
		instanceID, err := getAwsMetadata(imdsClient, ctx, "instance-id")
		if err != nil {
			return nil, err
		}

		// Create new EC2 client from config
		ec2Client := ec2.NewFromConfig(cfg)

		// Select current instance
		filters := []types.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []string{
					instanceID,
				},
			},
		}

		// Set up the DescribeTags input
		input := &ec2.DescribeTagsInput{

			// Do not bother with NextToken because the
			// maximum number of tags per resource is 50

			Filters:    filters,
			MaxResults: aws.Int32(1000),
		}

		// Attempt to call DescribeTags
		result, err := ec2Client.DescribeTags(ctx, input)
		if err != nil {
			return nil, err
		}

		output := make(Tags)
		// Convert result into tags slice
		for _, tag := range result.Tags {
			output[*tag.Key] = *tag.Value
		}

		return output, nil
	}
}
