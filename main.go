package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Use this exact path:
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// 1. Configure the connection
	opts := MQTT.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("go-client-listener")

	// 2. Create the client
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// 3. Subscribe to a topic
	topic := "my/test/topic"
	client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})

	fmt.Println("Listening for messages on", topic, "... Press Ctrl+C to stop.")

	// 4. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	client.Disconnect(250)
	fmt.Println("Disconnected.")
}
