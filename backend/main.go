package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/questdb/go-questdb-client/v3"
	"github.com/segmentio/kafka-go"
	"github.com/tidwall/gjson"
)

type SensorData struct {
	Serial    string  `json:"Serial"`
	Type      string  `json:"Type"`
	Timestamp string  `json:"Timestamp"`
	Reading1  float64 `json:"Reading1"`
	Reading2  float64 `json:"Reading2"`
	Reading3  float64 `json:"Reading3"`
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Global slice to hold active WebSocket connections
var clients = make(map[*websocket.Conn]bool)

// Broadcast channel for sending messages
var broadcast = make(chan []byte)

func handleWebSocket(c *gin.Context) {
	// Upgrade the connection to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Register the new client
	clients[ws] = true
	defer delete(clients, ws)

	// Handle incoming messages
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received: %s", message)

		// Optionally, we can send message right here by sending the message to the broadcast channel
		// broadcast <- message
	}
}

// Broadcast messages to all connected clients
func broadcastMessages() {
	for {
		message := <-broadcast // Wait for a message to broadcast

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error writing message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {

	server := gin.Default()

	//Apply the CORS middleware
	server.Use(cors.Default())

	//GET endpoints
	server.GET("/api/get_all_sensors", getAllSensors)                                     //connect to QuestDB to get all sensors info
	server.GET("/api/get_sensor_historical_data", getHistoricalData)                      //connect to QuestDB to get historical data
	server.GET("/api/create_kafka_topic_and_subscribe", create_kafka_topic_and_subscribe) //create kafka topic through Kafka rest api

	// Set up WebSocket endpoint
	server.GET("/ws", handleWebSocket)

	// Start the broadcast goroutine, url is ws://localhost:8080/ws
	go broadcastMessages()

	//Start the backend server listening on port 8080
	log.Fatal(server.Run(":8080"))
}

// /////////////// Database related functions /////////////////////
// Execute a SQL query
func executeQuery(query string) string {

	connStr := "postgres://admin:quest@questdb:8812/qdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Open database error: " + err.Error())
	}
	defer db.Close()

	// Ensure the connection is alive
	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database: " + err.Error())
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Prepare query error: " + err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("Execute query error: " + err.Error())
	}
	defer rows.Close()

	dataList := []SensorData{}
	for rows.Next() {

		var serialNumber string
		var timestamp string
		var sensorType string
		var reading1 float64
		var reading2 float64
		var reading3 float64

		// Scan the data into variables
		if err := rows.Scan(&serialNumber, &sensorType, &timestamp, &reading1, &reading2, &reading3); err != nil {
			log.Println("Failed to scan row: " + err.Error())
		}

		//add each record to a list to return
		var data SensorData
		data.Serial = serialNumber
		data.Type = sensorType
		data.Timestamp = timestamp
		data.Reading1 = reading1
		data.Reading2 = reading2
		data.Reading3 = reading3

		dataList = append(dataList, data)

		// Check for any errors encountered during iteration
		if err := rows.Err(); err != nil {
			log.Println("Error during row iteration: " + err.Error())
		}
	}

	// Convert the slice to JSON
	jsonData, err := json.Marshal(dataList)
	if err != nil {
		log.Println("Failed to marshal data: " + err.Error())

	}

	return string(jsonData)
}

// Get all current sensors
func getAllSensors(ctx *gin.Context) {

	// Define the SQL query using REST HTTP API
	query := `
			SELECT serial_number, sensor_type, timestamp, reading1, reading2, reading3
			FROM sensor_historical_data
			LATEST BY serial_number ORDER BY serial_number ASC;
			`
	result := executeQuery(query)

	ctx.String(http.StatusOK, result)
}

// Get individual sensor historical data based on datetime range
func getHistoricalData(ctx *gin.Context) {

	query := `SELECT * FROM sensor_historical_data
				WHERE serial_number = '#SERIAL#'
				AND timestamp BETWEEN '#SDT#' AND '#EDT#'
				ORDER BY timestamp ASC;`

	//Get sensor serial number and the date range
	serial := ctx.Query("serial_number")
	startDT := ctx.Query("start_dt")
	endDT := ctx.Query("end_dt")

	//modify query
	query = strings.ReplaceAll(query, "#SERIAL#", serial)
	query = strings.ReplaceAll(query, "#SDT#", startDT)
	query = strings.ReplaceAll(query, "#EDT#", endDT)

	result := executeQuery(query)

	ctx.String(http.StatusOK, result)
}

