import { Component, type ErrorInfo, type ReactNode } from 'react';

interface Props {
  children?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="flex min-h-screen items-center justify-center p-4 text-center">
          <div className="space-y-4 max-w-md w-full">
            <h1 className="text-4xl font-bold tracking-tight text-destructive">System Error</h1>
            <p className="text-muted-foreground">
              An unexpected error occurred in the data harvest application.
            </p>
            <div className="bg-destructive/10 text-destructive p-4 rounded-md font-mono text-sm text-left overflow-auto">
              {this.state.error?.message}
            </div>
            <button
              onClick={() => window.location.href = '/'}
              className="mt-4 px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium"
            >
              Return to Dashboard
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
