# Dataset Catalog API

This is a simple Go-based API server that fetches datasets from Open Data Hub and serves them in different formats: **DCAT**, **ODPS v1.0**, **ODPS v3.0**, and **ODPS v3.1**.

## Prerequisites

Ensure you have Go installed on your machine (Go 1.16+ recommended). You can download it from [golang.org](https://golang.org/dl/).

## Installation & Running the Server

1. Clone this repository:
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```

2 Start the server:
   ```sh
   go run main.go
   ```

The server will run on `http://localhost:8878`.

## Available Endpoints

### 1. DCAT Endpoint
- **URL:** `http://localhost:8878/dcat`
- **Description:** Returns dataset metadata in DCAT format.
- **Optional Query Parameters:**
  - `format=yaml` (returns YAML format instead of JSON)
  - `page=<number>` (fetches a specific page of datasets)

### 2. ODPS v1.0 Endpoint
- **URL:** `http://localhost:8878/odps`
- **Description:** Returns dataset metadata in ODPS v1.0 format.
- **Optional Query Parameters:**
  - `page=<number>` (fetches a specific page of datasets)

### 3. ODPS v3.1 Endpoint
- **URL:** `http://localhost:8878/odps31`
- **Description:** Returns dataset metadata in ODPS v3.1 format.
- **Optional Query Parameters:**
  - `page=<number>` (fetches a specific page of datasets)

### 4. ODPS v3.0 (dev) Endpoint
- **URL:** `http://localhost:8878/odps30`
- **Description:** Returns dataset metadata in ODPS v3.0 (dev) format.
- **Optional Query Parameters:**
  - `page=<number>` (fetches a specific page of datasets)

## Example Requests

- Fetch **DCAT** data in JSON:
  ```sh
  curl http://localhost:8878/dcat
  ```

- Fetch **DCAT** data in YAML:
  ```sh
  curl "http://localhost:8878/dcat?format=yaml"
  ```

- Fetch **ODPS v1.0** data:
  ```sh
  curl http://localhost:8878/odps
  ```

- Fetch **ODPS v3.1** data:
  ```sh
  curl http://localhost:8878/odps31
  ```

## License
This project is licensed under the [GNU General Public License v3.0](LICENSE).
