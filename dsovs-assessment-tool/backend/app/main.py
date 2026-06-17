from __future__ import annotations

import asyncio
import logging

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.routers import catalogue, projects, assessments, reports

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="DSOVS Assessment Tool", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=[settings.frontend_origin, "http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(catalogue.router)
app.include_router(projects.router)
app.include_router(assessments.router)
app.include_router(reports.router)


@app.on_event("startup")
async def startup_event() -> None:
    if settings.auto_sync_catalogue:
        logger.info("AUTO_SYNC_CATALOGUE=true – syncing on startup …")
        try:
            from app.database import SessionLocal
            from app.routers.catalogue import sync_catalogue

            db = SessionLocal()
            try:
                result = await sync_catalogue(db)
                logger.info("Auto-sync result: %s", result.message)
            finally:
                db.close()
        except Exception as exc:
            logger.warning("Auto-sync failed: %s", exc)


@app.get("/health")
def health() -> dict:
    return {"status": "ok"}
