# Data Observability Frontend

This is the frontend component of a tiny data-observability application. It is built using React with Typescript and utilizes Recharts for data visualization.

## Running

To run the frontend, ensure that you have Docker and Docker Compose installed. Then, navigate to the project directory and execute the following command:

```bash
docker compose up --build
```

After that, open your web browser and go to `http://localhost:3000` to access the dashboard.

## Overview

It consists of 3 components:

1. Event producer (Golang, MQTT)
2. Backend (Golang, MongoDB with timeseries collections, Server-Sent Events (SSE))
3. Frontend (React with Typescript, Recharts)

### Features

- Displays a live plot of multiple telemetry signals
- Connects to the backend's SSE endpoint for real-time updates
- Provides a form to query historical data for the last hour, showing raw data and statistical summaries (min, max, average, 95 percentile)
- **Important**: the frontend does not poll the database directly; it relies on the backend to provide data via SSE. Events are pushed from the event producer to the client through the backend, and the client does not make queries to the DB to fetch data for the charts.

## Shortcuts taken

- No authentication or authorization mechanisms are implemented.
- The UI is kept minimalistic without advanced styling or responsiveness.
- Error handling and edge cases are not thoroughly addressed.
- The code is structured for simplicity and clarity rather than scalability or maintainability.

## Reasoning for technology choices:

### Backend

- Golang: Chosen for its performance, simplicity, and strong support for concurrency, making it suitable for handling real-time data streams.
- MongoDB with timeseries collections: Selected for its ability to efficiently store and query time-series data, which is essential for telemetry signals.
- Server-Sent Events (SSE): Used for real-time data updates to the frontend, as it provides a simple and efficient way to push updates from the server to the client with automatic reconnection. Plus, it's built into modern browsers without requiring additional libraries.
- MQTT: A lightweight messaging protocol ideal for IoT and telemetry data, enabling efficient communication between the event producer and backend.

Technically, instead of SSE and MQTT, WebSockets could have used for both pushing the events from the producer to the backend and from the backend to the frontend. However, SSE and MQTT were chosen not because of any technical advantage but to showcase a wider range of technologies. WebSockets would have simplified the architecture by using a single protocol for both communication channels.

### Frontend
- React with Typescript: React is a popular library for building user interfaces, and Typescript adds type safety to JavaScript.
- Recharts: Selected for its simplicity and ease of integration with React for creating charts and visual