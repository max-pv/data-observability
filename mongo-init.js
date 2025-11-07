print("Starting mongo-init.js script...");

db = db.getSiblingDB("fourier"); // Replace with your database name

// timeseries helps optimize storage and queries for time-based data
db.createCollection("telemetry", {
  timeseries: {
    timeField: "timestamp",
    granularity: "seconds"
  }
});

print("Creating indexes...");

// we will be storing telemetry data with fields: timestamp, value, type
// where type is a string indicating the type of telemetry
db.telemetry.createIndex({ type: 1 });
db.telemetry.createIndex({ type: 1, timestamp: 1 });

print("Finished mongo-init.js script.");