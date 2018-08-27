package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	credentials = flag.String("credentials", "", "absolute path to the Google Credentials JSON file")
	projectID   = flag.String("project", "", "ID of your Google Project where the topic is created")
	topic       = flag.String("subscription", "", "name of the subscription")
)

func isFlagGiven(flag *string) bool {
	return len(*flag) > 0
}

func main() {

	flag.Parse()

	if !(isFlagGiven(credentials) || isFlagGiven(projectID) || isFlagGiven(topic)) {
		log.Fatalf("unable to start pub-slack: -credentials, -project or -subscription is missing. --help to get informations.")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, *projectID, option.WithCredentialsFile(*credentials))
	if err != nil {
		log.Fatalf("unable to create pub sub client: %v", err)
	}

	sub := client.Subscription(*topic)
	if err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		m.Ack()
	}); err != nil {
		log.Fatalf("unable to get message from topic: %v", err)
	}

}
