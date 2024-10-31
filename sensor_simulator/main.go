package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/segmentio/kafka-go"
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
	num_of_sensors := 5
	serial_characters := "ABCD"
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

// post sensor information
func postSensorInfo(serial string, sensor_type string, delay_interval int) {

	// Initialize an empty Payload struct
	var payload Payload

	// Define Kafka writer configuration
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
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
