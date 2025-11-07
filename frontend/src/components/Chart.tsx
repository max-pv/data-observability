import React from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts";
import { DataPoint } from "../types/DataPoint";

interface ChartProps {
  title: string;
  data: DataPoint[];
  color: string;
}

const Chart: React.FC<ChartProps> = ({ title, data, color }) => {
  // Format data for Recharts
  const chartData = data.map((point) => ({
    time: new Date(point.timestamp).toLocaleTimeString(),
    value: point.value,
  }));

  return (
    <div className="chart-container">
      <h3>{title}</h3>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart 
        margin={{ top: 5, right: 20, left: -25, bottom: 5 }}
        data={chartData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#333" />
          <XAxis 
            dataKey="time" 
            stroke="#888"
            tick={{ fontSize: 12 }}
          />
          <YAxis 
            stroke="#888"
            tick={{ fontSize: 12 }}
          />
          <Tooltip 
            contentStyle={{ 
              backgroundColor: "#1a1a1a", 
              border: "1px solid #333",
              borderRadius: "4px",
              padding: "8px 8px 5px 8px"
            }}
          />
          <Line 
            type="monotone" 
            dataKey="value" 
            stroke={color} 
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default Chart;
