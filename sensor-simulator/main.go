package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/segmentio/kafka-go"
	"github.com/tidwall/gjson"
)

// Struct to match the JSON structure
type Payload struct {
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

func main() {

	//Read the config file
	file, err := os.Open("config.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)
	var lines []string
	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for any error during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// //Get post url, Kafka URL for posting sensor to a topic channel
	// post_url, err := getKafkaPublishURLEndpointandStartListener()
	// fmt.Println(post_url)
	// if err != nil {
	// 	fmt.Println("Unable to start Kafka Rest endpoint listener for sensor data!, try again!")
	// }

	//splitStr := strings.Split(lines[0], "=")
	//post_url := strings.TrimSpace(splitStr[1])

	//Get number of sensor
	splitStr := strings.Split(lines[1], "=")
	number_of_sensors := strings.TrimSpace(splitStr[1])

	//Get sensor serial start with characters
	splitStr = strings.Split(lines[2], "=")
	serial_characters := strings.TrimSpace(splitStr[1])

	num_of_sensors, err := strconv.Atoi(number_of_sensors)
	list_of_sensors := []string{}
	//Generate sensor list
	for i := 0; i < num_of_sensors; i++ {
		list_of_sensors = append(list_of_sensors, serial_characters+fmt.Sprintf("%s%s", strings.Repeat(string('0'), 4-len(strconv.Itoa(i))), strconv.Itoa(i)))
	}

	// Define the range
	min := 5
	max := 10

	for _, serial := range list_of_sensors {
		// Create a new random number generator
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Generate a random integer
		randomInt := r.Intn(max-min) + min // Generates a random integer between min and max
		fmt.Println("RadomInt: " + strconv.Itoa(randomInt))

		time.Sleep(time.Second * 2)

		go postSensorInfo(serial, "Incubator", randomInt)
	}

	for {
		time.Sleep(time.Hour)
	}

}

// This return the Kafka endpoint and start the topic listener
func getKafkaPublishURLEndpointandStartListener() (string, error) {

	//Send GET request
	resp, err := http.Get("http://localhost:8080/api/create_kafka_topic_and_subscribe?topic_name=sensor_data")
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
	value := gjson.Get(string(body), "Kafka_Endpoint")

	return value.String(), nil
}

// post sensor information
func postSensorInfo(serial string, sensor_type string, delay_interval int) {

	// Initialize an empty Payload struct
	var payload Payload

	// Define Kafka writer configuration
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "sensor_data",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	//Keep looping
	for {
		// Create a new random number generator
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Define the range
		var min = -10.0
		var max = 30.0
		// Generate a random temperature value within the range [min, max)
		var randomTempValue = min + (r.Float64() * (max - min))

		// Define the range for humidity
		min = 0.0
		max = 100.0
		// Generate a random humidity value within the range [min, max)
		var randomHumValue = min + (r.Float64() * (max - min))

		// Generate a random Co2 value within the range [min, max)
		var randomCo2Value = min + (r.Float64() * (max - min))

		// Assign values to the struct fields
		payload.Value.Data.Reading1 = randomTempValue
		payload.Value.Data.Reading2 = randomHumValue
		payload.Value.Data.Reading3 = randomCo2Value
		payload.Value.Type = "JSON"
		payload.Value.Data.Serial = serial
		payload.Value.Data.Type = sensor_type
		payload.Value.Data.Timestamp = time.Now().UTC().Format(time.RFC3339)

		// Convert the struct to a JSON string
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshaling struct:", err)
			return
		}

		// Convert byte slice to a string and print
		jsonString := string(jsonData)
		fmt.Println("JSON String:", jsonString)

		// Define message
		message := kafka.Message{
			Value: []byte(jsonString),
		}

		// Write message to Kafka
		err = writer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Println("failed to write message to Kafka: " + err.Error())
			continue
		}

		log.Println("Message published successfully")

		time.Sleep(time.Duration(delay_interval) * time.Second)
	}

}
