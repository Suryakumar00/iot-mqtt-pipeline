package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("GoBackend_Surya")

	// Subscription Handler
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to Broker!")
		c.Subscribe("surya/iot/temperature", 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("Received Temperature: %s °C\n", string(msg.Payload()))
		})
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down backend...")
	client.Disconnect(250)
}
