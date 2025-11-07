import React, { useState } from "react";
import { fetchHistoricalData } from "../services/api";
import { DataPoint, DATA_TYPES, DATA_TYPE_LABELS, Statistics } from "../types/DataPoint";

type DataTypeKey = typeof DATA_TYPES[keyof typeof DATA_TYPES];

const HistoricalData: React.FC = () => {
  const [data, setData] = useState<DataPoint[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [type, setType] = useState<DataTypeKey>(DATA_TYPES.POWER_INPUT);

  const calculateStatistics = (dataPoints: DataPoint[]): Statistics | null => {
    if (dataPoints.length === 0) return null;

    const values = dataPoints.map((dp) => dp.value).sort((a, b) => a - b);
    const sum = values.reduce((acc, val) => acc + val, 0);
    const min = values[0];
    const max = values[values.length - 1];
    const average = sum / values.length;
    const index95 = Math.floor(values.length * 0.95);
    const percentile95 = values[index95];

    return { min, max, average, percentile95 };
  };

  const handleFetchData = async () => {
    setLoading(true);
    setError(null);

    try {
      const end = new Date();
      const start = new Date(end.getTime() - 60 * 60 * 1000); // 1 hour ago
      
      const result = await fetchHistoricalData(
        type,
        start.toISOString(),
        end.toISOString()
      );
      
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch data");
    } finally {
      setLoading(false);
    }
  };

  const stats = calculateStatistics(data);

  return (
    <div className="historical-data">
      <header>
        <h1>📈 Historical Data</h1>
        <p>Query and analyze data from the last hour</p>
      </header>

      <div className="query-form">
        <div className="form-group">
          <label htmlFor="type">Data Type:</label>
          <select
            id="type"
            value={type}
            onChange={(e) => setType(e.target.value as DataTypeKey)}
          >
            {Object.entries(DATA_TYPE_LABELS).map(([key, label]) => (
              <option key={key} value={key}>
                {label}
              </option>
            ))}
          </select>
        </div>

        <button onClick={handleFetchData} disabled={loading}>
          {loading ? "Loading..." : "Fetch Data"}
        </button>
      </div>

      {error && <div className="error">Error: {error}</div>}

      {stats && (
        <div className="statistics">
          <h2>Statistics</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <div className="stat-label">Minimum</div>
              <div className="stat-value">{stats.min.toFixed(2)}</div>
            </div>
            <div className="stat-card">
              <div className="stat-label">Maximum</div>
              <div className="stat-value">{stats.max.toFixed(2)}</div>
            </div>
            <div className="stat-card">
              <div className="stat-label">Average</div>
              <div className="stat-value">{stats.average.toFixed(2)}</div>
            </div>
            <div className="stat-card">
              <div className="stat-label">95th Percentile</div>
              <div className="stat-value">{stats.percentile95.toFixed(2)}</div>
            </div>
          </div>
        </div>
      )}

      {data.length > 0 && (
        <div className="raw-data">
          <h2>Raw Data ({data.length} points)</h2>
          <div className="data-table-container">
            <table>
              <thead>
                <tr>
                  <th>Timestamp</th>
                  <th>Value</th>
                  <th>Type</th>
                </tr>
              </thead>
              <tbody>
                {data.map((point, index) => (
                  <tr key={index}>
                    <td>{new Date(point.timestamp).toLocaleString()}</td>
                    <td>{point.value.toFixed(2)}</td>
                    <td>{point.type}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default HistoricalData;
