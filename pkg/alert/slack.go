package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	net "github.com/pinkgorilla/go-sample/pkg/http/client"
)

const apiURL = "https://slack.com/api"
const defaultHTTPTimeout = 80 * time.Second

// SlackResponse slack response ...
type SlackResponse struct {
	Ok        bool                   `json:"ok"`
	Error     string                 `json:"error"`
	Channel   string                 `json:"channel"`
	Timestamp float64                `json:"ts"`
	Message   map[string]interface{} `json:"message"`
}

// SlackAlertConfig represent the config needed when creating a new slack notifier
type SlackAlertConfig struct {
	Token   string
	Channel string
}

// SlackAlert alert implementation over slack channel
type SlackAlert struct {
	Token      string
	Channel    string
	httpClient *net.Client
}

func (sn *SlackAlert) Error(err error) error {
	return sn.Alert(Message{
		Text:  err.Error(),
		Error: err,
		Trace: nil,
	})
}

// Alert alerts message to a slack channel
func (sn *SlackAlert) Alert(message Message) error {
	/*
		Examples of calling the slack API:

			curl -X POST -H 'Authorization: Bearer xoxb-1234-56789abcdefghijklmnop' \
			-H 'Content-type: application/json' \
			--data '{
				"channel":"C061EG9SL",
				"text":"I hope the tour went well, Mr. Wonka.",
				"attachments": [{
					"text":"Who wins the lifetime supply of chocolate?",
					"fallback":"You could be telling the computer exactly what it can do with a lifetime supply of chocolate.",
					"color":"#3AA3E3",
					"attachment_type":"default",
					"callback_id":"select_simple_1234",
					"actions":[{
						"name":"winners_list",
						"text":"Who should win?",
						"type":"select",
						"data_source":"users"
					}]
				}]
			}' \
			https://slack.com/api/chat.postMessage
	*/

	payload := map[string]interface{}{
		"channel": sn.Channel,
		"text":    message.Text,
	}

	if message.Error != nil || len(message.Trace) > 0 {
		var errMessage string
		var traceMessage string

		if message.Error != nil {
			errMessage = message.Error.Error()
		}
		if len(message.Trace) > 0 {
			traceMessage = string(message.Trace)
		}
		payload["attachments"] = []interface{}{
			map[string]interface{}{
				"text": strings.Join([]string{errMessage, traceMessage}, "\n"),
			},
		}
	}

	if len(message.Title) > 0 {
		payload["username"] = message.Title
	}
	if len(message.Icon) > 0 {
		payload["icon_emoji"] = message.Icon
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat.postMessage", apiURL), bytes.NewBuffer(bs))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sn.Token))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
		return err
	}

	res, err := sn.httpClient.Do(req)
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
		return err
	}
	if res.StatusCode != http.StatusOK {
		e := fmt.Errorf("error: %v was occured while trying send message to slack.\nMessage was: %v", errors.New(http.StatusText(res.StatusCode)), message)
		log.Println(e)
		return e
	}
	var slackResponse SlackResponse
	err = json.NewDecoder(res.Body).Decode(&slackResponse)
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
		return err
	}
	if !slackResponse.Ok {
		e := fmt.Errorf("error: %v was occured while trying send message to slack.\nMessage was: %v", errors.New(slackResponse.Error), message)
		log.Println(e)
		return e
	}
	return nil
}

// NewSlackAlert creates a new slack notifier
func NewSlackAlert(token, channel string) *SlackAlert {

	return &SlackAlert{
		Token:      token,
		Channel:    channel,
		httpClient: net.NewClient(net.DefaultClientConfig()),
	}
}
