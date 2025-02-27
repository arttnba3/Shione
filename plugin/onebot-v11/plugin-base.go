package onebot_v11_plugin

import (
	"Shione/config"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
	"sync"
)

type OneBotV11PluginSystem struct {
	PluginLock sync.Mutex
	Plugins    *[]OneBotV11PluginAPI
}

type OneBotV11EventHandlers interface {
	HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) int64
	HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) int64
	HandlePrivateRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.FriendMessageRecalled) int64
	HandleGroupRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.GroupMessageRecalled) int64
	// TODO: add more handler for more different events
}

const (
	EVENT_IGNORE    = 0x10000
	EVENT_INTERCEPT = 0x10001
)

type OneBotV11PluginOperations interface {
	Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error
	IsEnabled() bool
	Enable()
	Disable()
	GetName() string
	GetHelpMsg() string
}

type OneBotV11PluginAPI interface {
	OneBotV11PluginOperations
	OneBotV11EventHandlers
}

type OneBotV11Plugin struct {
	OneBotV11PluginAPI
	name      string
	command   string
	isEnabled bool
	lock      sync.Mutex
}

func (this *OneBotV11Plugin) Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error {
	return nil
}

func (this *OneBotV11Plugin) IsEnabled() bool {
	return this.isEnabled
}

func (this *OneBotV11Plugin) Enable() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.isEnabled = true
}

func (this *OneBotV11Plugin) Disable() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.isEnabled = false
}

func (this *OneBotV11Plugin) GetName() string {
	return this.name
}

func (this *OneBotV11Plugin) GetHelpMsg() string {
	return "(not defined)"
}

func (this *OneBotV11Plugin) HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) int64 {
	return EVENT_IGNORE
}

func (this *OneBotV11Plugin) HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) int64 {
	return EVENT_IGNORE
}

func (this *OneBotV11Plugin) HandlePrivateRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.FriendMessageRecalled) int64 {
	return EVENT_IGNORE
}
func (this *OneBotV11Plugin) HandleGroupRecall(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.GroupMessageRecalled) int64 {
	return EVENT_IGNORE
}

func (this *OneBotV11PluginSystem) PluginSystemInit(botConfig config.BotConfig) error {
	pluginConfig := GetDefaultOneBotV11PluginsConfig()

	for _, plugin := range *pluginConfig {
		err := plugin.Init(botConfig, this)
		if err != nil {
			return err
		}
	}

	this.Plugins = pluginConfig

	return nil
}

func (this *OneBotV11PluginSystem) PrivateMessageHandler(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) {
	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		if plugin.HandlePrivateMessage(logger, bot, privateMsgEvent) != EVENT_IGNORE {
			break
		}
	}
}

func (this *OneBotV11PluginSystem) GroupMessageHandler(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) {
	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		if plugin.HandleGroupMessage(logger, bot, groupMsgEvent) != EVENT_IGNORE {
			break
		}
	}
}

func (this *OneBotV11PluginSystem) PrivateRecallHandler(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.FriendMessageRecalled) {
	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		if plugin.HandlePrivateRecall(logger, bot, recalledEvent) != EVENT_IGNORE {
			break
		}
	}
}

func (this *OneBotV11PluginSystem) GroupRecallHandler(logger func(...any), bot *onebot_v11_impl.V11Bot, recalledEvent *onebot_v11_api_event.GroupMessageRecalled) {
	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		if plugin.HandleGroupRecall(logger, bot, recalledEvent) != EVENT_IGNORE {
			break
		}
	}
}
