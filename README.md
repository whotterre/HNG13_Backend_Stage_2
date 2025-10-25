# Country Currency & Exchange API

A RESTful API that fetches country data from external APIs, stores it in a MySQL database, and provides CRUD operations with exchange rate calculations and GDP estimations.

## Overview

This API integrates with external services to:
- Fetch country information from REST Countries API
- Retrieve real-time exchange rates from Open Exchange Rates API
- Calculate estimated GDP based on population and exchange rates
- Generate summary images with country statistics
- Provide filtered and sorted country data

##  Features

- **Data Synchronization**: Fetch and cache country data with exchange rates
- **CRUD Operations**: Create, Read, Update, and Delete country records
- **Advanced Filtering**: Query by region, currency, and sort by GDP
- **Image Generation**: Automatic generation of summary statistics images
- **Validation**: Comprehensive input validation with detailed error responses
- **Error Handling**: Graceful handling of external API failures

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0
- **Image Processing**: golang.org/x/image
- **Containerization**: Docker & Docker Compose

## 📦 Dependencies

```go
github.com/gin-gonic/gin          // Web framework
gorm.io/gorm                      // ORM
gorm.io/driver/mysql              // MySQL driver
github.com/joho/godotenv          // Environment variables
golang.org/x/image                // Image generation
```

##  Project Structure

```
.
├── cmd/
│   └── main.go                   # Application entry point
├── config/
│   └── config.go                 # Configuration management
├── dto/
│   └── dto.go                    # Data Transfer Objects
├── handlers/
│   └── handlers.go               # HTTP request handlers
├── models/
│   └── country.go                # Database models
├── repository/
│   └── country.go                # Database operations
├── services/
│   └── services.go               # Business logic
├── routes/
│   └── routes.go                 # Route definitions
├── utils/
│   ├── currency.go               # GDP calculations
│   └── image.go                  # Image generation
├── clients/
│   └── clients.go                # External API clients
├── initializer/
│   └── connectDB.go              # Database initialization
├── cache/                        # Generated images
├── .env.example                  # Example environment variables
├── Dockerfile                    # Docker configuration
├── docker-compose.yml            # Docker Compose setup
└── README.md                     # This file
```

## 🔧 Installation & Setup

### Prerequisites

- Go 1.25+ installed
- MySQL 8.0+ installed (or use Docker)
- Docker & Docker Compose (optional but recommended)

### Local Setup (Without Docker)

1. **Clone the repository**
   ```bash
   git clone https://github.com/whotterre/HNG13_Backend_Stage_2.git
   cd HNG13_Backend_Stage_2
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   ```env
   PORT=8080
   MYSQL_HOST=localhost
   MYSQL_PORT=3306
   MYSQL_USER=root
   MYSQL_PASSWORD=yourpassword
   MYSQL_DATABASE=countries_db
   DB_STRING=root:yourpassword@tcp(localhost:3306)/countries_db?charset=utf8mb4&parseTime=True&loc=Local
   ```

4. **Create MySQL database**
   ```sql
   CREATE DATABASE countries_db;
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

The API will be available at `http://localhost:8080`

### Docker Setup (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/whotterre/HNG13_Backend_Stage_2.git
   cd HNG13_Backend_Stage_2
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Start the services**
   ```bash
   docker-compose up -d
   ```

The API will be available at `http://localhost:8080`

4. **Stop the services**
   ```bash
   docker-compose down
   ```

## 🌐 API Endpoints

### 1. Refresh Countries Data
**POST** `/countries/refresh`

Fetches country data and exchange rates from external APIs, then stores/updates in the database.

**Response (200 OK):**
```json
{
  "status": "Successfully refreshed countries"
}
```

**Error (503 Service Unavailable):**
```json
{
  "error": "External data source unavailable",
  "details": "failed to fetch country data from external API"
}
```

---

### 2. Get All Countries
**GET** `/countries`

Retrieve all countries with optional filtering and sorting.

**Query Parameters:**
- `region` - Filter by region (e.g., `Africa`, `Europe`)
- `currency` - Filter by currency code (e.g., `NGN`, `USD`)
- `sort` - Sort order: `gdp_desc` or `gdp_asc`

**Examples:**
```
GET /countries?region=Africa
GET /countries?currency=NGN
GET /countries?sort=gdp_desc
GET /countries?region=Africa&sort=gdp_desc
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Nigeria",
    "capital": "Abuja",
    "region": "Africa",
    "population": 206139589,
    "currency_code": "NGN",
    "exchange_rate": 1600.23,
    "estimated_gdp": 25767448125.2,
    "flag_url": "https://flagcdn.com/ng.svg",
    "last_refreshed_at": "2025-10-25T18:00:00Z"
  }
]
```

---

### 3. Get Country by Name
**GET** `/countries/:name`

Retrieve a specific country by name (case-insensitive).

**Example:**
```
GET /countries/Nigeria
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Nigeria",
  "capital": "Abuja",
  "region": "Africa",
  "population": 206139589,
  "currency_code": "NGN",
  "exchange_rate": 1600.23,
  "estimated_gdp": 25767448125.2,
  "flag_url": "https://flagcdn.com/ng.svg",
  "last_refreshed_at": "2025-10-25T18:00:00Z"
}
```

