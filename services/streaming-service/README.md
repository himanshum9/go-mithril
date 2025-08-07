# Streaming Service Documentation

## Overview
The Streaming Service is a microservice responsible for processing and streaming location data submitted by tenant users to a third-party application in real-time. It ensures that the data is handled securely and efficiently, adhering to the multi-tenant architecture of the overall system.

## Features
- Real-time streaming of location data to third-party applications.
- Robust error handling and connection management for streaming processes.
- Integration with AWS Cognito for secure user authentication and authorization.

## Setup Instructions
1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd github.com/himanshum9/go-mithril/services/streaming-service
   ```

2. **Install Dependencies**
   Ensure you have Go installed and run:
   ```bash
   go mod tidy
   ```

3. **Configuration**
   Update the configuration settings in `configs/config.yaml` to include the necessary details for connecting to the third-party application and any other required parameters.

4. **Run the Service**
   Start the streaming service by executing:
   ```bash
   go run main.go
   ```

## API Endpoints
- **POST /stream**
  - Description: Initiates the streaming of location data to the third-party application.
  - Request Body: JSON object containing location data (latitude, longitude, tenant ID).
  
- **GET /status**
  - Description: Checks the status of the streaming service and its connection to the third-party application.

## Error Handling
The streaming service includes mechanisms to handle failures in the streaming process. If a connection to the third-party application fails, the service will attempt to re-establish the connection automatically.

## Contribution
For contributions, please follow the standard Git workflow:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a clear description of your changes.

## License
This project is licensed under the MIT License - see the LICENSE file for details.