package main

import "github.com/arttnba3/Shigure-Bot/bot"

func CreateBot(botType string, configJson []byte, logger func(...any), handlers map[string]func(...any)) (*shigure.ShigureBot, error) {
	return shigure.NewShigureBot(botType, configJson, logger, handlers)
}
