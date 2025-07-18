package config

type BotConfig struct {
	Type      string      `json:"type"`
	AdminList []int64     `json:"admin_list"`
	BotQQ     string      `json:"bot_qq"`
	Config    interface{} `json:"config"`
}

type ShioneConfig struct {
	FileLogOutput    bool        `json:"file_log_output"`
	LogDir           string      `json:"log_dir"`
	LogFlushFreq     uint64      `json:"log_flush_freq"`
	BotConfig        []BotConfig `json:"bot_config"`
	ConsoleLogOutput bool        `json:"console_log_output"`
}

var AppVersion string = "0.0.1"
var DefaultConfigFile string = "./config.json"
var DefaultLogDir string = "./local_logs"
var DefaultLogFlushFreq uint64 = 4096
var DefaultConsoleLogOutput = true
