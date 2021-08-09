package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	log "github.com/taudelta/nanolog"

	"github.com/taudelta/notifyhub/internal/notifyhub"
	"github.com/taudelta/notifyhub/internal/notifyhub/config"
	"github.com/taudelta/notifyhub/internal/notifyhub/db"
	"github.com/taudelta/umq"
)

func main() {

	log.Init(log.Options{
		Level: log.DebugLevel,
	})

	var configFilePath string

	flag.StringVar(&configFilePath, "config", "./config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Fatal().Println(err)
	}

	db.CreateDefaultConnection(db.GetConnString(cfg.Database))

	brokerConfig := cfg.Broker

	queueOptions := umq.Options{
		Dial: umq.DialOptions{
			User:     brokerConfig.User,
			Password: brokerConfig.Password,
			Host:     brokerConfig.Host,
			Port:     brokerConfig.Port,
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
	}

	consumer, err := umq.NewConsumer(&umq.ConsumeOptions{
		Options:       queueOptions,
		NumOfWorkers:  10,
		PrefetchCount: 10,
	})

	if err != nil {
		log.Fatal().Printf("consumer error: %s\n", err)
	}

	notifyhub.SetEmailConfig(cfg.Email)
	if err := notifyhub.InitTelegram(cfg.Telegram.BotToken); err != nil {
		log.Fatal().Println("telegram client setup error: ", err)
	}
	notifyhub.InitSms(cfg.Sms)

	consumer.Run(&umq.WorkerOptions{
		Handler: &umq.QueueHandler{
			Manual: true,
			Apply:  umq.ManualHandler(SendHandler),
		},
		WaitTimeSeconds: 1,
		Logger:          log.DefaultLogger(),
	}, 0)

}

func SendHandler(delivery []amqp.Delivery) error {

	for _, mqMessage := range delivery {
		msg, err := notifyhub.ParseMessage(mqMessage.Body)
		if err != nil {
			log.Error().Println(err)
			continue
		}

		conn := db.Connection()
		err = notifyhub.SendMessage(*msg)
		if err != nil {
			log.Error().Println(err)
			_, err := conn.Exec(
				"update notification set status = $1, error = $2, processing_timestamp = $3 "+
					"where id = $4",
				"failed", fmt.Sprintf("%s", err), time.Now(),
				msg.DbID,
			)
			if err != nil {
				log.Error().Println(err)
			}
			mqMessage.Nack(false, false)
		} else {
			_, err := conn.Exec(
				"update notification set status = $1, processing_timestamp = $2 "+
					"where id = $3",
				"success", time.Now(), msg.DbID,
			)
			if err != nil {
				log.Error().Println(err)
			}
			mqMessage.Ack(false)
			log.Debug().Println("send notification ok")
		}

	}

	return nil
}
