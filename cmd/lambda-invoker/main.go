package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

func doTask(ctx context.Context, queueURL string, client *sqs.Client, message *types.Message) error {
	// 呼び出し元でキャンセルされた場合、処理しない
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	//log.Printf("Run Start: %v", *message.Body)
	//
	//time.Sleep(30 * time.Second)

	log.Printf("Run Task: %v", *message.Body)

	// TODO: ここにメイン処理を書く
	// メッセージの処理、	AWS リソースとの通信、コンテナ API サーバーとの通信など

	// メッセージを削除する
	_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("error deleting messqge: %w", err)
	}

	return nil
}

func main() {
	log.Println("Start999")

	ctx := context.Background()

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               "http://elasticmq:9324",
			HostnameImmutable: true,
		}, nil
	})

	log.Println("Start1")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(resolver))
	if err != nil {
		fmt.Printf("Error loading AWS config: %v", err)
		return
	}

	client := sqs.NewFromConfig(cfg)
	queueURL := "http://elasticmq:9324/000000000000/sample1"

	// セマフォ用のバッファチャネル
	semaphore := make(chan struct{}, 4)

	log.Println("Start2")

	for {
		log.Println("Start Polling...")

		// ポーリング
		result, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(queueURL),
			MaxNumberOfMessages:   1,
			WaitTimeSeconds:       20,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			fmt.Printf("Error receiving message: %v", err)
			return
		}

		log.Println("Polling...")

		//var wg sync.WaitGroup

		// メッセージを処理する
		for _, message := range result.Messages {
			// 子ゴルーチン用のシャドーイング
			message := message

			// 子ゴルーチンを起動
			semaphore <- struct{}{} // セマフォを取得。足りなければ、ここてブロックされる
			//wg.Add(1)

			go func() {
				defer func() {
					//wg.Done()
					<-semaphore
				}()
				_ = doTask(ctx, queueURL, client, &message)
			}()
		}

		//wg.Wait()
	}
}
