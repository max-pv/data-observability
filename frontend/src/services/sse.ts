import { DataPoint, SSEPayload } from "../types/DataPoint";

const sseEndpoint = "http://localhost:8080/events";

export const connectToSSE = (onMessage: (data: DataPoint) => void) => {
  const eventSource = new EventSource(sseEndpoint);

  eventSource.addEventListener("message", (event) => {
    try {
      const data: SSEPayload = JSON.parse(event.data);
      const initialData = data.payload as DataPoint[];
      initialData.forEach((dp) => onMessage(dp));
    } catch (error) {
      console.error("Error parsing SSE message:", error);
      return;
    }
  })

  eventSource.addEventListener("open", (ev) => {
    console.log("SSE connection opened:", ev);
  })

  return eventSource;
}