// Shared TypeScript types for the DSOVS Assessment Tool

export interface MaturityLevel {
  id: number;
  level: number;
  title: string | null;
  description: string | null;
  evidence_json: unknown[] | null;
}

export interface Control {
  id: number;
  control_id: string;
  code: string | null;
  title: string;
  phase: string | null;
  slug: string | null;
  status: string | null;
  type: string | null;
  summary: string | null;
  doc_url: string | null;
  maturity_levels: MaturityLevel[];
}

export interface Phase {
  id: number;
  name: string;
  sort_order: number;
}

export interface Standard {
  id: number;
  name: string;
  abbreviation: string;
  version: string;
  source_url: string | null;
  retrieved_at: string;
  raw_hash: string;
  phases: Phase[];
  controls: Control[];
}

export interface SyncResult {
  version: string;
  control_count: number;
  phase_count: number;
  changed: boolean;
  message: string;
}

export interface Project {
  id: number;
  name: string;
  client_name: string | null;
  owner: string | null;
  description: string | null;
  created_at: string;
  updated_at: string;
}

export interface ProjectCreate {
  name: string;
  client_name?: string;
  owner?: string;
  description?: string;
}

export interface Assessment {
  id: number;
  project_id: number;
  standard_id: number;
  name: string;
  assessment_date: string | null;
  assessor: string | null;
  scope: string | null;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface AssessmentCreate {
  name: string;
  assessment_date?: string;
  assessor?: string;
  scope?: string;
  status?: string;
}

export interface EvidenceLink {
  id: number;
  label: string | null;
  url: string | null;
  notes: string | null;
}

export interface Score {
  id: number;
  assessment_id: number;
  control_id: number;
  current_level: number | null;
  target_level: number | null;
  not_applicable: boolean;
  confidence: string | null;
  priority: string | null;
  evidence_notes: string | null;
  action_notes: string | null;
  created_at: string;
  updated_at: string;
  evidence_links: EvidenceLink[];
}

export interface ScoreUpsert {
  current_level?: number | null;
  target_level?: number | null;
  not_applicable?: boolean;
  confidence?: string | null;
  priority?: string | null;
  evidence_notes?: string | null;
  action_notes?: string | null;
}

export interface PhaseScore {
  phase: string;
  current_score: number;
  target_score: number;
  control_count: number;
  completed_count: number;
}

export interface ControlGap {
  control_id: number;
  code: string | null;
  title: string;
  phase: string | null;
  current_level: number | null;
  target_level: number;
  gap: number;
  priority: string | null;
  action_notes: string | null;
}

export interface AssessmentResults {
  overall_score: number;
  phase_scores: PhaseScore[];
  control_gaps: ControlGap[];
  top_risks: ControlGap[];
  completed_count: number;
  total_controls: number;
  completion_percentage: number;
}

export interface TrendPoint {
  assessment_id: number;
  assessment_name: string;
  assessment_date: string | null;
  overall_score: number;
  phase_scores: Record<string, number>;
}

export interface Trends {
  project_id: number;
  trend_points: TrendPoint[];
}

export interface ActionItem {
  control_id: number;
  code: string | null;
  title: string;
  phase: string | null;
  current_level: number | null;
  target_level: number;
  gap: number;
  priority: string | null;
  action_notes: string | null;
}

export interface ActionPlan {
  days_30: ActionItem[];
  days_60: ActionItem[];
  days_90: ActionItem[];
}

export interface ReportData {
  project: Project;
  assessment: Assessment;
  standard: Standard;
  results: AssessmentResults;
  action_plan: ActionPlan;
  scores: Score[];
}
