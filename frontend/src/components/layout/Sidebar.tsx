import React from 'react';

export const Sidebar: React.FC = () => {
  return (
    <div className="pb-12 border-r bg-background w-64 min-h-[calc(100vh-3.5rem)] hidden md:block">
      <div className="space-y-4 py-4">
        <div className="px-3 py-2">
          <h2 className="mb-2 px-4 text-lg font-semibold tracking-tight">
            Navigation
          </h2>
          <div className="space-y-1">
            <a href="/" className="flex items-center rounded-md px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground">
              Dashboard
            </a>
            <a href="/tasks" className="flex items-center rounded-md px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground">
              Tasks
            </a>
            <a href="/record" className="flex items-center rounded-md px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground">
              Record
            </a>
            <a href="/uploads" className="flex items-center rounded-md px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground">
              Uploads
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};
