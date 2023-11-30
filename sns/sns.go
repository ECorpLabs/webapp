package sns

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func PublishToSNS(userDetails string) error {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")), // Replace with your AWS region
		// Add other AWS credentials or configurations as needed
	})
	if err != nil {
		return err
	}

	// Create an SNS service client
	svc := sns.New(sess)

	// Specify the SNS topic ARN
	topicARN := os.Getenv("SNS_TOPIC_ARN") // Replace with your SNS topic ARN

	// Define the message you want to publish
	message := userDetails // This should be the user details you want to send

	// Publish the message to the SNS topic
	_, err = svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicARN),
	})
	if err != nil {
		return err
	}

	return nil
}