// //////////////// Other helper functions ////////////////////////////////////////////////////////////////
// Get Kafka cluster id using Kafka REST Proxy instead of native client, just demonstrating its flexibility
func getKafkaClusterID(url string) (string, error) {

	//Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Get cluster id string from the json response
	value := gjson.Get(string(body), "data.0.cluster_id")

	return value.String(), nil
}

// Get an instance Kafka cluster ID based on URL
func create_kafka_topic_and_subscribe(ctx *gin.Context) {

	url := "http://kafka-rest-proxy:8082/v3/clusters"

	//Get Kafka local cluster id using Kafka REST API proxy, demonstrating REST works as well
	cluster_id, err := getKafkaClusterID(url)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"Status": "Server error, unable to get Kafka cluster id.", "Error": err.Error()})
	}

	// Get topic name from the post
	topic := ctx.Query("topic_name")
	broker := "kafka:9092"
	partitions := 1
	replicationFactor := 1

	// Initialize a Kafka writer using native Kafka client library, demonstrating native client works as well
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer writer.Close()

	// Create a new admin client
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "Connect to Kafka failed!", "error": err.Error()})
	}
	defer conn.Close()

	// Create the topic
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "Failed to create topic!", "error": err.Error()})
		return
	}

	//add path to url for creating publish message endpoint
	new_url := strings.ReplaceAll(url, "kafka-rest-proxy", "localhost") + "/" + cluster_id + "/topics"

	//Kafka topic is created, now start the listening for messages on this topic by subscribing to it
	go subscribeToKafkaTopic(topic)

	//Return messages to user
	ctx.JSON(http.StatusOK, gin.H{"Status": "OK", "Kafka_Endpoint": new_url + "/" + topic + "/records"})

}

// Subscribe to Kafka
func subscribeToKafkaTopic(topic_name string) {
	// Define the Kafka broker addresses and topic
	brokers := []string{"kafka:9092"}   // Replace with your Kafka broker addresses
	topic := topic_name                 // Replace with your topic name
	groupID := topic_name + "-consumer" // Replace with your consumer group ID

	// Set up the Kafka reader (consumer)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID, // Consumer group ID, ensures message offsets are tracked
	})

	// Close the reader when the function returns
	defer r.Close()

	log.Println("Subscribed to topic: " + topic)

	// Loop and read messages from the topic
	var rcv_msg string
	for {
		// Read a message from the topic
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message: " + err.Error())
		}

		rcv_msg = string(msg.Value)
		// Print the message key and value to the console
		log.Println("Received Message (" + topic_name + ") : " + rcv_msg)

		// Send the message over websocket for any clients are currently connected.
		broadcast <- []byte(rcv_msg)

		// Parse the JSON data into a struct
		var data SensorData
		err = json.Unmarshal([]byte(rcv_msg), &data)
		if err != nil {
			log.Println("Error unmarshalling JSON:" + err.Error())
			return
		}

		log.Println("Serial: " + data.Serial)
		log.Println("Type: " + data.Type)
		log.Print("LastUpdate: ")
		log.Println(data.Timestamp)
		log.Print("Value 1: ")
		log.Println(data.Reading1)
		log.Print("Value 2: ")
		log.Println(data.Reading2)
		log.Print("Value 3: ")
		log.Println(data.Reading3)

		go insertSensorData(data)
	}
}

func insertSensorData(data SensorData) {

	// Connect to QuestDB
	ctx := context.TODO()
	sender, err := questdb.LineSenderFromConf(ctx, "http::addr=questdb:9000;")
	if err != nil {
		log.Println("Error connecting to QuestDB: " + err.Error())
	}

	// Make sure to close the sender on exit to release resources.
	defer sender.Close(ctx)

	// Create the row and insert into QuestDB
	last_updated, _ := time.Parse(time.RFC3339, data.Timestamp)
	err = sender.
		Table("sensor_historical_data").
		Symbol("serial_number", data.Serial).
		StringColumn("sensor_type", data.Type).
		TimestampColumn("timestamp", last_updated).
		Float64Column("reading1", data.Reading1).
		Float64Column("reading2", data.Reading2).
		Float64Column("reading3", data.Reading3).
		At(ctx, last_updated)

	if err != nil {
		log.Println("Insert data error: " + err.Error())
		sender.Close(ctx)
		return
	}

	// Make sure that the messages are sent over the network.
	err = sender.Flush(ctx)
	if err != nil {
		log.Println("Flush context error: " + err.Error())
		sender.Close(ctx)
		return
	}

	log.Println("Insert data OK.")
	//sender.Close(ctx)

}
