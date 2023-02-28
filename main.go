package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Panicf("error loading config, %v", err)
	}
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Panicf("error creating new discord session, %v", err)
	}
	dg.AddHandler(onReady)
	webserver := NewWebServer(":8080")
	webserver.mux.HandleFunc("/", webHandler)
	webserver.mux.Handle("/settings", settingsHandler(config))
	go func() {
		log.Printf("web server listening at %v", webserver.server.Addr)
		log.Fatal(webserver.server.ListenAndServe())
	}()
	if config.Token != "" {
		err = dg.Open()
		if err != nil {
			log.Panicf("error opening discord session, %v", err)
		}
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	err = dg.Close()
	if err != nil {
		log.Printf("discord session closed with error: %v", err)
	}
	log.Println("disconnected")
	err = webserver.server.Shutdown(context.Background())
	if err != nil {
		log.Printf("web server shutdown with error: %v", err)
	}
}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	log.Printf("connected as %v (%v)", event.User.Username, event.User.ID)
}

type config struct {
	Token         string `json:"token"`
	DefaultPrefix string `json:"default_prefix"`
}

// will handle initial creating config.json and setting global currentConfig
func loadConfig() (*config, error) {
	var c *config
	log.Println("checking for config.json")
	_, err := os.Stat("config.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	//handle config.json file not created
	if errors.Is(err, os.ErrNotExist) {
		log.Println("did not find config.json")
		newConfigFile, err := os.Create("config.json")
		if err != nil {
			return nil, err
		}
		c = &config{Token: ""}
		b, err := json.MarshalIndent(c, "", "\t")
		if err != nil {
			return nil, err
		}
		_, err = newConfigFile.Write(b)
		if err != nil {
			return nil, err
		}
		log.Println("created config.json")
		return c, nil
	}
	//handle existing config.json
	log.Println("reading config.json")
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	log.Println("loading current config")
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*
webserver chan - should send bot token, default prefix
discord chan - wants to receive bot token, default prefix

*/
