package onebot_v11_plugin

func GetDefaultOneBotV11PluginsConfig() *[]OneBotV11Plugin {
	return &[]OneBotV11Plugin{
		&PluginManagerPlugin{
			OneBotV11PluginConfig: OneBotV11PluginConfig{
				Name:      "PluginManager",
				Command:   "/plugin",
				isEnabled: true,
			},
		},
		&RepeaterPlugin{
			OneBotV11PluginConfig: OneBotV11PluginConfig{
				Name:      "RepeaterPlugin",
				Command:   "",
				isEnabled: true,
			},
		},
	}
}
