package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/katesclau/opensearch-test/models"
)

func main() {
	ctx := context.Background()

	// AWS Config
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	dynClient := dynamodb.NewFromConfig(awsCfg)

	addItem := func(svc *dynamodb.Client) {
		w := &models.TimeseriesString{
			TS:  time.Now(),
			Val: "100",
		}
		buf := &bytes.Buffer{}
		err = json.NewEncoder(buf).Encode(w)
		if err != nil {
			log.Fatalf("Failed to encode JSON: %s", err)
			return
		}

		item, err := attributevalue.MarshalMap(map[string]interface{}{
			"ID":        uuid.New().String(),
			"Timestamp": time.Now().Unix(),
			"Data":      buf.Bytes(),
			"model":     "timeseries",
		})
		if err != nil {
			log.Fatalf("Failed to marshal item: %s", err)
			return
		}

		input := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String("service-name"),
		}

		_, err = svc.PutItem(ctx, input)
		if err != nil {
			log.Fatalf("Failed to put item in DynamoDB: %s", err)
		}
	}

	// Graceful shutdown
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done, breaking the loop")
			return
		case <-stopCh:
			fmt.Println("Received signal to stop, breaking the loop")
			return
		default:
			time.Sleep(5 * time.Second)
			addItem(dynClient)
			fmt.Println("Successfully added item to DynamoDB")
		}
	}
}
