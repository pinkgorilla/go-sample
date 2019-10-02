package managed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// NewGStream ...
func NewGStream(topicID, subscriptionID, projectID, credential string) Stream {
	ctx, cancel := context.WithCancel(context.Background())
	credentials, err := google.CredentialsFromJSON(ctx, []byte(credential), pubsub.ScopePubSub)
	if err != nil {
		panic(err)
	}

	client, err := pubsub.NewClient(
		ctx,
		projectID,
		option.WithCredentials(credentials),
	)
	if err != nil {
		panic(err)
	}
	// ensure topic exists
	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			panic(err)
		}
	}

	// ensure subscription
	subscription := client.Subscription(subscriptionID)
	exists, err = subscription.Exists(ctx)
	config := pubsub.SubscriptionConfig{
		Topic: topic,
	}
	if err != nil {
		panic(err)
	}
	if !exists {
		subscription, err = client.CreateSubscription(ctx, subscriptionID, config)
		if err != nil {
			panic(err)
		}
	} else {
		con, err := subscription.Config(ctx)
		if err != nil {
			panic(err)
		}
		boundTopicID := con.Topic.ID()
		if boundTopicID != config.Topic.ID() {
			panic(fmt.Errorf("failed to bind subscription:'%s' to topic:'%s', already bound to topic:'%s'", subscription.ID(), config.Topic.ID(), boundTopicID))
		}
	}

	return &GStream{
		client:       client,
		ctx:          ctx,
		cancel:       cancel,
		topic:        topic,
		subscription: subscription,
		ch:           make(chan interface{}, 10),
		once:         sync.Once{},
	}
}

// GStream implementations of Stream using google pub/sub
type GStream struct {
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	ctx          context.Context
	cancel       context.CancelFunc
	client       *pubsub.Client
	ch           chan interface{}
	once         sync.Once
}

// Push pushes data to stream
func (s *GStream) Push(data interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	result := s.topic.Publish(ctx, &pubsub.Message{
		Data: bs,
	})
	_, err = result.Get(ctx)
	return err
}

func (s *GStream) receive() error {
	return s.subscription.Receive(s.ctx, func(ctx context.Context, msg *pubsub.Message) {
		s.ch <- msg.Data
		msg.Ack()
		log.Println("pull")
	})
}

// Pull pulls data from stream, returns []byte,error
func (s *GStream) Pull() (interface{}, error) {
	s.once.Do(func() {
		go s.receive()
	})

	select {
	case data := <-s.ch:
		return data, nil
	default:
		return nil, nil
	}
	// var bs []byte
	// // ctx, _ := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	// ctx := context.Background()
	// err := s.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
	// 	bs = msg.Data
	// 	msg.Ack()
	// 	log.Println("pull")
	// })
	// // cancel()
	// if err != nil {
	// 	log.SetFlags(log.LstdFlags)
	// 	log.Print("Google PubSub Stream, error occured when receiving message - ", err, "\n")
	// }
	// if len(bs) == 0 {
	// 	return nil, err
	// }
	// return bs, err
}

// Dispose disposes instance
func (s *GStream) Dispose() {
	s.topic.Stop()
	s.cancel()
}
