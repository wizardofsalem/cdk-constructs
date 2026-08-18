// Package dynamoutil contains wrapper functions for DynamoDB SDK operations.
package dynamoutil

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoItem represents a DynamoDB item with its table name and attribute map.
type DynamoItem struct {
	Table        string
	AttributeMap map[string]types.AttributeValue
}

var ErrItemNotFound = errors.New("dynamo item not found")

func NewClient(ctx context.Context, cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func PutItem(ctx context.Context, client *dynamodb.Client, item DynamoItem) error {
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(item.Table), Item: item.AttributeMap,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func DeleteItem(ctx context.Context, client *dynamodb.Client, item DynamoItem) error {
	out, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:          item.AttributeMap,
		TableName:    aws.String(item.Table),
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		log.Printf("Couldn't delete item from table. Here's why: %v\n", err)
		return err
	}

	if len(out.Attributes) == 0 {
		return ErrItemNotFound
	}

	return nil
}

func GetItem(ctx context.Context, client *dynamodb.Client, item DynamoItem) (map[string]types.AttributeValue, error) {
	out, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(item.Table),
		Key:       item.AttributeMap,
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, ErrItemNotFound
	}

	return out.Item, nil
}

func Query(ctx context.Context, client *dynamodb.Client, tableName string, keyCondition string, attributeNames map[string]string, attributeValues map[string]types.AttributeValue, filter string) ([]map[string]types.AttributeValue, error) {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeNames:  attributeNames,
		ExpressionAttributeValues: attributeValues,
	}

	if filter != "" {
		input.FilterExpression = aws.String(filter)
	}

	out, err := client.Query(ctx, input)
	if err != nil {
		log.Println("error performing dynamo query", err)
		return nil, err
	}

	return out.Items, nil
}

func Scan(ctx context.Context, client *dynamodb.Client, tableName string, projection string, propertyName string, value string) ([]map[string]types.AttributeValue, error) {
	var startKey map[string]types.AttributeValue
	var scannedItems []map[string]types.AttributeValue
	for {
		out, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(tableName),
			FilterExpression: aws.String("#a = :v"),
			ExpressionAttributeNames: map[string]string{
				"#a": propertyName,
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberS{Value: value},
			},
			ProjectionExpression: aws.String(projection),
			ExclusiveStartKey:    startKey,
		})
		if err != nil {
			return nil, err
		}

		scannedItems = append(scannedItems, out.Items...)
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		startKey = out.LastEvaluatedKey
	}

	return scannedItems, nil
}

func ScanAll(ctx context.Context, client *dynamodb.Client, table string) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	var lastKey map[string]types.AttributeValue

	for {
		out, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(table),
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, out.Items...)
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		lastKey = out.LastEvaluatedKey
	}

	return items, nil
}
