package watermill

import "github.com/ThreeDotsLabs/watermill/message"

type Agent struct {
	SubscriberFactory SubscriberFactory
	Publisher         message.Publisher
}
