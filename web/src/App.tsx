import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import ProjectView from './pages/ProjectView';
import SuiteView from './pages/SuiteView';
import RunView from './pages/RunView';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Dashboard />} />
        <Route path="projects/:projectId" element={<ProjectView />} />
        <Route path="suites/:suiteId" element={<SuiteView />} />
        <Route path="runs/:runId" element={<RunView />} />
      </Route>
    </Routes>
  );
}

export default App;
