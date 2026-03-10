import { create } from 'zustand';
import { wsClient } from '../websocket/client';

interface Pose {
  x: number; y: number; z: number;
  qw: number; qx: number; qy: number; qz: number;
}

interface TelemetryState {
  timestampNs: number;
  robotId: string;
  jointPositions: number[];
  batteryLevel: number;
  errorFlags: string[];
  pose: Pose | null;
  framesCaptured: number;
  updateTelemetry: (payload: Uint8Array) => void;
  reset: () => void;
}

export const useTelemetryStore = create<TelemetryState>((set) => {
  // Setup the WS listener once
  wsClient.addMessageHandler((_data: ArrayBuffer) => {
    // In a real implementation we would decode protobuf here
    // For now we'll simulate update since we don't have the generated ts files
    set((state) => ({
      framesCaptured: state.framesCaptured + 1,
      // Decode protobuf and update...
    }));
  });

  return {
    timestampNs: 0,
    robotId: '',
    jointPositions: [],
    batteryLevel: 100,
    errorFlags: [],
    pose: null,
    framesCaptured: 0,
    
    updateTelemetry: (_payload: Uint8Array) => {
      // Decode protobuf here
      set({ framesCaptured: useTelemetryStore.getState().framesCaptured + 1 });
    },
    
    reset: () => {
      set({
        timestampNs: 0,
        jointPositions: [],
        batteryLevel: 100,
        errorFlags: [],
        pose: null,
        framesCaptured: 0
      });
    }
  };
});
