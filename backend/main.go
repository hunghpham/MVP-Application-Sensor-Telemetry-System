package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func main() {
    broker := os.Getenv("KAFKA_BROKER")
    topic := "sensor-data"

    producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
    if err != nil {
        panic(err)
    }

    // Kafka consumer for receiving data from producers
    go consumeKafkaMessages(broker, topic)

    // Websocket handler for frontend
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal(err)
        }
        go handleWebSocketConnection(ws)
    })

    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func consumeKafkaMessages(broker string, topic string) {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": broker,
        "group.id":          "go-consumer-group",
        "auto.offset.reset": "earliest",
    })

    if err != nil {
        log.Fatal(err)
    }

    c.SubscribeTopics([]string{topic}, nil)

    for {
        msg, err := c.ReadMessage(-1)
        if err == nil {
            log.Printf("Received: %s\n", string(msg.Value))
            // TODO: Store the message to QuestDB
        } else {
            log.Printf("Consumer error: %v (%v)\n", err, msg)
        }
    }

    c.Close()
}

func handleWebSocketConnection(ws *websocket.Conn) {
    defer ws.Close()
    for {
        _, message, err := ws.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }

        // Send message to Kafka producer
        log.Println("Received from client:", string(message))
        // TODO: Process the message and send to Kafka
    }
}
