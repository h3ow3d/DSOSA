from __future__ import annotations

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from app import crud
from app.database import get_db
from app.schemas import ProjectCreate, ProjectSchema, ProjectUpdate

router = APIRouter(prefix="/api/projects", tags=["projects"])


@router.get("", response_model=list[ProjectSchema])
def list_projects(db: Session = Depends(get_db)) -> list[ProjectSchema]:
    return crud.list_projects(db)


@router.post("", response_model=ProjectSchema, status_code=201)
def create_project(data: ProjectCreate, db: Session = Depends(get_db)) -> ProjectSchema:
    return crud.create_project(db, data)


@router.get("/{project_id}", response_model=ProjectSchema)
def get_project(project_id: int, db: Session = Depends(get_db)) -> ProjectSchema:
    project = crud.get_project(db, project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")
    return project


@router.put("/{project_id}", response_model=ProjectSchema)
def update_project(
    project_id: int, data: ProjectUpdate, db: Session = Depends(get_db)
) -> ProjectSchema:
    project = crud.get_project(db, project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")
    return crud.update_project(db, project, data)


@router.delete("/{project_id}", status_code=204)
def delete_project(project_id: int, db: Session = Depends(get_db)) -> None:
    project = crud.get_project(db, project_id)
    if project is None:
        raise HTTPException(status_code=404, detail="Project not found")
    crud.delete_project(db, project)
