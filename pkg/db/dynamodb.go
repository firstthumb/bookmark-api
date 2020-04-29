package db

import (
	"bookmark-api/pkg/utils"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

func GetTableBookmark() string {
	v := os.Getenv("DYNAMO_TABLE_BOOKMARK")
	if v == "" {
		return "dev-bookmark-api-bookmark-table"
	}
	return v
}

func GetTableTag() string {
	v := os.Getenv("DYNAMO_TABLE_TAG")
	if v == "" {
		return "dev-bookmark-api-tag-table"
	}
	return v
}

func GetTableUser() string {
	v := os.Getenv("DYNAMO_TABLE_USER")
	if v == "" {
		return "dev-bookmark-api-user-table"
	}
	return v
}

func GenerateID() string {
	return uuid.New().String()
}

func GetDynamoDb() *dynamo.DB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if utils.IsOffline() {
		cred := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_ACCESS_SECRET"), "")
		return dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_ACCESS_REGION")).WithCredentials(cred))
	} else {
		return dynamo.New(sess)
	}
}
