import axios from "axios";
import type {
  Assessment,
  AssessmentCreate,
  AssessmentResults,
  Project,
  ProjectCreate,
  ReportData,
  Score,
  ScoreUpsert,
  Standard,
  SyncResult,
  Trends,
} from "../types";

const BASE = import.meta.env.VITE_API_BASE_URL || "";

const api = axios.create({ baseURL: BASE });

// Catalogue
export const syncCatalogue = (): Promise<SyncResult> =>
  api.post("/api/catalogue/sync").then((r) => r.data);

export const getCurrentCatalogue = (): Promise<Standard> =>
  api.get("/api/catalogue/current").then((r) => r.data);

// Projects
export const listProjects = (): Promise<Project[]> =>
  api.get("/api/projects").then((r) => r.data);

export const createProject = (data: ProjectCreate): Promise<Project> =>
  api.post("/api/projects", data).then((r) => r.data);

export const getProject = (id: number): Promise<Project> =>
  api.get(`/api/projects/${id}`).then((r) => r.data);

export const updateProject = (
  id: number,
  data: Partial<ProjectCreate>
): Promise<Project> => api.put(`/api/projects/${id}`, data).then((r) => r.data);

export const deleteProject = (id: number): Promise<void> =>
  api.delete(`/api/projects/${id}`).then(() => undefined);

// Assessments
export const listAssessments = (projectId: number): Promise<Assessment[]> =>
  api.get(`/api/projects/${projectId}/assessments`).then((r) => r.data);

export const createAssessment = (
  projectId: number,
  data: AssessmentCreate
): Promise<Assessment> =>
  api.post(`/api/projects/${projectId}/assessments`, data).then((r) => r.data);

export const getAssessment = (id: number): Promise<Assessment> =>
  api.get(`/api/assessments/${id}`).then((r) => r.data);

export const updateAssessment = (
  id: number,
  data: Partial<AssessmentCreate>
): Promise<Assessment> =>
  api.put(`/api/assessments/${id}`, data).then((r) => r.data);

export const deleteAssessment = (id: number): Promise<void> =>
  api.delete(`/api/assessments/${id}`).then(() => undefined);

// Scores
export const upsertScore = (
  assessmentId: number,
  controlId: number,
  data: ScoreUpsert
): Promise<Score> =>
  api
    .put(`/api/assessments/${assessmentId}/scores/${controlId}`, data)
    .then((r) => r.data);

// Results
export const getResults = (assessmentId: number): Promise<AssessmentResults> =>
  api.get(`/api/assessments/${assessmentId}/results`).then((r) => r.data);

// Trends
export const getTrends = (projectId: number): Promise<Trends> =>
  api.get(`/api/projects/${projectId}/trends`).then((r) => r.data);

// Report
export const getReportData = (assessmentId: number): Promise<ReportData> =>
  api.get(`/api/assessments/${assessmentId}/report-data`).then((r) => r.data);
