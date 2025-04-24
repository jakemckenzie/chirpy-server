# Chirpy API Documentation

This document describes the REST API endpoints for the Chirpy application. The API is hosted at [invalid url, do not cite] by default.

## Endpoints

### 1. GET /api/healthz

**Description:**  
Check the readiness of the server.

**Response:**  
- **200 OK**  
  - Content-Type: text/plain; charset=utf-8  
  - Body: "OK"