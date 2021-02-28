package main

import (
	"log"

	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

const telegramToken string = "1601846360:AAHPRgAazXY-bX-fZI5NAh0ffUWGbPmH0-I"

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("😍 Да, с удовольствием"),
		tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
	),
)

var startTestKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("😍 Да, давайте начнем тест!"),
		tgbotapi.NewKeyboardButton("🙄 Хочу узнать больше о школе"),
	),
)

var endKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonContact("Отправить номер телефона"),
	),
)

type userProfile struct {
	Points          int
	ChatID          int64
	inTest          bool
	currentQuestion int
	levelAfterTest  string
}

type question struct {
	QuestionText string   `yaml:"QuestionText"`
	Answers      []string `yaml:"Answers"`
	RightAnswer  string   `yaml:"RightAnswer"`
	Points       int      `yaml:"Points"`
}

type questionsGroup struct {
	Questions []question `yaml:"Questions"`
}

func userContainsIn(a []userProfile, u userProfile) bool {
	for _, n := range a {
		if u.ChatID == n.ChatID {
			return true
		}
	}
	return false
}

func (q *questionsGroup) getQuestions() *questionsGroup {
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

func getQuestion(chatID int64, questionNumber int) tgbotapi.MessageConfig {

	var questions questionsGroup

	qArray := questions.getQuestions().Questions

	buttons := make([][]tgbotapi.KeyboardButton, len(qArray[questionNumber].Answers))
	for i, a := range qArray[questionNumber].Answers {
		buttons[i] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(a))
	}

	msg := tgbotapi.NewMessage(chatID, "")
	msg.Text = qArray[questionNumber].QuestionText
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)

	return msg
}

func makeBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func main() {

	var knownUsers map[int64]userProfile
	var currentUser userProfile

	bot := makeBot()

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("From: [%s] Message is: %s", update.Message.From.UserName, update.Message.Text)

		if currentUser.ChatID == 0 {
			currentUser = userProfile{0, update.Message.Chat.ID, false, 0, ""}
		} else {
			currentUser.ChatID = update.Message.Chat.ID
		}

		if knownUsers == nil {
			knownUsers = make(map[int64]userProfile)
		}
		if knownUsers[update.Message.Chat.ID].ChatID == 0 {
			knownUsers[update.Message.Chat.ID] = currentUser
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				msg.Text = "🌴 Привет! Мы школа английского языка English Island!\n\n🔥 С нами ты заговоришь на английском уже с первого занятия которое ты можешь посетить совершенно бесплатно!\n\nХочешь прийти на бесплатный урок?"
				msg.ReplyMarkup = startKeyboard
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}

		if update.Message.Text == "😍 Да, с удовольствием" {
			msg.Text = "🔥 Прежде, чем мы запишем тебя на бесплатный пробный урок, мы предлагаем пройти тест на определение твоего уровня английского! 🔥\n\nЭто нужно для того, чтобы изучение языка было легким и комфортным для тебя.\n\nТест займет не более 10ти минут\nТы готов пройти тест?"
			msg.ReplyMarkup = startTestKeyboard
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.Message.Text == "🙄 Хочу узнать больше о школе" {
			msg.Text = "Школа иностранных языков\n\n🌴English Island🌴\n\n🔥Это уютная атмосфера, современный подход и уроки с носителями языка.\n\n🔥Забудьте о нудной зубрежке и скучных домашних заданиях.\n\n🔥Приходи к нам в\n🌴English Island School🌴\nИ получи опыт живого языка, на котором действительно говорят."
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.Message.Text == "😍 Да, давайте начнем тест!" {
			user := knownUsers[update.Message.Chat.ID]
			if user.inTest == false {
				user.inTest = true
				user.Points = 0 //If test completed several times
			}
			qMsg := getQuestion(user.ChatID, user.currentQuestion)

			if _, err := bot.Send(qMsg); err != nil {
				log.Panic(err)
			}

			user.currentQuestion++
			knownUsers[update.Message.Chat.ID] = user

		} else if update.Message.Contact != nil {
			msgToSchool := tgbotapi.NewMessage(418634811, "Пользователь "+update.Message.Contact.FirstName+" прошел тест и прислал номер телефона:"+update.Message.Contact.PhoneNumber+"\nЕго уровень по результатам теста: "+knownUsers[update.Message.Chat.ID].levelAfterTest)
			msg.Text = "Спасибо, наш менеджер свяжется с тобой в ближайшее время. До встречи в English Island School.🔥 \n P.S. нажми /start если хочешь начать заново"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			if _, err := bot.Send(msgToSchool); err != nil {
				log.Panic(err)
			}
		} else {
			user := knownUsers[update.Message.Chat.ID]
			if user.inTest == false {
				//Exit when user not in test
				continue
			}
			var questions questionsGroup
			qArray := questions.getQuestions().Questions
			numOfQuestions := len(qArray)

			if user.currentQuestion < numOfQuestions {
				question := qArray[user.currentQuestion-1] //Because we are reading answers of the previous question
				if question.RightAnswer == update.Message.Text {
					user.Points += question.Points
				}

				qMsg := getQuestion(user.ChatID, user.currentQuestion)
				if _, err := bot.Send(qMsg); err != nil {
					log.Panic(err)
				}
				user.currentQuestion++
				knownUsers[update.Message.Chat.ID] = user
			} else {
				var level string
				if user.Points < 20 {
					level = "Elementary"
				} else if user.Points >= 20 && user.Points < 45 {
					level = "Intermediate"
				} else if user.Points >= 45 {
					level = "Upper Intermediate"
				}
				user.levelAfterTest = level
				user.inTest = false
				knownUsers[update.Message.Chat.ID] = user
				msg.ReplyMarkup = endKeyboard
				msg.Text = "Твой уровень языка: " + level + "\n\nПоздравляем тебя, ты успешно прошел тест на определение уровня языка🔥\n\n Для того, чтобы мы могли записать тебя на бесплатный урок, тебе надо оставить свой номер телефона 🌴"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

		}
	}
}
