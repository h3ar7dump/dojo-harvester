import React, { useState } from 'react';
import { api } from '@/services/api';

export interface PrerequisiteItem {
  item_id: string;
  name: string;
  description: string;
  category: string;
  verification_type: string;
  status: string;
  is_required: boolean;
  verification_data?: {
    free_space_mb?: number;
    [key: string]: any;
  };
}

interface ChecklistProps {
  sessionId: string;
  items: PrerequisiteItem[];
  onItemVerified: () => void;
}

export const PrerequisiteChecklist: React.FC<ChecklistProps> = ({ sessionId, items, onItemVerified }) => {
  const [verifying, setVerifying] = useState<string | null>(null);

  const handleVerify = async (itemId: string) => {
    setVerifying(itemId);
    try {
      await api.post(`/prerequisites/${sessionId}/${itemId}/verify`);
      onItemVerified();
    } catch (err) {
      console.error('Failed to verify item:', err);
    } finally {
      setVerifying(null);
    }
  };

  return (
    <div className="space-y-4">
      {items.map((item) => (
        <div key={item.item_id} className="flex items-start space-x-4 p-4 border rounded-lg bg-card text-card-foreground shadow-sm">
          <div className="flex h-6 items-center">
            {item.status === 'verified' ? (
              <div className="h-5 w-5 rounded-full bg-green-500 text-white flex items-center justify-center">✓</div>
            ) : item.status === 'failed' ? (
              <div className="h-5 w-5 rounded-full bg-red-500 text-white flex items-center justify-center">✗</div>
            ) : (
              <div className="h-5 w-5 rounded-full border-2 border-muted-foreground" />
            )}
          </div>
          <div className="flex-1 space-y-1">
            <p className="font-medium leading-none">{item.name}</p>
            <p className="text-sm text-muted-foreground">{item.description}</p>
            {item.verification_data?.free_space_mb && (
              <p className="text-sm font-mono text-blue-500 mt-1">
                Available Space: {(item.verification_data.free_space_mb / 1024).toFixed(2)} GB
              </p>
            )}
          </div>
          <div className="flex flex-col items-end">
            <span className="text-xs uppercase tracking-wider text-muted-foreground mb-2">
              {item.category} • {item.verification_type}
            </span>
            {item.status !== 'verified' && (
              <button
                onClick={() => handleVerify(item.item_id)}
                disabled={verifying === item.item_id}
                className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2"
              >
                {verifying === item.item_id ? 'Verifying...' : item.verification_type === 'auto' ? 'Auto Check' : 'Verify'}
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};
