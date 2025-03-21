package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration of the bot.
type Config struct {
	TelegramToken string
	AdminChatID   int64
	QuestionsFile string
}

// BotWrapper encapsulates bot functionality and state.
type BotWrapper struct {
	Bot         *tgbotapi.BotAPI
	KnownUsers  map[int64]userProfile
	AdminChatID int64
	Questions   []question
}

// userProfile represents information about a user's session.
type userProfile struct {
	Points          int
	ChatID          int64
	InTest          bool
	CurrentQuestion int
	LevelAfterTest  string
}

// question represents a single question in the test.
type question struct {
	QuestionText string   `yaml:"QuestionText"`
	Answers      []string `yaml:"Answers"`
	RightAnswer  string   `yaml:"RightAnswer"`
	Points       int      `yaml:"Points"`
}

// questionsGroup represents a collection of questions parsed from YAML.
type questionsGroup struct {
	Questions []question `yaml:"Questions"`
}

const (
	// Messages
	welcomeMessage = "🌴 Привет! Мы школа английского языка English Island!\n\n🔥 С нами ты заговоришь на английском уже с первого занятия, которое ты можешь посетить совершенно бесплатно!\n\nХочешь прийти на бесплатный урок?"
	testIntro      = "🔥 Прежде, чем мы запишем тебя на бесплатный пробный урок, мы предлагаем пройти тест на определение твоего уровня английского! 🔥\n\nЭто нужно для того, чтобы изучение языка было легким и комфортным для тебя.\n\nТест займет не более 10 минут.\nТы готов пройти тест?"
	schoolInfo     = "Школа иностранных языков\n\n🌴English Island🌴\n\n🔥 Это уютная атмосфера, современный подход и уроки с носителями языка.\n\n🔥 Забудьте о нудной зубрежке и скучных домашних заданиях.\n\n🔥 Приходи к нам в\n🌴English Island School🌴\nИ получи опыт живого языка, на котором действительно говорят."
	thankYou       = "Спасибо, наш менеджер свяжется с тобой в ближайшее время. До встречи в English Island School.🔥 \nP.S. нажми /start если хочешь начать заново."
)

// loadConfig loads the configuration from environment variables.
func loadConfig() Config {
	adminChatID, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		log.Panicf("Invalid ADMIN_CHAT_ID: %v", err)
	}

	return Config{
		TelegramToken: os.Getenv("TELETOKEN"),
		AdminChatID:   adminChatID,
		QuestionsFile: "questions.yml",
	}
}

// loadQuestions reads questions from the YAML file.
func loadQuestions(filePath string) ([]question, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var group questionsGroup
	if err = yaml.Unmarshal(yamlFile, &group); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return group.Questions, nil
}

// NewBotWrapper initializes a new BotWrapper.
func NewBotWrapper(config Config) *BotWrapper {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panicf("Failed to initialize bot: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	questions, err := loadQuestions(config.QuestionsFile)
	if err != nil {
		log.Panicf("Failed to load questions: %v", err)
	}

	return &BotWrapper{
		Bot:         bot,
		KnownUsers:  make(map[int64]userProfile),
		AdminChatID: config.AdminChatID,
		Questions:   questions,
	}
}

// getUserProfile retrieves or creates a user profile.
func (wrapper *BotWrapper) getUserProfile(chatID int64) userProfile {
	user, exists := wrapper.KnownUsers[chatID]
	if !exists {
		user = userProfile{ChatID: chatID}
		wrapper.KnownUsers[chatID] = user
	}
	return user
}

// makeStartKeyboard dynamically creates the start keyboard.
func makeStartKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("😍 Да, с удовольствием"),
			tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
		),
	)
}

// makeStartTestKeyboard dynamically creates the test start keyboard.
func makeStartTestKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("😍 Да, давайте начнем тест!"),
			tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
		),
	)
}

// makeContactKeyboard dynamically creates the contact sharing keyboard.
func makeContactKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Отправить номер телефона"),
		),
	)
}

