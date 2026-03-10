import { LoginPage } from './pages/LoginPage';
import { TaskDashboardPage } from './pages/TaskDashboardPage';
import { RecordingPreparationPage } from './pages/RecordingPreparationPage';
import { RecordingDashboardPage } from './pages/RecordingDashboardPage';
import { ConversionStatusPage } from './pages/ConversionStatusPage';
import { UploadProgressPage } from './pages/UploadProgressPage';

function App() {
  const path = window.location.pathname;

  if (path === '/login') return <LoginPage />;
  if (path === '/prepare') return <RecordingPreparationPage />;
  if (path.startsWith('/record/')) return <RecordingDashboardPage />;
  if (path.startsWith('/convert/')) return <ConversionStatusPage />;
  if (path.startsWith('/upload/')) return <UploadProgressPage />;

  // Default homepage (e.g. '/' or '/tasks')
  return <TaskDashboardPage />;
}

export default App;
