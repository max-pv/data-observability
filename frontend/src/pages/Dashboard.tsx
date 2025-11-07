import React, { useEffect, useState } from "react";
import { connectToSSE } from "../services/sse";
import { DataPoint, DATA_TYPES, DATA_TYPE_LABELS, DATA_TYPE_COLORS } from "../types/DataPoint";
import Chart from "../components/Chart";

const MAX_DATA_POINTS = 60; // Keep last 60 data points per chart

const Dashboard: React.FC = () => {
  const [dataByType, setDataByType] = useState<Record<string, DataPoint[]>>({
    [DATA_TYPES.POWER_INPUT]: [],
    [DATA_TYPES.WATER_FLOW_RATE]: [],
    [DATA_TYPES.TEMPERATURE]: [],
    [DATA_TYPES.HYDROGEN_PRODUCTION_RATE]: [],
    [DATA_TYPES.EFFICIENCY]: [],
  });

  useEffect(() => {
    const cb = function (dataPoint: DataPoint) {
      setDataByType((prevData) => {
        const typeData = prevData[dataPoint.type] || [];
        const newData = [...typeData, dataPoint].slice(-MAX_DATA_POINTS);

        return {
          ...prevData,
          [dataPoint.type]: newData,
        };
      });
    };

    const eventSource = connectToSSE(cb);

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="dashboard">
      <header>
        <h1>📊 Live Telemetry Dashboard</h1>
        <p>Real-time data from the hydrogen production system</p>
      </header>

      <div className="charts-grid">
        {Object.values(DATA_TYPES).map((type) => (
          <Chart
            key={type}
            title={DATA_TYPE_LABELS[type]}
            data={dataByType[type]}
            color={DATA_TYPE_COLORS[type]}
          />
        ))}
      </div>
    </div>
  );
};

export default Dashboard;
