package services

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// SnsActions encapsulates the Amazon Simple Notification Service (Amazon SNS) actions
// used in the examples.
type SnsActions struct {
	SnsClient *sns.Client
}

func (actor SnsActions) CreateTopic(ctx context.Context, topicName string, isFifoTopic bool, contentBasedDeduplication bool) (string, error) {
	var topicArn string
	topicAttributes := map[string]string{}
	if isFifoTopic {
		topicAttributes["FifoTopic"] = "true"
	}
	if contentBasedDeduplication {
		topicAttributes["ContentBasedDeduplication"] = "true"
	}
	topic, err := actor.SnsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name:       aws.String(topicName),
		Attributes: topicAttributes,
	})
	if err != nil {
		log.Printf("Couldn't create topic %v. Here's why: %v\n", topicName, err)
	} else {
		topicArn = *topic.TopicArn
	}

	return topicArn, err
}

func (actor SnsActions) DeleteTopic(ctx context.Context, topicArn string) error {
	_, err := actor.SnsClient.DeleteTopic(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn)})
	if err != nil {
		log.Printf("Couldn't delete topic %v. Here's why: %v\n", topicArn, err)
	}
	return err
}

// CreateUserEndpoint registers a device token (FCM) to SNS and returns an EndpointARN
func (actor SnsActions) CreateUserEndpoint(ctx context.Context, platformAppArn string, fcmToken string) (string, error) {
	input := &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(platformAppArn),
		Token:                  aws.String(fcmToken),
	}

	result, err := actor.SnsClient.CreatePlatformEndpoint(ctx, input)
	if err != nil {
		log.Printf("Couldn't create endpoint for token. Error: %v", err)
		return "", err
	}

	return *result.EndpointArn, nil
}

func (actor SnsActions) SubscribeToTopic(ctx context.Context, topicArn string, endpointArn string) error {
	_, err := actor.SnsClient.Subscribe(ctx, &sns.SubscribeInput{
		Protocol: aws.String("application"), // Use "application" for Push/FCM
		TopicArn: aws.String(topicArn),
		Endpoint: aws.String(endpointArn),
	})
	return err
}

// Modified Publish to support TargetArn (Direct to one user)
func (actor SnsActions) PublishToUser(ctx context.Context, endpointArn string, message string) error {
	publishInput := sns.PublishInput{
		TargetArn: aws.String(endpointArn), // <--- Use EndpointArn here
		Message:   aws.String(message),
	}
	_, err := actor.SnsClient.Publish(ctx, &publishInput)
	return err
}

func (actor SnsActions) Publish(ctx context.Context, topicArn string, message string, groupId string, dedupId string, filterKey string, filterValue string) error {
	publishInput := sns.PublishInput{TopicArn: aws.String(topicArn), Message: aws.String(message)}
	if groupId != "" {
		publishInput.MessageGroupId = aws.String(groupId)
	}
	if dedupId != "" {
		publishInput.MessageDeduplicationId = aws.String(dedupId)
	}
	if filterKey != "" && filterValue != "" {
		publishInput.MessageAttributes = map[string]types.MessageAttributeValue{
			filterKey: {DataType: aws.String("String"), StringValue: aws.String(filterValue)},
		}
	}
	_, err := actor.SnsClient.Publish(ctx, &publishInput)
	if err != nil {
		log.Printf("Couldn't publish message to topic %v. Here's why: %v", topicArn, err)
	}
	return err
}
