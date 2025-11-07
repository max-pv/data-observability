import { DataPoint } from "../types/DataPoint";

const API_BASE_URL = "http://localhost:8080";

export const fetchHistoricalData = async (
  type: string,
  start: string,
  end: string
): Promise<DataPoint[]> => {
  const url = new URL(`${API_BASE_URL}/historical`);
  url.searchParams.append("type", type);
  url.searchParams.append("start", start);
  url.searchParams.append("end", end);

  const response = await fetch(url.toString());
  
  if (!response.ok) {
    throw new Error(`Failed to fetch historical data: ${response.statusText}`);
  }

  return response.json();
};
