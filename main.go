package main

import (
	"Shione/config"
	"Shione/handlers"
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/arttnba3/Shigure-Bot/bot"
	"os"
	"sync"
	"time"
)

func EarlyLoggerLn(params ...any) {
	fmt.Println("["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprint(params...))
}

func LogoPrint() {
	EarlyLoggerLn(
		"\n",
		"                                                                    \n",
		"                                                                    \n",
		"       █████╗    ██╗    ██╗  ██╗   ██████╗   ██╗   ██╗  ██████╗     \n",
		"      ██╔═══╝    ██║    ██║  ██║  ██╔═══██║  ████╗ ██║  ██╔═══╝     \n",
		"       ███████╗  █████████║  ██║  ██║   ██║  ██╔██╗██║  █████╗      \n",
		"       ╚════██║  ██╔════██║  ██║  ██║   ██║  ██║╚████║  ██╔══╝      \n",
		"      ███████╔╝  ██║    ██║  ██║  ╚██████╔╝  ██║ ╚═██║  ███████╗    \n",
		"      ╚══════╝   ╚═╝    ╚═╝  ╚═╝   ╚═════╝   ╚═╝   ╚═╝  ╚══════╝    \n",
		"                                                                    \n",
		"                                                                    \n",
		"                                                                    \n",
		"  Copyright(c) 2025 arttnba3                   「汐音(シオネ)」製作委員会\n",
		"\n")
}

func LoggerInit(config *config.ShioneConfig) (func(params ...any), error) {
	var ConsoleLogger func(params ...any)
	var fileLogLock sync.Mutex
	var consoleLogLock sync.Mutex

	if config.ConsoleLogOutput {
		ConsoleLogger = func(params ...any) {
			consoleLogLock.Lock()
			defer consoleLogLock.Unlock()

			logPrefix := fmt.Sprintf("[%v] ", time.Now().Format("2006-01-02 15:04:05"))
			logStr := logPrefix + fmt.Sprint(params...) + "\n"

			fmt.Print(logStr)
		}
	} else {
		ConsoleLogger = func(params ...any) {
			// do nothing
		}
	}

	// no directories or files will be chosen and created, only output to console
	if !config.FileLogOutput {
		return func(params ...any) {
			ConsoleLogger(params...)
		}, nil
	}

	logDirInfo, err := os.Stat(config.LogDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(config.LogDir, os.ModePerm)
			if err != nil {
				EarlyLoggerLn("Unable to create log directory, error:", err.Error())
				return nil, err
			}

			logDirInfo, err = os.Stat(config.LogDir)
		}

		if err != nil {
			EarlyLoggerLn("Unable to stat log directory, error:", err.Error())
			return nil, err
		}
	}

	if !logDirInfo.IsDir() {
		EarlyLoggerLn("Log directory path has been occupied and is not a directory")
		return nil, errors.New("log directory path has been occupied and is not a directory")
	}

	loggerStartTime := time.Now()
	logFilePath := config.LogDir + "/" + loggerStartTime.Format("2006-01-02_15-04-05") + ".log"

	logFile, err := os.Create(logFilePath)
	if err != nil {
		EarlyLoggerLn("Unable to create log file, error: ", err.Error())
		return nil, err
	}

	logFileWriter := bufio.NewWriter(logFile)

	var logCounter uint64 = 0
	logFlushFreq := config.LogFlushFreq

	return func(params ...any) {
		fileLogLock.Lock()
		defer fileLogLock.Unlock()

		logPrefix := fmt.Sprintf("[%v] ", time.Now().Format("2006-01-02 15:04:05"))
		logStr := logPrefix + fmt.Sprint(params...) + "\n"

		ConsoleLogger(params...)
		_, err = logFileWriter.WriteString(logStr)
		if err != nil {
			ConsoleLogger("Unable to write to log file, error:", err.Error())
		}

		logCounter++
		if logCounter%logFlushFreq == 0 {
			err = logFileWriter.Flush()
			if err != nil {
				ConsoleLogger("Unable to write to log file, error:", err.Error())
			}
			logCounter = 0
		}
	}, nil
}

