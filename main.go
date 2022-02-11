package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/codeedu/imersao6-go/email"
	"github.com/codeedu/imersao6-go/kafka"
	ckafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	gomail "gopkg.in/mail.v2"
)

func main() {
	var emailCh = make(chan email.Email)
	var msgChan = make(chan *ckafka.Message)

	d := gomail.NewDialer("sntp.mailgun.org", 587, "exemplo@schoolofnet.com", "dsçjfhesuiogfçeuh389750469hgfoegojr-=wegopj")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	es := email.NewMailSender()
	es.From = "exemplo@schoolofnet.com"
	es.Dailer = d

	go es.Send(emailCh)

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9094",
		"client.id":         "emailapp",
		"group.id":          "emailapp",
	}
	topics := []string{"emails"}
	consumer := kafka.NewConsumer(configMap, topics)
	go consumer.Consume(msgChan)

	for msg := range msgChan {
		var input email.Email
		json.Unmarshal(msg.Value, &input)
		fmt.Println("Recebendo mensagem")
		emailCh <- input
	}

}
