import { create } from 'zustand';

interface AppState {
  isConnected: boolean;
  activeSessionId: string | null;
  setConnectionStatus: (status: boolean) => void;
  setActiveSession: (sessionId: string | null) => void;
}

export const useAppStore = create<AppState>((set) => ({
  isConnected: false,
  activeSessionId: null,
  setConnectionStatus: (status) => set({ isConnected: status }),
  setActiveSession: (sessionId) => set({ activeSessionId: sessionId }),
}));
