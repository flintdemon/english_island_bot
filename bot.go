package main

import (
	"log"

	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("😍 Да, с удовольствием"),
		tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
	),
)

//Struct used for questions in english level test
type Question struct {
	QuestionText string
	Answers      []string
	RightAnswer  string
	Points       int
}

func (q *Question) getQuestions() *Question {
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

func main() {

	var q Question
	q.getQuestions()
	log.Println(q.QuestionText)

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
		//if !update.Message.IsCommand() { //ignore any non command Messages
		//	continue
		//}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "🌴 Привет! Мы школа английского языка English Island!\n\n🔥 С нами ты заговоришь на английском уже с первого занятия которое ты можешь посетить совершенно бесплатно!"
			msg.ReplyMarkup = keyboard
		case "info":
			msg.Text = "Школа иностранных языков\n\n🌴English Island🌴\n\n🔥Это уютная атмосфера, современный подход и уроки с носителями языка.\n\n🔥Забудьте о нудной зубрежке и скучных домашних заданиях.\n\n🔥Приходи к нам в\n🌴English Island School🌴\nИ получи опыт живого языка, на котором действительно говорят."
		case "test":
			msg.Text = "Тут будет тест"
		default:
			msg.Text = "Я не знаю такой команды :("
		}

		switch update.Message.Text {
		case "😍 Да, с удовольствием":
			msg.Text = "Ok"
		case "🙄 Хочу узнать больше о школе":
			msg.Text = "Школа иностранных языков\n\n🌴English Island🌴\n\n🔥Это уютная атмосфера, современный подход и уроки с носителями языка.\n\n🔥Забудьте о нудной зубрежке и скучных домашних заданиях.\n\n🔥Приходи к нам в\n🌴English Island School🌴\nИ получи опыт живого языка, на котором действительно говорят."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
