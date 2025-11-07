export interface DataPoint {
  timestamp: string; // ISO 8601 format
  value: number;
  type: string;
}

export interface SSEPayload {
  kind: string;
  payload: DataPoint[];
}

export interface Statistics {
  min: number;
  max: number;
  average: number;
  percentile95: number;
}

export const DATA_TYPES = {
  POWER_INPUT: "PowerInput",
  WATER_FLOW_RATE: "WaterFlowRate",
  TEMPERATURE: "Temperature",
  HYDROGEN_PRODUCTION_RATE: "HydrogenProductionRate",
  EFFICIENCY: "Efficiency",
} as const;

export const DATA_TYPE_LABELS: Record<string, string> = {
  [DATA_TYPES.POWER_INPUT]: "Power Input (kW)",
  [DATA_TYPES.WATER_FLOW_RATE]: "Water Flow Rate (L/min)",
  [DATA_TYPES.TEMPERATURE]: "Temperature (°C)",
  [DATA_TYPES.HYDROGEN_PRODUCTION_RATE]: "Hydrogen Production Rate (kg/h)",
  [DATA_TYPES.EFFICIENCY]: "Efficiency (%)",
};

export const DATA_TYPE_COLORS: Record<string, string> = {
  [DATA_TYPES.POWER_INPUT]: "#8884d8",
  [DATA_TYPES.WATER_FLOW_RATE]: "#82ca9d",
  [DATA_TYPES.TEMPERATURE]: "#ffc658",
  [DATA_TYPES.HYDROGEN_PRODUCTION_RATE]: "#ff7c7c",
  [DATA_TYPES.EFFICIENCY]: "#8dd1e1",
};
