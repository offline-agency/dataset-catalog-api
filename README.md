<!--
SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>

SPDX-License-Identifier: CC0-1.0
-->

# Dataset Catalog API

This is a simple Go-based API server that fetches datasets from Open Data Hub and serves them in different formats: **DCAT**, **ODPS v1.0**, **ODPS v3.0 (dev)**, and **ODPS v3.1**.

> **Note:** Pagination always starts at page 1. A request with `page=0` or any page number greater than the total number of pages will return a "No data found" response.

## Prerequisites

Ensure you have Go installed on your machine (Go 1.16+ recommended). You can download it from [golang.org](https://golang.org/dl/).

## Installation & Running the Server

1. Clone this repository:
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Start the server:
   ```sh
   cd src
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

### 3. ODPS v3.1 Endpoints
- **Listing Endpoint**
  - **URL:** `http://localhost:8878/odps31`
  - **Description:** Returns a paginated list of dataset endpoints in ODPS v3.1 format.
  - **Optional Query Parameters:**
    - `page=<number>` (fetches a specific page of datasets)
  - **Pagination Details:**  
    The response includes `current_page`, `total_pages`, and `totalRecord` fields so you can verify the complete dataset list and navigate through pages.
- **Detail Endpoint**
  - **URL:** `http://localhost:8878/odps31/{uuid}`
  - **Description:** Returns detailed information for a specific dataset in ODPS v3.1 format.
  - **Path Parameter:**
    - `{uuid}` – The unique identifier of the dataset.

### 4. ODPS v3.0 (dev) Endpoints
- **Listing Endpoint**
  - **URL:** `http://localhost:8878/odps30`
  - **Description:** Returns a paginated list of dataset endpoints in ODPS v3.0 (dev) format.
  - **Optional Query Parameters:**
    - `page=<number>` (fetches a specific page of datasets)
  - **Pagination Details:**  
    Similar to ODPS v3.1, the response includes `current_page`, `total_pages`, and `totalRecord` fields.
- **Detail Endpoint**
  - **URL:** `http://localhost:8878/odps30/{uuid}`
  - **Description:** Returns detailed information for a specific dataset in ODPS v3.0 (dev) format.
  - **Path Parameter:**
    - `{uuid}` – The unique identifier of the dataset.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).
```
