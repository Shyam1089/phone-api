# Phone Number Lookup API

A Go REST API built with Gin that validates and parses phone numbers according to E.164 format specifications.


## 📁 Project Structure
```
phone-api/
├── api/                  # Core API package
│   ├── handlers.go       # HTTP handlers
│   └── validator.go      # Phone validation logic
├── cmd/api/              # Main application
│   └── main.go           # Application entry point
├── tests/                # Test suite
│   ├── handlers_test.go  # API endpoint tests
│   └── validator_test.go # Validation logic tests
├── deploy/               # Docker configuration
├── go.mod               # Go module dependencies
└── Makefile             # Development commands
```

  

## 🚀 Quick Start

  
### Prerequisites

- Go 1.21+ (optional - only for local development)
- Docker (for development and testing)

  

### Local Development


### Using Docker

```bash
make app
```


**Test the API:**

```bash

curl  "http://localhost:8000/v1/phone-numbers/?phoneNumber=%2B12125690123"

```

  

## 📋 API Usage

  
### Endpoints


-  `GET /health/` - Health check

-  `GET /v1/phone-numbers/` - Phone number lookup


### Parameters

-  `phoneNumber` (required): Phone number in E.164 format

-  `countryCode` (optional): ISO 3166-1 alpha-2 country code

  

### Examples


```bash

# Valid requests

curl  "http://localhost:8000/v1/phone-numbers/?phoneNumber=%2B12125690123"

curl  "http://localhost:8000/v1/phone-numbers/?phoneNumber=%2B52%20631%203118150"

curl  "http://localhost:8000/v1/phone-numbers/?phoneNumber=2125690123&countryCode=US"

```

  

**Success Response:**

```json

{

"phoneNumber": "+12125690123",

"countryCode": "US",

"areaCode": "212",

"localPhoneNumber": "5690123"

}

```

  

**Error Response:**

```json

{

"phoneNumber": "25690123",

"error": {

"countryCode": "required value is missing"

}

}

```

  

## 🧪 Testing

```bash
make test
```

**Features:**
- 21 comprehensive tests covering all API endpoints and validation logic
- 88.7% code coverage with detailed reporting
- Function-by-function coverage breakdown
- Coverage reports show missing lines for easy improvement

*Note: Tests run inside Docker containers, no local Go setup required.*

  

## 📝 Available Commands

  

```bash
make app   # Build and run with Docker
make test  # Run all tests in Docker
```

  

## ✅ Validation Rules

  

- E.164 format: `[+][country code][area code][local phone number]`

-  `+` sign is optional

- Only digits and spaces allowed

- Spaces allowed between country, area code, and local number

- 4 space-separated parts are invalid (e.g., `351 21 094 2000`)

- Invalid characters rejected (letters, hyphens, etc.)

  

## 🌍 Supported Countries

  

US, CA, MX, ES, PT, GB, FR, DE, IT, BR

  

## 🛠️ Technology Choices

**Go + Gin Framework**

- **Go**: High performance, excellent concurrency, strong typing, fast compilation
- **Gin Framework**: Lightweight, fast HTTP router with middleware support
- **Minimal Dependencies**: Only essential packages (gin, cors, testify for testing)

  

## 🚀 Production Deployment

  

**Recommended: Docker + Cloud Platform**

  

1.  **Build for production:**

```bash

docker build -t phone-api .

```


2.  **Production settings:**

- Set `GIN_MODE=release` environment variable
- Configure appropriate `PORT` (defaults to 8000)
- Use `/health` endpoint for health checks
- Add SSL at load balancer level
- Set resource limits in production containers

  

##  Assumptions Made

  

1.  **Phone Format**: E.164 standard with spacing rules from provided examples

2.  **Countries**: Limited to 10 common countries for scope

3.  **Length Limits**: Applied realistic phone number lengths per country

4.  **No Database**: Stateless validation service, no data storage needed

5.  **Public API**: No authentication required for phone validation

6.  **Error Format**: Matched exact JSON structure from requirements

7.  **Input Rules**: Only digits, spaces, and + symbols allowed

  

## 📈 Potential Improvements

- Add rate limiting to prevent abuse

- Add authentication if deployed to PROD 

- Redis caching for performance

- API metrics and monitoring

- Swagger documentation

- Batch validation (multiple numbers at once)