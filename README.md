# Packify

Packify is a Go application that calculates the optimal packs to fulfill customer orders based on specific business rules.

## Business Rules

1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1, send out the least amount of items to fulfill the order.
3. Within the constraints of Rules 1 & 2, send out as few packs as possible to fulfill each order.

Note: Rule #2 takes precedence over rule #3.

## Project Structure

The project follows a standard Go project layout:

```
packify/
├── cmd/
│   └── api/            # API application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Database models
│   └── services/       # Business logic services
├── pkg/
│   └── calculator/     # Pack calculation algorithm
├── .env                # Environment variables
├── docker-compose.yaml # Docker Compose configuration
├── Dockerfile          # Docker build configuration
├── go.mod              # Go module file
└── main.go             # Application entry point
```

## Requirements

- Go 1.23 or higher
- PostgreSQL database
- Docker and Docker Compose (optional)

## Setup

### Using Docker Compose

1. Clone the repository
2. Run the application:

```bash
docker-compose up
```

### Manual Setup

1. Clone the repository
2. Set up a PostgreSQL database
3. Configure the `.env` file with your database credentials
4. Run the application:

```bash
go run main.go
```

## Web UI

Packify includes a web-based user interface built with HTMX and Go templates. The UI provides a user-friendly way to:

1. Calculate optimal packs for orders
2. View and manage available pack sizes
3. See examples of how the pack calculation works

### Accessing the Web UI

Simply navigate to the root URL of the application in your web browser:

```
http://localhost:8080/
```

### Features

- **Home Page**: Overview of the application with a quick calculate form and examples
- **Calculate Packs**: Full page for calculating optimal packs for orders
- **Manage Pack Sizes**: Page for viewing, adding, activating/deactivating, and deleting pack sizes

## API Documentation

### Calculate Packs

Calculates the optimal packs to fulfill an order.

**Endpoint:** `POST /api/calculate`

**Request:**

```json
{
  "itemsOrdered": 501
}
```

**Response:**

```json
{
  "packs": [
    {
      "size": 500,
      "count": 1
    },
    {
      "size": 250,
      "count": 1
    }
  ],
  "totalPacks": 2,
  "totalItems": 750,
  "excessItems": 249
}
```

### Get Pack Sizes

Returns all available pack sizes.

**Endpoint:** `GET /api/pack-sizes`

**Response:**

```json
[
  {
    "ID": 1,
    "CreatedAt": "2023-01-01T00:00:00Z",
    "UpdatedAt": "2023-01-01T00:00:00Z",
    "DeletedAt": null,
    "Size": 250,
    "IsAvailable": true
  },
  {
    "ID": 2,
    "CreatedAt": "2023-01-01T00:00:00Z",
    "UpdatedAt": "2023-01-01T00:00:00Z",
    "DeletedAt": null,
    "Size": 500,
    "IsAvailable": true
  },
  {
    "ID": 3,
    "CreatedAt": "2023-01-01T00:00:00Z",
    "UpdatedAt": "2023-01-01T00:00:00Z",
    "DeletedAt": null,
    "Size": 1000,
    "IsAvailable": true
  }
]
```

### Add Pack Size

Adds a new pack size.

**Endpoint:** `POST /api/pack-sizes`

**Request:**

```json
{
  "size": 300
}
```

**Response:**

```json
{
  "message": "Pack size added successfully"
}
```

### Update Pack Size

Updates a pack size availability.

**Endpoint:** `PUT /api/pack-sizes/:id`

**Request:**

```json
{
  "isAvailable": false
}
```

**Response:**

```json
{
  "message": "Pack size updated successfully"
}
```

### Delete Pack Size

Deletes a pack size.

**Endpoint:** `DELETE /api/pack-sizes/:id`

**Response:**

```json
{
  "message": "Pack size deleted successfully"
}
```

## Examples

Here are some examples of how the pack calculation works:

| Items ordered | Correct number of packs | Explanation |
|---------------|-------------------------|-------------|
| 1             | 1 x 250                 | Smallest pack that can fulfill the order |
| 250           | 1 x 250                 | Exact match with a single pack |
| 251           | 1 x 500                 | Next pack size up (minimizes excess) |
| 501           | 1 x 500, 1 x 250        | Combination that minimizes excess |
| 12001         | 2 x 5000, 1 x 2000, 1 x 250 | Combination that minimizes excess |

## Flexibility

The application is designed to be flexible:
- Pack sizes can be added, updated, or removed through the API
- The algorithm automatically adapts to the available pack sizes
- No code changes are required when pack sizes change
