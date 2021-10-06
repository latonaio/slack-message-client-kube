package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slack-message-client-kube/cmd"
	"slack-message-client-kube/config"
	"syscall"

	"github.com/streadway/amqp"
)

const msName = "slack-message-client-kube"

func failOnError(err error, errCh chan error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		errCh <- err
	}
}

func main() {
	errCh := make(chan error, 1)
	quiteCh := make(chan syscall.Signal, 1)

	_, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()

	// 接続
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, errCh, "Failed to open a channel")
	defer conn.Close()

	// チャンネル生成
	ch, err := conn.Channel()
	failOnError(err, errCh, "Failed to open a channel")
	defer ch.Close()

	// 参照する queue の存在確認
	q, err := ch.QueueInspect(os.Getenv("QUEUE_FROM"))
	failOnError(err, errCh, "queue does not exist")

	// queue からメッセージの受け取り
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, errCh, "Failed to register a consumer")

	if err = ch.Qos(1, 0, false); err != nil {
		log.Printf("failed to set prefetch count: %v", err)
	}

	go func() {
		for d := range msgs {
			jsonData := map[string]interface{}{}
			err := json.Unmarshal(d.Body, &jsonData)
			if err != nil {
				log.Printf("JSON decode error: %v\n", err)
			}
			log.Printf("message : %v\n", jsonData)

			if level := jsonData["level"].(string); level != "warning" {
				d.Ack(false)
				continue
			}

			fmt.Printf("send message to channel %v\n", cfg.Slack.ChannelId)

			if err := cmd.SendToSlack(
				cfg.Slack.Token,
				cfg.Slack.ChannelId,
				"#ff6347",
				"*Podの異常を検知*",
				fmt.Sprintf("pod %s has failure in running. reason is %s", jsonData["pod_name"].(string), jsonData["status"].(string)),
			); err != nil {
				errCh <- err
				d.Nack(false, false)
				continue
			}

			d.Ack(false)
		}
	}()

	for {
		select {
		case err := <-errCh:
			log.Print(err)
		case <-quiteCh:
			cancel()
		}
	}
}
