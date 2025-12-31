# Hillview Links API Documentation

## Overview

The Hillview Links API is a URL shortening and redirection service built with Go. It provides endpoints to look up short link routes and redirect users to their destinations while optionally tracking click analytics.

**Base URL:** `/links/v1.1`

**Version:** 1.1

---

## Table of Contents

- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Check Link Route](#check-link-route)
- [Data Models](#data-models)
- [Environment Variables](#environment-variables)
- [CORS Configuration](#cors-configuration)

---

## Authentication

The API supports JWT Bearer token authentication. While some endpoints may not require authentication, the API includes middleware that processes authorization tokens when provided.

### Authorization Header

```
Authorization: Bearer <jwt_token>
```

When a valid JWT token is provided, the API will:

- Log the request with the associated user ID
- Track user activity

### Token Types

| Type           | Description                                      |
| -------------- | ------------------------------------------------ |
| `access_token` | Standard access token for authenticated requests |

---

## Response Format

### Success Response

All successful responses follow this structure:

```json
{
  "success": true,
  "message": "request was successful",
  "data": { ... }
}
```

| Field     | Type                      | Description                           |
| --------- | ------------------------- | ------------------------------------- |
| `success` | `boolean`                 | Always `true` for successful requests |
| `message` | `string`                  | Human-readable success message        |
| `data`    | `object \| array \| null` | The response payload                  |

### Error Response

All error responses follow this structure:

```json
{
  "error": "error details or null",
  "error_message": "human readable error message",
  "error_code": 1000
}
```

| Field           | Type             | Description                              |
| --------------- | ---------------- | ---------------------------------------- |
| `error`         | `string \| null` | Detailed error information (may be null) |
| `error_message` | `string`         | Human-readable error description         |
| `error_code`    | `integer`        | Application-specific error code          |

---

## Error Handling

### HTTP Status Codes

| Status Code | Description                                 |
| ----------- | ------------------------------------------- |
| `200`       | Success                                     |
| `400`       | Bad Request - Missing or invalid parameters |
| `401`       | Unauthorized - Invalid or expired token     |
| `404`       | Not Found - Resource does not exist         |
| `409`       | Conflict - Resource conflict                |
| `500`       | Internal Server Error                       |

### Common Error Types

| Error Type              | Status | Description                                            |
| ----------------------- | ------ | ------------------------------------------------------ |
| Missing Route Parameter | `400`  | The `route` parameter is required but was not provided |
| Route Not Found         | `404`  | The specified route does not exist or is inactive      |
| Internal Error          | `500`  | Server-side error occurred during processing           |

---

## Endpoints

### Health Check

Check if the API service is running and healthy.

**Endpoint:** `GET /healthcheck`

> **Note:** This endpoint is at the root level, not under the `/links/v1.1` prefix.

#### Request

```http
GET /healthcheck HTTP/1.1
Host: api.example.com
```

#### Response

**Success (200 OK)**

Returns an empty response with status code `200`.

#### Example

```bash
curl -X GET https://api.example.com/healthcheck
```

---

### Check Link Route

Look up a short link route and retrieve its destination URL. Optionally record a click for analytics.

**Endpoint:** `GET /links/v1.1/check/{route}`

#### Path Parameters

| Parameter | Type     | Required | Description                                |
| --------- | -------- | -------- | ------------------------------------------ |
| `route`   | `string` | Yes      | The short link route identifier to look up |

#### Query Parameters

| Parameter     | Type     | Required | Default | Description                                     |
| ------------- | -------- | -------- | ------- | ----------------------------------------------- |
| `recordClick` | `string` | No       | `false` | Set to `"true"` to record a click for analytics |

#### Request

```http
GET /links/v1.1/check/my-short-link?recordClick=true HTTP/1.1
Host: api.example.com
Content-Type: application/json
```

#### Response

**Success (200 OK)**

```json
{
  "success": true,
  "message": "successfully found route",
  "data": {
    "id": 1,
    "route": "my-short-link",
    "destination": "https://example.com/very/long/url/path",
    "created_by": 123,
    "active": true,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Not Found (404)**

Returns an empty response with status code `404` when the route doesn't exist or is inactive.

**Bad Request (400)**

```json
{
  "error": null,
  "error_message": "missing required body field: route",
  "error_code": 1000
}
```

**Internal Server Error (500)**

```json
{
  "error": "database connection failed",
  "error_message": "failed to lookup route",
  "error_code": 1000
}
```

#### Examples

**Basic lookup:**

```bash
curl -X GET "https://api.example.com/links/v1.1/check/my-link"
```

**Lookup with click tracking:**

```bash
curl -X GET "https://api.example.com/links/v1.1/check/my-link?recordClick=true"
```

**With authentication:**

```bash
curl -X GET "https://api.example.com/links/v1.1/check/my-link" \
  -H "Authorization: Bearer <your_jwt_token>"
```

---

## Data Models

### Route

Represents a short link route configuration.

```json
{
  "id": 1,
  "route": "my-short-link",
  "destination": "https://example.com/destination",
  "created_by": 123,
  "active": true,
  "created_at": "2024-01-15T10:30:00Z"
}
```

| Field         | Type                | Description                                |
| ------------- | ------------------- | ------------------------------------------ |
| `id`          | `integer`           | Unique identifier for the route            |
| `route`       | `string`            | The short link identifier (slug)           |
| `destination` | `string`            | The target URL to redirect to              |
| `created_by`  | `integer`           | User ID of the creator                     |
| `active`      | `boolean`           | Whether the route is active and accessible |
| `created_at`  | `string (ISO 8601)` | Timestamp when the route was created       |

### User

Represents a user in the system (used internally for authentication).

```json
{
  "id": 1,
  "username": "johndoe",
  "name": "John Doe",
  "email": "john@example.com",
  "profile_image_url": "https://example.com/avatar.jpg",
  "authentication": {
    "id": 1,
    "name": "Standard",
    "short_name": "std"
  },
  "inserted_at": "2024-01-01T00:00:00Z",
  "last_active": "2024-01-15T12:00:00Z"
}
```

| Field               | Type                        | Description                     |
| ------------------- | --------------------------- | ------------------------------- |
| `id`                | `integer`                   | Unique user identifier          |
| `username`          | `string \| null`            | Optional username               |
| `name`              | `string`                    | User's display name             |
| `email`             | `string`                    | User's email address            |
| `profile_image_url` | `string`                    | URL to user's profile image     |
| `authentication`    | `object`                    | Authentication type information |
| `inserted_at`       | `string (ISO 8601)`         | Account creation timestamp      |
| `last_active`       | `string (ISO 8601) \| null` | Last activity timestamp         |

---

## Environment Variables

The API requires the following environment variables to be configured:

| Variable           | Required | Default | Description                           |
| ------------------ | -------- | ------- | ------------------------------------- |
| `PORT`             | No       | `8000`  | Port number for the API server        |
| `DATABASE_DSN`     | Yes      | -       | Database connection string            |
| `JWT_SIGNING_KEY`  | Yes      | -       | Secret key for JWT token validation   |
| `HEALTH_CHECK_URL` | Yes      | -       | URL for external health check polling |

### Example `.env` File

```env
PORT=8000
DATABASE_DSN=user:password@tcp(localhost:3306)/hillview_links
JWT_SIGNING_KEY=your-secret-signing-key
HEALTH_CHECK_URL=https://your-health-check-endpoint.com/ping
```

---

## CORS Configuration

The API is configured with the following CORS settings:

### Allowed Origins

```
* (all origins)
```

### Allowed Headers

- `X-Requested-With`
- `Content-Type`
- `Origin`
- `Authorization`
- `Accept`
- `X-CSRF-Token`
- `Accept-Encoding`
- `Connection`
- `Content-Length`

### Allowed Methods

- `GET`
- `HEAD`
- `POST`
- `PUT`
- `OPTIONS`
- `DELETE`

---

## Middleware

The API applies the following middleware to all `/links/v1.1/*` routes:

### 1. Logging Middleware

Logs all incoming requests with method and URI.

**Log Format:**

```
<METHOD> <REQUEST_URI>
```

### 2. Header Middleware

Automatically sets response headers:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Credentials: true`
- `Content-Type: application/json`
- `Server: Go`

### 3. Token Handler

When an `Authorization` header with a valid Bearer token is provided:

- Parses and validates the JWT token
- Logs the request with the user ID for analytics
- Tracks user activity in the request log

---

## Rate Limiting

Currently, the API does not implement rate limiting. Consider implementing rate limiting for production deployments.

---

## Database Schema

### `links` Table

| Column        | Type        | Description                  |
| ------------- | ----------- | ---------------------------- |
| `id`          | `INT`       | Primary key, auto-increment  |
| `route`       | `VARCHAR`   | Unique short link identifier |
| `destination` | `VARCHAR`   | Target URL                   |
| `active`      | `BOOLEAN`   | Whether the link is active   |
| `created_by`  | `INT`       | Foreign key to users table   |
| `created_at`  | `TIMESTAMP` | Creation timestamp           |

### `link_clicks` Table

| Column       | Type        | Description                 |
| ------------ | ----------- | --------------------------- |
| `id`         | `INT`       | Primary key, auto-increment |
| `link_id`    | `INT`       | Foreign key to links table  |
| `created_at` | `TIMESTAMP` | Click timestamp             |

---

## Background Services

### Health Check Polling

The API runs a background service that polls an external health check endpoint every minute. This can be used for uptime monitoring or keep-alive functionality.

**Interval:** 1 minute  
**Timeout:** 10 seconds per request

---

## Example Integration

### JavaScript/TypeScript

```typescript
interface Route {
  id: number;
  route: string;
  destination: string;
  created_by: number;
  active: boolean;
  created_at: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

async function checkRoute(
  route: string,
  recordClick: boolean = false
): Promise<Route | null> {
  const url = `https://api.example.com/links/v1.1/check/${route}?recordClick=${recordClick}`;

  const response = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    throw new Error("Failed to check route");
  }

  const data: ApiResponse<Route> = await response.json();
  return data.data;
}

// Usage
const route = await checkRoute("my-short-link", true);
if (route) {
  window.location.href = route.destination;
}
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Route struct {
    ID          int    `json:"id"`
    Route       string `json:"route"`
    Destination string `json:"destination"`
    CreatedBy   int    `json:"created_by"`
    Active      bool   `json:"active"`
    CreatedAt   string `json:"created_at"`
}

type ApiResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    Route  `json:"data"`
}

func checkRoute(route string, recordClick bool) (*Route, error) {
    url := fmt.Sprintf("https://api.example.com/links/v1.1/check/%s?recordClick=%t", route, recordClick)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil, nil
    }

    var apiResp ApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }

    return &apiResp.Data, nil
}
```

### cURL

```bash
# Check a route
curl -X GET "https://api.example.com/links/v1.1/check/my-link"

# Check a route and record click
curl -X GET "https://api.example.com/links/v1.1/check/my-link?recordClick=true"

# With authentication
curl -X GET "https://api.example.com/links/v1.1/check/my-link" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Changelog

### v1.1 (Current)

- Initial documented version
- Link route lookup endpoint
- Click tracking analytics
- JWT authentication support
- Health check polling background service
