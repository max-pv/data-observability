# Fourier Frontend

A React-based data observability dashboard for real-time telemetry monitoring.

## Features

- **Live Dashboard**: Real-time visualization of 5 telemetry streams using Server-Sent Events (SSE)
  - Power Input (kW)
  - Water Flow Rate (L/min)
  - Temperature (°C)
  - Hydrogen Production Rate (kg/h)
  - Efficiency (%)

- **Historical Data Viewer**: Query and analyze historical data
  - Fetch data from the last hour by type
  - View statistical summaries (min, max, average, 95th percentile)
  - Browse raw data in a table format

## Architecture

- **SSE Connection**: A single EventSource connection receives all data types and routes them to the appropriate charts
- **Fetch API**: Used for historical data queries (no external HTTP libraries)
- **Minimal CSS**: Clean, dark-themed interface with no CSS frameworks

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- Backend server running on `http://localhost:8080`

### Installation

```bash
npm install
```

### Running the App

```bash
npm start
```

The app will open at `http://localhost:3000`.

### Building for Production

```bash
npm build
```

## Project Structure

```
src/
├── components/
│   └── Chart.tsx          # Reusable chart component
├── pages/
│   ├── Dashboard.tsx      # Live data dashboard
│   └── HistoricalData.tsx # Historical data viewer
├── services/
│   ├── api.ts             # API service using fetch
│   └── sse.ts             # SSE connection handler
├── types/
│   └── DataPoint.ts       # TypeScript type definitions
├── App.tsx                # Main app with routing
├── App.css                # Global styles
└── index.tsx              # Entry point
```

## API Endpoints

- **SSE Endpoint**: `GET /events` - Receives live telemetry data
- **Historical Data**: `GET /historical?type={type}&start={ISO8601}&end={ISO8601}` - Queries historical data

## Technologies

- React 18
- TypeScript
- React Router 6
- Recharts (for data visualization)
- Native Fetch API
- EventSource (SSE)
