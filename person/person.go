package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/streadway/amqp"
)

const (
	URL_ENV     = "url"
	DEFAULT_URL = "amqp://guest:guest@localhost:5672/"
)

var (
	version, build, buildDate string
	showVersion               bool
	id                        string
	url                       = DEFAULT_URL
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func personRPC(id string) (res string, err error) {
	//log.Printf("Connecting to %s", url)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(id),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			failOnError(err, "Failed to get response body")
			break
		}
	}

	if res=="" {
		os.Exit(1)
	}

	return
}

func main() {
	env()
	flags()

	if showVersion {
		banner()
		os.Exit(0)
	}
	if id == "" {
		showUsage()
		os.Exit(1)
	}

	//log.Printf(" [x] Requesting customer(%s)", id)
	res, err := personRPC(id)
	failOnError(err, "Failed to handle RPC request")

	fmt.Println(res)
}

func env() {
	urlEnv, exist := os.LookupEnv(URL_ENV)
	if exist {
		url = urlEnv
	}
}

func flags() {
	flag.BoolVar(&showVersion, "version", false, "Shows the program version and exit")
	flag.Parse()

	if len(flag.Args()) == 1 {
		id = flag.Args()[0]
	}
}

func banner() {
	fmt.Println(" ______         __                                                              ")
	fmt.Println("|      |.-----.|  |.-----.-----.---.-.-----.                                    ")
	fmt.Println("|   ---||  _  ||  ||  -__|  _  |  _  |__ --|                                    ")
	fmt.Println("|______||_____||__||_____|___  |___._|_____|                                    ")
	fmt.Println("                         |_____|                                                ")
	fmt.Println("                                                                                ")
	fmt.Println(".-----.-----.----.-----.-----.-----.                                            ")
	fmt.Println("|  _  |  -__|   _|__ --|  _  |     |                                            ")
	fmt.Println("|   __|_____|__| |_____|_____|__|__|                                            ")
	fmt.Println("|__|                                                                            ")
	fmt.Printf("\n\tVersion: %s\tBuild: %s\tBuild date: %s\n\n", version, build, buildDate)
}

func showUsage() {
	fmt.Println("person [options] <id>")
	fmt.Println("Where options:")
	fmt.Println("    -version: Shows the program version and exit")
	// TODO complete the program usage
}
