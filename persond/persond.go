package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const (
	URL_ENV     = "url"
	DEFAULT_URL = "amqp://guest:guest@localhost:5672/"
)

type Person struct {
	id, name, surname string
}

var (
	version, build, buildDate string
	verbose                   bool
	url                       = DEFAULT_URL
)

func person(id string) Person {
	switch id {
	case "jomoespe":
		return Person{id: "jomoespe", name: "Jose", surname: "Moreno Esteban"}
	case "illescas":
		return Person{id: "illescas", name: "Jose Antonio", surname: "Illescas del Olmo"}
	case "juergas":
		return Person{id: "juergas", name: "Juan Ernesto", surname: "Roldan Garcia"}
	case "jbbarquero":
		return Person{id: "jbbarquero", name: "Javier", surname: "Beneito Barquero"}
	case "dparra":
		return Person{id: "dparra", name: "David", surname: "Parra Catalan"}
	}
	return Person{}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	env()
	flags()
	banner()

	if verbose {
		log.Printf("Establishing connection to command bus on %s", url)
	}
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			id := string(d.Body)
			failOnError(err, "Failed to get the id from the message body")

			if verbose {
				log.Printf("request for person %s\n", id)

			}
			response := person(id)
			if verbose {
				log.Printf("person found for %s is %s\n", id, response.name)
			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(response.name),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf("persond service started successfully. Awaiting requests.")
	<-forever
}

func env() {
	urlEnv, exist := os.LookupEnv(URL_ENV)
	if exist {
		url = urlEnv
	}
}

func flags() {
	flag.BoolVar(&verbose, "verbose", false, "Start service in verbose mode")
	flag.Parse()
}

func banner() {
	fmt.Println(" ______         __                                                              ")
	fmt.Println("|      |.-----.|  |.-----.-----.---.-.-----.                                    ")
	fmt.Println("|   ---||  _  ||  ||  -__|  _  |  _  |__ --|                                    ")
	fmt.Println("|______||_____||__||_____|___  |___._|_____|                                    ")
	fmt.Println("                         |_____|                                                ")
	fmt.Println("                                                                 __             ")
	fmt.Println(".-----.-----.----.-----.-----.-----.    .-----.-----.----.--.--.|__|.----.-----.")
	fmt.Println("|  _  |  -__|   _|__ --|  _  |     |    |__ --|  -__|   _|  |  ||  ||  __|  -__|")
	fmt.Println("|   __|_____|__| |_____|_____|__|__|    |_____|_____|__|  \\___/ |__||____|_____|")
	fmt.Println("|__|                                                                            ")
	fmt.Printf("\nVersion: %s\tBuild: %s\tBuild date: %s\n\n", version, build, buildDate)
}
