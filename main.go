package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config, err := parseConfig("config.json") //store settings in config.json
	if err != nil {
		log.Panicf("error parsing config, %v", err)
	}
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Panicf("error creating new discord session, %v", err)
	}
	dg.AddHandler(onReady)
	err = dg.Open()
	if err != nil {
		log.Panicf("error opening discord session, %v", err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
	log.Println("disconnected")
}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	log.Printf("connected as %v (%v)", event.User.Username, event.User.ID)
}

type config struct {
	Token string `json:"token"`
}

func parseConfig(filename string) (*config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
