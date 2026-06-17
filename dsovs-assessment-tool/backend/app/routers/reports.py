from __future__ import annotations

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from app import crud, scoring
from app.database import get_db
from app.schemas import ReportDataSchema, StandardSchema

router = APIRouter(prefix="/api/assessments", tags=["reports"])


@router.get("/{assessment_id}/report-data", response_model=ReportDataSchema)
def get_report_data(assessment_id: int, db: Session = Depends(get_db)) -> ReportDataSchema:
    assessment = crud.get_assessment(db, assessment_id)
    if assessment is None:
        raise HTTPException(status_code=404, detail="Assessment not found")

    project = crud.get_project(db, assessment.project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")

    std = crud.get_latest_standard(db)
    if std is None:
        raise HTTPException(status_code=404, detail="No catalogue loaded")

    controls = crud.get_controls_for_standard(db, assessment.standard_id)
    scores = crud.list_scores(db, assessment_id)
    scores_by_control = {s.control_id: s for s in scores}
    results = scoring.compute_results(controls, scores_by_control)
    action_plan = scoring.build_action_plan(results)

    return ReportDataSchema(
        project=project,
        assessment=assessment,
        standard=std,
        results=results,
        action_plan=action_plan,
        scores=scores,
    )
