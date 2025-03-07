package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/olivere/elastic/v7"
)

var (
	esClient *elastic.Client
)

const mapping = `
{
	"settings":{
		"number_of_shards": 2,
		"number_of_replicas": 0
	},
	"mappings":{
		"timeseries":{
			"properties":{
				"ts":{
					"type":"date"
				},
				"val":{
					"type":"text",
					"store": true,
					"fielddata": true
				}
			}
		}
	}
}`

func init() {
	var err error
	esClient, err = elastic.NewClient(
		elastic.SetURL(os.Getenv("OPENSEARCH_URL")),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	// Ensure the index exists
	exists, err := esClient.IndexExists("listings").Do(ctx)
	if err != nil {
		log.Fatalf("Error checking if the index exists: %s", err)
		return err
	}
	if !exists {
		createIndex, err := esClient.CreateIndex("listings").BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatalf("Error creating the index: %s", err)
		}
		if !createIndex.Acknowledged {
			log.Fatalf("CreateIndex was not acknowledged. %v", createIndex)
		}
	}

	for _, record := range event.Records {
		if record.EventName == "INSERT" || record.EventName == "MODIFY" {
			item := record.Change.NewImage
			tmp := eventStreamToMap(item)
			var i any
			if err := attributevalue.UnmarshalMap(tmp, &i); err != nil {
				log.Fatalf("Error unmarshalling record: %s", err)
				return err
			}
			_, err = esClient.Index().
				Index("listings").
				Id(record.Change.Keys["ID"].String()).
				BodyJson(i).
				Do(ctx)
			if err != nil {
				log.Printf("Error indexing record: %s", err)
			}
		} else if record.EventName == "REMOVE" {
			_, err := esClient.Delete().
				Index("listings").
				Id(record.Change.Keys["ID"].String()).
				Do(ctx)
			if err != nil {
				log.Printf("Error deleting record: %s", err)
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

// ugly hack because the types
// events.DynamoDBAttributeValue != *dynamodb.AttributeValue
func eventStreamToMap(attribute any) map[string]types.AttributeValue {
	// Map to be returned
	m := make(map[string]types.AttributeValue)

	tmp := make(map[string]events.DynamoDBAttributeValue)

	switch t := attribute.(type) {
	case map[string]events.DynamoDBAttributeValue:
		tmp = t
	case events.DynamoDBAttributeValue:
		tmp = t.Map()
	}

	for k, v := range tmp {
		switch v.DataType() {
		case events.DataTypeString:
			s := v.String()
			m[k] = &types.AttributeValueMemberS{
				Value: s,
			}
		case events.DataTypeBoolean:
			b := v.Boolean()
			m[k] = &types.AttributeValueMemberBOOL{
				Value: b,
			}
		case events.DataTypeMap:
			m[k] = &types.AttributeValueMemberM{
				Value: eventStreamToMap(v),
			}
		case events.DataTypeNumber:
			n := v.Number()
			m[k] = &types.AttributeValueMemberN{
				Value: n,
			}
		case events.DataTypeList:
			m[k] = &types.AttributeValueMemberL{
				Value: eventStreamToList(v),
			}
		}
	}
	return m
}

// ugly hack because the types
// events.DynamoDBAttributeValue != *dynamodb.AttributeValue
func eventStreamToList(attribute any) []types.AttributeValue {
	// List to be returned
	l := make([]types.AttributeValue, 0)

	var tmp []events.DynamoDBAttributeValue

	switch t := attribute.(type) {
	case []events.DynamoDBAttributeValue:
		tmp = t
	case events.DynamoDBAttributeValue:
		tmp = t.List()
	}

	for _, v := range tmp {
		switch v.DataType() {
		case events.DataTypeString:
			s := v.String()
			l = append(l, &types.AttributeValueMemberS{
				Value: s,
			})
		case events.DataTypeBoolean:
			b := v.Boolean()
			l = append(l, &types.AttributeValueMemberBOOL{
				Value: b,
			})
		case events.DataTypeMap:
			l = append(l, &types.AttributeValueMemberM{
				Value: eventStreamToMap(v),
			})
		case events.DataTypeNumber:
			n := v.Number()
			l = append(l, &types.AttributeValueMemberN{
				Value: n,
			})
		case events.DataTypeList:
			l = append(l, &types.AttributeValueMemberL{
				Value: eventStreamToList(v),
			})
		}
	}
	return l
}
