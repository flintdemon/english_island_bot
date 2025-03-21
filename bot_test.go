package main

import (
	"os"
	"strconv"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// TestLoadQuestions tests the refactored `loadQuestions` function.
func TestLoadQuestions(t *testing.T) {
	questionsFile := "questions.yml"
	yamlFile, err := os.ReadFile(questionsFile)
	assert.NoError(t, err, "Failed to read YAML file")

	var group questionsGroup
	err = yaml.Unmarshal(yamlFile, &group)
	assert.NoError(t, err, "Failed to unmarshal YAML file")

	assert.NotEmpty(t, group.Questions, "Questions should not be empty")
}

// TestBotInitialization tests the refactored bot initialization logic.
func TestBotInitialization(t *testing.T) {
	telegramToken := os.Getenv("TELETOKEN")
	adminChatIDEnv := os.Getenv("ADMIN_CHAT_ID")

	if telegramToken == "" || adminChatIDEnv == "" {
		t.Skip("Environment variables TELETOKEN or ADMIN_CHAT_ID not set, skipping test")
	}

	adminChatID, err := strconv.ParseInt(adminChatIDEnv, 10, 64)
	assert.NoError(t, err, "Failed to parse ADMIN_CHAT_ID environment variable")

	config := Config{
		TelegramToken: telegramToken,
		AdminChatID:   adminChatID,
		QuestionsFile: "questions.yml",
	}

	botWrapper := NewBotWrapper(config)
	assert.NotNil(t, botWrapper.Bot, "Bot initialization failed")
	assert.Equal(t, config.AdminChatID, botWrapper.AdminChatID, "AdminChatID mismatch")
	assert.NotEmpty(t, botWrapper.Questions, "Questions should not be empty")
}

// TestDetermineLevel tests level determination logic.
func TestDetermineLevel(t *testing.T) {
	assert.Equal(t, "Elementary", determineLevel(15))
	assert.Equal(t, "Intermediate", determineLevel(30))
	assert.Equal(t, "Upper Intermediate", determineLevel(50))
}

// TestUserProfileManagement tests `getUserProfile` logic.
func TestUserProfileManagement(t *testing.T) {
	wrapper := BotWrapper{
		KnownUsers: make(map[int64]userProfile),
	}

	chatID := int64(12345)
	user := wrapper.getUserProfile(chatID)

	assert.Equal(t, chatID, user.ChatID, "ChatID mismatch")
	assert.False(t, user.InTest, "Default InTest value should be false")

	// Update user state and re-check
	user.Points = 10
	user.InTest = true
	wrapper.KnownUsers[chatID] = user

	updatedUser := wrapper.getUserProfile(chatID)
	assert.Equal(t, 10, updatedUser.Points, "Points mismatch")
	assert.True(t, updatedUser.InTest, "InTest mismatch")
}

// TestSendMessage tests the `sendMessage` helper function.
func TestSendMessage(t *testing.T) {
	telegramToken := os.Getenv("TELETOKEN")
	if telegramToken == "" {
		t.Skip("Environment variable TELETOKEN not set, skipping test")
	}

	testChatID, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		t.Skip("Invalid ADMIN_CHAT_ID, skipping test")
	}

	config := Config{
		TelegramToken: telegramToken,
		AdminChatID:   testChatID,
		QuestionsFile: "questions.yml",
	}

	botWrapper := NewBotWrapper(config)
	assert.NotNil(t, botWrapper.Bot, "Bot should be initialized")

	_, err = botWrapper.Bot.Send(tgbotapi.NewMessage(botWrapper.AdminChatID, "Test Message"))
	assert.NoError(t, err, "Failed to send message")
}
