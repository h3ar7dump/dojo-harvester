import React, { useEffect, useState } from 'react';
import { MainLayout } from '@/components/layout/MainLayout';
import { api } from '@/services/api';

export interface Task {
  task_id: string;
  title: string;
  description: string;
  required_episodes: number;
  completed_episodes: number;
  status: string;
  priority: string;
  robot_configuration: any;
  objectives: string[];
}

export const TaskDashboardPage: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    try {
      const data = await api.get<Task[]>('/tasks');
      setTasks(data || []);
    } catch (err) {
      console.error('Failed to fetch tasks', err);
    } finally {
      setLoading(false);
    }
  };

  const claimTask = async (taskId: string) => {
    try {
      await api.post(`/tasks/${taskId}/claim`);
      fetchTasks();
    } catch (err) {
      console.error('Failed to claim task', err);
    }
  };

  return (
    <MainLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Task Dashboard</h1>
          <p className="text-muted-foreground mt-2">
            View and claim recording tasks assigned from the platform.
          </p>
        </div>

        {loading ? (
          <div className="text-center p-8">Loading tasks...</div>
        ) : tasks.length === 0 ? (
          <div className="text-center p-12 border border-dashed rounded-lg">
            <p className="text-muted-foreground mb-4">No tasks found in local cache.</p>
            <button 
              onClick={() => api.post('/tasks/seed').then(fetchTasks)}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium"
            >
              Seed Mock Task
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {tasks.map(task => (
              <div key={task.task_id} className="bg-card border rounded-lg overflow-hidden shadow-sm flex flex-col">
                <div className="p-5 flex-1">
                  <div className="flex justify-between items-start mb-4">
                    <h3 className="font-semibold text-lg line-clamp-2">{task.title}</h3>
                    <span className={`text-xs px-2 py-1 rounded-full uppercase font-medium tracking-wide
                      ${task.priority === 'high' ? 'bg-red-100 text-red-700' : 'bg-blue-100 text-blue-700'}`}>
                      {task.priority}
                    </span>
                  </div>
                  <p className="text-sm text-muted-foreground mb-4 line-clamp-3">
                    {task.description}
                  </p>
                  
                  <div className="space-y-2 mb-4">
                    <div className="flex justify-between text-sm">
                      <span>Progress</span>
                      <span>{task.completed_episodes} / {task.required_episodes} eps</span>
                    </div>
                    <div className="h-2 bg-secondary rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-primary transition-all"
                        style={{ width: `${Math.min(100, (task.completed_episodes / task.required_episodes) * 100)}%` }}
                      />
                    </div>
                  </div>
                </div>
                
                <div className="bg-muted/50 p-4 border-t flex justify-between items-center">
                  <span className="text-xs text-muted-foreground uppercase font-medium tracking-wide">
                    {task.status}
                  </span>
                  
                  {task.status === 'pending' ? (
                    <button 
                      onClick={() => claimTask(task.task_id)}
                      className="px-4 py-1.5 bg-primary text-primary-foreground rounded text-sm font-medium hover:bg-primary/90"
                    >
                      Claim Task
                    </button>
                  ) : task.status === 'in_progress' ? (
                    <a 
                      href="/prepare"
                      className="px-4 py-1.5 bg-green-600 text-white rounded text-sm font-medium hover:bg-green-700"
                    >
                      Continue
                    </a>
                  ) : (
                    <span className="px-4 py-1.5 border border-input rounded text-sm font-medium text-muted-foreground bg-background">
                      Completed
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </MainLayout>
  );
};
