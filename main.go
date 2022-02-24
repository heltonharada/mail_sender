package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/codeedu/imersao6-go/email"
	"github.com/codeedu/imersao6-go/kafka"
	ckafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	gomail "gopkg.in/mail.v2"
)

func main() {
	var emailCh = make(chan email.Email)
	var msgChan = make(chan *ckafka.Message)

	//incluido alteracao aula kubernetes

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	//substituido na aula kubernetes

	// d := gomail.NewDialer("smtp.sendgrid.net", 465, "apikey", "**************************")

	d := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		port,
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL PASSWORD"),
	)

	//

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	es := email.NewMailSender()
	// es.From = "emaildohelton@gmail.com"
	es.From = os.Getenv("MAIL FROM")
	es.Dailer = d

	// for i := 1; i <= 2; i++ {
	// 	go es.Send(emailCh, i)

	// }

	go es.Send(emailCh)

	configMap := &ckafka.ConfigMap{
		//"bootstrap.servers": "kafka:9094",
		//"bootstrap.servers": "host.docker.internal:9094",
		"client.id": "emailapp",
		"group.id":  "emailapp",
		//incluido abaixo
		"bootstrap.servers":  os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol":  os.Getenv("SECURITY_PROTOCOL"),
		"sasl.mechanisms":    os.Getenv("SASL_MECHANISMS"),
		"sasl.username":      os.Getenv("SASL_USERNAME"),
		"sasl.password":      os.Getenv("SASL_PASSWORD"),
		"session.timeout.ms": 45000,
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
