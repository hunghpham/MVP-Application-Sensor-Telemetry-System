package models

import (
    "time"
)

type Sensor struct {
    Serial      string      `binding: required`
    DateTime    time.Time
    Value       float64     `binding: required`
}

var sensors = []Sensor{}

func (s Sensor) Save() {
    //Later add it to the database
    sensors = append(sensors, s)
}

func GetAllSensors() []Sensor {
    return sensors
}