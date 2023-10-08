package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() {
	mu.Lock()
	for k := range sensors {
		delete(sensors, k)
	}
	mu.Unlock()
}

func TestStoreSensorHandler(t *testing.T) {
	// Mock a request to our endpoint
	sensor := Sensor{Name: "TestSensor", Location: GPS{Latitude: 10.0, Longitude: 10.0}}
	body, _ := json.Marshal(sensor)
	req, err := http.NewRequest("POST", "/sensor", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeSensorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetSensorByName(t *testing.T) {
	// Add a test sensor for retrieval
	testSensor := Sensor{Name: "TestSensor", Location: GPS{Latitude: 10.0, Longitude: 10.0}}
	sensors[testSensor.Name] = testSensor

	req, err := http.NewRequest("GET", "/sensor/TestSensor", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sensorByNameHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Further checks can be made on the response body
}

func TestUpdateSensor(t *testing.T) {
	setup()

	// Add a test sensor for updating
	testSensor := Sensor{Name: "TestSensorToUpdate", Location: GPS{Latitude: 20.0, Longitude: 20.0}}
	sensors[testSensor.Name] = testSensor

	updatedSensor := Sensor{Location: GPS{Latitude: 30.0, Longitude: 30.0}}
	body, _ := json.Marshal(updatedSensor)

	req, err := http.NewRequest("PUT", "/sensor/TestSensorToUpdate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sensorByNameHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check if the sensor actually updated in the global map
	mu.Lock()
	updated, exists := sensors["TestSensorToUpdate"]
	mu.Unlock()

	if !exists || updated.Location.Latitude != 30.0 || updated.Location.Longitude != 30.0 {
		t.Error("Failed to update the sensor")
	}
}

func TestNearestSensorHandler(t *testing.T) {
	setup()

	// Add some sensors
	sensors["Sensor1"] = Sensor{Name: "Sensor1", Location: GPS{Latitude: 10.0, Longitude: 10.0}}
	sensors["Sensor2"] = Sensor{Name: "Sensor2", Location: GPS{Latitude: 20.0, Longitude: 20.0}}
	sensors["Sensor3"] = Sensor{Name: "Sensor3", Location: GPS{Latitude: 30.0, Longitude: 30.0}}

	req, err := http.NewRequest("GET", "/nearest?lat=15&lng=15", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nearestSensorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Decode the response to check if the correct sensor is returned
	var nearest Sensor
	err = json.NewDecoder(rr.Body).Decode(&nearest)
	if err != nil {
		t.Fatal(err)
	}
	if nearest.Name != "Sensor1" {
		t.Errorf("Expected nearest sensor to be Sensor1, got %s", nearest.Name)
	}
}

func TestStoreSensorHandlerBadRequest(t *testing.T) {
	setup()

	// Invalid JSON
	body := []byte(`{"name": "BadSensor" "location": {"latitude": 10.0, "longitude": 10.0}}`) // Missing comma
	req, err := http.NewRequest("POST", "/sensor", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeSensorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for bad request: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestStoreSensorHandlerWrongMethod(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/sensor", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeSensorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for wrong method: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetNonexistentSensor(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/sensor/NonexistentSensor", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sensorByNameHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code for nonexistent sensor: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUnsupportedMethodForSensorByName(t *testing.T) {
	setup()

	req, err := http.NewRequest("DELETE", "/sensor/SomeSensor", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sensorByNameHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code for unsupported method: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestNearestSensorWithEmptyData(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/nearest?lat=15&lng=15", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nearestSensorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code when no sensors are available: got %v want %v", status, http.StatusNotFound)
	}
}
