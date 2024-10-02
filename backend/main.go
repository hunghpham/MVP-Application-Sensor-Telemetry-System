package main

import (
//	"context"
	"log"
	"net/http"
	"io"
 	"time"
 	"errors"
 	//"encoding/json"
 	//"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
//	"github.com/questdb/go-questdb-client/v3"
	"mesa.com/backend/models"
)


func main() {

 	server := gin.Default()

    //GET endpoints
 	server.GET("/api/get_all_sensors", getAllSensors)                   //connect to QuestDB to get all sensors info
    server.GET("/api/get_sensor_historical_data", getHistoricalData)    //connect to QuestDB to get historical data
    server.GET("/api/get_message_broker_topics", getKafkaTopics)        //get Kafka topics

    //POST endpoints
    server.POST("/api/create_kafka_topic", createKafkaTopic)            //Create Kafka topic through Kafka Go client

    //Start the backend server listening on port 8080
 	log.Fatal(server.Run(":8080"))
}


///////////////// Database related functions /////////////////////
//Get all current sensors
func getAllSensors(context *gin.Context) {
    sensors := models.GetAllSensors()
    context.JSON(http.StatusOK, sensors)


}

//Get individual sensor historical data based on datetime range
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


///////////////// Kafka related functions /////////////////////
//Create a Kafka topic
func createKafkaTopic(context *gin.Context) {

    topic := "my_messages"
    broker := "kafka:9092"
    partitions := 1
    replicationFactor := 1

    // Get topic name
	topic = context.Query("topic_name")

    // Initialize a Kafka writer
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{broker},
        Topic:   topic,
    })
    defer writer.Close()

    // Create a new admin client
    conn, err := kafka.Dial("tcp", broker)
    if err != nil {
        context.JSON(http.StatusOK, gin.H{"status": "Connect to Kafka failed!"})
    }
    defer conn.Close()

    // Create the topic
    err = conn.CreateTopics(kafka.TopicConfig{
        Topic:             topic,
        NumPartitions:     partitions,
        ReplicationFactor: replicationFactor,
    })
    if err != nil {
        context.JSON(http.StatusOK, gin.H{"status": "Failed to create topic!"})
        return
    }

    //Return http status OK 200
    context.JSON(http.StatusOK, gin.H{"status": "OK"})
}

//Get an instance Kafka cluster ID based on URL
func getKafkaTopics(context *gin.Context) {

	url := kafka_url
	//cluster_id := ""
	//topics := []string{}

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to make GET request")
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Error HTTP Status Code")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read response body")
	}

	// Unmarshal the JSON response into the struct
// 	var response map[string]interface{}
//
// 	err = json.Unmarshal([]byte(body), &response)
// 	if err != nil {
// 		return "", errors.New("Failed to parse JSON")
// 	}



	return "", "", nil
}

// func produceMsgToKafka(context *gin.Context) {
//
//     // Create a Kafka writer (producer)
// 	writer := kafka.NewWriter(kafka.WriterConfig{
// 		Brokers:  []string{"kafka:9092"},  // Kafka broker address
// 		Topic:    "hung_topic",                  // Topic name
// 		Balancer: &kafka.LeastBytes{},         // Message distribution strategy
// 	})
// 	defer writer.Close()
//
// 	// Produce messages
// 	for i := 0; i < 10; i++ {
// 		msg := kafka.Message{
// 			Key:   []byte(fmt.Sprintf("Key-%d", i)),   // Optional key
// 			Value: []byte(fmt.Sprintf("Message %d", i)), // Message content
// 		}
//
// 		err := writer.WriteMessages(context.Background(), msg)
// 		if err != nil {
// 			fmt.Println("Failed to write message: %v", err)
// 		}
//
// 		fmt.Printf("Produced message: %s\n", msg.Value)
// 		time.Sleep(1 * time.Second) // Simulate some work
// 	}
//
// }



