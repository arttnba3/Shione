package handlers

import (
	"Shione/config"
	"Shione/plugin/onebot-v11"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
)

func NewOneBotV11Handlers(logger func(...any), config config.BotConfig) (map[string]func(params ...any), error) {
	newBotSystem := &onebot_v11_plugin.OneBotV11PluginSystem{}

	err := newBotSystem.PluginSystemInit(config)
	if err != nil {
		logger("Unable to initialize OneBotV11 plugin system")
		return nil, err
	}

	handlers := make(map[string]func(params ...any))

	// PrivateMessage Handler
	handlers["message_private"] = func(params ...any) {
		if len(params) < 2 {
			logger("Error: insufficient parameters for PrivateMessageHandler")
			return
		}

		bot, ok1 := params[0].(*onebot_v11_impl.V11Bot)
		message, ok2 := params[1].(onebot_v11_api_event.PrivateMessage)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in PrivateMessageHandler")
			return
		}

		logger("Receive private message from user ", message.UserId, " : ", message.RawMessage)

		newBotSystem.PrivateMessageHandler(logger, bot, &message)
	}

	//
	handlers["message_group"] = func(params ...any) {
		if len(params) < 2 {
			logger("Error: insufficient parameters for GroupMessageHandler")
			return
		}

		bot, ok1 := params[0].(*onebot_v11_impl.V11Bot)
		message, ok2 := params[1].(onebot_v11_api_event.GroupMessage)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in GroupMessageHandler")
			return
		}

		logger("Receive group message from user ", message.UserId, " in group ", message.GroupId, " : ", message.RawMessage)

		newBotSystem.GroupMessageHandler(logger, bot, &message)
	}

	return handlers, nil
}
