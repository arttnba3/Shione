package onebot_v11_plugin

func GetDefaultOneBotV11PluginsConfig() *[]OneBotV11PluginAPI {
	return &[]OneBotV11PluginAPI{
		&PluginManagerPlugin{
			OneBotV11Plugin: OneBotV11Plugin{
				name:      "PluginManager",
				command:   "/plugin",
				isEnabled: true,
			},
		},
		&ChatPlugin{
			OneBotV11Plugin: OneBotV11Plugin{
				name:      "ChatPlugin",
				command:   "/chat",
				isEnabled: true,
			},
		},
		&RepeaterPlugin{
			OneBotV11Plugin: OneBotV11Plugin{
				name:      "RepeaterPlugin",
				command:   "",
				isEnabled: true,
			},
		},
		&AntiRecallPlugin{
			OneBotV11Plugin: OneBotV11Plugin{
				OneBotV11PluginAPI: nil,
				name:               "AntiRecallPlugin",
				command:            "",
				isEnabled:          true,
			},
		},
	}
}
