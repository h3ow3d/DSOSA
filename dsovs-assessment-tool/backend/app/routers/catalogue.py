from __future__ import annotations

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from app import crud, dsovs_client
from app.database import get_db
from app.schemas import StandardSchema, SyncResultSchema

router = APIRouter(prefix="/api/catalogue", tags=["catalogue"])


@router.post("/sync", response_model=SyncResultSchema)
async def sync_catalogue(db: Session = Depends(get_db)) -> SyncResultSchema:
    try:
        data, raw_hash = await dsovs_client.fetch_dsovs()
    except Exception as exc:
        raise HTTPException(status_code=502, detail=f"Failed to fetch DSOVS: {exc}") from exc

    existing = crud.get_standard_by_hash(db, raw_hash)
    if existing:
        phase_count = len(existing.phases) if existing.phases else 0
        # Count controls for existing
        from sqlalchemy import select
        from app.models import Control
        ctrl_count = db.execute(
            select(Control).where(Control.standard_id == existing.id)
        ).fetchall()
        return SyncResultSchema(
            version=existing.version,
            control_count=len(ctrl_count),
            phase_count=phase_count,
            changed=False,
            message="Catalogue already up to date.",
        )

    name = data.get("standard", "OWASP DSOVS")
    abbreviation = data.get("abbreviation", "DSOVS")
    version = data.get("version", "unknown")
    source_url = data.get("source")

    std = crud.create_standard(
        db,
        name=name,
        abbreviation=abbreviation,
        version=version,
        source_url=source_url,
        raw_hash=raw_hash,
        raw_json=data,
    )

    phases_raw = data.get("phases", [])
    for i, phase_name in enumerate(phases_raw):
        crud.create_phase(db, standard_id=std.id, name=phase_name, sort_order=i)

    controls_raw = data.get("controls", [])
    for ctrl_data in controls_raw:
        ctrl = crud.create_control(
            db,
            standard_id=std.id,
            control_id=str(ctrl_data.get("id", ctrl_data.get("code", ""))),
            code=ctrl_data.get("code"),
            title=ctrl_data.get("title", ""),
            phase=ctrl_data.get("phase"),
            slug=ctrl_data.get("slug"),
            status=ctrl_data.get("status"),
            type=ctrl_data.get("type"),
            summary=ctrl_data.get("summary"),
            doc_url=ctrl_data.get("doc_url"),
        )
        for lvl_data in ctrl_data.get("levels", []):
            crud.create_maturity_level(
                db,
                control_id=ctrl.id,
                level=int(lvl_data.get("level", 0)),
                title=lvl_data.get("title"),
                description=lvl_data.get("description"),
                evidence_json=lvl_data.get("evidence"),
            )

    db.commit()

    return SyncResultSchema(
        version=version,
        control_count=len(controls_raw),
        phase_count=len(phases_raw),
        changed=True,
        message=f"Synced DSOVS v{version} with {len(controls_raw)} controls.",
    )


@router.get("/current", response_model=StandardSchema)
def get_current_catalogue(db: Session = Depends(get_db)) -> StandardSchema:
    std = crud.get_latest_standard(db)
    if std is None:
        raise HTTPException(status_code=404, detail="No catalogue loaded yet.")
    return std
