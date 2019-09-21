package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/anaskhan96/soup"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tbot.BotAPI

func start(ID int64) {
	bot.Send(tbot.NewMessage(ID, "Hello , I'm OSDC Bot, Use /help to know more. To join OSDC group: https://t.me/jiitosdc"))
}

func github(ID int64) {
	bot.Send(tbot.NewMessage(ID, "https://github.com/osdc"))
}

func telegram(ID int64) {
	bot.Send(tbot.NewMessage(ID, "https://t.me/jiitosdc"))
}

func website(ID int64) {
	bot.Send(tbot.NewMessage(ID, "https://osdc.netlify.com"))
}

func blog(ID int64) {
	bot.Send(tbot.NewMessage(ID, "https://osdc.github.io/blog"))
}

func irc(ID int64) {
	bot.Send(tbot.NewMessage(ID, "Join us on IRC server of Freenode at #jiit-lug. To get started refer our IRC wiki- https://github.com/osdc/community-committee/wiki/IRC ."))
}

func xkcd(ID int64) {
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 2000
	randomnum := (rand.Intn(max-min+1) + min)
	fmt.Println(randomnum)
	xkcdurl := "https://xkcd.com/" + strconv.Itoa(randomnum)
	resp, err := soup.Get(xkcdurl)
	if err == nil {
		doc := soup.HTMLParse(resp)
		links := doc.Find("div", "id", "comic").FindAll("img")
		for _, link := range links {
			linkimg := link.Attrs()["src"]
			fullurl := "https:" + linkimg
			fmt.Println(fullurl)
			bot.Send(tbot.NewMessage(ID, fullurl))
		}
	} else {
		fullurl := "https://imgs.xkcd.com/comics/operating_systems.png"
		bot.Send(tbot.NewMessage(ID, fullurl))
	}

}

func help(ID int64) {
	msg := ` Use one of the following commands
	/github - to get a link to OSDC's Github page.
	/telegram - to get an invite link for OSDC's Telegram Group.
	/website - to get the link of the official website of OSDC.
	/blog - to get the link of the OSDC blog.
	/irc - to find us on IRC.
	To contribute to|modify this bot : https://github.com/vaibhavk/osdc-bots
	`
	bot.Send(tbot.NewMessage(ID, msg))
}

func welcome(user tbot.User, ID int64) {
	User := fmt.Sprintf("[%v](tg://user?id=%v)", user.FirstName, user.ID)
	reply := tbot.NewMessage(ID, "**Welcome** "+User+", please introduce yourself")
	reply.ParseMode = "markdown"
	bot.Send(reply)
}

func kickUser(user int, ID int64) {
	bot.KickChatMember(tbot.KickChatMemberConfig{
		ChatMemberConfig: tbot.ChatMemberConfig{
			ChatID: ID,
			UserID: user,
		},
		UntilDate: time.Now().Add(time.Hour * 24).Unix(),
	})
}

func main() {
	var err error
	bot, err = tbot.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tbot.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ID := update.Message.Chat.ID
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			switch update.Message.Command() {

			case "start":
				start(ID)
			case "help":
				help(ID)
			case "github":
				github(ID)
			case "telegram":
				telegram(ID)
			case "website":
				website(ID)
			case "blog":
				blog(ID)
			case "irc":
				irc(ID)
			case "xkcd":
				xkcd(ID)
			default:
				bot.Send(tbot.NewMessage(ID, "I don't know that command"))
			}
		}
		if update.Message.NewChatMembers != nil {
			for _, user := range *(update.Message.NewChatMembers) {
				if user.IsBot && user.UserName != "osdcbot" {
					go kickUser(user.ID, ID)
				} else {
					go welcome(user, ID)
				}
			}
		}
	}
}
