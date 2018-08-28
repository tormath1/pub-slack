package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	credentials = flag.String("credentials", "", "absolute path to the Google Credentials JSON file")
	projectID   = flag.String("project", "", "ID of your Google Project where the topic is created")
	topic       = flag.String("subscription", "", "name of the subscription")
	slackURL    = flag.String("slackURL", "", "Incoming webhook URL for Slack")
)

type slackMessage struct {
	Text string `json:"text"`
	Md   bool   `json:"mrkdwn"`
}

type pubSubMessage struct {
	BudgetDisplayName string  `json:"budgetDisplayName"`
	CostAmount        float64 `json:"costAmount"`
	CostIntervalStart string  `json:"costIntervalStart"`
	BudgetAmount      float64 `json:"budgetAmount"`
	BudgetAmountType  string  `json:"budgetAmountType"`
	CurrencyCode      string  `json:"currencyCode"`
}

func isFlagGiven(flag *string) bool {
	return len(*flag) > 0
}

func newSlackMessage(name, costAmount, budgetAmount, currencyCode string) *slackMessage {
	return &slackMessage{
		Text: fmt.Sprintf("*%s*\n_Cost amount:_ %s\n_Budget Amount:_ %s %s", name, costAmount, budgetAmount, currencyCode),
		Md:   true,
	}
}

func main() {

	flag.Parse()

	if !(isFlagGiven(credentials) || isFlagGiven(projectID) || isFlagGiven(topic) || isFlagGiven(slackURL)) {
		log.Fatalf("unable to start pub-slack: -credentials, -project, -slackURL or -subscription is missing. --help to get informations.")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, *projectID, option.WithCredentialsFile(*credentials))
	if err != nil {
		log.Fatalf("unable to create pub sub client: %v", err)
	}

	sub := client.Subscription(*topic)
	if err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		if err = handleResponse(m.Data); err != nil {
			log.Printf("unable to handle response from pub/sub: %v", err)
		} else {
			m.Ack()
		}
	}); err != nil {
		log.Fatalf("unable to get message from topic: %v", err)
	}

}

func handleResponse(data []byte) error {

	var pubsub pubSubMessage
	if err := json.Unmarshal(data, &pubsub); err != nil {
		return fmt.Errorf("unable to extract json from google response: %v", err)
	}

	msg := newSlackMessage(
		pubsub.BudgetDisplayName,
		strconv.FormatFloat(pubsub.CostAmount, 'f', -1, 64),
		strconv.FormatFloat(pubsub.BudgetAmount, 'f', -1, 64),
		pubsub.CurrencyCode,
	)

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("unable to create payload: %v", err)
	}

	_, err = http.Post(*slackURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("unable to post payload on slack URL: %v", err)
	}
	return nil
}
