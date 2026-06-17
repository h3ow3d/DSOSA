import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getAssessment, getResults } from "../api/client";
import GapBarChart from "../components/GapBarChart";
import PhaseRadarChart from "../components/PhaseRadarChart";
import type { Assessment, AssessmentResults } from "../types";

const PRIORITY_COLORS: Record<string, string> = {
  critical: "bg-red-100 text-red-700",
  high: "bg-orange-100 text-orange-700",
  medium: "bg-yellow-100 text-yellow-700",
  low: "bg-green-100 text-green-700",
};

export default function AssessmentResultsPage() {
  const { assessmentId } = useParams<{ assessmentId: string }>();
  const id = Number(assessmentId);

  const [assessment, setAssessment] = useState<Assessment | null>(null);
  const [results, setResults] = useState<AssessmentResults | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    Promise.all([getAssessment(id), getResults(id)])
      .then(([a, r]) => {
        setAssessment(a);
        setResults(r);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="text-gray-500">Loading…</div>;
  if (error) return <div className="text-red-500">Error: {error}</div>;
  if (!assessment || !results)
    return <div className="text-gray-400">Not found.</div>;

  const scoreColor =
    results.overall_score >= 2
      ? "text-green-600"
      : results.overall_score >= 1
      ? "text-yellow-600"
      : "text-red-600";

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">
            {assessment.name}
          </h1>
          <p className="text-sm text-gray-500">Assessment Results</p>
        </div>
        <div className="flex gap-2">
          <Link
            to={`/assessments/${id}`}
            className="border text-sm px-4 py-2 rounded hover:bg-gray-50"
          >
            ← Back to Wizard
          </Link>
          <Link
            to={`/report/${id}`}
            target="_blank"
            className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded"
          >
            Print / Save PDF
          </Link>
        </div>
      </div>

      {/* Summary cards */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
        <div className="bg-white border rounded-lg p-4 text-center">
          <div className={`text-3xl font-bold ${scoreColor}`}>
            {results.overall_score.toFixed(1)}
          </div>
          <div className="text-xs text-gray-500 mt-1">Overall Score / 3</div>
        </div>
        <div className="bg-white border rounded-lg p-4 text-center">
          <div className="text-3xl font-bold text-blue-600">
            {results.completion_percentage}%
          </div>
          <div className="text-xs text-gray-500 mt-1">Completion</div>
        </div>
        <div className="bg-white border rounded-lg p-4 text-center">
          <div className="text-3xl font-bold text-gray-700">
            {results.completed_count}
          </div>
          <div className="text-xs text-gray-500 mt-1">Controls Scored</div>
        </div>
        <div className="bg-white border rounded-lg p-4 text-center">
          <div className="text-3xl font-bold text-gray-700">
            {results.total_controls}
          </div>
          <div className="text-xs text-gray-500 mt-1">Total Controls</div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Radar */}
        <div className="bg-white border rounded-lg p-5">
          <h2 className="font-semibold text-gray-700 mb-3">Phase Maturity</h2>
          <PhaseRadarChart phaseScores={results.phase_scores} />
        </div>

        {/* Gap bar */}
        <div className="bg-white border rounded-lg p-5">
          <h2 className="font-semibold text-gray-700 mb-3">Top Gaps</h2>
          <GapBarChart gaps={results.control_gaps} />
        </div>
      </div>

      {/* Phase table */}
      <div className="bg-white border rounded-lg p-5 mb-6">
        <h2 className="font-semibold text-gray-700 mb-3">Phase Summary</h2>
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500 border-b">
              <th className="pb-2 pr-4">Phase</th>
              <th className="pb-2 pr-4">Current</th>
              <th className="pb-2 pr-4">Target</th>
              <th className="pb-2">Progress</th>
            </tr>
          </thead>
          <tbody>
            {results.phase_scores.map((ps) => (
              <tr key={ps.phase} className="border-b last:border-0">
                <td className="py-2 pr-4 font-medium">{ps.phase}</td>
                <td className="py-2 pr-4">{ps.current_score.toFixed(2)}</td>
                <td className="py-2 pr-4">{ps.target_score.toFixed(2)}</td>
                <td className="py-2">
                  {ps.completed_count}/{ps.control_count}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Control matrix */}
      <div className="bg-white border rounded-lg p-5">
        <h2 className="font-semibold text-gray-700 mb-3">Control Matrix</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500 border-b">
                <th className="pb-2 pr-3">Code</th>
                <th className="pb-2 pr-3">Title</th>
                <th className="pb-2 pr-3">Phase</th>
                <th className="pb-2 pr-3">Current</th>
                <th className="pb-2 pr-3">Target</th>
                <th className="pb-2 pr-3">Gap</th>
                <th className="pb-2 pr-3">Priority</th>
                <th className="pb-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              {results.control_gaps.map((g) => (
                <tr key={g.control_id} className="border-b last:border-0">
                  <td className="py-1.5 pr-3 font-mono text-xs">
                    {g.code ?? "—"}
                  </td>
                  <td className="py-1.5 pr-3 max-w-xs truncate">{g.title}</td>
                  <td className="py-1.5 pr-3 text-xs text-gray-500">
                    {g.phase}
                  </td>
                  <td className="py-1.5 pr-3">
                    {g.current_level ?? "—"}
                  </td>
                  <td className="py-1.5 pr-3">{g.target_level}</td>
                  <td className="py-1.5 pr-3">
                    {g.gap > 0 ? (
                      <span className="text-red-500 font-medium">
                        -{g.gap}
                      </span>
                    ) : (
                      <span className="text-green-500">✓</span>
                    )}
                  </td>
                  <td className="py-1.5 pr-3">
                    {g.priority && (
                      <span
                        className={`text-xs px-2 py-0.5 rounded-full font-medium ${
                          PRIORITY_COLORS[g.priority] ??
                          "bg-gray-100 text-gray-600"
                        }`}
                      >
                        {g.priority}
                      </span>
                    )}
                  </td>
                  <td className="py-1.5 text-xs text-gray-500 max-w-xs truncate">
                    {g.action_notes ?? "—"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
