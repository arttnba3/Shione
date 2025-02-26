package onebot_v11_plugin

import (
	"Shione/config"
	"Shione/tools"
	"errors"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	onebot_v11_api_message "github.com/arttnba3/Shigure-Bot/api/onebot/v11/message"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
	"strconv"
	"sync"
	"time"
)

var (
	DEFAULT_MIN_REQ_INTERVAL = 10
)

type ChatPlugin struct {
	OneBotV11Plugin
	botConfig        config.BotConfig
	lastChat         time.Time
	reqLock          sync.Mutex
	model            string
	provider         string
	url              string
	prompt           string
	headers          map[string]interface{}
	min_req_interval time.Duration
}

func (this *ChatPlugin) Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error {
	this.botConfig = botConfig
	this.lastChat = time.Unix(0, 0)
	this.reqLock = sync.Mutex{}

	chatConfig, ok := botConfig.Config.(map[string]interface{})["chat"].(map[string]interface{})
	if !ok {
		return errors.New("chat config is invalid")
	}

	this.provider, ok = chatConfig["provider"].(string)
	if !ok {
		return errors.New("chat provider is invalid")
	}

	this.model, ok = chatConfig["model"].(string)
	if !ok {
		return errors.New("chat model is invalid")
	}

	this.url, ok = chatConfig["url"].(string)
	if !ok {
		return errors.New("chat url is invalid")
	}

	this.prompt, ok = chatConfig["prompt"].(string)
	if !ok {
		return errors.New("chat prompt is invalid")
	}

	this.headers, ok = chatConfig["headers"].(map[string]interface{})
	if !ok {
		this.headers = nil
	}

	minReqInterval, ok := chatConfig["min_req_interval"].(float64)
	if !ok {
		this.min_req_interval = time.Duration(DEFAULT_MIN_REQ_INTERVAL) * time.Second
	} else {
		this.min_req_interval = time.Duration(minReqInterval) * time.Second
	}

	return nil
}

func (this *ChatPlugin) GetHelpMsg() string {
	return "Usage:\n" +
		"\n" +
		"    /chat [what you'd like to say]\n" +
		"\n" +
		"Redundant arguments will be automatically ignored\n"
}

func (this *ChatPlugin) ChatOperator(rawMessage string) (interface{}, bool) {
	var replyMsg string
	var err error

	params := tools.TextParser(rawMessage)
	if len(params) == 0 {
		return nil, false
	}

	switch params[0] {
	case this.command:
		if len(params) == 1 {
			replyMsg = this.GetHelpMsg()
			break
		}

		if time.Now().Sub(this.lastChat) < this.min_req_interval {
			replyMsg = "Request frequency is too fast."
			break
		}

		this.reqLock.Lock()
		defer this.reqLock.Unlock()
		this.lastChat = time.Now()

		replyMsg, err = tools.ChatWithAIText(this.provider, this.url, this.model, this.prompt, this.headers, rawMessage[len(params[0]):])
		if err != nil {
			replyMsg = "Error occur while requesting the model:  " + err.Error()
			break
		}

		break
	default:
		return nil, false
	}

	return onebot_v11_api_message.MessageArray{
		onebot_v11_api_message.MessageSegment{
			Type: "text",
			Data: onebot_v11_api_message.MessageSegmentDataText{
				Text: replyMsg,
			},
		},
	}, true
}

func (this *ChatPlugin) HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) int64 {
	replyMsg, intercept := this.ChatOperator(privateMsgEvent.RawMessage)
	if intercept {
		_, err := bot.SendPrivateMsg(privateMsgEvent.UserId, replyMsg, false)
		if err != nil {
			logger("Unable to send message to private: ", privateMsgEvent.UserId, ", error:", err.Error())
			return EVENT_IGNORE
		}
		return EVENT_INTERCEPT
	}

	return EVENT_IGNORE
}

func (this *ChatPlugin) HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) int64 {
	var replyMsg onebot_v11_api_message.MessageArray

	ollamaReplyMsgRaw, intercept := this.ChatOperator(groupMsgEvent.RawMessage)
	ollamaReplyMsg, _ := ollamaReplyMsgRaw.(onebot_v11_api_message.MessageArray)
	if intercept {
		if groupMsgEvent.Anonymous == nil {
			replyMsg = onebot_v11_api_message.MessageArray{
				onebot_v11_api_message.MessageSegment{
					Type: "at",
					Data: onebot_v11_api_message.MessageSegmentDataAt{
						QQ: strconv.FormatInt(groupMsgEvent.UserId, 10),
					},
				},
			}
			replyMsg = append(replyMsg, ollamaReplyMsg...)
		} else {
			replyMsg = ollamaReplyMsg
		}

		_, err := bot.SendGroupMsg(groupMsgEvent.GroupId, replyMsg, false)
		if err != nil {
			logger("Unable to send message to group: ", groupMsgEvent.GroupId, ", error:", err.Error())
			return EVENT_IGNORE
		}
		return EVENT_INTERCEPT
	}

	return EVENT_IGNORE
}
