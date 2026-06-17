import { useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getAssessment, getCurrentCatalogue } from "../api/client";
import ControlCard from "../components/ControlCard";
import type { Assessment, Control, Score, Standard } from "../types";

// Scores are fetched per-save; we start from an empty map and the ControlCard
// updates its own local state, so here we just pass initialScore once.
const PHASE_ORDER = [
  "Organisation",
  "Requirements",
  "Design",
  "Code",
  "Build",
  "Test",
  "Release",
  "Deploy",
  "Operate",
  "Monitor",
];

function phaseSort(name: string): number {
  const idx = PHASE_ORDER.findIndex((p) =>
    name.toLowerCase().includes(p.toLowerCase())
  );
  return idx === -1 ? 99 : idx;
}

export default function AssessmentWizardPage() {
  const { assessmentId } = useParams<{ assessmentId: string }>();
  const id = Number(assessmentId);

  const [assessment, setAssessment] = useState<Assessment | null>(null);
  const [catalogue, setCatalogue] = useState<Standard | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [activePhase, setActivePhase] = useState<string>("");

  useEffect(() => {
    Promise.all([getAssessment(id), getCurrentCatalogue()])
      .then(([a, cat]) => {
        setAssessment(a);
        setCatalogue(cat);
        const phases = getPhases(cat.controls);
        if (phases.length > 0) setActivePhase(phases[0]);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="text-gray-500">Loading…</div>;
  if (error) return <div className="text-red-500">Error: {error}</div>;
  if (!assessment || !catalogue)
    return <div className="text-gray-400">Not found.</div>;

  const phases = getPhases(catalogue.controls);
  const controlsByPhase = groupByPhase(catalogue.controls);
  const currentControls = controlsByPhase[activePhase] ?? [];

  return (
    <div className="flex gap-6">
      {/* Left navigation */}
      <div className="w-48 flex-shrink-0">
        <div className="sticky top-4">
          <div className="text-xs font-semibold text-gray-500 uppercase mb-2">
            Phases
          </div>
          <nav className="flex flex-col gap-1">
            {phases.map((phase) => (
              <button
                key={phase}
                onClick={() => setActivePhase(phase)}
                className={`text-left text-sm px-3 py-2 rounded transition-colors ${
                  activePhase === phase
                    ? "bg-blue-600 text-white font-medium"
                    : "text-gray-600 hover:bg-gray-100"
                }`}
              >
                {phase}
                <span className="ml-1 text-xs opacity-70">
                  ({controlsByPhase[phase]?.length ?? 0})
                </span>
              </button>
            ))}
          </nav>

          <div className="mt-4 pt-4 border-t">
            <Link
              to={`/assessments/${id}/results`}
              className="text-sm text-green-600 hover:underline"
            >
              → View Results
            </Link>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-xl font-bold text-gray-800">
              {assessment.name}
            </h1>
            {assessment.assessor && (
              <p className="text-sm text-gray-500">
                Assessor: {assessment.assessor}
              </p>
            )}
          </div>
          <Link
            to={`/assessments/${id}/results`}
            className="bg-green-600 hover:bg-green-700 text-white text-sm font-medium px-4 py-2 rounded"
          >
            View Results
          </Link>
        </div>

        <h2 className="text-lg font-semibold text-gray-700 mb-3">
          {activePhase} ({currentControls.length} controls)
        </h2>

        <div className="flex flex-col gap-4">
          {currentControls.map((ctrl) => (
            <ControlCard
              key={ctrl.id}
              control={ctrl}
              assessmentId={id}
              initialScore={undefined}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

function getPhases(controls: Control[]): string[] {
  const phases = [...new Set(controls.map((c) => c.phase ?? "Unknown"))];
  return phases.sort((a, b) => phaseSort(a) - phaseSort(b));
}

function groupByPhase(controls: Control[]): Record<string, Control[]> {
  return controls.reduce<Record<string, Control[]>>((acc, ctrl) => {
    const phase = ctrl.phase ?? "Unknown";
    if (!acc[phase]) acc[phase] = [];
    acc[phase].push(ctrl);
    return acc;
  }, {});
}
