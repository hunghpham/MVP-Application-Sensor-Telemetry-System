package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	//"encoding/json"
	//"bytes"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/tidwall/gjson"

	//	"github.com/questdb/go-questdb-client/v3"
	"mesa.com/backend/models"
)

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
func getAllSensors(context *gin.Context) {
	sensors := models.GetAllSensors()
	context.JSON(http.StatusOK, sensors)

}

// Get individual sensor historical data based on datetime range
func getHistoricalData(context *gin.Context) {

}

func createSensorRecord(context *gin.Context) {
	var sensor models.Sensor
	err := context.ShouldBindJSON(&sensor)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create sensor!"})
		return
	}

	sensor.Serial = "9988"
	sensor.DateTime = time.Now().UTC()
	sensor.Value = 25.7
	sensor.Save()

	context.JSON(http.StatusCreated, gin.H{"message": "Sensor Created", "Sensor": sensor})
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

	//     // Create the data structure to send in the POST request body, using REST API proxy here, works as well
	//     data := map[string]string {
	// 		"topic_name": string(topic),
	// 	}
	//
	//     // Convert the data to JSON
	// 	jsonData, err := json.Marshal(data)
	// 	if err != nil {
	// 		context.JSON(http.StatusOK, gin.H{"Status": "Server error, encoding JSON failed.", "Error": err.Error()})
	// 		return
	// 	}
	//
	//     // Create a new POST request with JSON data
	//     url = url + "/" + cluster_id + "/topics"
	// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	// 	if err != nil {
	// 		context.JSON(http.StatusOK, gin.H{"Status": "Server error, creating request failed.", "Error": err.Error()})
	// 		return
	// 	}
	//
	//     // Set the appropriate headers
	// 	req.Header.Set("Content-Type", "application/json")
	//
	//     // Send the POST request
	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		context.JSON(http.StatusOK, gin.H{"Status": "Server error, sending request failed.", "Error": err.Error()})
	// 		return
	// 	}
	// 	defer resp.Body.Close()
	//
	// 	// Check the response status
	// 	if resp.StatusCode == http.StatusOK {
	// 		context.JSON(http.StatusOK, gin.H{"Status": "OK", "Post_endpoint": url + "/" + topic + "/records"})
	// 		return
	// 	} else {
	// 		context.JSON(http.StatusOK, gin.H{"Status": "Server error, sending request failed.", "Error": err.Error()})
	// 		return
	// 	}

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
	for {
		// Read a message from the topic
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message: " + err.Error())
		}

		// Print the message key and value to the console
		log.Println("Received Message (" + topic_name + ") : " + string(msg.Value))
	}
}
