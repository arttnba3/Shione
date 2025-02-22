package onebot_v11_plugin

import (
	"Shione/config"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
	"sync"
)

type OneBotV11PluginSystem struct {
	PluginLock sync.Mutex
	Plugins    *[]OneBotV11Plugin
}

type OneBotV11EventHandlers interface {
	HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) bool
	HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) bool
	// TODO: add more handler for more different events
}

type OneBotV11PluginOperations interface {
	Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error
	IsEnabled() bool
	Enable()
	Disable()
	GetName() string
}

type OneBotV11Plugin interface {
	OneBotV11PluginOperations
	OneBotV11EventHandlers
}

type OneBotV11PluginConfig struct {
	Name      string
	Command   string
	isEnabled bool
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
	this.PluginLock.Lock()
	defer this.PluginLock.Unlock()

	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		passToNext := plugin.HandlePrivateMessage(logger, bot, privateMsgEvent)
		if !passToNext {
			break
		}
	}
}

func (this *OneBotV11PluginSystem) GroupMessageHandler(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) {
	this.PluginLock.Lock()
	defer this.PluginLock.Unlock()

	for _, plugin := range *this.Plugins {
		if !plugin.IsEnabled() {
			continue
		}

		passToNext := plugin.HandleGroupMessage(logger, bot, groupMsgEvent)
		if !passToNext {
			break
		}
	}
}
