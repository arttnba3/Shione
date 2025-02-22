package onebot_v11_plugin

import (
	"Shione/config"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
	"sync"
)

type repeatPrivate struct {
	HasRepeated bool
	Message     *onebot_v11_api_event.PrivateMessage
}

type repeatGroup struct {
	HasRepeated bool
	Message     *onebot_v11_api_event.GroupMessage
}

type RepeaterPlugin struct {
	OneBotV11PluginConfig
	LastPrivateMessage map[int64]*repeatPrivate
	PrivateMessageLock sync.Mutex
	LastGroupMessage   map[int64]*repeatGroup
	GroupMessageLock   sync.Mutex
}

func (this *RepeaterPlugin) Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error {
	this.LastPrivateMessage = make(map[int64]*repeatPrivate)
	this.LastGroupMessage = make(map[int64]*repeatGroup)
	return nil
}

func (this *RepeaterPlugin) IsEnabled() bool {
	return this.isEnabled
}

func (this *RepeaterPlugin) Enable() {
	this.isEnabled = true
}

func (this *RepeaterPlugin) Disable() {
	this.isEnabled = false
}

func (this *RepeaterPlugin) GetName() string {
	return this.Name
}

func (this *RepeaterPlugin) HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) bool {
	this.PrivateMessageLock.Lock()
	defer this.PrivateMessageLock.Unlock()

	lastMsg, ok := this.LastPrivateMessage[privateMsgEvent.UserId]
	if !ok {
		this.LastPrivateMessage[privateMsgEvent.UserId] = &repeatPrivate{HasRepeated: false, Message: privateMsgEvent}
		return true
	}

	if lastMsg.Message.RawMessage == privateMsgEvent.RawMessage {
		if lastMsg.HasRepeated {
			return false
		}

		_, err := bot.SendPrivateMsg(privateMsgEvent.UserId, privateMsgEvent.Message, false)
		if err != nil {
			logger("Unable to send message to private: ", privateMsgEvent.UserId, ", error:", err.Error())
			return true
		}

		lastMsg.HasRepeated = true
	} else {
		this.LastPrivateMessage[privateMsgEvent.UserId] = &repeatPrivate{HasRepeated: false, Message: privateMsgEvent}
	}

	return true
}

func (this *RepeaterPlugin) HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) bool {
	this.GroupMessageLock.Lock()
	defer this.GroupMessageLock.Unlock()

	lastMsg, ok := this.LastGroupMessage[groupMsgEvent.GroupId]
	if !ok {
		this.LastGroupMessage[groupMsgEvent.GroupId] = &repeatGroup{HasRepeated: false, Message: groupMsgEvent}
		return true
	}

	if lastMsg.Message.RawMessage == groupMsgEvent.RawMessage {
		if lastMsg.HasRepeated {
			return false
		}

		_, err := bot.SendGroupMsg(groupMsgEvent.GroupId, groupMsgEvent.Message, false)
		if err != nil {
			logger("Unable to send message to group: ", groupMsgEvent.GroupId, ", error:", err.Error())
			return true
		}

		lastMsg.HasRepeated = true
	} else {
		this.LastGroupMessage[groupMsgEvent.GroupId] = &repeatGroup{HasRepeated: false, Message: groupMsgEvent}
	}

	return true
}
