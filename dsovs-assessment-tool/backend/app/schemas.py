from __future__ import annotations

from datetime import datetime
from typing import Any

from pydantic import BaseModel, ConfigDict


# ── Maturity Level ────────────────────────────────────────────────────────────

class MaturityLevelSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    level: int
    title: str | None = None
    description: str | None = None
    evidence_json: list[Any] | None = None


# ── Control ───────────────────────────────────────────────────────────────────

class ControlSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    control_id: str
    code: str | None = None
    title: str
    phase: str | None = None
    slug: str | None = None
    status: str | None = None
    type: str | None = None
    summary: str | None = None
    doc_url: str | None = None
    maturity_levels: list[MaturityLevelSchema] = []


# ── Phase ─────────────────────────────────────────────────────────────────────

class PhaseSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    name: str
    sort_order: int


# ── Standard ──────────────────────────────────────────────────────────────────

class StandardSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    name: str
    abbreviation: str
    version: str
    source_url: str | None = None
    retrieved_at: datetime
    raw_hash: str
    phases: list[PhaseSchema] = []
    controls: list[ControlSchema] = []


# ── Catalogue sync ────────────────────────────────────────────────────────────

class SyncResultSchema(BaseModel):
    version: str
    control_count: int
    phase_count: int
    changed: bool
    message: str


# ── Project ───────────────────────────────────────────────────────────────────

class ProjectCreate(BaseModel):
    name: str
    client_name: str | None = None
    owner: str | None = None
    description: str | None = None


class ProjectUpdate(BaseModel):
    name: str | None = None
    client_name: str | None = None
    owner: str | None = None
    description: str | None = None


class ProjectSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    name: str
    client_name: str | None = None
    owner: str | None = None
    description: str | None = None
    created_at: datetime
    updated_at: datetime


# ── Assessment ────────────────────────────────────────────────────────────────

class AssessmentCreate(BaseModel):
    name: str
    assessment_date: datetime | None = None
    assessor: str | None = None
    scope: str | None = None
    status: str = "draft"


class AssessmentUpdate(BaseModel):
    name: str | None = None
    assessment_date: datetime | None = None
    assessor: str | None = None
    scope: str | None = None
    status: str | None = None


class AssessmentSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    project_id: int
    standard_id: int
    name: str
    assessment_date: datetime | None = None
    assessor: str | None = None
    scope: str | None = None
    status: str
    created_at: datetime
    updated_at: datetime


# ── Score ─────────────────────────────────────────────────────────────────────

class EvidenceLinkSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    label: str | None = None
    url: str | None = None
    notes: str | None = None


class ScoreUpsert(BaseModel):
    current_level: int | None = None
    target_level: int | None = None
    not_applicable: bool = False
    confidence: str | None = None
    priority: str | None = None
    evidence_notes: str | None = None
    action_notes: str | None = None


class ScoreSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    assessment_id: int
    control_id: int
    current_level: int | None = None
    target_level: int | None = None
    not_applicable: bool
    confidence: str | None = None
    priority: str | None = None
    evidence_notes: str | None = None
    action_notes: str | None = None
    created_at: datetime
    updated_at: datetime
    evidence_links: list[EvidenceLinkSchema] = []


# ── Results ───────────────────────────────────────────────────────────────────

class PhaseScoreSchema(BaseModel):
    phase: str
    current_score: float
    target_score: float
    control_count: int
    completed_count: int


class ControlGapSchema(BaseModel):
    control_id: int
    code: str | None
    title: str
    phase: str | None
    current_level: int | None
    target_level: int
    gap: int
    priority: str | None
    action_notes: str | None


class AssessmentResultsSchema(BaseModel):
    overall_score: float
    phase_scores: list[PhaseScoreSchema]
    control_gaps: list[ControlGapSchema]
    top_risks: list[ControlGapSchema]
    completed_count: int
    total_controls: int
    completion_percentage: float


# ── Trends ────────────────────────────────────────────────────────────────────

class TrendPointSchema(BaseModel):
    assessment_id: int
    assessment_name: str
    assessment_date: datetime | None
    overall_score: float
    phase_scores: dict[str, float]


class TrendsSchema(BaseModel):
    project_id: int
    trend_points: list[TrendPointSchema]


# ── Report ────────────────────────────────────────────────────────────────────

class ActionItemSchema(BaseModel):
    control_id: int
    code: str | None
    title: str
    phase: str | None
    current_level: int | None
    target_level: int
    gap: int
    priority: str | None
    action_notes: str | None


class ActionPlanSchema(BaseModel):
    days_30: list[ActionItemSchema]
    days_60: list[ActionItemSchema]
    days_90: list[ActionItemSchema]


class ReportDataSchema(BaseModel):
    project: ProjectSchema
    assessment: AssessmentSchema
    standard: StandardSchema
    results: AssessmentResultsSchema
    action_plan: ActionPlanSchema
    scores: list[ScoreSchema]
