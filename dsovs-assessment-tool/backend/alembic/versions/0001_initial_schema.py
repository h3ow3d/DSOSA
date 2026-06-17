"""initial schema

Revision ID: 0001
Revises:
Create Date: 2024-01-01 00:00:00.000000

"""
from __future__ import annotations

from typing import Any

import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

from alembic import op

# revision identifiers
revision: str = "0001"
down_revision: str | None = None
branch_labels: str | None = None
depends_on: str | None = None


def upgrade() -> None:
    # standards
    op.create_table(
        "standards",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("abbreviation", sa.String(50), nullable=False),
        sa.Column("version", sa.String(50), nullable=False),
        sa.Column("source_url", sa.String(500)),
        sa.Column(
            "retrieved_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
        sa.Column("raw_hash", sa.String(64), nullable=False, unique=True),
        sa.Column("raw_json", postgresql.JSONB, nullable=False),
    )

    # phases
    op.create_table(
        "phases",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "standard_id",
            sa.Integer,
            sa.ForeignKey("standards.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("sort_order", sa.Integer, nullable=False, server_default="0"),
    )

    # controls
    op.create_table(
        "controls",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "standard_id",
            sa.Integer,
            sa.ForeignKey("standards.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("control_id", sa.String(100), nullable=False),
        sa.Column("code", sa.String(50)),
        sa.Column("title", sa.String(500), nullable=False),
        sa.Column("phase", sa.String(255)),
        sa.Column("slug", sa.String(255)),
        sa.Column("status", sa.String(50)),
        sa.Column("type", sa.String(50)),
        sa.Column("summary", sa.Text),
        sa.Column("doc_url", sa.String(500)),
        sa.UniqueConstraint("standard_id", "control_id", name="uq_control_standard"),
    )

    # maturity_levels
    op.create_table(
        "maturity_levels",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "control_id",
            sa.Integer,
            sa.ForeignKey("controls.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("level", sa.Integer, nullable=False),
        sa.Column("title", sa.String(500)),
        sa.Column("description", sa.Text),
        sa.Column("evidence_json", postgresql.JSONB),
        sa.UniqueConstraint("control_id", "level", name="uq_maturity_level"),
    )

    # projects
    op.create_table(
        "projects",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("client_name", sa.String(255)),
        sa.Column("owner", sa.String(255)),
        sa.Column("description", sa.Text),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
        sa.Column(
            "updated_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
    )

    # assessments
    op.create_table(
        "assessments",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "project_id",
            sa.Integer,
            sa.ForeignKey("projects.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column(
            "standard_id",
            sa.Integer,
            sa.ForeignKey("standards.id"),
            nullable=False,
        ),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("assessment_date", sa.DateTime(timezone=True)),
        sa.Column("assessor", sa.String(255)),
        sa.Column("scope", sa.Text),
        sa.Column("status", sa.String(50), server_default="draft"),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
        sa.Column(
            "updated_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
    )

    # assessment_scores
    op.create_table(
        "assessment_scores",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "assessment_id",
            sa.Integer,
            sa.ForeignKey("assessments.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column(
            "control_id",
            sa.Integer,
            sa.ForeignKey("controls.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("current_level", sa.Integer),
        sa.Column("target_level", sa.Integer),
        sa.Column("not_applicable", sa.Boolean, server_default="false"),
        sa.Column("confidence", sa.String(50)),
        sa.Column("priority", sa.String(50)),
        sa.Column("evidence_notes", sa.Text),
        sa.Column("action_notes", sa.Text),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
        sa.Column(
            "updated_at",
            sa.DateTime(timezone=True),
            server_default=sa.func.now(),
        ),
        sa.UniqueConstraint(
            "assessment_id", "control_id", name="uq_assessment_control_score"
        ),
        sa.CheckConstraint(
            "current_level IS NULL OR (current_level >= 0 AND current_level <= 3)",
            name="chk_current_level",
        ),
        sa.CheckConstraint(
            "target_level IS NULL OR (target_level >= 0 AND target_level <= 3)",
            name="chk_target_level",
        ),
    )

    # evidence_links
    op.create_table(
        "evidence_links",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "assessment_score_id",
            sa.Integer,
            sa.ForeignKey("assessment_scores.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("label", sa.String(255)),
        sa.Column("url", sa.String(500)),
        sa.Column("notes", sa.Text),
    )


def downgrade() -> None:
    op.drop_table("evidence_links")
    op.drop_table("assessment_scores")
    op.drop_table("assessments")
    op.drop_table("projects")
    op.drop_table("maturity_levels")
    op.drop_table("controls")
    op.drop_table("phases")
    op.drop_table("standards")
