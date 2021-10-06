package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendToSlack(token, channel, color, title, message string) error {
	messageInterface := map[string]interface{}{
		"channel": channel,
		"text":    title,
		"attachments": []interface{}{
			map[string]interface{}{
				"color": color,
				"blocks": []interface{}{
					map[string]interface{}{
						"type": "section",
						"text": map[string]interface{}{
							"type": "plain_text",
							"text": message,
						},
					},
				},
			},
		},
	}
	messageJSON, err := json.Marshal(messageInterface)
	if err != nil {
		return fmt.Errorf("json marshal error %w", err)
	}

	endpoint := "https://slack.com/api/chat.postMessage"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("Failed to create new request.err = %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execugte request.err = %w", err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to parse response.err = %w", err)
	}

	resBody := map[string]interface{}{}

	err = json.Unmarshal(res, &resBody)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal error %w", err)
	}
	if _, exists := resBody["ok"]; !exists {
		return errors.New("key ok is not found in response json")
	}
	if _, success := resBody["ok"].(bool); !success {
		return errors.New("key ok is not boolean")
	}
	if isOK := resBody["ok"].(bool); !isOK {
		return errors.New("ok is not true")
	}

	return nil

}
