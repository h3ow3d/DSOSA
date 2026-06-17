from datetime import datetime
from typing import Any

from sqlalchemy import (
    Boolean,
    CheckConstraint,
    DateTime,
    ForeignKey,
    Integer,
    String,
    Text,
    UniqueConstraint,
    func,
)
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.database import Base


class Standard(Base):
    __tablename__ = "standards"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    abbreviation: Mapped[str] = mapped_column(String(50), nullable=False)
    version: Mapped[str] = mapped_column(String(50), nullable=False)
    source_url: Mapped[str | None] = mapped_column(String(500))
    retrieved_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now()
    )
    raw_hash: Mapped[str] = mapped_column(String(64), nullable=False, unique=True)
    raw_json: Mapped[dict[str, Any]] = mapped_column(JSONB, nullable=False)

    phases: Mapped[list["Phase"]] = relationship(
        "Phase", back_populates="standard", cascade="all, delete-orphan"
    )
    controls: Mapped[list["Control"]] = relationship(
        "Control", back_populates="standard", cascade="all, delete-orphan"
    )
    assessments: Mapped[list["Assessment"]] = relationship(
        "Assessment", back_populates="standard"
    )


class Phase(Base):
    __tablename__ = "phases"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    standard_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("standards.id", ondelete="CASCADE"), nullable=False
    )
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    sort_order: Mapped[int] = mapped_column(Integer, nullable=False, default=0)

    standard: Mapped["Standard"] = relationship("Standard", back_populates="phases")


class Control(Base):
    __tablename__ = "controls"
    __table_args__ = (
        UniqueConstraint("standard_id", "control_id", name="uq_control_standard"),
    )

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    standard_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("standards.id", ondelete="CASCADE"), nullable=False
    )
    control_id: Mapped[str] = mapped_column(String(100), nullable=False)
    code: Mapped[str | None] = mapped_column(String(50))
    title: Mapped[str] = mapped_column(String(500), nullable=False)
    phase: Mapped[str | None] = mapped_column(String(255))
    slug: Mapped[str | None] = mapped_column(String(255))
    status: Mapped[str | None] = mapped_column(String(50))
    type: Mapped[str | None] = mapped_column(String(50))
    summary: Mapped[str | None] = mapped_column(Text)
    doc_url: Mapped[str | None] = mapped_column(String(500))

    standard: Mapped["Standard"] = relationship("Standard", back_populates="controls")
    maturity_levels: Mapped[list["MaturityLevel"]] = relationship(
        "MaturityLevel", back_populates="control", cascade="all, delete-orphan"
    )
    scores: Mapped[list["AssessmentScore"]] = relationship(
        "AssessmentScore", back_populates="control"
    )


class MaturityLevel(Base):
    __tablename__ = "maturity_levels"
    __table_args__ = (
        UniqueConstraint("control_id", "level", name="uq_maturity_level"),
    )

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    control_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("controls.id", ondelete="CASCADE"), nullable=False
    )
    level: Mapped[int] = mapped_column(Integer, nullable=False)
    title: Mapped[str | None] = mapped_column(String(500))
    description: Mapped[str | None] = mapped_column(Text)
    evidence_json: Mapped[list[Any] | None] = mapped_column(JSONB)

    control: Mapped["Control"] = relationship(
        "Control", back_populates="maturity_levels"
    )


class Project(Base):
    __tablename__ = "projects"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    client_name: Mapped[str | None] = mapped_column(String(255))
    owner: Mapped[str | None] = mapped_column(String(255))
    description: Mapped[str | None] = mapped_column(Text)
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now()
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )

    assessments: Mapped[list["Assessment"]] = relationship(
        "Assessment", back_populates="project", cascade="all, delete-orphan"
    )


class Assessment(Base):
    __tablename__ = "assessments"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    project_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("projects.id", ondelete="CASCADE"), nullable=False
    )
    standard_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("standards.id"), nullable=False
    )
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    assessment_date: Mapped[datetime | None] = mapped_column(DateTime(timezone=True))
    assessor: Mapped[str | None] = mapped_column(String(255))
    scope: Mapped[str | None] = mapped_column(Text)
    status: Mapped[str] = mapped_column(String(50), default="draft")
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now()
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )

    project: Mapped["Project"] = relationship("Project", back_populates="assessments")
    standard: Mapped["Standard"] = relationship(
        "Standard", back_populates="assessments"
    )
    scores: Mapped[list["AssessmentScore"]] = relationship(
        "AssessmentScore", back_populates="assessment", cascade="all, delete-orphan"
    )


class AssessmentScore(Base):
    __tablename__ = "assessment_scores"
    __table_args__ = (
        UniqueConstraint(
            "assessment_id", "control_id", name="uq_assessment_control_score"
        ),
        CheckConstraint(
            "current_level IS NULL OR (current_level >= 0 AND current_level <= 3)",
            name="chk_current_level",
        ),
        CheckConstraint(
            "target_level IS NULL OR (target_level >= 0 AND target_level <= 3)",
            name="chk_target_level",
        ),
    )

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    assessment_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("assessments.id", ondelete="CASCADE"), nullable=False
    )
    control_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("controls.id", ondelete="CASCADE"), nullable=False
    )
    current_level: Mapped[int | None] = mapped_column(Integer)
    target_level: Mapped[int | None] = mapped_column(Integer)
    not_applicable: Mapped[bool] = mapped_column(Boolean, default=False)
    confidence: Mapped[str | None] = mapped_column(String(50))
    priority: Mapped[str | None] = mapped_column(String(50))
    evidence_notes: Mapped[str | None] = mapped_column(Text)
    action_notes: Mapped[str | None] = mapped_column(Text)
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now()
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )

    assessment: Mapped["Assessment"] = relationship(
        "Assessment", back_populates="scores"
    )
    control: Mapped["Control"] = relationship("Control", back_populates="scores")
    evidence_links: Mapped[list["EvidenceLink"]] = relationship(
        "EvidenceLink", back_populates="score", cascade="all, delete-orphan"
    )


class EvidenceLink(Base):
    __tablename__ = "evidence_links"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    assessment_score_id: Mapped[int] = mapped_column(
        Integer,
        ForeignKey("assessment_scores.id", ondelete="CASCADE"),
        nullable=False,
    )
    label: Mapped[str | None] = mapped_column(String(255))
    url: Mapped[str | None] = mapped_column(String(500))
    notes: Mapped[str | None] = mapped_column(Text)

    score: Mapped["AssessmentScore"] = relationship(
        "AssessmentScore", back_populates="evidence_links"
    )
