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
		event, ok2 := params[1].(onebot_v11_api_event.PrivateMessage)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in PrivateMessageHandler")
			return
		}

		logger("Receive private message from user ", event.UserId, " : ", event.RawMessage)

		newBotSystem.PrivateMessageHandler(logger, bot, &event)
	}

	// GroupMessage Handler
	handlers["message_group"] = func(params ...any) {
		if len(params) < 2 {
			logger("Error: insufficient parameters for GroupMessageHandler")
			return
		}

		bot, ok1 := params[0].(*onebot_v11_impl.V11Bot)
		event, ok2 := params[1].(onebot_v11_api_event.GroupMessage)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in GroupMessageHandler")
			return
		}

		logger("Receive group message from user ", event.UserId, " in group ", event.GroupId, " : ", event.RawMessage)

		newBotSystem.GroupMessageHandler(logger, bot, &event)
	}

	// PrivateRecall Handler
	handlers["notice_friend_recall"] = func(params ...any) {
		if len(params) < 2 {
			logger("Error: insufficient parameters for PrivateRecallHandler")
			return
		}

		bot, ok1 := params[0].(*onebot_v11_impl.V11Bot)
		event, ok2 := params[1].(onebot_v11_api_event.FriendMessageRecalled)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in PrivateRecallHandler")
			return
		}

		logger("Receive friend recalled event from user ", event.UserID, ", recalled message id : ", event.MessageID)

		newBotSystem.PrivateRecallHandler(logger, bot, &event)
	}

	// GroupRecall Handler
	handlers["notice_group_recall"] = func(params ...any) {
		if len(params) < 2 {
			logger("Error: insufficient parameters for GroupRecallHandler")
			return
		}

		bot, ok1 := params[0].(*onebot_v11_impl.V11Bot)
		event, ok2 := params[1].(onebot_v11_api_event.GroupMessageRecalled)
		if !ok1 || !ok2 {
			logger("Error: parameter type mismatch in GroupRecallHandler")
			return
		}

		logger("Receive group recalled event operated by user ", event.OperatorID, " in group ", event.GroupID, ", recalled message ", event.MessageID, " originally sent by ", event.UserID)

		newBotSystem.GroupRecallHandler(logger, bot, &event)
	}

	return handlers, nil
}
