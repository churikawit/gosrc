package linebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

type BotConfig struct {
	BotID              int    // id register to database
	BotUserId          string // destination from LineApi
	ChannelSecret      string
	ChannelAccessToken string
	Client             *linebot.Client
	Handler            func(w http.ResponseWriter, req *http.Request, config BotConfig)
}

var BotList = []BotConfig{}

func init() {

}

func RegisterBot(bot_id int, bot_userid, channel_secret, channel_access_token string,
	handler func(w http.ResponseWriter, req *http.Request, config BotConfig)) {
	client, err := linebot.New(
		channel_secret,
		channel_access_token,
	)
	if err != nil {
		log.Fatal(err)
	}

	BotList = append(BotList, BotConfig{
		BotID: bot_id, BotUserId: bot_userid,
		ChannelSecret: channel_secret, ChannelAccessToken: channel_access_token,
		Client: client, Handler: handler})
	fmt.Printf("register bot: %s\n", bot_userid)
}

func GetHandler() func(w http.ResponseWriter, req *http.Request) {
	return _handler
}

func _handler(w http.ResponseWriter, req *http.Request) {
	// PreProcess ------------------------------------
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(400)
	}

	// re-create body
	temp := ioutil.NopCloser(bytes.NewBuffer(body))
	req.Body = temp

	request := &struct {
		Events      []*linebot.Event `json:"events"`
		Destination string           `json:"destination"`
	}{}
	err = json.Unmarshal(body, request)
	fmt.Printf("handle destination: \"%s\"\n", string(request.Destination))
	// endPreProcess ---------------------------------

	// Bot Routing
	dest := string(request.Destination)
	for _, bot := range BotList {
		if dest == bot.BotUserId {
			bot.Handler(w, req, bot)
			return
		}
	}
	log.Printf("No handler for %s\n", dest)
}
