package main

import (
	"context"
	"encoding/json"
	"fmt"
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
)

// Sensor data structure, use for marhaling JSON and vice versa
type SensorData struct {
	Serial    string  `json:"Serial"`
	Type      string  `json:"Type"`
	Timestamp string  `json:"Timestamp"`
	Reading1  float64 `json:"Reading1"`
	Reading2  float64 `json:"Reading2"`
	Reading3  float64 `json:"Reading3"`
}

// data structure, use for marhaling JSON for min max report
type MinMaxData struct {
	Serial      string  `json:"Serial"`
	Type        string  `json:"Type"`
	Timestamp   string  `json:"Timestamp"`
	Reading1Min float64 `json:"Reading1Min"`
	Reading1Max float64 `json:"Reading1Max"`
	Reading1Avg float64 `json:"Reading1Avg"`
	Reading2Min float64 `json:"Reading2Min"`
	Reading2Max float64 `json:"Reading2Max"`
	Reading2Avg float64 `json:"Reading2Avg"`
	Reading3Min float64 `json:"Reading3Min"`
	Reading3Max float64 `json:"Reading3Max"`
	Reading3Avg float64 `json:"Reading3Avg"`
}

// //////////////////////////// WEB SOCKET SECTION ///////////////////////////////////////
// Global slice to hold active WebSocket connections
var clients = make(map[*websocket.Conn]bool)

// Broadcast channel for sending messages
var broadcast = make(chan []byte)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origin clients
	},
}

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

// //////////////////////////// Database related functions //////////////////////////////////////////
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

