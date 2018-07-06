package utility

import (
	"fmt"
	"github.com/nlopes/slack"

	"time"
	"github.com/richardsang2008/BotManager/model"

)

type SlackUtility struct {
}

func (s *SlackUtility) SendMessage(msg model.BotMessage, api *slack.Client) string {
	if msg.ChannelID == "" {
		_, _, channelID, err := api.OpenIMChannel(msg.UserID)
		if err != nil {
			panic(err)
		}
		msg.ChannelID = channelID
	}
	channelID, timestamp, err := api.PostMessage(msg.ChannelID, msg.Message, slack.PostMessageParameters{Username: "lisa"})
	if err != nil {
		panic(err)
	}
	ret := fmt.Sprintf("Message successfully sent to channel %s at %s", channelID, timestamp)
	MLog.Debug(ret)
	return ret
}
func (s *SlackUtility) DeleteMessage(channel string, ts string, api *slack.Client) error {
	_, _, err := api.DeleteMessage(channel, ts)
	if err != nil {
		MLog.Error(err)

	}
	time.Sleep(1 * time.Second)
	return err
}
func (s *SlackUtility) DeleteFile(fileId string, api *slack.Client) error {
	err := api.DeleteFile(fileId)
	if err != nil {
		MLog.Error(err)

	}
	time.Sleep(1 * time.Second)
	return err
}

func (s *SlackUtility) GetUserInfo(userId string, api *slack.Client) (*model.SlackUser, error) {
	user, err := api.GetUserInfo(userId)
	if err != nil {
		MLog.Error(err)
		return nil, err
	} else {
		slackUser := model.SlackUser{}
		slackUser.ReferenceID = user.ID
		slackUser.Deleted = user.Deleted
		slackUser.IsAdmin = user.IsAdmin
		slackUser.IsBot = user.IsBot
		slackUser.IsOwner = user.IsOwner
		slackUser.Name = user.Name
		slackUser.DisplayName = user.Profile.DisplayName
		slackUser.Email = user.Profile.Email
		slackUser.FirstName = user.Profile.FirstName
		slackUser.LastName = user.Profile.LastName
		slackUser.Phone = user.Profile.Phone
		slackUser.RealName = user.RealName
		slackUser.StatusText = user.Profile.StatusText
		return &slackUser, nil
	}
}
