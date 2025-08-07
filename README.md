# Multi-Tenant Location Streaming System

A scalable, real-time location tracking system built with Go microservices. Enables multi-tenant location data collection and streaming with secure authentication and WebSocket support.

## 🏗️ Project Structure

```
go-mithril/
├── services/                 # Microservices
│   ├── auth-service/         # User authentication with AWS Cognito
│   ├── tenant-service/       # Multi-tenant management
│   ├── location-service/     # Location data collection
│   └── streaming-service/    # Real-time WebSocket streaming
├── deployments/              # Docker and Kubernetes configs
├── migrations/               # Database migrations
├── scripts/                  # Utility scripts
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build and dev commands
└── .env.example             # Environment config template
```

## 🛠️ Technologies Used

- **Go 1.24.5** - Core programming language
- **AWS Cognito** - Secure authentication
- **PostgreSQL** - Multi-tenant database
- **Kafka** - Message streaming
- **WebSockets** - Real-time updates
- **Docker** - Containerization

## ✨ Features

- **Multi-Tenant Architecture**: Isolated tenant data with shared schema
- **Real-Time Streaming**: WebSocket-based location updates
- **Secure Authentication**: JWT-based auth with AWS Cognito
- **Scalable**: Microservices design for horizontal scaling
- **Containerized**: Easy deployment with Docker

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- AWS Account (for Cognito)
- PostgreSQL 13+
- Kafka (included in Docker Compose)

### 1. Environment Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Configure `.env` with your settings:
   ```bash
   # AWS Configuration
   AWS_REGION=your-aws-region
   COGNITO_USER_POOL_ID=your-user-pool-id
   
   # Database
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your-password
   DB_NAME=multi_tenant_db
   
   # Kafka
   KAFKA_BROKER=kafka:9092
   ```

### 2. Start Services

Using Docker Compose (recommended for development):

```bash
# Start all services
docker-compose -f deployments/docker-compose.yml up -d

# Run database migrations
docker-compose -f deployments/docker-compose.yml run --rm migrate

# View logs
docker-compose -f deployments/docker-compose.yml logs -f
```

### 3. Development Workflow

The project includes a Makefile for common tasks:

```bash
# Start all services
make up

# Run database migrations
make migrate

# Access database console
make psql

# Stop all services
make down
```

## 🏃‍♂️ Running Services Individually

For development, you can run services individually:

```bash
# Auth Service (port 8080)
cd services/auth-service && go run main.go

# Tenant Service (port 8080)
cd services/tenant-service && go run main.go

# Location Service (port 8080)
cd services/location-service && go run main.go

# Streaming Service (port 8080)
cd services/streaming-service && go run main.go
```

> **Note**: Ensure all dependencies (PostgreSQL, Kafka) are running when starting services individually.

## 🌐 API Endpoints

All services run on port 8080 by default (Docker maps them to different ports).

### Auth Service
- `POST /api/auth/register` - Register new user
  ```bash
  curl -X POST http://localhost:8080/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"password123"}'
  ```
- `POST /api/auth/login` - Get JWT token
  ```bash
  curl -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"password123"}'
  ```
- `GET /api/protected/me` - Get user info (requires JWT)

### Tenant Service
- `POST /tenants` - Create tenant
- `GET /tenants/{id}` - Get tenant details
- `GET /tenants` - List all tenants

### Location Service
- `POST /location` - Submit location data
  ```bash
  curl -X POST http://localhost:8080/location \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"latitude":37.7749,"longitude":-122.4194}'
  ```

### Streaming Service
- `POST /stream` - Send location data
- `GET /ws` - Connect via WebSocket for real-time updates

## ⚙️ Configuration

Configuration is managed through environment variables in the `.env` file. Here are the key settings:

### Database
```bash
DB_HOST=db                 # Database host (use 'db' for Docker)
DB_PORT=5432              # PostgreSQL port
DB_USER=postgres          # Database user
DB_PASSWORD=yourpassword  # Database password
DB_NAME=multi_tenant_db   # Database name
```

### AWS Cognito
```bash
AWS_REGION=us-east-1                    # AWS region
COGNITO_USER_POOL_ID=your-user-pool-id  # Cognito User Pool ID
COGNITO_CLIENT_ID=your-client-id        # Cognito App Client ID
```

### Kafka
```bash
KAFKA_BROKER=kafka:9092    # Kafka broker address
KAFKA_TOPIC=locations      # Topic for location updates
```

### Service Ports (Docker)
```bash
AUTH_SERVICE_PORT=8080
TENANT_SERVICE_PORT=8080
LOCATION_SERVICE_PORT=8080
STREAMING_SERVICE_PORT=8080
```

> **Note**: In Docker Compose, services communicate using internal port 8080, while Docker maps them to different external ports.

## 🚦 Usage Examples

### 1. Authentication Flow

1. **Register a new user**
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123"}'
   ```

2. **Login to get JWT token**
   ```bash
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123"}'
   ```
   Save the `access_token` from the response for authenticated requests.

### 2. Location Tracking

**Submit location data**
```bash
curl -X POST http://localhost:8080/location \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 37.7749,
    "longitude": -122.4194,
    "device_id": "device-123",
    "accuracy": 15.5
  }'
```

### 3. Real-time Streaming with WebSocket

**JavaScript Example**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket server');  
  // Subscribe to location updates
  ws.send(JSON.stringify({
    type: 'subscribe',
    device_id: 'device-123'  // Optional: filter by device
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Location update:', data);
  // Update UI with new location
};

// Handle errors
ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};
```

**Python Example**
```python
import websockets
import asyncio
import json

async def listen():
    async with websockets.connect('ws://localhost:8080/ws') as ws:
        # Subscribe to updates
        await ws.send(json.dumps({
            'type': 'subscribe',
            'tenant_id': 'your-tenant-id'  # Optional: filter by tenant
        }))
        
        while True:
            data = await ws.recv()
            print(f"Received: {data}")

asyncio.get_event_loop().run_until_complete(listen())
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.