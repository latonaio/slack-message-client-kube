package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"slack-message-client-kube/config"
)

type SlackMessage struct {
	Channel     string       `json:"channel"`
	Message     string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Color string `json:"color"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (msg *SlackMessage) SendSlack(cfg *config.Config) error {
	jsonString, err := json.Marshal(msg)
	if err != nil {
		log.Printf("encoding error.err = %v", err)
		return err
	}

	endpoint := "https://slack.com/api/chat.postMessage"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonString))
	if err != nil {
		log.Printf("Failed to create new request.err = %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Slack.Token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execugte request.err = %v", err)
		return err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to parse response.err = %v", err)
		return err
	}

	resBinder := struct {
		Ok bool `json:"ok"`
		ResponseMetadata struct {
			Warnings []string `json:"warnings"`
		} `json:"response_metadata"`
	}{}

	err = json.Unmarshal(res,&resBinder)
	if err != nil {
		log.Println("Failed to unmarshal error")
		return err
	}

	return nil
}

func BuildMessage(podName, status,level string, channel string) *SlackMessage {
	var title string
	if level == "warning" {
		title = "Podの異常を検知"
	} else {
		return nil
	}
	ats := []Attachment{}
	a := Attachment{
		Color: "#ff6347",
		Title: title,
		Text:  fmt.Sprintf("pod %s has failure in running. reaon is %s", podName, status),
	}
	ats = append(ats, a)

	return &SlackMessage{
		Channel:     "#" + channel,
		Message:     "pod running failure!",
		Attachments: ats,
	}
}
