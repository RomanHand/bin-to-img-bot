package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	tele "gopkg.in/telebot.v3"
	"romanhand.ru/bin-to-img/internal/config"
	"romanhand.ru/bin-to-img/internal/imgen"
	"romanhand.ru/bin-to-img/internal/logging"
)



var userStates = make(map[int64]string)

func main() {
	logging.SetupLogging()

	cfgPathFlag := flag.String("config", "", "Path to the configuration file")
	flag.Parse() 

	cfgPath := *cfgPathFlag
	if cfgPath == "" {
		cfgPath = "/etc/tg-bti-bot/config.yml"
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatal("Error config read")
	}

	token := os.Getenv("BTP_TELEGRAM_BOT_TOKEN")
	if token == "" {
		token = cfg.Token
	} 
	

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(e error, c tele.Context) {
			log.Println("Error:", e)
		},
	}
	saveDir := cfg.BinsDir
	if err != nil {
		log.Fatal("Error config read (BinsDir)")
	}
	imgDir := cfg.CompliteDir
	if err != nil {
		log.Fatal("Error config read (BinsDir)")
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Creating bot: %v", err)
	}

	bot.Handle("/start", func(c tele.Context) error {
		userStates[c.Chat().ID] = "waiting_for_bin"
		log.Printf("User %d started interaction with the bot (username: %s)", c.Chat().ID, c.Chat().Username)
		return c.Send(cfg.WelcomeMsg)
	})

	bot.Handle("/help", func(c tele.Context) error {
		userStates[c.Chat().ID] = "waiting_for_target"
		log.Printf("User %d entered help command (username: %s)", c.Chat().ID, c.Chat().Username)
		return c.Send("Тут это, короче, могу только морально поддрежать. А вообще тыка на клавиатуре /start , а после по подсказкам.")
	})

	bot.Handle(tele.OnDocument, func(c tele.Context) error {
		log.Printf("User %d started creating an image (username: %s)", c.Chat().ID, c.Chat().Username)

		filePath := filepath.Join(saveDir, c.Chat().Username)
		imgPath := filepath.Join(imgDir, c.Chat().Username)
		imgPath = imgPath + ".png"

		c.Bot().Download(&c.Message().Document.File, filePath)
		
		err = imgen.GenerateImg(filePath, imgPath)
		if err != nil {
			return c.Send("Во время конвертации произошла ошибка, попробуйте снова.")
		}

		photo := &tele.Photo{
			File: tele.FromDisk(imgPath), 
		}

		_, err = photo.Send(c.Bot(), c.Chat(), nil)
		if err != nil {
			return c.Send("Ошибка при отправке изображения.")
		}else if err == nil{
			log.Printf("Image for user %d successfully sented (username: %s)", c.Chat().ID, c.Chat().Username)
			
		}
		return c.Send("Твой картинка готов!") 
	})

	bot.Start()
}







