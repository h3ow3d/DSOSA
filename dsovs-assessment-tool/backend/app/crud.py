from __future__ import annotations

from typing import Any

from sqlalchemy import select
from sqlalchemy.orm import Session, selectinload

from app.models import (
    Assessment,
    AssessmentScore,
    Control,
    EvidenceLink,
    MaturityLevel,
    Phase,
    Project,
    Standard,
)
from app.schemas import (
    AssessmentCreate,
    AssessmentUpdate,
    ProjectCreate,
    ProjectUpdate,
    ScoreUpsert,
)


# ── Standards / Catalogue ─────────────────────────────────────────────────────

def get_standard_by_hash(db: Session, raw_hash: str) -> Standard | None:
    return db.scalar(select(Standard).where(Standard.raw_hash == raw_hash))


def get_latest_standard(db: Session) -> Standard | None:
    return db.scalar(
        select(Standard)
        .options(
            selectinload(Standard.phases),
            selectinload(Standard.controls).selectinload(Control.maturity_levels),
        )
        .order_by(Standard.retrieved_at.desc())
    )


def create_standard(
    db: Session,
    *,
    name: str,
    abbreviation: str,
    version: str,
    source_url: str | None,
    raw_hash: str,
    raw_json: dict[str, Any],
) -> Standard:
    std = Standard(
        name=name,
        abbreviation=abbreviation,
        version=version,
        source_url=source_url,
        raw_hash=raw_hash,
        raw_json=raw_json,
    )
    db.add(std)
    db.flush()
    return std


def create_phase(
    db: Session, *, standard_id: int, name: str, sort_order: int
) -> Phase:
    phase = Phase(standard_id=standard_id, name=name, sort_order=sort_order)
    db.add(phase)
    return phase


def create_control(
    db: Session,
    *,
    standard_id: int,
    control_id: str,
    code: str | None,
    title: str,
    phase: str | None,
    slug: str | None,
    status: str | None,
    type: str | None,
    summary: str | None,
    doc_url: str | None,
) -> Control:
    ctrl = Control(
        standard_id=standard_id,
        control_id=control_id,
        code=code,
        title=title,
        phase=phase,
        slug=slug,
        status=status,
        type=type,
        summary=summary,
        doc_url=doc_url,
    )
    db.add(ctrl)
    db.flush()
    return ctrl


def create_maturity_level(
    db: Session,
    *,
    control_id: int,
    level: int,
    title: str | None,
    description: str | None,
    evidence_json: list[Any] | None,
) -> MaturityLevel:
    ml = MaturityLevel(
        control_id=control_id,
        level=level,
        title=title,
        description=description,
        evidence_json=evidence_json,
    )
    db.add(ml)
    return ml


# ── Projects ──────────────────────────────────────────────────────────────────

def list_projects(db: Session) -> list[Project]:
    return list(db.scalars(select(Project).order_by(Project.created_at.desc())))


def get_project(db: Session, project_id: int) -> Project | None:
    return db.scalar(select(Project).where(Project.id == project_id))


def create_project(db: Session, data: ProjectCreate) -> Project:
    project = Project(**data.model_dump())
    db.add(project)
    db.commit()
    db.refresh(project)
    return project


def update_project(db: Session, project: Project, data: ProjectUpdate) -> Project:
    for field, value in data.model_dump(exclude_none=True).items():
        setattr(project, field, value)
    db.commit()
    db.refresh(project)
    return project


def delete_project(db: Session, project: Project) -> None:
    db.delete(project)
    db.commit()


# ── Assessments ───────────────────────────────────────────────────────────────

def list_assessments(db: Session, project_id: int) -> list[Assessment]:
    return list(
        db.scalars(
            select(Assessment)
            .where(Assessment.project_id == project_id)
            .order_by(Assessment.assessment_date.asc().nulls_last(), Assessment.created_at.asc())
        )
    )


def get_assessment(db: Session, assessment_id: int) -> Assessment | None:
    return db.scalar(select(Assessment).where(Assessment.id == assessment_id))


def create_assessment(
    db: Session, project_id: int, standard_id: int, data: AssessmentCreate
) -> Assessment:
    assessment = Assessment(
        project_id=project_id,
        standard_id=standard_id,
        **data.model_dump(),
    )
    db.add(assessment)
    db.commit()
    db.refresh(assessment)
    return assessment


def update_assessment(
    db: Session, assessment: Assessment, data: AssessmentUpdate
) -> Assessment:
    for field, value in data.model_dump(exclude_none=True).items():
        setattr(assessment, field, value)
    db.commit()
    db.refresh(assessment)
    return assessment


def delete_assessment(db: Session, assessment: Assessment) -> None:
    db.delete(assessment)
    db.commit()


# ── Scores ────────────────────────────────────────────────────────────────────

def get_score(
    db: Session, assessment_id: int, control_id: int
) -> AssessmentScore | None:
    return db.scalar(
        select(AssessmentScore)
        .options(selectinload(AssessmentScore.evidence_links))
        .where(
            AssessmentScore.assessment_id == assessment_id,
            AssessmentScore.control_id == control_id,
        )
    )


def upsert_score(
    db: Session, assessment_id: int, control_id: int, data: ScoreUpsert
) -> AssessmentScore:
    score = get_score(db, assessment_id, control_id)
    if score is None:
        score = AssessmentScore(assessment_id=assessment_id, control_id=control_id)
        db.add(score)

    for field, value in data.model_dump().items():
        setattr(score, field, value)

    db.commit()
    db.refresh(score)
    return score


def list_scores(db: Session, assessment_id: int) -> list[AssessmentScore]:
    return list(
        db.scalars(
            select(AssessmentScore)
            .options(selectinload(AssessmentScore.evidence_links))
            .where(AssessmentScore.assessment_id == assessment_id)
        )
    )


def get_controls_for_standard(db: Session, standard_id: int) -> list[Control]:
    return list(
        db.scalars(
            select(Control)
            .options(selectinload(Control.maturity_levels))
            .where(Control.standard_id == standard_id)
        )
    )
