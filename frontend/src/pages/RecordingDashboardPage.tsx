import React, { useState, useEffect } from 'react';
import { MainLayout } from '@/components/layout/MainLayout';
import { LiveTelemetryPanel } from '@/components/recording/LiveTelemetryPanel';
import { RobotVisualization } from '@/components/three/RobotVisualization';
import { api } from '@/services/api';
import { useAppStore } from '@/stores/useAppStore';

export const RecordingDashboardPage: React.FC = () => {
  const sessionId = useAppStore(state => state.activeSessionId);
  const [status, setStatus] = useState<'recording' | 'converting' | 'completed' | 'error'>('recording');
  const [duration, setDuration] = useState(0);

  useEffect(() => {
    if (status !== 'recording') return;

    const timer = setInterval(() => {
      setDuration(prev => prev + 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [status]);

  const handleStopRecording = async () => {
    if (!sessionId) return;
    try {
      await api.post(`/sessions/${sessionId}/stop`);
      setStatus('converting');
      // Typically we'd redirect to conversion page here or update UI
      setTimeout(() => {
        window.location.href = `/convert/${sessionId}`;
      }, 2000);
    } catch (err) {
      console.error('Failed to stop recording:', err);
      setStatus('error');
    }
  };

  const formatDuration = (seconds: number) => {
    const m = Math.floor(seconds / 60).toString().padStart(2, '0');
    const s = (seconds % 60).toString().padStart(2, '0');
    return `${m}:${s}`;
  };

  if (!sessionId) {
    return (
      <MainLayout>
        <div className="p-8 text-center text-muted-foreground">
          No active session. Please start a recording first.
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="flex flex-col space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Recording Session</h1>
            <p className="text-muted-foreground mt-2">
              Session ID: <span className="font-mono">{sessionId.substring(0, 8)}...</span>
            </p>
          </div>
          
          <div className="flex flex-col items-end space-y-2">
            <div className="text-4xl font-mono tabular-nums tracking-tight">
              {formatDuration(duration)}
            </div>
            {status === 'recording' && (
              <div className="flex items-center space-x-2 text-red-500 font-medium animate-pulse">
                <div className="w-3 h-3 rounded-full bg-red-500" />
                <span>REC</span>
              </div>
            )}
          </div>
        </div>

        {status === 'error' && (
          <div className="bg-destructive/15 text-destructive p-4 rounded-md">
            An error occurred while stopping the recording. Check system logs.
          </div>
        )}

        {status === 'converting' && (
          <div className="bg-blue-500/15 text-blue-500 p-4 rounded-md flex items-center space-x-2">
            <span className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
            <span>Transitioning to conversion phase...</span>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-6">
            <div className="bg-card border rounded-lg p-4">
              <h3 className="font-semibold text-lg mb-4">Robot Visualization</h3>
              <RobotVisualization />
            </div>
          </div>
          
          <div className="space-y-6">
            <LiveTelemetryPanel />
            
            <div className="pt-4">
              <button
                onClick={handleStopRecording}
                disabled={status !== 'recording'}
                className="w-full inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-destructive text-destructive-foreground hover:bg-destructive/90 h-14 text-lg"
              >
                {status === 'recording' ? 'Stop Recording' : 'Stopping...'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
};
