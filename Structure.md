```
The Challenge: Build a tiny data-observability app
A mock telemetry producer that emits time-series values for multiple signals
A backend that ingests, stores and serves that data
A React dashboard that shows the live data

The Details
1. Producer
Simulates a handful of signals (sine waves, random noise, etc)
Streams values at a steady rate
 
2. Backend
Accepts data from the producer
Stores data in any database of your choice
Provides an API for querying new values in real time
 
3. Frontend
Written in React
Displays a plot of the last hour of all of the signals
Updates the plot live as new data arrives
```

# Project Structure

## Technologies to use

Producer: Golang, MQTT
Backend: Golang, MongoDB with timeseries collections, Server-Sent Events (SSE). Extra endpoint for historical data queries.
Frontend: React with Typescript, Recharts

We will strart with a Docker Compose file. That will contain 3 services: producer, backend, and frontend. There will be no authenication or authorization in this app for the sake of simplicity. It will also contain a MongoDB service for the backend to use as its database.

## Producer

The producer will run multiple goroutines, each simulating a different signal. Each goroutine will use the publisher client to send data to the backend. Here is the telemetry we will focus on:

1. Power Input (kW)
2. Water Flow Rate (L/min)
3. Temperature (C)
4. Hydrogen Production Rate (kg/h)
5. Efficiency (%)

## Backend

The backend will have the following endpoints:
1. MQTT endpoint to receive data from the producer
2. SSE endpoint to stream live data to the frontend
3. REST endpoint to query historical data

The backend will use MongoDB's time-series collections to store the telemetry data efficiently.

## Frontend

The frontend will have a dashboard that displays the live data using Recharts. It will connect to the backend's SSE endpoint to receive real-time updates and will also have a form to query historical data for the last hour, show it in the raw format, and provide statistical summaries (min, max, average, 95 percentile).