package onebot_v11_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/message"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
)

type AntiRecallPlugin struct {
	OneBotV11Plugin
}

func (this *AntiRecallPlugin) HandlePrivateRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.FriendMessageRecalled) int64 {
	recalledMessageData, err := bot.GetMsg(int32(recalledEvent.MessageID))
	if err != nil {
		logger("Unable to get recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	recalledMsgJson, err := json.Marshal(recalledMessageData)
	if err != nil {
		logger("Unable to marshal recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	var recalledMessage onebot_v11_api_message.MessageArray
	err = json.Unmarshal(recalledMsgJson, &recalledMessage)
	if err != nil {
		logger("Unable to unmarshal recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	replyMessage := onebot_v11_api_message.MessageArray{
		onebot_v11_api_message.MessageSegment{
			Type: "text",
			Data: onebot_v11_api_message.MessageSegmentDataText{
				Text: fmt.Sprintf(
					"Detected a message has been recalled by %v\nOriginal message: \n",
					recalledEvent.UserID),
			},
		},
	}

	// Note that if there's picture in the msg, it'll fail to send, maybe it's the problem of v11 impl backend :(
	replyMessage = append(replyMessage, recalledMessage...)

	_, err = bot.SendPrivateMsg(recalledEvent.UserID, replyMessage, false)
	if err != nil {
		logger(fmt.Sprintf("Unable to send message to user %v, error: %v", recalledEvent.UserID, err))
		return EVENT_IGNORE
	}

	return EVENT_INTERCEPT
}

func (this *AntiRecallPlugin) HandleGroupRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.GroupMessageRecalled) int64 {
	recalledMessageData, err := bot.GetMsg(int32(recalledEvent.MessageID))
	if err != nil {
		logger("Unable to get recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	recalledMsgJson, err := json.Marshal(recalledMessageData)
	if err != nil {
		logger("Unable to marshal recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	var recalledMessage onebot_v11_api_message.MessageArray
	err = json.Unmarshal(recalledMsgJson, &recalledMessage)
	if err != nil {
		logger("Unable to unmarshal recalled message, id: ", recalledEvent.MessageID, ", Error: ", err)
		return EVENT_IGNORE
	}

	replyMessage := onebot_v11_api_message.MessageArray{
		onebot_v11_api_message.MessageSegment{
			Type: "text",
			Data: onebot_v11_api_message.MessageSegmentDataText{
				Text: fmt.Sprintf(
					"Detected a message has been recalled by %v, originally sent by %v\nOriginal message: \n",
					recalledEvent.OperatorID,
					recalledEvent.UserID),
			},
		},
	}

	// Note that if there's picture in the msg, it'll fail to send, maybe it's the problem of v11 impl backend :(
	replyMessage = append(replyMessage, recalledMessage...)

	_, err = bot.SendGroupMsg(recalledEvent.GroupID, replyMessage, false)
	if err != nil {
		logger(fmt.Sprintf("Unable to send message to group %v, error: %v", recalledEvent.GroupID, err))
		return EVENT_IGNORE
	}

	return EVENT_INTERCEPT
}
