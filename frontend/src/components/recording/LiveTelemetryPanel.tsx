import React from 'react';
import { useTelemetryStore } from '@/stores/telemetryStore';

export const LiveTelemetryPanel: React.FC = () => {
  const { framesCaptured, batteryLevel, errorFlags } = useTelemetryStore();

  return (
    <div className="bg-card text-card-foreground rounded-lg border p-4 space-y-4 shadow-sm">
      <h3 className="font-semibold text-lg">Live Telemetry</h3>
      
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground">Frames Captured</p>
          <p className="font-mono text-2xl">{framesCaptured}</p>
        </div>
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground">Battery Level</p>
          <div className="flex items-center space-x-2">
            <div className="flex-1 h-2 bg-secondary rounded-full overflow-hidden">
              <div 
                className={`h-full ${batteryLevel > 20 ? 'bg-green-500' : 'bg-red-500'}`} 
                style={{ width: `${batteryLevel}%` }}
              />
            </div>
            <span className="font-mono text-sm">{batteryLevel}%</span>
          </div>
        </div>
      </div>

      {errorFlags.length > 0 && (
        <div className="p-3 bg-destructive/10 text-destructive rounded-md mt-4">
          <p className="font-semibold text-sm mb-1">Errors Detected:</p>
          <ul className="list-disc pl-4 text-sm">
            {errorFlags.map((flag, idx) => (
              <li key={idx}>{flag}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};
