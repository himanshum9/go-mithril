# Location Service

This document provides an overview of the Location Service within the multi-tenant location streaming system.

## Overview

The Location Service is responsible for collecting and processing geographical location data submitted by tenant users. It ensures that each tenant's data is isolated and securely handled. The service streams the collected location data to a third-party application in real-time.

## Features

- Multi-tenant architecture with data isolation using tenant identifiers.
- Real-time location data submission and processing.
- Integration with third-party applications for data streaming.

## Setup Instructions

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd github.com/himanshum9/go-mithrilservices/location-service
   ```

2. **Install Dependencies**
   Ensure you have Go installed and run:
   ```bash
   go mod tidy
   ```

3. **Configuration**
   Update the configuration file located at `../../configs/config.yaml` with the necessary settings, including database connection details and AWS Cognito settings.

4. **Run the Service**
   Start the Location Service:
   ```bash
   go run main.go
   ```

## API Endpoints

- **Submit Location Data**
  - **Endpoint:** `POST /locations`
  - **Description:** Submits location data (latitude, longitude) for the tenant user.
  - **Request Body:**
    ```json
    {
      "latitude": <float>,
      "longitude": <float>
    }
    ```

- **Stream Location Data**
  - **Endpoint:** `GET /locations/stream`
  - **Description:** Streams location data to a third-party application.

## Error Handling

The service includes mechanisms to handle errors during data submission and streaming. In case of a failure, the service will attempt to re-establish connections and ensure data integrity.

## License

This project is licensed under the MIT License. See the LICENSE file for details.