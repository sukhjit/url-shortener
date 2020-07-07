package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sukhjit/url-shortener/model"
	"github.com/sukhjit/url-shortener/repo"
)

type svc struct {
	db        *dynamodb.DynamoDB
	tableName string
}

// NewShortener func
func NewShortener(awsRegion, tableName string) repo.Shortener {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	return &svc{
		db:        dynamodb.New(session),
		tableName: tableName,
	}
}

func (s *svc) Add(item *model.Shortener) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	})

	return err
}

func (s *svc) Info(slug string) (*model.Shortener, error) {
	result, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"slug": {
				S: aws.String(slug),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	item := model.Shortener{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
