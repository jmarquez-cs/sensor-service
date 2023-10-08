#!/bin/bash

# Ask the user for the number of sensors to generate
read -p "Enter the number of sensors you want to generate: " SENSOR_COUNT

# Validates if the input is a number
if ! [[ $SENSOR_COUNT =~ ^[0-9]+$ ]]; then
   echo "Error: Please enter a valid number."
   exit 1
fi

# Defines the base curl command with a 10-second timeout
BASE_CMD="curl --max-time 15 -X POST http://localhost:8080/sensor -H \"Content-Type: application/json\""

# Defines sensor data payloads
declare -a PAYLOADS=()

# Function generates a random float number within a range
generate_random_float() {
  printf "%0.6f\n" "$(echo "$1 + ( $RANDOM * ($2 - $1) / 32767 )" | bc)"
}

# Start generating payloads
echo "Generating random sensor payloads..."

# Generates random payloads based on the user's input
for i in $(seq 1 $SENSOR_COUNT); do
  RAND_LAT=$(generate_random_float 30.788990 40.788990)
  RAND_LNG=$(generate_random_float -70.948499 -80.908499)
  PAYLOADS+=("{\"name\":\"sensor_$i\", \"location\":{\"latitude\":$RAND_LAT, \"longitude\":$RAND_LNG}, \"tags\":[\"outdoor\", \"temperature\"]}")
done

# Confirm payloads have been generated
echo "Payloads generated successfully."

# Loop over payloads and send the data
echo "Sending payloads to the server..."

for payload in "${PAYLOADS[@]}"; do
  echo "Sending payload: $payload"
  
  # Send payload and print the server response
  RESPONSE=$($BASE_CMD -d "$payload")
  echo "Server response: $RESPONSE"
done

echo "Sensor data creation complete."
