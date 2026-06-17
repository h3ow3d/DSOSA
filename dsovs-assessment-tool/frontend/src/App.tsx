import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Layout from "./components/Layout";
import DashboardPage from "./pages/DashboardPage";
import ProjectsPage from "./pages/ProjectsPage";
import ProjectDetailPage from "./pages/ProjectDetailPage";
import AssessmentWizardPage from "./pages/AssessmentWizardPage";
import AssessmentResultsPage from "./pages/AssessmentResultsPage";
import ReportPage from "./pages/ReportPage";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="projects" element={<ProjectsPage />} />
          <Route path="projects/:projectId" element={<ProjectDetailPage />} />
          <Route
            path="projects/:projectId/assessments/new"
            element={<AssessmentWizardPage />}
          />
          <Route
            path="assessments/:assessmentId"
            element={<AssessmentWizardPage />}
          />
          <Route
            path="assessments/:assessmentId/results"
            element={<AssessmentResultsPage />}
          />
        </Route>
        <Route path="/report/:assessmentId" element={<ReportPage />} />
      </Routes>
    </BrowserRouter>
  );
}