func ConfigInit() (*config.ShioneConfig, error) {
	var isHelp bool
	var configFilePath string
	var shioneConfig config.ShioneConfig

	EarlyLoggerLn("Parsing commandline args...")

	flag.BoolVar(&isHelp, "h", false, "print help message and exit")
	flag.StringVar(&configFilePath, "c", "", "`path` of config file (default: "+config.DefaultConfigFile+")")
	flag.Parse()

	if isHelp {
		fmt.Println("Usage: ")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if configFilePath == "" {
		EarlyLoggerLn("`-c` param is not specified, use default value: ", config.DefaultConfigFile)
		configFilePath = config.DefaultConfigFile
	}

	EarlyLoggerLn(fmt.Sprintf("Loading config from file: %v", configFilePath))

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		EarlyLoggerLn(fmt.Sprintf("Error while reading config file: %v", err))
		return nil, err
	}

	err = json.Unmarshal(configData, &shioneConfig)
	if err != nil {
		EarlyLoggerLn(fmt.Sprintf("Error while parsing config file: %v", err))
		return nil, err
	}

	if shioneConfig.FileLogOutput {
		if shioneConfig.LogDir == "" {
			EarlyLoggerLn("`log_dir` param is not specified in config file, use default value: ", config.DefaultLogDir)
			shioneConfig.LogDir = config.DefaultLogDir
		} else {
			EarlyLoggerLn("Use log directory: ", shioneConfig.LogDir)
		}

		if shioneConfig.LogFlushFreq == 0 {
			EarlyLoggerLn("`log_flush_freq` param is not specified, use default value: ", config.DefaultLogFlushFreq)
			shioneConfig.LogFlushFreq = config.DefaultLogFlushFreq
		} else {
			EarlyLoggerLn("Use log flush frequency: ", shioneConfig.LogFlushFreq)
		}
	} else {
		EarlyLoggerLn("Log will not be output to log files")
	}

	if shioneConfig.ConsoleLogOutput {
		EarlyLoggerLn("Log will be output to console concurrently")
	} else {
		EarlyLoggerLn("Log will not be output to console")
	}

	return &shioneConfig, nil
}

func BotInit(logger func(params ...any), config config.BotConfig) (*shigure.ShigureBot, error) {
	var bot *shigure.ShigureBot
	var botHandlers map[string]func(params ...any)
	var err error
	var configJson []byte

	configJson, err = json.Marshal(config.Config)
	if err != nil {
		logger("Unable to marshal configuration into JSON, error:", err.Error())
		return nil, err
	}

	switch config.Type {
	case "OneBot-V11":
		botHandlers, err = handlers.NewOneBotV11Handlers(logger, config)
		if err != nil {
			break
		}
		bot, err = shigure.NewShigureBot(config.Type, configJson, logger, botHandlers)
		break
	default:
		logger("Unknown Bot Type: " + config.Type)
		return nil, errors.New("Unknown Bot Type: " + config.Type)
	}

	if err != nil {
		logger("Unable to create ShigureBot, error:", err.Error())
		return nil, err
	}

	return bot, nil
}

func BotGroupInit(logger func(params ...any), config *config.ShioneConfig) (*[]*shigure.ShigureBot, error) {
	var botGroup []*shigure.ShigureBot

	for idx, botConfig := range config.BotConfig {
		bot, err := BotInit(logger, botConfig)
		if err != nil {
			logger(fmt.Sprintf("Error while creating %v bot, error: %v", idx, err))
			return nil, err
		}
		botGroup = append(botGroup, bot)
	}

	return &botGroup, nil
}

func main() {
	LogoPrint()

	EarlyLoggerLn("Shione chat bot application start, version:", config.AppVersion)

	shoneConfig, err := ConfigInit()
	if err != nil {
		EarlyLoggerLn(fmt.Sprintf("Error initializing config: %v", err.Error()))
		return
	}

	logger, err := LoggerInit(shoneConfig)
	if err != nil {
		EarlyLoggerLn(fmt.Sprintf("Error initializing logger: %v", err.Error()))
		return
	}

	_, err = BotGroupInit(logger, shoneConfig)
	if err != nil {
		EarlyLoggerLn(fmt.Sprintf("Error initializing bot group: %v", err.Error()))
		return
	}

	for {
		// infinite loop
		// TODO: add a commandline console
		time.Sleep(time.Second * 1)
	}
}
