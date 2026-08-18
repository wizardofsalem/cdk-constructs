// Package cdkutil contains shared CDK helper functions for ARN formatting.
package cdkutil

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

// DynamoTableArn formats a DynamoDB table ARN from the stack context.
func DynamoTableArn(stack awscdk.Stack, tableName string) *string {
	return stack.FormatArn(&awscdk.ArnComponents{
		Service:      jsii.String("dynamodb"),
		Resource:     jsii.String("table"),
		ResourceName: jsii.String(tableName),
	})
}

// S3BucketArn formats an S3 bucket ARN with an optional key path.
func S3BucketArn(stack awscdk.Stack, bucket string, key string) *string {
	return stack.FormatArn(&awscdk.ArnComponents{
		Service:      jsii.String("s3"),
		Region:       jsii.String(""),
		Account:      jsii.String(""),
		Resource:     jsii.String(bucket),
		ResourceName: jsii.String(key),
	})
}