**Error (404 Not Found):**
```json
{
  "error": "Country not found"
}
```

---

### 4. Delete Country
**DELETE** `/countries/:name`

Delete a country record by name.

**Example:**
```
DELETE /countries/Nigeria
```

**Response:** `204 No Content`

---

### 5. Get Statistics
**GET** `/status`

Get total countries count and last refresh timestamp.

**Response (200 OK):**
```json
{
  "total_countries": 250,
  "last_refreshed_at": "2025-10-25T18:00:00Z"
}
```

---

### 6. Get Summary Image
**GET** `/countries/image`

Retrieve the auto-generated summary image containing:
- Total number of countries
- Top 5 countries by estimated GDP
- Last refresh timestamp

**Response:** PNG image file

**Error (404 Not Found):**
```json
{
  "error": "Summary image not found. Please refresh countries first."
}
```

---

##  Data Model

### Country Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | uint | Auto | Primary key |
| `name` | string | Yes | Country name |
| `capital` | string | No | Capital city |
| `region` | string | No | Geographic region |
| `population` | int64 | Yes | Population count |
| `currency_code` | string | Conditional | ISO currency code (required if currencies array is not empty) |
| `exchange_rate` | float64 | No | Exchange rate to USD |
| `estimated_gdp` | float64 | Computed | `population × random(1000–2000) ÷ exchange_rate` |
| `flag_url` | string | No | Country flag URL |
| `last_refreshed_at` | timestamp | Auto | ISO 8601 timestamp |
| `created_at` | timestamp | Auto | Record creation time |
| `updated_at` | timestamp | Auto | Last update time |

##  Validation Rules

- `name` - Required, non-empty string
- `population` - Required, must be ≥ 0
- `currency_code` - Required only if the country has currencies in the external API response

**Validation Error Response (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "details": {
    "currency_code": "is required",
    "population": "must be non-negative"
  }
}
```

##  External APIs

1. **REST Countries API**
   - URL: `https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies`
   - Purpose: Fetch country information

2. **Open Exchange Rates API**
   - URL: `https://open.er-api.com/v6/latest/USD`
   - Purpose: Fetch real-time exchange rates

## 🔄 Refresh Behavior

### Currency Handling
- **Multiple currencies**: Only the first currency code is stored
- **Empty currencies array**: 
  - `currency_code` → `null`
  - `exchange_rate` → `null`
  - `estimated_gdp` → `0`
- **Currency not in exchange rates API**:
  - `exchange_rate` → `null`
  - `estimated_gdp` → `null`

### Update vs Insert Logic
- Countries are matched by **name** (case-insensitive)
- **Existing country**: All fields updated, including new `estimated_gdp` with fresh random multiplier
- **New country**: Inserted with validation
- **Random multiplier**: Generated fresh (1000-2000) for each country on every refresh

### Image Generation
After successful refresh:
1. Queries database for total countries and top 5 by GDP
2. Generates PNG image at `cache/summary.png`
3. Includes timestamp, total count, and top countries
4. Accessible via `GET /countries/image`

##  Error Handling

| Status Code | Response |
|-------------|----------|
| 400 | `{ "error": "Validation failed", "details": {...} }` |
| 404 | `{ "error": "Country not found" }` |
| 500 | `{ "error": "Internal server error", "details": "..." }` |
| 503 | `{ "error": "External data source unavailable", "details": "..." }` |

## Testing

### Manual Testing

1. **Refresh data:**
   ```bash
   curl -X POST http://localhost:8080/countries/refresh
   ```

2. **Get all countries:**
   ```bash
   curl http://localhost:8080/countries
   ```

3. **Filter by region:**
   ```bash
   curl http://localhost:8080/countries?region=Africa
   ```

4. **Get statistics:**
   ```bash
   curl http://localhost:8080/status
   ```

5. **Get summary image:**
   ```bash
   curl http://localhost:8080/countries/image --output summary.png
   ```

##  Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=yourpassword
MYSQL_DATABASE=countries_db
DB_STRING=root:yourpassword@tcp(localhost:3306)/countries_db?charset=utf8mb4&parseTime=True&loc=Local
```

## 🐳 Docker Commands

```bash
# Build and start containers
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop containers
docker-compose down

# Rebuild after code changes
docker-compose up -d --build
```

## 📝 Development Notes

- Database connection includes automatic retry with exponential backoff
- GORM AutoMigrate creates/updates tables on startup
- All timestamps follow ISO 8601 format (RFC3339)
- Image generation uses Go's standard `image` package
- Transaction-based refresh ensures atomicity

## Deployment

The application can be deployed to:
- Railway
- Heroku
- AWS (EC2, ECS, Lambda)
- DigitalOcean App Platform
- Any platform supporting Docker or Go applications

##  Acknowledgments

- [HNG Internship](https://hng.tech) for the challenge
- [REST Countries API](https://restcountries.com)
- [Open Exchange Rates API](https://open.er-api.com)

---
