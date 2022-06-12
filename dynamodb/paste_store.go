package dynamodb

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jsmithdenverdev/microbin"
)

type PasteStore struct {
	Table  string
	Client *dynamodb.Client
}

func (s *PasteStore) Create(ctx context.Context, p microbin.Paste) error {
	item, err := attributevalue.MarshalMap(p)

	if err != nil {
		return err
	}

	_, err = s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.Table),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *PasteStore) Read(ctx context.Context, id int) (microbin.Paste, error) {
	paste := microbin.Paste{}

	result, err := s.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{
				Value: strconv.Itoa(id),
			},
		},
	})

	if err != nil {
		return paste, err
	}

	if result.Item == nil {
		return paste, microbin.ErrPasteNotFound
	}

	err = attributevalue.UnmarshalMap(result.Item, &paste)
	if err != nil {
		return paste, err
	}

	return paste, nil
}

func (s *PasteStore) Delete(ctx context.Context, id int) error {
	_, err := s.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{
				Value: strconv.Itoa(id),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *PasteStore) List(ctx context.Context) ([]microbin.Paste, error) {
	result, err := s.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.Table),
	}, nil)

	if err != nil {
		return []microbin.Paste{}, err
	}

	pastes := []microbin.Paste{}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &pastes)

	if err != nil {
		return []microbin.Paste{}, err
	}

	return pastes, nil
}
