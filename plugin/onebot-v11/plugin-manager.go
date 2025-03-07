package onebot_v11_plugin

import (
	"Shione/config"
	"Shione/tools"
	"encoding/json"
	"fmt"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/event"
	"github.com/arttnba3/Shigure-Bot/api/onebot/v11/message"
	"github.com/arttnba3/Shigure-Bot/bot/onebot/v11"
)

type PluginManagerPlugin struct {
	OneBotV11Plugin
	config.BotConfig
	pluginSystem *OneBotV11PluginSystem
}

func (this *PluginManagerPlugin) Init(botConfig config.BotConfig, pluginSystem *OneBotV11PluginSystem) error {
	this.BotConfig = botConfig
	this.pluginSystem = pluginSystem

	return nil
}

func (this *PluginManagerPlugin) GetHelpMsg() string {
	return fmt.Sprintf(
		"Usage:\n\n"+
			"    %v [operations] [plugin_Name]\n"+
			"\nRedundant arguments will be automatically ignored\n",
		this.command,
	)
}

func ParsePureTextMessage(logger func(...any), msgArray onebot_v11_api_message.MessageArray) (string, bool) {
	var pureMsgText string

	// While one of which segment is not text, ignore
	for _, msgSeg := range msgArray {
		if msgSeg.Type != "text" {
			return "", false
		}

		messageText, ok := msgSeg.Data.(map[string]interface{})["text"].(string)
		if !ok {
			logger("[Internal Error] Invalid message from backend: ", msgArray)
			return "", false
		}

		pureMsgText += messageText
	}

	return pureMsgText, true
}

func (this *PluginManagerPlugin) PermissionCheck(user int64) bool {
	for _, admin := range this.AdminList {
		if admin == user {
			return true
		}
	}

	return false
}

func (this *PluginManagerPlugin) PluginOperationExecutor(logger func(...any), pureMsgText string, user int64) (string, bool) {
	var messageToSend string

	if pureMsgText == this.command {
		if !this.PermissionCheck(user) {
			messageToSend = "Permission denied."
			goto out
		} else {
			messageToSend = this.GetHelpMsg()
		}
	} else if len(pureMsgText) > len(this.command) &&
		pureMsgText[:len(this.command)] == this.command &&
		pureMsgText[len(this.command)] == ' ' {
		if !this.PermissionCheck(user) {
			messageToSend = "Permission denied."
			goto out
		}

		command := tools.TextParser(pureMsgText)
		if len(command) > 1 {
			switch command[1] {
			case "list":
				if len(command) > 2 {
					if command[2] == "all" {
						messageToSend = "All Plugins:\n"
						for idx, plugin := range *this.pluginSystem.Plugins {
							messageToSend += fmt.Sprintf("  %v. %v", idx+1, plugin.GetName())
							if plugin.IsEnabled() {
								messageToSend += "[Enabled]\n"
							}
						}
					} else {
						messageToSend = "Unrecognized second-level sub-command: " + command[2]
					}
				} else {
					messageToSend = "Loaded Plugins:\n"
					enableCount := 1
					for _, plugin := range *this.pluginSystem.Plugins {
						if plugin.IsEnabled() {
							messageToSend += fmt.Sprintf("  %v. %v\n", enableCount, plugin.GetName())
							enableCount++
						}
					}
				}

				break
			case "help":
				messageToSend = this.GetHelpMsg()
				break
			case "load":
				fallthrough
			case "unload":
				if len(command) == 2 {
					messageToSend = this.GetHelpMsg()
					break
				}

				pluginName := command[2]
				var plugin *OneBotV11PluginAPI = nil
				for _, candidate := range *this.pluginSystem.Plugins {
					if candidate.GetName() == pluginName {
						plugin = &candidate
						break
					}
				}

				if plugin == nil {
					messageToSend = "Unable to find plugin " + pluginName
					break
				}

				if command[1] == "load" {
					if (*plugin).IsEnabled() {
						messageToSend = "Plugin has already been loaded"
					} else {
						(*plugin).Enable()
						messageToSend = "Operation complete."
					}
				} else if command[1] == "unload" {
					if (*plugin).IsEnabled() {
						(*plugin).Disable()
						messageToSend = "Operation complete."
					} else {
						messageToSend = "Plugin has already been unloaded"
					}
				} else {
					logger(fmt.Sprintf(
						"[Internal Error] command for %v got changed in the procedure, new value: %v",
						this.command,
						command[1]),
					)
					messageToSend = "Internal error, please check the log."
				}

				break
			default:
				messageToSend = "Unrecognized sub-command: " + command[1]
			}
		} else {
			messageToSend = this.GetHelpMsg()
		}
	} else {
		return "", false
	}

out:

	return messageToSend, true
}

func (this *PluginManagerPlugin) HandlePrivateMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, privateMsgEvent *onebot_v11_api_event.PrivateMessage) int64 {
	var pureMsgText string
	var messageArray onebot_v11_api_message.MessageArray
	var messageString onebot_v11_api_message.MessageString
	var ok bool

	messageJson, err := json.Marshal(privateMsgEvent.Message)
	if err != nil {
		logger("Unable to marshal private message: ", privateMsgEvent.Message)
		return EVENT_IGNORE
	}

	err = json.Unmarshal(messageJson, &messageArray)
	if err == nil {
		pureMsgText, ok = ParsePureTextMessage(logger, messageArray)
		if !ok {
			return EVENT_IGNORE
		}
	} else {
		err = json.Unmarshal(messageJson, &messageString)
		if err != nil {
			logger("Unable to re-unmarshal private message: ", privateMsgEvent.Message)
			return EVENT_IGNORE
		}
		pureMsgText = messageString.Message
	}

	messageToSend, sendMessage := this.PluginOperationExecutor(logger, pureMsgText, privateMsgEvent.UserId)

	if sendMessage {
		_, err := bot.SendPrivateMsg(privateMsgEvent.UserId, messageToSend, false)
		if err != nil {
			logger("Unable to send message to user ", privateMsgEvent.UserId, " , error: ", err.Error())
		}

		return EVENT_INTERCEPT
	}

	return EVENT_IGNORE
}

func (this *PluginManagerPlugin) HandleGroupMessage(logger func(...any), bot *onebot_v11_impl.V11Bot, groupMsgEvent *onebot_v11_api_event.GroupMessage) int64 {
	var pureMsgText string
	var messageArray onebot_v11_api_message.MessageArray
	var messageString onebot_v11_api_message.MessageString
	var ok bool

	messageJson, err := json.Marshal(groupMsgEvent.Message)
	if err != nil {
		logger("Unable to marshal group message: ", groupMsgEvent.Message)
		return EVENT_IGNORE
	}

	err = json.Unmarshal(messageJson, &messageArray)
	if err == nil {
		pureMsgText, ok = ParsePureTextMessage(logger, messageArray)
		if !ok {
			return EVENT_IGNORE
		}
	} else {
		err = json.Unmarshal(messageJson, &messageString)
		if err != nil {
			logger("Unable to re-unmarshal group message: ", groupMsgEvent.Message)
			return EVENT_IGNORE
		}
		pureMsgText = messageString.Message
	}

	messageToSend, sendMessage := this.PluginOperationExecutor(logger, pureMsgText, groupMsgEvent.UserId)

	if sendMessage {
		_, err = bot.SendGroupMsg(groupMsgEvent.GroupId, messageToSend, false)
		if err != nil {
			logger("Unable to send message to group ", groupMsgEvent.GroupId, " , error: ", err.Error())
		}

		return EVENT_INTERCEPT
	}

	return EVENT_IGNORE
}
