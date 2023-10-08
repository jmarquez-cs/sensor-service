package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"sync"
)

type Sensor struct {
	Name     string   `json:"name"`
	Location GPS      `json:"location"`
	Tags     []string `json:"tags"`
}

type GPS struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var sensors = make(map[string]Sensor)
var mu sync.Mutex

func main() {
	http.HandleFunc("/sensor", storeSensorHandler)
	http.HandleFunc("/sensor/", sensorByNameHandler)
	http.HandleFunc("/nearest", nearestSensorHandler)

	http.ListenAndServe(":8080", nil)
}

func storeSensorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusBadRequest)
		return
	}

	var s Sensor
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	sensors[s.Name] = s
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func sensorByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/sensor/"):]

	switch r.Method {
	case http.MethodGet:
		getSensorByName(name, w)
	case http.MethodPut:
		updateSensor(name, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func nearestSensorHandler(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)

	var nearestSensor Sensor
	var minDist float64 = math.MaxFloat64

	mu.Lock()
	for _, sensor := range sensors {
		dist := distance(sensor.Location, GPS{lat, lng})
		if dist < minDist {
			minDist = dist
			nearestSensor = sensor
		}
	}
	mu.Unlock()

	if minDist == math.MaxFloat64 {
		http.Error(w, "No sensors found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(nearestSensor)
}

func distance(p1, p2 GPS) float64 {
	// This is a simplified calculation and may not be accurate for large distances.
	// Consider using the Haversine formula for a more accurate calculation.
	return math.Sqrt(math.Pow(p2.Latitude-p1.Latitude, 2) + math.Pow(p2.Longitude-p1.Longitude, 2))
}

func getSensorByName(name string, w http.ResponseWriter) {
	mu.Lock()
	sensor, exists := sensors[name]
	mu.Unlock()

	if !exists {
		http.Error(w, "Sensor not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(sensor)
}

func updateSensor(name string, w http.ResponseWriter, r *http.Request) {
	var updatedSensor Sensor
	if err := json.NewDecoder(r.Body).Decode(&updatedSensor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := sensors[name]
	if !exists {
		mu.Unlock()
		http.Error(w, "Sensor not found", http.StatusNotFound)
		return
	}

	updatedSensor.Name = name
	sensors[name] = updatedSensor
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSensor)
}
