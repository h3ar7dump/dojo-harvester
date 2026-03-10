import React, { useState } from 'react';
import { MainLayout } from '@/components/layout/MainLayout';
import { PrerequisiteChecklist, type PrerequisiteItem } from '@/components/recording/PrerequisiteChecklist';
import { api } from '@/services/api';
import { useAppStore } from '@/stores/useAppStore';

export const RecordingPreparationPage: React.FC = () => {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [prereqs, setPrereqs] = useState<PrerequisiteItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const setActiveSession = useAppStore(state => state.setActiveSession);

  const fetchPrereqs = async (sid: string) => {
    try {
      const items = await api.get<PrerequisiteItem[]>(`/sessions/${sid}/prerequisites`);
      setPrereqs(items);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch prerequisites');
    }
  };

  const handleStartSession = async () => {
    setLoading(true);
    setError(null);
    try {
      // Mocked TaskID and RobotID for MVP
      const session = await api.post<{ session_id: string }>('/sessions', {
        task_id: "task-123",
        robot_id: "robot-alpha",
        operator_id: "op-1"
      });
      setSessionId(session.session_id);
      setActiveSession(session.session_id);
      await fetchPrereqs(session.session_id);
    } catch (err: any) {
      setError(err.message || 'Failed to start session');
    } finally {
      setLoading(false);
    }
  };

  const allRequiredVerified = prereqs.length > 0 && prereqs
    .filter(p => p.is_required)
    .every(p => p.status === 'verified');

  return (
    <MainLayout>
      <div className="flex flex-col space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Recording Preparation</h1>
          <p className="text-muted-foreground mt-2">
            Complete the checklist to ensure the robot and environment are ready for data collection.
          </p>
        </div>

        {error && (
          <div className="bg-destructive/15 text-destructive p-4 rounded-md">
            {error}
          </div>
        )}

        {!sessionId ? (
          <div className="flex items-center justify-center p-12 border border-dashed rounded-lg">
            <button
              onClick={handleStartSession}
              disabled={loading}
              className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-11 px-8"
            >
              {loading ? 'Initializing...' : 'Start New Recording Session'}
            </button>
          </div>
        ) : (
          <div className="space-y-8">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold">Prerequisites</h2>
              <span className="text-sm bg-secondary text-secondary-foreground px-3 py-1 rounded-full">
                Session: {sessionId.substring(0, 8)}...
              </span>
            </div>

            <PrerequisiteChecklist 
              sessionId={sessionId} 
              items={prereqs} 
              onItemVerified={() => fetchPrereqs(sessionId)} 
            />

            <div className="flex justify-end pt-4 border-t">
              <button
                disabled={!allRequiredVerified}
                onClick={() => {
                  window.location.href = `/record/${sessionId}`;
                }}
                className={`inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors h-11 px-8 ${
                  allRequiredVerified 
                    ? 'bg-green-600 text-white hover:bg-green-700' 
                    : 'bg-muted text-muted-foreground cursor-not-allowed'
                }`}
              >
                Launch Recording
              </button>
            </div>
          </div>
        )}
      </div>
    </MainLayout>
  );
};
