package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/jsii-runtime-go"
)

func handleRequest(ctx context.Context, event events.CodePipelineEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Println("Failed to create bucket", err)
		return "error", err
	}

	// Create a CloudFront client with additional configuration
	cloudfrontClient := cloudfront.New(sess)
	var quantity int64
	var items []*string
	paths := "/*"
	items = append(items, &paths)
	quantity = 1

	// Create a CodePipeline client
	codepipelineClient := codepipeline.New(sess)

	result, erraws := cloudfrontClient.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: jsii.String(event.CodePipelineJob.Data.ActionConfiguration.Configuration.UserParameters),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: jsii.String(time.Now().Format("20060102150405")),
			Paths: &cloudfront.Paths{
				Quantity: &quantity,
				Items:    items,
			},
		},
	})
	if erraws != nil {
		if aerr, ok := erraws.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				log.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeMissingBody:
				log.Println(cloudfront.ErrCodeMissingBody, aerr.Error())
			case cloudfront.ErrCodeInvalidArgument:
				log.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			case cloudfront.ErrCodeNoSuchDistribution:
				log.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodeBatchTooLarge:
				log.Println(cloudfront.ErrCodeBatchTooLarge, aerr.Error())
			case cloudfront.ErrCodeTooManyInvalidationsInProgress:
				log.Println(cloudfront.ErrCodeTooManyInvalidationsInProgress, aerr.Error())
			case cloudfront.ErrCodeInconsistentQuantities:
				log.Println(cloudfront.ErrCodeInconsistentQuantities, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
			_, errcode := codepipelineClient.PutJobFailureResult(&codepipeline.PutJobFailureResultInput{
				JobId: jsii.String(event.CodePipelineJob.ID),
				FailureDetails: &codepipeline.FailureDetails{
					Message: jsii.String(aerr.Error()),
					Type:    jsii.String("JobFailed"),
				},
			})
			if errcode != nil {
				log.Println(errcode)
			}
		}
		return string("Fail"), nil
	}
	log.Println(result)
	// insert codepipeline success here
	_, errcode := codepipelineClient.PutJobSuccessResult(&codepipeline.PutJobSuccessResultInput{
		JobId: jsii.String(event.CodePipelineJob.ID),
	})
	if errcode != nil {
		log.Println(errcode)
	}
	return string(eventJson), nil
}

func main() {
	runtime.Start(handleRequest)
}
