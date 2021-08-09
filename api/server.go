package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/taudelta/nanolog"
	"github.com/taudelta/notifyhub/app/notifications"
	"github.com/taudelta/notifyhub/app/recipients"
	"github.com/taudelta/notifyhub/internal/notifyhub/config"
	"github.com/taudelta/notifyhub/internal/notifyhub/db"
	"github.com/taudelta/umq"
)

func StartServer(port int, messageProducer *umq.RabbitProducer) {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		notifications.SendNotification(messageProducer, w, r)
	})

	http.HandleFunc("/recipients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			recipients.CreateRecipient(w, r)
		} else if r.Method == "GET" {
			recipients.GetRecipients(w, r)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal().Println(err)
	}

}

func main() {

	var configFilePath string

	flag.StringVar(&configFilePath, "config", "./config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Fatal().Println(err)
	}

	err = db.CreateDefaultConnection(db.GetConnString(cfg.Database))
	if err != nil {
		log.Fatal().Println(err)
	}

	messageProducer, err := umq.NewProducer(umq.Options{
		Dial: umq.DialOptions{
			User:     cfg.Broker.User,
			Password: cfg.Broker.Password,
			Host:     cfg.Broker.Host,
			Port:     cfg.Broker.Port,
		},
		Queue: umq.QueueOptions{
			Name:    "notifications",
			Durable: true,
		},
		Exchange: umq.ExchangeOptions{
			Name:    "notifications",
			Kind:    umq.DirectKind,
			Durable: true,
		},
		Bind: umq.BindOptions{
			RoutingKey: "notifications",
		},
	})

	if err != nil {
		log.Fatal().Println(err)
	}

	StartServer(cfg.Server.Port, messageProducer)

}
