import React, { useState, useEffect } from 'react';
import { MainLayout } from '@/components/layout/MainLayout';
import { api } from '@/services/api';

interface UploadJob {
  job_id: string;
  dataset_id: string;
  status: string;
  total_bytes: number;
  uploaded_bytes: number;
  progress_percentage: number;
  retry_count: number;
  last_error?: string;
  platform_endpoint: string;
}

export const UploadProgressPage: React.FC = () => {
  // In a real app, extract this from routing context
  const datasetId = window.location.pathname.split('/').pop();
  const [job, setJob] = useState<UploadJob | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (datasetId) {
      startUpload();
    }
  }, [datasetId]);

  const startUpload = async () => {
    setError(null);
    try {
      // POST to start or get existing upload
      const data = await api.post<UploadJob>(`/datasets/${datasetId}/upload`);
      setJob(data);
      pollStatus(data.job_id);
    } catch (err: any) {
      setError(err.message || 'Failed to start upload');
    }
  };

  const pollStatus = (jobId: string) => {
    const interval = setInterval(async () => {
      try {
        const data = await api.get<UploadJob>(`/uploads/${jobId}`);
        setJob(data);
        if (data.status === 'completed' || data.status === 'failed') {
          clearInterval(interval);
        }
      } catch (err) {
        console.error('Failed to poll upload status', err);
      }
    }, 1000);
  };

  const formatBytes = (bytes: number, decimals = 2) => {
    if (!+bytes) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  };

  if (!datasetId) {
    return (
      <MainLayout>
        <div className="p-8 text-center text-muted-foreground">Invalid dataset ID</div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="max-w-2xl mx-auto space-y-8 mt-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dataset Upload</h1>
          <p className="text-muted-foreground mt-2">
            Uploading converted dataset to data platform
          </p>
        </div>

        {error && (
          <div className="bg-destructive/15 text-destructive p-4 rounded-md">
            {error}
          </div>
        )}

        {job && (
          <div className="bg-card border rounded-lg p-6 space-y-6 shadow-sm">
            <div className="flex justify-between items-center border-b pb-4">
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                <p className="font-semibold capitalize">{job.status.replace('_', ' ')}</p>
              </div>
              <div className="text-right">
                <p className="text-sm text-muted-foreground">Platform</p>
                <p className="font-mono text-sm">{job.platform_endpoint}</p>
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>{formatBytes(job.uploaded_bytes)} / {formatBytes(job.total_bytes)}</span>
                <span className="font-mono">{job.progress_percentage.toFixed(1)}%</span>
              </div>
              <div className="h-3 bg-secondary rounded-full overflow-hidden">
                <div 
                  className={`h-full transition-all duration-300 ${
                    job.status === 'failed' ? 'bg-red-500' :
                    job.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'
                  }`}
                  style={{ width: `${job.status === 'completed' ? 100 : job.progress_percentage}%` }}
                />
              </div>
            </div>

            {job.retry_count > 0 && job.status === 'in_progress' && (
              <div className="text-sm text-amber-500 bg-amber-500/10 p-3 rounded flex items-center space-x-2">
                <span className="w-2 h-2 rounded-full bg-amber-500 animate-pulse" />
                <span>Recovering from network interruption (Retry {job.retry_count}/5)...</span>
              </div>
            )}

            {job.status === 'failed' && (
              <div className="bg-destructive/10 text-destructive p-4 rounded-md space-y-2">
                <p className="font-semibold">Upload Failed</p>
                <p className="text-sm font-mono">{job.last_error || 'Unknown error occurred'}</p>
              </div>
            )}

            {job.status === 'completed' && (
              <div className="bg-green-500/10 text-green-600 p-4 rounded-md flex items-center space-x-2">
                <span>✓</span>
                <span className="font-medium">Dataset successfully uploaded to platform.</span>
              </div>
            )}

            <div className="flex justify-end pt-4 gap-4">
              {job.status === 'failed' && (
                <button
                  onClick={() => pollStatus(job.job_id)}
                  className="px-4 py-2 border rounded-md hover:bg-accent text-sm font-medium"
                >
                  Retry Upload
                </button>
              )}
              {job.status === 'completed' && (
                <a
                  href="/"
                  className="px-6 py-2 rounded-md font-medium bg-primary text-primary-foreground hover:bg-primary/90 text-sm inline-flex items-center justify-center"
                >
                  Return to Dashboard
                </a>
              )}
            </div>
          </div>
        )}
      </div>
    </MainLayout>
  );
};
