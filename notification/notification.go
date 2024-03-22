package notifications

import (
	novu "github.com/novuhq/go-novu/lib"

	"github.com/joho/godotenv"

	"context"
	"log"
	"os"
)

var novuClient *novu.APIClient

var ctx = context.Background()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	novuApiKey := os.Getenv("NOVU_API_KEY")
	if novuApiKey == "" {
		log.Fatal("Error loading novu api key")
	}
	novuClient = novu.NewAPIClient(novuApiKey, &novu.Config{})
}

func CreateSubscriber(userDetails Subscriber) (*novu.SubscriberResponse, error) {
	subscriber := novu.SubscriberPayload{
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Email:     userDetails.Email,
		Avatar:    userDetails.Avatar,
		Data:      userDetails.Data,
		Phone:     userDetails.Phone,
	}
	resp, err := novuClient.SubscriberApi.Identify(ctx, userDetails.SubscriberID, subscriber)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func UpdateSubscriber(subscriberID string, userDetails Subscriber) (*novu.SubscriberResponse, error) {
	updateSubscriber := novu.SubscriberPayload{
		FirstName: userDetails.FirstName,
		LastName: userDetails.LastName,
		Email:    userDetails.Email,
		Avatar:   userDetails.Avatar,
		Data:     userDetails.Data,
		Phone:    userDetails.Phone,
	}
	resp, err := novuClient.SubscriberApi.Update(ctx, subscriberID, updateSubscriber)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func SendNotification(payload Trigger) (*novu.EventResponse, error) {
	resp, err := novuClient.EventApi.Trigger(ctx, payload.EventID, novu.ITriggerPayloadOptions{
		To:      payload.To,
		Payload: payload.Data, // dynamic data
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func CreateTopic(topicKey string, topicName string) error {
	err := novuClient.TopicsApi.Create(ctx, topicKey, topicName)
	if err != nil {
		return err
	}
	return nil
}

func AddSubscriber(topicKey string, subscribers []string) error {
	err := novuClient.TopicsApi.AddSubscribers(ctx, topicKey, subscribers)
	if err != nil {
		return err
	}
	return nil
}

func RemoveSubscriber(topicKey string, subscribers []string) error {
	err := novuClient.TopicsApi.RemoveSubscribers(ctx, topicKey, subscribers)
	if err != nil {
		return err
	}
	return nil
}

func SendTopicNotification(arg TriggerTopic) (*novu.EventResponse, error) {
	to := map[string]interface{}{
		"type":     "Topic",
		"topicKey": arg.TopicKey,
	}
	payload := map[string]interface{}{
		"name": arg.Title,
		"organization": map[string]interface{}{
			"logo": arg.Logo,
		},
	}

	resp, err := novuClient.EventApi.Trigger(ctx, arg.EventID, novu.ITriggerPayloadOptions{
		To:      to,
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}
