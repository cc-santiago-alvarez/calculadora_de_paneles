import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AppShell from './components/Layout/AppShell';
import ProjectsListPage from './pages/ProjectsList/ProjectsListPage';
import ProjectFormPage from './pages/ProjectForm/ProjectFormPage';
import MapPage from './pages/MapView/MapPage';
import ResultsPage from './pages/Results/ResultsPage';
import ComparatorPage from './pages/Comparator/ComparatorPage';
import ReportPage from './pages/Report/ReportPage';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppShell />}>
          <Route path="/" element={<ProjectsListPage />} />
          <Route path="/new" element={<ProjectFormPage />} />
          <Route path="/map" element={<MapPage />} />
          <Route path="/results" element={<ResultsPage />} />
          <Route path="/compare" element={<ComparatorPage />} />
          <Route path="/report" element={<ReportPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
