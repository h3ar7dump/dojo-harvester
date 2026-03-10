import React from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useAppStore } from '@/stores/useAppStore';

export const Header: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();
  const isConnected = useAppStore(state => state.isConnected);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center justify-between">
        <div className="flex items-center gap-4">
          <a href="/" className="flex items-center space-x-2">
            <span className="hidden font-bold sm:inline-block">
              Dojo Harvester
            </span>
          </a>
        </div>
        <div className="flex flex-1 items-center justify-end space-x-4">
          <div className="flex items-center space-x-2">
            <div className={`h-3 w-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} title={isConnected ? "Connected to Robot" : "Disconnected"} />
            <span className="text-sm text-muted-foreground hidden sm:inline-block">
              {isConnected ? 'Connected' : 'Offline'}
            </span>
          </div>
          <nav className="flex items-center space-x-1">
            {isAuthenticated ? (
              <button
                onClick={logout}
                className="text-sm font-medium transition-colors hover:text-primary"
              >
                Logout
              </button>
            ) : null}
          </nav>
        </div>
      </div>
    </header>
  );
};
