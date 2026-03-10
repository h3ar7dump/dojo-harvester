import React, { useState, useEffect } from 'react';
import { MainLayout } from '@/components/layout/MainLayout';
import { api } from '@/services/api';

interface Dataset {
  dataset_id: string;
  status: string;
  validation_errors: string[];
}

export const ConversionStatusPage: React.FC = () => {
  // Extract session ID from URL in a real app (e.g., using React Router)
  const sessionId = window.location.pathname.split('/').pop();
  const [dataset, setDataset] = useState<Dataset | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (sessionId) {
      startConversion();
    }
  }, [sessionId]);

  const startConversion = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.post<Dataset>(`/datasets/${sessionId}/convert`);
      setDataset(data);
      pollStatus(data.dataset_id);
    } catch (err: any) {
      setError(err.message || 'Failed to start conversion');
      setLoading(false);
    }
  };

  const pollStatus = async (datasetId: string) => {
    const interval = setInterval(async () => {
      try {
        const data = await api.get<Dataset>(`/datasets/${datasetId}`);
        setDataset(data);
        if (data.status === 'converted' || data.status === 'invalid' || data.status === 'failed') {
          clearInterval(interval);
          setLoading(false);
        }
      } catch (err) {
        console.error('Failed to fetch status', err);
      }
    }, 2000);
  };

  if (!sessionId) {
    return (
      <MainLayout>
        <div className="p-8 text-center text-muted-foreground">
          Invalid session ID
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="max-w-2xl mx-auto space-y-8 mt-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Data Conversion</h1>
          <p className="text-muted-foreground mt-2">
            Converting raw recording to LeRobot V3.0 format
          </p>
        </div>

        {error && (
          <div className="bg-destructive/15 text-destructive p-4 rounded-md">
            {error}
          </div>
        )}

        {dataset && (
          <div className="bg-card border rounded-lg p-6 space-y-6">
            <div className="flex justify-between items-center">
              <span className="font-semibold text-lg">Status</span>
              <span className="px-3 py-1 bg-secondary rounded-full uppercase text-sm font-medium tracking-wider">
                {dataset.status}
              </span>
            </div>

            <div className="space-y-2">
              <div className="flex justify-between text-sm text-muted-foreground">
                <span>Progress</span>
                <span>{dataset.status === 'converted' ? '100%' : 'Processing...'}</span>
              </div>
              <div className="h-2 bg-secondary rounded-full overflow-hidden">
                <div 
                  className={`h-full transition-all duration-500 ${
                    dataset.status === 'failed' || dataset.status === 'invalid' ? 'bg-red-500' :
                    dataset.status === 'converted' ? 'bg-green-500' : 'bg-blue-500'
                  }`}
                  style={{ width: dataset.status === 'converted' ? '100%' : dataset.status === 'converting' || dataset.status === 'validating' ? '60%' : '0%' }}
                />
              </div>
            </div>

            {dataset.status === 'invalid' && dataset.validation_errors && (
              <div className="bg-destructive/10 text-destructive p-4 rounded-md mt-4 space-y-2">
                <p className="font-semibold">Validation Errors:</p>
                <ul className="list-disc pl-5 text-sm space-y-1">
                  {dataset.validation_errors.map((err, i) => (
                    <li key={i}>{err}</li>
                  ))}
                </ul>
              </div>
            )}

            <div className="flex justify-end pt-4 border-t gap-4">
              {(dataset.status === 'failed' || dataset.status === 'invalid') && (
                <button
                  onClick={startConversion}
                  disabled={loading}
                  className="px-4 py-2 border rounded-md hover:bg-accent"
                >
                  Retry Conversion
                </button>
              )}
              <button
                disabled={dataset.status !== 'converted'}
                onClick={() => {
                  window.location.href = `/upload/${dataset.dataset_id}`;
                }}
                className={`px-6 py-2 rounded-md font-medium text-white ${
                  dataset.status === 'converted' ? 'bg-green-600 hover:bg-green-700' : 'bg-muted text-muted-foreground cursor-not-allowed'
                }`}
              >
                Proceed to Upload
              </button>
            </div>
          </div>
        )}
      </div>
    </MainLayout>
  );
};
