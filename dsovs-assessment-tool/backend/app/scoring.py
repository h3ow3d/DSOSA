from __future__ import annotations

from app.models import AssessmentScore, Control
from app.schemas import (
    ActionItemSchema,
    ActionPlanSchema,
    AssessmentResultsSchema,
    ControlGapSchema,
    PhaseScoreSchema,
)

DEFAULT_TARGET = 2


def _effective_target(score: AssessmentScore | None) -> int:
    if score and score.target_level is not None:
        return score.target_level
    return DEFAULT_TARGET


def _suggest_priority(current: int | None, target: int) -> str:
    if current is None:
        return "medium"
    gap = target - current
    if current == 0 and target == 3:
        return "critical"
    if gap >= 2:
        return "high"
    if gap == 1:
        return "medium"
    return "low"


def compute_results(
    controls: list[Control],
    scores_by_control: dict[int, AssessmentScore],
) -> AssessmentResultsSchema:
    phase_data: dict[str, dict] = {}
    all_applicable: list[tuple[Control, AssessmentScore | None]] = []

    for ctrl in controls:
        score = scores_by_control.get(ctrl.id)
        if score and score.not_applicable:
            continue

        all_applicable.append((ctrl, score))
        phase_name = ctrl.phase or "Unknown"
        if phase_name not in phase_data:
            phase_data[phase_name] = {
                "current_levels": [],
                "target_levels": [],
                "control_count": 0,
                "completed_count": 0,
            }

        pd = phase_data[phase_name]
        pd["control_count"] += 1
        if score and score.current_level is not None:
            pd["current_levels"].append(score.current_level)
            pd["completed_count"] += 1
        if score and score.target_level is not None:
            pd["target_levels"].append(score.target_level)
        else:
            pd["target_levels"].append(DEFAULT_TARGET)

    # Phase scores
    phase_scores: list[PhaseScoreSchema] = []
    for phase_name, pd in sorted(phase_data.items()):
        current_avg = (
            sum(pd["current_levels"]) / len(pd["current_levels"])
            if pd["current_levels"]
            else 0.0
        )
        target_avg = (
            sum(pd["target_levels"]) / len(pd["target_levels"])
            if pd["target_levels"]
            else float(DEFAULT_TARGET)
        )
        phase_scores.append(
            PhaseScoreSchema(
                phase=phase_name,
                current_score=round(current_avg, 2),
                target_score=round(target_avg, 2),
                control_count=pd["control_count"],
                completed_count=pd["completed_count"],
            )
        )

    # Overall
    all_current = [
        s.current_level
        for _, s in all_applicable
        if s and s.current_level is not None
    ]
    overall_score = round(sum(all_current) / len(all_current), 2) if all_current else 0.0
    completed_count = len(all_current)
    total_controls = len(all_applicable)
    completion_pct = (
        round(completed_count / total_controls * 100, 1) if total_controls else 0.0
    )

    # Gaps
    control_gaps: list[ControlGapSchema] = []
    for ctrl, score in all_applicable:
        target = _effective_target(score)
        current = score.current_level if score else None
        gap = target - (current if current is not None else 0)
        priority = (score.priority if score and score.priority else None) or _suggest_priority(
            current, target
        )
        control_gaps.append(
            ControlGapSchema(
                control_id=ctrl.id,
                code=ctrl.code,
                title=ctrl.title,
                phase=ctrl.phase,
                current_level=current,
                target_level=target,
                gap=gap,
                priority=priority,
                action_notes=score.action_notes if score else None,
            )
        )

    control_gaps.sort(key=lambda g: -g.gap)
    top_risks = [g for g in control_gaps if g.gap > 0][:10]

    return AssessmentResultsSchema(
        overall_score=overall_score,
        phase_scores=phase_scores,
        control_gaps=control_gaps,
        top_risks=top_risks,
        completed_count=completed_count,
        total_controls=total_controls,
        completion_percentage=completion_pct,
    )


def build_action_plan(results: AssessmentResultsSchema) -> ActionPlanSchema:
    def to_action(g: ControlGapSchema) -> ActionItemSchema:
        return ActionItemSchema(
            control_id=g.control_id,
            code=g.code,
            title=g.title,
            phase=g.phase,
            current_level=g.current_level,
            target_level=g.target_level,
            gap=g.gap,
            priority=g.priority,
            action_notes=g.action_notes,
        )

    days_30 = [
        to_action(g)
        for g in results.control_gaps
        if g.priority in ("critical", "high")
        and g.gap >= 2
    ]
    days_60 = [
        to_action(g)
        for g in results.control_gaps
        if g.priority in ("high", "medium")
        and g.gap >= 1
        and to_action(g) not in days_30
    ]
    # Avoid duplicates in days_60
    days_30_ids = {a.control_id for a in days_30}
    days_60 = [a for a in [
        to_action(g)
        for g in results.control_gaps
        if g.priority in ("high", "medium") and g.gap >= 1
    ] if a.control_id not in days_30_ids]

    days_30_60_ids = days_30_ids | {a.control_id for a in days_60}
    days_90 = [
        to_action(g)
        for g in results.control_gaps
        if g.gap > 0 and g.control_id not in days_30_60_ids
    ]

    return ActionPlanSchema(days_30=days_30, days_60=days_60, days_90=days_90)
