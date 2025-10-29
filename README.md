# Image Metadata API

A RESTful API built with Go for extracting and managing image metadata (width, height, format) from image URLs. The API includes a background worker that asynchronously processes images and stores their metadata in a PostgreSQL database.

## Features

- **Secure Authentication** - Bearer token authentication for API endpoints
- **Image Processing** - Asynchronous background worker extracts metadata from images
- **PostgreSQL Database** - Persistent storage for image metadata
- **Docker Support** - Easy deployment with Docker Compose
- **Health Checks** - Built-in health check endpoint
- **Comprehensive Tests** - Test coverage for all handlers and workers
- **CORS Enabled** - Cross-origin resource sharing support

## Architecture

The application consists of several components:

- **API Layer**: REST endpoints for managing images
- **Background Worker**: Processes images asynchronously to extract metadata
- **Database**: PostgreSQL for storing image records and metadata
- **Authentication Middleware**: Bearer token-based authentication

### Image Processing Flow

1. Submit an image URL via POST `/v1/image`
2. Image is stored with status "pending"
3. Background worker picks up pending images
4. Worker fetches the image and extracts width, height, and format
5. Image status is updated to "done" with extracted metadata

## Prerequisites

- Go 1.25.2 or higher
- Docker and Docker Compose
- PostgreSQL 16 (or use Docker)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/Giammarco-Ferranti/image-metadata-api.git
cd image-metadata-api
```

2. Create a `.env` file in the root directory:

Generate a secure API key using OpenSSL:

```bash
openssl rand -hex 32
```

Create a `.env` file with the following contents:

```env
PORT=8080
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=image_metadata
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
API_KEY=<paste_your_generated_key_here>
```

Replace `<paste_your_generated_key_here>` with the output from the OpenSSL command above.

3. Run with Docker Compose:

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Set up PostgreSQL database and configure environment variables.

   Generate a secure API key using OpenSSL:

   ```bash
   openssl rand -hex 32
   ```

   Create a `.env` file with your database credentials and the generated API key.

3. Run the application:

```bash
go run cmd/main.go
```

## API Endpoints

All API endpoints (except health check) require authentication using a Bearer token.

### Health Check

```
GET /healthz
```

Returns the health status of the API.

**Response:** HTTP 200 OK

### Get All Images

```
GET /v1/images
Authorization: Bearer <API_KEY>
```

Retrieves all images from the database.

**Response:**

```json
{
  "data": [
    {
      "id": "uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "url": "https://example.com/image.jpg",
      "status": "done",
      "width": 1920,
      "height": 1080,
      "format": "jpeg"
    }
  ]
}
```

### Add Image

```
POST /v1/image
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

Submits an image URL for processing.

**Request Body:**

```json
{
  "url": "https://example.com/image.jpg"
}
```

**Response:**

```json
{
  "data": {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "url": "https://example.com/image.jpg",
    "status": "pending"
  }
}
```

### Get Image by ID

```
GET /v1/image/{id}
Authorization: Bearer <API_KEY>
```

Retrieves a specific image by its UUID.

**Response:**

```json
{
  "data": {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "url": "https://example.com/image.jpg",
    "status": "done",
    "width": 1920,
    "height": 1080,
    "format": "jpeg"
  }
}
```

### Delete Image

```
DELETE /v1/image/{id}
Authorization: Bearer <API_KEY>
```

Deletes an image from the database.

**Response:**

```json
"Successfully deleted image"
```

HTTP Status: 200 OK

## Image Status

Images progress through the following statuses:

- **pending**: Image has been submitted but not yet processed
- **in process**: Image is currently being processed by the worker
- **done**: Image has been successfully processed and metadata extracted
- **failed**: Image processing failed (invalid URL, unsupported format, etc.)

## Background Worker

The application includes a background worker that:

- Runs every minute (configurable)
- Processes up to 5 pending images concurrently (configurable)
- Fetches images from their URLs
- Extracts width, height, and format (supports PNG, JPEG, and GIF)
- Updates the database with extracted metadata

Worker configuration can be adjusted in `cmd/main.go`:

```go
go worker.StartExtract(db, time.Minute, 5)
```

Parameters:

- First parameter: Time between worker runs
- Second parameter: Number of images to process concurrently

## Testing

Run all tests:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./pkg/api/...
go test ./pkg/worker/...
```

Run tests with coverage:

```bash
go test ./... -cover
```

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── pkg/
│   ├── api/                     # HTTP handlers
│   │   ├── handler_add_image.go
│   │   ├── handler_delete_image.go
│   │   ├── handler_get_image.go
│   │   ├── handler_get_images.go
│   │   └── routes.go
│   ├── health/                  # Health check handler
│   ├── middleware/              # Authentication middleware
│   ├── models/                  # Data models
│   ├── responses/               # Response utilities
│   └── worker/                  # Background worker
├── tests/                       # Test files
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                   # Docker image definition
├── docs/                        # API documentation (OpenAPI)
└── go.mod                       # Go dependencies

```

## Configuration

### Environment Variables

| Variable            | Description                     | Default  |
| ------------------- | ------------------------------- | -------- |
| `PORT`              | Server port                     | 8080     |
| `POSTGRES_USER`     | PostgreSQL username             | Required |
| `POSTGRES_PASSWORD` | PostgreSQL password             | Required |
| `POSTGRES_DB`       | Database name                   | Required |
| `POSTGRES_HOST`     | PostgreSQL host                 | Required |
| `POSTGRES_PORT`     | PostgreSQL port                 | Required |
| `API_KEY`           | Bearer token for authentication | Required |

## Supported Image Formats

The worker currently supports:

- PNG
- JPEG
- GIF

## Development

### Adding New Endpoints

1. Create a handler in `pkg/api/`
2. Add the route in `pkg/api/routes.go`
3. Write tests in `tests/api/`

### Database Migrations

The application uses GORM's auto-migration. The database schema is automatically created/updated on startup. See `cmd/main.go`:

```go
db.AutoMigrate(&models.ImageProcess{})
```

## Author

Giammarco Ferranti
