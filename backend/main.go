package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	//"bytes"

	"strings"

	"github.com/gin-gonic/gin"
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
}

func main() {

	server := gin.Default()

	//GET endpoints
	server.GET("/api/get_all_sensors", getAllSensors)                                     //connect to QuestDB to get all sensors info
	server.GET("/api/get_sensor_historical_data", getHistoricalData)                      //connect to QuestDB to get historical data
	server.GET("/api/create_kafka_topic_and_subscribe", create_kafka_topic_and_subscribe) //create kafka topic through Kafka rest api

	//Start the backend server listening on port 8080
	log.Fatal(server.Run(":8080"))
}

// /////////////// Database related functions /////////////////////
// Get all current sensors
func getAllSensors(ctx *gin.Context) {

	// Define the SQL query
	query := `
			SELECT serial_number, timestamp, sensor_type, reading1, reading2
			FROM sensor_historical_data
			LATEST BY serial_number;
			`

	u, err := url.Parse("http://localhost:9000")
	if err != nil {
		log.Println("Parse URL error: " + err.Error())
	}

	u.Path += "exec"
	params := url.Values{}
	params.Add("query", query)
	u.RawQuery = params.Encode()
	url := fmt.Sprintf("%v", u)

	res, err := http.Get(url)
	if err != nil {
		log.Println("Query through QUESTDB Rest API error: " + err.Error())
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Read response body error: " + err.Error())
	}

	log.Println(string(body))

	ctx.JSON(http.StatusCreated, gin.H{"message": "OK"})

}

// Get individual sensor historical data based on datetime range
func getHistoricalData(context *gin.Context) {
	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

// //////////////// Other helper functions
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
func create_kafka_topic_and_subscribe(context *gin.Context) {

	url := "http://kafka-rest-proxy:8082/v3/clusters"

	//Get Kafka local cluster id using Kafka REST API proxy, demonstrating REST works as well
	cluster_id, err := getKafkaClusterID(url)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"Status": "Server error, unable to get Kafka cluster id.", "Error": err.Error()})
	}

	// Get topic name from the post
	topic := context.Query("topic_name")
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
		context.JSON(http.StatusOK, gin.H{"status": "Connect to Kafka failed!", "error": err.Error()})
	}
	defer conn.Close()

	// Create the topic
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"status": "Failed to create topic!", "error": err.Error()})
		return
	}

	//add path to url for creating publish message endpoint
	new_url := strings.ReplaceAll(url, "kafka-rest-proxy", "localhost") + "/" + cluster_id + "/topics"

	//Kafka topic is created, now start the listening for messages on this topic by subscribing to it
	go subscribeToKafkaTopic(topic)

	//Return messages to user
	context.JSON(http.StatusOK, gin.H{"Status": "OK", "Kafka_Endpoint": new_url + "/" + topic + "/records"})

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
	//defer sender.Close(ctx)

	// Create the row and insert into QuestDB
	last_updated, _ := time.Parse(time.RFC3339, data.Timestamp)
	err = sender.
		Table("sensor_historical_data").
		Symbol("serial_number", data.Serial).
		StringColumn("sensor_type", data.Type).
		TimestampColumn("timestamp", last_updated).
		Float64Column("reading1", data.Reading1).
		Float64Column("reading2", data.Reading2).
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
	sender.Close(ctx)

}
