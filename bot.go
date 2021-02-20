package main

import (
	"log"

	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

var buttonTouch bool = false
var userPoint int = 0

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("😍 Да, с удовольствием"),
		tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
	),
)

//question struct using for questions in english level test
type question struct {
	QuestionText string   `yaml:"QuestionText"`
	Answers      []string `yaml:"Answers"`
	RightAnswer  string   `yaml:"RightAnswer"`
	Points       int      `yaml:"Points"`
}

//questionGroup struct contains array of questions
type questionGroup struct {
	Questions []question `yaml:"Questions"`
}

//getQuestions method for import questions
func (q *questionGroup) getQuestions() *questionGroup {
	yamlFile, err := ioutil.ReadFile("questions.yml")
	if err != nil {
		log.Printf("yamlFile.Get err  #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, q)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return q
}

//getQuestion function gets question by uts number and makes reply keyboard with answers
func getQuestion(msg *tgbotapi.MessageConfig, questionNumber int) {
	var questions questionGroup

	qArray := questions.getQuestions().Questions

	buttons := make([][]tgbotapi.KeyboardButton, len(qArray[questionNumber].Answers))
	for i, a := range qArray[questionNumber].Answers {
		buttons[i] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(a))
	}

	msg.Text = qArray[questionNumber].QuestionText
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)
}

func main() {

	bot, err := tgbotapi.NewBotAPI("1601846360:AAHPRgAazXY-bX-fZI5NAh0ffUWGbPmH0-I")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "🌴 Привет! Мы школа английского языка English Island!\n\n🔥 С нами ты заговоришь на английском уже с первого занятия которое ты можешь посетить совершенно бесплатно!\n\nХочешь прийти на бесплатный урок?"
			msg.ReplyMarkup = startKeyboard
		default:
			msg.Text = "Напиши /start или нажимай на кнопки внизу."
		}

		switch update.Message.Text {
		case "😍 Да, с удовольствием":
			if buttonTouch == false {
				msg.Text = "🔥 Прежде, чем мы запишем тебя на бесплатный пробный урок, мы предлагаем пройти тест на определение твоего уровня английского! 🔥\n\nЭто нужно для того, чтобы изучение языка было легким и комфортным для тебя.\n\nТест займет не более 10ти минут\nТы готов пройти тест?"
				buttonTouch = true
			} else {
				getQuestion(&msg, 6)

				buttonTouch = false
			}

		case "🙄 Хочу узнать больше о школе":
			msg.Text = "Школа иностранных языков\n\n🌴English Island🌴\n\n🔥Это уютная атмосфера, современный подход и уроки с носителями языка.\n\n🔥Забудьте о нудной зубрежке и скучных домашних заданиях.\n\n🔥Приходи к нам в\n🌴English Island School🌴\nИ получи опыт живого языка, на котором действительно говорят."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
