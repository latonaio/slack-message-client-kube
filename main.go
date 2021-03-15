package main

import (
	"context"
	"log"
	"slack-message-client-kube/cmd"
	"slack-message-client-kube/config"
	"syscall"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
)

const msName = "slack-message-client-kube"

func main() {
	errCh := make(chan error, 1)
	quiteCh := make(chan syscall.Signal, 1)

	ctx, cancel := context.WithCancel(context.Background())

	c, err := msclient.NewKanbanClient(ctx)
	if err != nil {
		errCh <- err
	}

	kCh, err := c.GetKanbanCh(msName, c.GetProcessNumber())
	if err != nil {
		errCh <- err
	}

	cfg, err := config.New()

	go func() {
		for kanban := range kCh {
			data, err := kanban.GetMetadataByMap()
			if err != nil {
				errCh <- err
			}

			msg := data["pod_name"].(string)
			status := data["status"].(string)
			level := data["level"].(string)
			sm := cmd.BuildMessage(msg, status,level,cfg.Slack.ChannelId)
			if err := sm.SendSlack(cfg); err != nil {
				errCh <- err
			}

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
