from __future__ import annotations

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from app import crud, scoring
from app.crud import get_latest_standard
from app.database import get_db
from app.schemas import (
    AssessmentCreate,
    AssessmentResultsSchema,
    AssessmentSchema,
    AssessmentUpdate,
    ScoreSchema,
    ScoreUpsert,
    TrendPointSchema,
    TrendsSchema,
)

router = APIRouter(tags=["assessments"])


# ── Assessments ───────────────────────────────────────────────────────────────

@router.get("/api/projects/{project_id}/assessments", response_model=list[AssessmentSchema])
def list_assessments(project_id: int, db: Session = Depends(get_db)):
    crud.get_project(db, project_id)  # raises 404 if not found via dependency
    return crud.list_assessments(db, project_id)


@router.post(
    "/api/projects/{project_id}/assessments",
    response_model=AssessmentSchema,
    status_code=201,
)
def create_assessment(
    project_id: int, data: AssessmentCreate, db: Session = Depends(get_db)
):
    project = crud.get_project(db, project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")
    std = get_latest_standard(db)
    if std is None:
        raise HTTPException(
            status_code=400, detail="No catalogue loaded. Sync DSOVS first."
        )
    return crud.create_assessment(db, project_id, std.id, data)


@router.get("/api/assessments/{assessment_id}", response_model=AssessmentSchema)
def get_assessment(assessment_id: int, db: Session = Depends(get_db)):
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")
    return assessment


@router.put("/api/assessments/{assessment_id}", response_model=AssessmentSchema)
def update_assessment(
    assessment_id: int, data: AssessmentUpdate, db: Session = Depends(get_db)
):
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")
    return crud.update_assessment(db, assessment, data)


@router.delete("/api/assessments/{assessment_id}", status_code=204)
def delete_assessment(assessment_id: int, db: Session = Depends(get_db)):
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")
    crud.delete_assessment(db, assessment)


# ── Scores ────────────────────────────────────────────────────────────────────

@router.put(
    "/api/assessments/{assessment_id}/scores/{control_id}",
    response_model=ScoreSchema,
)
def upsert_score(
    assessment_id: int,
    control_id: int,
    data: ScoreUpsert,
    db: Session = Depends(get_db),
):
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")
    return crud.upsert_score(db, assessment_id, control_id, data)


# ── Results ───────────────────────────────────────────────────────────────────

@router.get(
    "/api/assessments/{assessment_id}/results",
    response_model=AssessmentResultsSchema,
)
def get_results(assessment_id: int, db: Session = Depends(get_db)):
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")

    controls = crud.get_controls_for_standard(db, assessment.standard_id)
    scores = crud.list_scores(db, assessment_id)
    scores_by_control = {s.control_id: s for s in scores}

    return scoring.compute_results(controls, scores_by_control)


# ── Trends ────────────────────────────────────────────────────────────────────

@router.get("/api/projects/{project_id}/trends", response_model=TrendsSchema)
def get_trends(project_id: int, db: Session = Depends(get_db)):
    project = crud.get_project(db, project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")

    assessments = crud.list_assessments(db, project_id)
    trend_points: list[TrendPointSchema] = []

    for a in assessments:
        controls = crud.get_controls_for_standard(db, a.standard_id)
        scores = crud.list_scores(db, a.id)
        scores_by_control = {s.control_id: s for s in scores}
        results = scoring.compute_results(controls, scores_by_control)

        phase_scores = {ps.phase: ps.current_score for ps in results.phase_scores}
        trend_points.append(
            TrendPointSchema(
                assessment_id=a.id,
                assessment_name=a.name,
                assessment_date=a.assessment_date,
                overall_score=results.overall_score,
                phase_scores=phase_scores,
            )
        )

    return TrendsSchema(project_id=project_id, trend_points=trend_points)
