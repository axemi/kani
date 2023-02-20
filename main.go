package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var currentConfig *config

func main() {
	config, err := parseConfig("config.json") //store settings in config.json
	if err != nil {
		log.Panicf("error parsing config, %v", err)
	}
	currentConfig = config
	dg, err := discordgo.New("Bot " + currentConfig.Token)
	if err != nil {
		log.Panicf("error creating new discord session, %v", err)
	}
	dg.AddHandler(onReady)
	err = dg.Open()
	if err != nil {
		log.Panicf("error opening discord session, %v", err)
	}

	webserver := webserver(":8080")
	go func() {
		log.Printf("web server listening at %v", webserver.Addr)
		log.Fatal(webserver.ListenAndServe())
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	err = dg.Close()
	if err != nil {
		log.Printf("discord session closed with error: %v", err)
	}
	log.Println("disconnected")
	err = webserver.Shutdown(context.Background())
	if err != nil {
		log.Printf("web server shutdown with error: %v", err)
	}
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
