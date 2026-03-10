import { useAppStore } from '../stores/useAppStore';

export class WSClient {
  private ws: WebSocket | null = null;
  private url: string;
  private reconnectTimer: number | null = null;
  private messageHandlers: Set<(data: ArrayBuffer) => void> = new Set();

  constructor(url: string = 'ws://localhost:8080/api/v1/ws') {
    this.url = url;
  }

  public connect() {
    if (this.ws?.readyState === WebSocket.OPEN) return;

    this.ws = new WebSocket(this.url);
    this.ws.binaryType = 'arraybuffer'; // Crucial for protobuf

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      useAppStore.getState().setConnectionStatus(true);
      if (this.reconnectTimer) {
        window.clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      useAppStore.getState().setConnectionStatus(false);
      this.ws = null;
      this.scheduleReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      // Let onclose handle the reconnection
    };

    this.ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        // T019 specifies [1 byte type][payload bytes]
        const data = event.data as ArrayBuffer;
        const typeArray = new Uint8Array(data, 0, 1);
        const type = typeArray[0];

        if (type === 0x01) {
          // Type 1 is TelemetryFrame, notify handlers
          const payload = data.slice(1);
          this.messageHandlers.forEach(handler => handler(payload));
        }
      }
    };
  }

  public disconnect() {
    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer);
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  public addMessageHandler(handler: (data: ArrayBuffer) => void) {
    this.messageHandlers.add(handler);
  }

  public removeMessageHandler(handler: (data: ArrayBuffer) => void) {
    this.messageHandlers.delete(handler);
  }

  private scheduleReconnect() {
    if (!this.reconnectTimer) {
      this.reconnectTimer = window.setTimeout(() => {
        console.log('Attempting to reconnect WebSocket...');
        this.connect();
      }, 5000);
    }
  }
}

export const wsClient = new WSClient();
