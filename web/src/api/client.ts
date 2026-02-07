import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// --- Types ---
export interface Project {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface TestSuite {
  id: number;
  project_id: number;
  name: string;
  yaml_content: string;
  created_at: string;
  updated_at: string;
}

export interface TestRun {
  id: number;
  suite_id: number;
  started_at: string;
  completed_at?: string;
  status: string;
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
}

export interface TestResult {
  id: number;
  run_id: number;
  test_name: string;
  passed: boolean;
  status_code: number;
  error_message: string;
  duration_ms: number;
  created_at: string;
}

export interface RunResponse {
  run_id: number;
  status: string;
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
  results: {
    Name: string;
    Passed: boolean;
    StatusCode: number;
    Error: string | null;
    Duration: number;
  }[];
}

// --- API Functions ---

// Health
export const getHealth = () => api.get('/health');

// Projects
export const listProjects = () =>
  api.get<{ projects: Project[] }>('/projects').then((r) => r.data.projects ?? []);

export const createProject = (name: string, description: string) =>
  api.post<Project>('/projects', { name, description }).then((r) => r.data);

export const getProject = (id: number) =>
  api.get<Project>(`/projects/${id}`).then((r) => r.data);

// Suites
export const listSuites = (projectId: number) =>
  api.get<{ suites: TestSuite[] }>(`/projects/${projectId}/suites`).then((r) => r.data.suites ?? []);

export const createSuite = (projectId: number, name: string, yamlContent: string) =>
  api.post<TestSuite>('/suites', { project_id: projectId, name, yaml_content: yamlContent }).then((r) => r.data);

export const getSuite = (id: number) =>
  api.get<TestSuite>(`/suites/${id}`).then((r) => r.data);

export const runSuite = (id: number) =>
  api.post<RunResponse>(`/suites/${id}/run`).then((r) => r.data);

export const listRuns = (suiteId: number) =>
  api.get<{ runs: TestRun[] }>(`/suites/${suiteId}/runs`).then((r) => r.data.runs ?? []);

// Runs
export const getTestRun = (id: number) =>
  api.get<TestRun>(`/runs/${id}`).then((r) => r.data);

export const getTestResults = (runId: number) =>
  api.get<{ results: TestResult[] }>(`/runs/${runId}/results`).then((r) => r.data.results ?? []);

export default api;
