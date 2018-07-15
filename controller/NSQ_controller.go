package controller

import(
	"github.com/bitly/go-nsq"
	"github.com/pkg/errors"
	"sync"
	"fmt"
	"github.com/richardsang2008/BotManager/utility"

	"encoding/json"
	"github.com/richardsang2008/BotManager/model"
)

type NSQController struct {
	Producer *nsq.Producer
	Consumer *nsq.Consumer

}
type SlackMessageHandler struct{}
func (h *SlackMessageHandler) HandleMessage(message *nsq.Message) error{
	if len(message.Body) ==0{
		return errors.New("body is blank re-enqueue message")
	}
	k:= string(message.Body)
	str:=fmt.Sprintf("NSQ message received: %s", k)
	var slackMessage  model.SlackDBMessage
	json.Unmarshal([]byte(k), &slackMessage)
	var regionId uint
	regionId =1
	err:=Data.InsertSlackMessage(regionId, slackMessage.ChannelId,slackMessage.Ts)
	if err != nil {
		utility.MLog.Error(err)
	}
	utility.MLog.Debug(str)
	return nil
}

func (s *NSQController) InitNSQ(producerAddress string, consumerLookupAddress string, consumerTopic string, consumerChannel string, wg *sync.WaitGroup ) {
	config := nsq.NewConfig()
	err := errors.New("")
	if s.Producer == nil {
		s.Producer,err = nsq.NewProducer(producerAddress, config)
		if err != nil  {
			utility.MLog.Panic(err)
		}
	}
	if s.Consumer == nil {
		s.Consumer,err = nsq.NewConsumer(consumerTopic, consumerChannel, config)
		if err != nil  {
			utility.MLog.Panic(err)
		}
		wg.Add(1)
		s.Consumer.ChangeMaxInFlight(200)
		s.Consumer.AddConcurrentHandlers(&SlackMessageHandler{},20)
		err = s.Consumer.ConnectToNSQLookupd(consumerLookupAddress)
		if err != nil {
			utility.MLog.Panic("Could not connect")
		}
		utility.MLog.Info("Awaiting message from topic ...")
		wg.Wait()
	}
}

func (s *NSQController) ProducerPublishMessage(body []byte, topic string){
	err := s.Producer.Publish(topic, body)
	if err != nil  {
		utility.MLog.Panic(err)
	}
}
