package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	appToken              = ""
	botToken              = ""
	botUserID             = ""
	gopherPing            = ""
	notificationChannelID = ""
	ws                    = ""
	debug                 = false

	threadReply = "thanks, added to @gopher backlog"
)

func quote(str string) string {
	x := strings.ReplaceAll(str, "\n", "\n> ")
	return fmt.Sprintf("> %v", x)
}

func ts(msg *slackevents.MessageEvent) string {
	if msg.ThreadTimeStamp != "" {
		return msg.ThreadTimeStamp
	}
	return msg.TimeStamp
}

func processMsgEvent(api *slack.Client, data interface{}) bool {
	ev, ok := data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("unexpected type for EventsAPIEvent %v %T\n", ev.Type, ev.Data)
		return true
	}
	msg, ok := ev.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		if debug {
			log.Printf("unexpected type for MessageEvent %v %T\n", ev.InnerEvent.Type, ev.InnerEvent.Data)
		}
		return true
	}
	if msg.User == botUserID {
		if debug {
			log.Println("skipping, message user == botUserID")
		}
		return true
	}
	if !strings.Contains(msg.Text, gopherPing) {
		if debug {
			log.Println("skipping, missing gopherPing", gopherPing)
			log.Println(msg.Text)
		}
		return true
	}
	backlogMsg := fmt.Sprintf("Message from <@%v>\n%v\n\n*link:* https://%v.slack.com/archives/%v/p%v", msg.User, quote(msg.Text), ws, msg.Channel, msg.TimeStamp)
	if _, _, err := api.PostMessage(notificationChannelID, slack.MsgOptionText(backlogMsg, false)); err != nil {
		log.Printf("failed posting reply: %v\n", err)
		return false
	}
	if _, _, err := api.PostMessage(msg.Channel, slack.MsgOptionText(threadReply, false), slack.MsgOptionTS(ts(msg))); err != nil {
		log.Printf("failed posting reply: %v\n", err)
		return false
	}
	return true
}

func main() {
	flag.StringVar(&appToken, "app-token", "", "Slack App Token")
	flag.StringVar(&botToken, "bot-token", "", "Slack Bot Token")
	flag.StringVar(&botUserID, "bot-user-id", "", "Slack App Bot user ID")
	flag.StringVar(&gopherPing, "gopher-ping", "", "Slack group slug expected to be in the message")
	flag.StringVar(&notificationChannelID, "notification-channel-id", "", "Slack channel ID where the bot should send notifications")
	flag.StringVar(&ws, "workspace", "", "Slack workspace")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	api := slack.New(
		botToken,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(debug),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for e := range client.Events {
			switch e.Type {
			case socketmode.EventTypeConnecting:
				log.Println("connecting to slack with socket mode")
			case socketmode.EventTypeConnectionError:
				log.Fatal(fmt.Errorf("received %q: %v", socketmode.EventTypeConnectionError, e.Data))
			case socketmode.EventTypeConnected:
				log.Println("connected")
			case socketmode.EventTypeEventsAPI:
				if ack := processMsgEvent(api, e.Data); ack {
					client.Ack(*e.Request)
				}
			}
		}
	}()
	err := client.Run()
	log.Fatal(err)
}
