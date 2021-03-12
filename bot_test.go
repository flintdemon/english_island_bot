package main

import (
	"io/ioutil"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

func TestYamlReadParse(t *testing.T) {
	var q questionsGroup
	yamlFile, err := ioutil.ReadFile(questionsFile)
	if err != nil {
		t.Error("yamlFile.Get error: ", err)
	}
	err = yaml.Unmarshal(yamlFile, &q)
	if err != nil {
		t.Error("Unmarshal: ", err)
	}
}

func TestRunBot(t *testing.T) {
	_, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		t.Error(err)
	}
}
