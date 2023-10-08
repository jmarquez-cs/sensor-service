# sensor-service

This is a simple API Go project serving a RESTful API for storing and querying sensor metadata.

## Prerequisites
- [Go v1.17 installed](https://go.dev/doc/install)
- Optional: [Docker installed](https://docs.docker.com/engine/install/)

## How to run sensor-service
1. `cd /cmd/app && go run main.go`

## How to manually exercise sensor-service endpoints 
1. Create data
    ```shell
    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_1", "location":{"latitude":40.730610, "longitude":-73.924240}, "tags":["outdoor", "temperature"]}'

    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_2", "location":{"latitude":40.730610, "longitude":-73.916112}, "tags":["outdoor", "temperature"]}'

    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_4", "location":{"latitude":40.730610, "longitude":-73.934999}, "tags":["outdoor", "temperature"]}'

    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_5", "location":{"latitude":40.730610, "longitude":-73.911444}, "tags":["outdoor", "temperature"]}'

    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_6", "location":{"latitude":40.730610, "longitude":-73.938442}, "tags":["outdoor", "temperature"]}'

    curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"name":"sensor_7", "location":{"latitude":40.730610, "longitude":-73.924241}, "tags":["outdoor", "temperature"]}'
    ```
2. Read previously created sensor by name
    ```shell
    curl -X GET http://localhost:8080/sensor/sensor_1
    ```
3. Retrieve nearest sensor
    ```shell
    curl -X GET "http://localhost:8080/nearest?lat=40.730600&lng=-73.935230"
    ```
4. Update an existing sensor by name
    ```shell
    curl -X PUT -H "Content-Type: application/json" \ -d '{"name": "sensor1", "location": {"latitude": 52.5200, "longitude": 13.4050}, "tags": ["indoor", "humidity' \
    http://localhost:8080/sensor/sensor1
    ```

### Optional: 

Automatically create sensor data
1. Run the `./create_sensors.sh` script located at `scripts/create_sensors.sh`  
    **Note:** curl timeout has hard coded value of 15 seconds.

2. Import `configs/sensor-service.postman_collection.json` into your Postman tool  
    **Note:** collection created and exported in Postman Version 10.16.5

3. Run unit tests on behavior of endpoints and code coverage
- Fom the root directory run `go test -cover ./cmd/app/`

## How to build and run binary file
1. Build the project using `go build -o bin/app ./cmd/app`
2. Run the compiled binary: `./bin/app`

## Install and run using Docker
1. Build the Docker image: `docker-compose build`
2. Run the Docker container: `docker-compose up`

## Features to implement

- Authentication & Authorization
   - Implement JWT (JSON Web Tokens) for stateless authentication.
   - Role-based access control for different endpoints.

- Rate Limiting
   - Introduce rate limiting to prevent abuse and protect the service.

- Logging
    - Choose/import a log package
    - Create an output destination for logs
    - Write logs using different levels of severity
    - Format logs
    - Set log flags
    - Handle panics
- Monitoring 
    - Integrate with monitoring tools like Prometheus and Grafana for real-time service metrics.
- Pagination and Filtering
    - Add pagination to endpoints returning multiple items.
    - Implement filtering options for more refined queries.
- Caching
    - Introduce caching mechanisms (e.g., Redis) to improve response times and reduce database load.
- Database Improvements
    - Implement database connection pooling.
    - Explore database replication and sharding for scalability.
    - Regularly back up the database and implement a recovery strategy.
- Error Handling
    - Comprehensive error handling and user-friendly error messages.
    - Implement a global error handler for uniform error responses.
- API Versioning (semantic versioning)
    - Introduce API versioning to ensure backward compatibility as the service evolves.
- Documentation & Testing
    - Use tools like Swagger for API documentation and interactive testing.
    - Implement unit and integration testing. Consider using a CI/CD pipeline for automated tests and deployments.
- WebSockets & Real-time Communication
    - Implement WebSockets for endpoints requiring real-time data updates.
    - End-to-End Encryption
    - Ensure data integrity and confidentiality with end-to-end encryption.
- Multi-Environment Setup
    - Establish separate environments for development, testing, staging, and production.
- Localization & Internationalization
    - Make the service adaptable to different languages and regions.
- Data Backup & Recovery
    - Regular automated backups of all data.
    - Implement a robust disaster recovery plan.
- Horizontal Scalability
    - Design the service to scale out, accommodating a larger number of users and requests.