func executeQueryMinMaxReport(query string) string {

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

	dataList := []MinMaxData{}
	for rows.Next() {

		var serialNumber string
		var timestamp string
		var sensorType string
		var reading1min float64
		var reading1max float64
		var reading1avg float64
		var reading2min float64
		var reading2max float64
		var reading2avg float64
		var reading3min float64
		var reading3max float64
		var reading3avg float64

		// Scan the data into variables
		if err := rows.Scan(&serialNumber, &sensorType, &reading1max, &reading1min, &reading1avg, &reading2max, &reading2min, &reading2avg, &reading3max, &reading3min, &reading3avg, &timestamp); err != nil {
			log.Println("Failed to scan row: " + err.Error())
		}

		//add each record to a list to return
		var data MinMaxData
		data.Serial = serialNumber
		data.Type = sensorType
		data.Timestamp = timestamp
		data.Reading1Min = reading1min
		data.Reading1Max = reading1max
		data.Reading1Avg = reading1avg

		data.Reading2Min = reading2min
		data.Reading2Max = reading2max
		data.Reading2Avg = reading2avg

		data.Reading3Min = reading3min
		data.Reading3Max = reading3max
		data.Reading3Avg = reading3avg

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

// Get min, max and avg values
func getMinMaxAvgData(ctx *gin.Context) {
	query := `SELECT
				serial_number,
				sensor_type,
				MAX(reading1) AS max_reading1,
				MIN(reading1) AS min_reading1,
				AVG(reading1) AS avg_reading1,
				MAX(reading2) AS max_reading2,
				MIN(reading2) AS min_reading2,
				AVG(reading2) AS avg_reading2,
				MAX(reading3) AS max_reading3,
				MIN(reading3) AS min_reading3,
				AVG(reading3) AS avg_reading3,
				date_trunc('#PERIOD#', timestamp) AS period  -- Use date_trunc to bucket by hour
			FROM
				sensor_historical_data
			WHERE
				serial_number = '#SERIAL#'
			GROUP BY
				serial_number,
				sensor_type,
				period
			ORDER BY
				period;`

	//Get sensor serial number and the date range
	serial := ctx.Query("serial_number")
	interval := ctx.Query("interval")

	//modify query
	query = strings.ReplaceAll(query, "#SERIAL#", serial)
	query = strings.ReplaceAll(query, "#PERIOD#", interval)

	fmt.Println(query)

	result := executeQueryMinMaxReport(query)

	ctx.String(http.StatusOK, result)
}

func insertSensorData(data Message) {

	// Connect to QuestDB
	ctx := context.TODO()
	sender, err := questdb.LineSenderFromConf(ctx, "http::addr=questdb:9000;")
	if err != nil {
		log.Println("Error connecting to QuestDB: " + err.Error())
	}

	// Make sure to close the sender on exit to release resources.
	defer sender.Close(ctx)

	// Create the row and insert into QuestDB
	last_updated, _ := time.Parse(time.RFC3339, data.Value.Data.Timestamp)
	err = sender.
		Table("sensor_historical_data").
		Symbol("serial_number", data.Value.Data.Serial).
		StringColumn("sensor_type", data.Value.Data.Type).
		TimestampColumn("timestamp", last_updated).
		Float64Column("reading1", data.Value.Data.Reading1).
		Float64Column("reading2", data.Value.Data.Reading2).
		Float64Column("reading3", data.Value.Data.Reading3).
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

// Create a kafka topic and then call another method to subscribe to Kafka topic using native client
func create_kafka_topic_then_subscribe(topic_name string) {
	topic := topic_name
	broker := "kafka:9092"
	partitions := 1
	replicationFactor := 1

	// Initialize a Kafka writer using native Kafka client library, demonstrating native client works as well
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer writer.Close()

	var conn *kafka.Conn
	var err error

	// Retry connecting to Kafka until successful
	for {
		conn, err = kafka.Dial("tcp", broker)
		if err != nil {
			log.Println("Connect to Kafka failed! Retrying... " + err.Error())
			time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
			continue
		}

		// Check if the Kafka connection is responsive by fetching broker metadata
		_, err = conn.Brokers()
		if err != nil {
			log.Println("Kafka connection established, but failed to retrieve broker metadata. Retrying... " + err.Error())
			conn.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to Kafka successfully and broker is responsive.")
		break
	}
	defer conn.Close()

	// Create the topic
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		log.Println("Failed to create topic! " + err.Error())
		return
	}

	//Kafka topic is created, now start the listening for messages on this topic by subscribing to it
	go subscribeToKafkaTopic(topic)
}

// Define the structure for the JSON data
type Message struct {
	Value struct {
		Type string `json:"type"`
		Data struct {
			Serial    string  `json:"Serial"`
			Type      string  `json:"Type"`
			Timestamp string  `json:"Timestamp"`
			Reading1  float64 `json:"Reading1"`
			Reading2  float64 `json:"Reading2"`
			Reading3  float64 `json:"Reading3"`
		} `json:"data"`
	} `json:"value"`
}

// Subscribe to Kafka
func subscribeToKafkaTopic(topic_name string) {
	// Define the Kafka broker addresses and topic
	brokers := []string{"kafka:9092"}
	topic := topic_name
	groupID := topic_name + "-consumer"

	// Set up the Kafka reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID, // Consumer group ID, ensures message offsets are tracked
	})

	// Close the reader when the function returns
	defer reader.Close()

	log.Println("Subscribed to topic: " + topic)

	// Loop and read messages from the topic
	var rcv_msg string
	for {
		// Read a message from the topic
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message: " + err.Error())
		}

		rcv_msg = string(msg.Value)
		// Print the message key and value to the console
		log.Println("Received Message (" + topic_name + ") : " + rcv_msg)

		// Send the message over websocket for any clients are currently connected.
		broadcast <- []byte(rcv_msg)

		// Parse the JSON data into a struct
		var kafka_msg Message
		err = json.Unmarshal([]byte(rcv_msg), &kafka_msg)
		if err != nil {
			log.Println("Error unmarshalling JSON:" + err.Error())
			return
		}

		// go routine insert the received data from kafka pipeline
		go insertSensorData(kafka_msg)
	}
}

// ///////////////////////////////////////// Main function, entry point ////////////////////////////////////////
func main() {
	server := gin.Default()

	//Apply the CORS middleware
	server.Use(cors.Default())

	//REST GET endpoints
	server.GET("/api/get_all_sensors", getAllSensors)                //connect to QuestDB to get all sensors info
	server.GET("/api/get_sensor_historical_data", getHistoricalData) //connect to QuestDB to get historical data
	server.GET("/api/get_sensor_min_max_avg_data", getMinMaxAvgData) //connect to QuestDB to get historical data

	// Set up WebSocket endpoint
	server.GET("/ws", handleWebSocket)

	// Start the broadcast goroutine, url is ws://localhost:8080/ws
	go broadcastMessages()

	// Create a topic and subcribe
	go create_kafka_topic_then_subscribe("sensor_data")

	//Start the backend server listening on port 8080
	log.Fatal(server.Run(":8080"))
}
