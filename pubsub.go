package pubsub

type chanMapStringList map[chan interface{}][]string
type stringMapChanList map[string][]chan interface{}

type Pubsub struct {
	capacity int // capacity for each chan buffer

	clientMapTopics chanMapStringList // map to store "chan -> topic List" for find subcription

	topicMapClients stringMapChanList // map to store "topic -> chan List" for publish
}

// Crate a pubsub with expect initial size, but the size could be extend.
func NewPubsub(initChanCapacity int) *Pubsub {
	initClientMapTopics := make(chanMapStringList)
	initTopicMapClients := make(stringMapChanList)

	server := Pubsub{capacity: initChanCapacity, clientMapTopics: initClientMapTopics, topicMapClients: initTopicMapClients}
	return &server
}

// Subscribe channels, the channels could be a list of channels name
// The channel name could be any, without define in server
func (p *Pubsub) Subscribe(topics ...string) chan interface{} {
	// initial new chan using capacity as channel buffer
	workChan := make(chan interface{}, p.capacity)
	p.updateTopicMapClient(workChan, topics)
	return workChan
}

// Publish a content to a list of channels
// The content could be any type.
func (p *Pubsub) Publish(content interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapClients[topic]; ok {
			// Someone has subscribed this topic.
			for _, channel := range chanList {
				channel <- content
			}
		}
	}
}

func (p *Pubsub) updateTopicMapClient(clientChan chan interface{}, topics []string) {
	var updateChanList []chan interface{}
	for _, topic := range topics {
		updateChanList = p.topicMapClients[topic]
		updateChanList = append(updateChanList, clientChan)
		p.topicMapClients[topic] = updateChanList
	}
	p.clientMapTopics[clientChan] = topics
}

// AddSubscription: Add a new topic subscribed to specific client channel.
func (p *Pubsub) AddSubscription(clientClient chan interface{}, topic ...string) {
	p.updateTopicMapClient(clientClient, topic)
}

// RemoveSubscription: Remove subscribed topic list on specific client channel.
func (p *Pubsub) RemoveSubscription(clientChan chan interface{}, topics ...string) {
	for _, topic := range topics {
		// Remove from topic->chan map
		if chanList, ok := p.topicMapClients[topic]; ok {
			var updateChanList []chan interface{}
			for _, client := range chanList {
				if client != clientChan {
					updateChanList = append(updateChanList, client)
				}
			}
			p.topicMapClients[topic] = updateChanList
		}

		// Remove from chan->topic map
		if topicList, ok := p.clientMapTopics[clientChan]; ok {
			var updateTopicList []string
			for _, pic := range topicList {
				if pic != topic {
					updateTopicList = append(updateTopicList, pic)
				}
			}
			p.clientMapTopics[clientChan] = updateTopicList
		}
	}
}