// getQuestion formats a question into a Telegram message.
func getQuestion(chatID int64, question question) tgbotapi.MessageConfig {
	buttons := make([][]tgbotapi.KeyboardButton, len(question.Answers))
	for i, answer := range question.Answers {
		buttons[i] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer))
	}

	msg := tgbotapi.NewMessage(chatID, question.QuestionText)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)
	return msg
}

// determineLevel calculates the user's English level based on points.
func determineLevel(points int) string {
	switch {
	case points < 20:
		return "Elementary"
	case points < 45:
		return "Intermediate"
	default:
		return "Upper Intermediate"
	}
}

// handleMessage processes a Telegram update.
func (wrapper *BotWrapper) handleMessage(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		wrapper.handleCommand(update.Message)
		return
	}

	if update.Message.Contact != nil {
		wrapper.handleContact(update.Message)
		return
	}

	wrapper.handleConversation(update.Message)
}

// handleCommand handles user commands.
func (wrapper *BotWrapper) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		wrapper.sendMessage(msg.Chat.ID, welcomeMessage, makeStartKeyboard())
	}
}

// handleContact handles contact sharing.
func (wrapper *BotWrapper) handleContact(msg *tgbotapi.Message) {
	user := wrapper.KnownUsers[msg.Chat.ID]
	msgToSchool := fmt.Sprintf(
		"Пользователь %s прошел тест и прислал номер телефона: %s.\nУровень: %s",
		msg.Contact.FirstName, msg.Contact.PhoneNumber, user.LevelAfterTest,
	)

	wrapper.sendMessage(msg.Chat.ID, thankYou, nil)
	wrapper.sendMessage(wrapper.AdminChatID, msgToSchool, nil)
}

// handleConversation processes normal conversation messages.
func (wrapper *BotWrapper) handleConversation(msg *tgbotapi.Message) {
	user := wrapper.getUserProfile(msg.Chat.ID)

	switch msg.Text {
	case "😍 Да, с удовольствием":
		wrapper.sendMessage(msg.Chat.ID, testIntro, makeStartTestKeyboard())

	case "🙄 Хочу узнать больше о школе":
		wrapper.sendMessage(msg.Chat.ID, schoolInfo, nil)

	case "😍 Да, давайте начнем тест!":
		user.InTest = true
		user.Points = 0
		user.CurrentQuestion = 0
		qMsg := getQuestion(user.ChatID, wrapper.Questions[user.CurrentQuestion])
		wrapper.sendMessageStruct(qMsg)
	default:
		if user.InTest {
			wrapper.handleTest(msg, &user)
		}
	}
	wrapper.KnownUsers[msg.Chat.ID] = user
}

// handleTest processes test-related messages.
func (wrapper *BotWrapper) handleTest(msg *tgbotapi.Message, user *userProfile) {
	if user.CurrentQuestion < len(wrapper.Questions) {
		question := wrapper.Questions[user.CurrentQuestion]
		if question.RightAnswer == msg.Text {
			user.Points += question.Points
		}
		user.CurrentQuestion++
	}

	if user.CurrentQuestion == len(wrapper.Questions) {
		user.InTest = false
		user.LevelAfterTest = determineLevel(user.Points)
		wrapper.sendMessage(msg.Chat.ID, "Твой уровень языка: "+user.LevelAfterTest+"\n\nОставь свой номер телефона для записи на бесплатный урок 🌴", makeContactKeyboard())
	} else {
		qMsg := getQuestion(user.ChatID, wrapper.Questions[user.CurrentQuestion])
		wrapper.sendMessageStruct(qMsg)
	}
}

// sendMessage sends a message to a user.
func (wrapper *BotWrapper) sendMessage(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	if _, err := wrapper.Bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// sendMessageStruct sends a fully constructed message.
func (wrapper *BotWrapper) sendMessageStruct(msg tgbotapi.MessageConfig) {
	if _, err := wrapper.Bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Main starts the bot and listens to updates.
func main() {
	config := loadConfig()
	botWrapper := NewBotWrapper(config)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botWrapper.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to initialize updates channel: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		botWrapper.handleMessage(update)
	}
}
