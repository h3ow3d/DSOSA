import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getReportData } from "../api/client";
import GapBarChart from "../components/GapBarChart";
import PhaseRadarChart from "../components/PhaseRadarChart";
import type { ActionItem, ReportData } from "../types";

export default function ReportPage() {
  const { assessmentId } = useParams<{ assessmentId: string }>();
  const id = Number(assessmentId);

  const [data, setData] = useState<ReportData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    getReportData(id)
      .then(setData)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading)
    return (
      <div className="flex items-center justify-center h-screen text-gray-500">
        Loading report…
      </div>
    );
  if (error)
    return (
      <div className="flex items-center justify-center h-screen text-red-500">
        Error: {error}
      </div>
    );
  if (!data) return null;

  const { project, assessment, standard, results, action_plan, scores } = data;

  const scoresByControlId = Object.fromEntries(
    scores.map((s) => [s.control_id, s])
  );

  return (
    <div className="max-w-4xl mx-auto px-8 py-10 font-sans text-gray-800">
      {/* Print button */}
      <div className="no-print mb-6 flex justify-end">
        <button
          onClick={() => window.print()}
          className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded"
        >
          Print / Save as PDF
        </button>
      </div>

      {/* Cover */}
      <div className="mb-10 border-b pb-8">
        <div className="text-xs text-gray-400 uppercase tracking-widest mb-2">
          DevSecOps Maturity Assessment Report
        </div>
        <h1 className="text-3xl font-bold mb-1">{project.name}</h1>
        {project.client_name && (
          <div className="text-gray-500">{project.client_name}</div>
        )}
        <div className="mt-4 grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">Assessment:</span>{" "}
            <span className="font-medium">{assessment.name}</span>
          </div>
          <div>
            <span className="text-gray-500">Assessor:</span>{" "}
            <span className="font-medium">{assessment.assessor ?? "—"}</span>
          </div>
          <div>
            <span className="text-gray-500">Date:</span>{" "}
            <span className="font-medium">
              {assessment.assessment_date
                ? new Date(assessment.assessment_date).toLocaleDateString()
                : "—"}
            </span>
          </div>
          <div>
            <span className="text-gray-500">Standard:</span>{" "}
            <span className="font-medium">
              {standard.name} v{standard.version}
            </span>
          </div>
          {assessment.scope && (
            <div className="col-span-2">
              <span className="text-gray-500">Scope:</span>{" "}
              <span>{assessment.scope}</span>
            </div>
          )}
        </div>
      </div>

      {/* Executive Summary */}
      <section className="mb-10">
        <h2 className="text-xl font-bold mb-4">Executive Summary</h2>
        <div className="grid grid-cols-3 gap-4 mb-4">
          <SummaryBox
            label="Overall Score"
            value={`${results.overall_score.toFixed(1)} / 3`}
          />
          <SummaryBox
            label="Completion"
            value={`${results.completion_percentage}%`}
          />
          <SummaryBox
            label="Controls Scored"
            value={`${results.completed_count} / ${results.total_controls}`}
          />
        </div>
        <p className="text-sm text-gray-600">
          This assessment evaluated{" "}
          <strong>{results.total_controls}</strong> applicable controls across{" "}
          <strong>{results.phase_scores.length}</strong> phases. The overall
          maturity score is{" "}
          <strong>{results.overall_score.toFixed(2)}</strong> out of 3.
          {results.top_risks.length > 0 && (
            <>
              {" "}
              Top risks have been identified and an action plan has been
              generated below.
            </>
          )}
        </p>
      </section>

      {/* Phase Table */}
      <section className="mb-10">
        <h2 className="text-xl font-bold mb-4">Phase Maturity</h2>
        <table className="w-full text-sm border">
          <thead>
            <tr className="bg-gray-50 text-left">
              <th className="border px-3 py-2">Phase</th>
              <th className="border px-3 py-2">Current</th>
              <th className="border px-3 py-2">Target</th>
              <th className="border px-3 py-2">Controls</th>
              <th className="border px-3 py-2">Scored</th>
            </tr>
          </thead>
          <tbody>
            {results.phase_scores.map((ps) => (
              <tr key={ps.phase}>
                <td className="border px-3 py-1.5 font-medium">{ps.phase}</td>
                <td className="border px-3 py-1.5">
                  {ps.current_score.toFixed(2)}
                </td>
                <td className="border px-3 py-1.5">
                  {ps.target_score.toFixed(2)}
                </td>
                <td className="border px-3 py-1.5">{ps.control_count}</td>
                <td className="border px-3 py-1.5">{ps.completed_count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      {/* Charts */}
      <section className="mb-10 no-print">
        <h2 className="text-xl font-bold mb-4">Maturity Radar</h2>
        <PhaseRadarChart phaseScores={results.phase_scores} />
      </section>

      <section className="mb-10 no-print">
        <h2 className="text-xl font-bold mb-4">Top Maturity Gaps</h2>
        <GapBarChart gaps={results.control_gaps} />
      </section>

      {/* Top Risks */}
      {results.top_risks.length > 0 && (
        <section className="mb-10">
          <h2 className="text-xl font-bold mb-4">Top Risks</h2>
          <table className="w-full text-sm border">
            <thead>
              <tr className="bg-gray-50 text-left">
                <th className="border px-3 py-2">Code</th>
                <th className="border px-3 py-2">Control</th>
                <th className="border px-3 py-2">Phase</th>
                <th className="border px-3 py-2">Current</th>
                <th className="border px-3 py-2">Target</th>
                <th className="border px-3 py-2">Gap</th>
                <th className="border px-3 py-2">Priority</th>
              </tr>
            </thead>
            <tbody>
              {results.top_risks.map((r) => (
                <tr key={r.control_id}>
                  <td className="border px-3 py-1.5 font-mono text-xs">
                    {r.code ?? "—"}
                  </td>
                  <td className="border px-3 py-1.5">{r.title}</td>
                  <td className="border px-3 py-1.5 text-xs">{r.phase}</td>
                  <td className="border px-3 py-1.5">{r.current_level ?? "—"}</td>
                  <td className="border px-3 py-1.5">{r.target_level}</td>
                  <td className="border px-3 py-1.5 text-red-600 font-medium">
                    {r.gap}
                  </td>
                  <td className="border px-3 py-1.5 capitalize">
                    {r.priority ?? "—"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}

      {/* Action Plan */}
      <section className="mb-10 page-break">
        <h2 className="text-xl font-bold mb-4">30 / 60 / 90 Day Action Plan</h2>
        <ActionPlanSection
          title="30 Days"
          description="Critical and High priority controls with significant gaps"
          items={action_plan.days_30}
        />
        <ActionPlanSection
          title="60 Days"
          description="High and Medium priority controls with remaining gaps"
          items={action_plan.days_60}
        />
        <ActionPlanSection
          title="90 Days"
          description="Remaining improvements and optimisation work"
          items={action_plan.days_90}
        />
      </section>

      {/* Full Control Matrix */}
      <section className="mb-10 page-break">
        <h2 className="text-xl font-bold mb-4">Full Control Matrix</h2>
        <table className="w-full text-xs border">
          <thead>
            <tr className="bg-gray-50 text-left">
              <th className="border px-2 py-1.5">Code</th>
              <th className="border px-2 py-1.5">Control</th>
              <th className="border px-2 py-1.5">Phase</th>
              <th className="border px-2 py-1.5">Current</th>
              <th className="border px-2 py-1.5">Target</th>
              <th className="border px-2 py-1.5">Gap</th>
              <th className="border px-2 py-1.5">Priority</th>
            </tr>
          </thead>
          <tbody>
            {results.control_gaps.map((g) => (
              <tr key={g.control_id}>
                <td className="border px-2 py-1 font-mono">{g.code ?? "—"}</td>
                <td className="border px-2 py-1">{g.title}</td>
                <td className="border px-2 py-1">{g.phase ?? "—"}</td>
                <td className="border px-2 py-1">{g.current_level ?? "—"}</td>
                <td className="border px-2 py-1">{g.target_level}</td>
                <td className="border px-2 py-1">{g.gap}</td>
                <td className="border px-2 py-1 capitalize">
                  {g.priority ?? "—"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      {/* Evidence Appendix */}
      <section className="mb-10 page-break">
        <h2 className="text-xl font-bold mb-4">Evidence Appendix</h2>
        {scores.filter((s) => s.evidence_notes || s.action_notes).length ===
        0 ? (
          <p className="text-sm text-gray-400">No evidence notes recorded.</p>
        ) : (
          <div className="space-y-4">
            {scores
              .filter((s) => s.evidence_notes || s.action_notes)
              .map((s) => {
                const ctrl = data.standard.controls.find(
                  (c) => c.id === s.control_id
                );
                return (
                  <div
                    key={s.id}
                    className="border rounded p-3 text-sm"
                  >
                    <div className="font-medium mb-1">
                      {ctrl?.code && (
                        <span className="font-mono text-xs mr-2">
                          {ctrl.code}
                        </span>
                      )}
                      {ctrl?.title ?? `Control #${s.control_id}`}
                    </div>
                    {s.evidence_notes && (
                      <div className="text-gray-600 text-xs mb-1">
                        <span className="font-semibold">Evidence:</span>{" "}
                        {s.evidence_notes}
                      </div>
                    )}
                    {s.action_notes && (
                      <div className="text-gray-600 text-xs">
                        <span className="font-semibold">Action:</span>{" "}
                        {s.action_notes}
                      </div>
                    )}
                  </div>
                );
              })}
          </div>
        )}
      </section>
    </div>
  );
}

function SummaryBox({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="border rounded p-3 text-center">
      <div className="text-2xl font-bold text-blue-700">{value}</div>
      <div className="text-xs text-gray-500 mt-1">{label}</div>
    </div>
  );
}

function ActionPlanSection({
  title,
  description,
  items,
}: {
  title: string;
  description: string;
  items: ActionItem[];
}) {
  return (
    <div className="mb-6">
      <h3 className="font-semibold text-gray-700 mb-1">
        {title}{" "}
        <span className="font-normal text-sm text-gray-500">
          — {description}
        </span>
      </h3>
      {items.length === 0 ? (
        <p className="text-sm text-gray-400 ml-4">No items.</p>
      ) : (
        <table className="w-full text-sm border ml-0">
          <thead>
            <tr className="bg-gray-50 text-left">
              <th className="border px-3 py-1.5">Code</th>
              <th className="border px-3 py-1.5">Control</th>
              <th className="border px-3 py-1.5">Gap</th>
              <th className="border px-3 py-1.5">Action</th>
            </tr>
          </thead>
          <tbody>
            {items.map((item) => (
              <tr key={item.control_id}>
                <td className="border px-3 py-1 font-mono text-xs">
                  {item.code ?? "—"}
                </td>
                <td className="border px-3 py-1">{item.title}</td>
                <td className="border px-3 py-1 text-red-600 font-medium">
                  {item.gap}
                </td>
                <td className="border px-3 py-1 text-gray-600">
                  {item.action_notes ?? "—"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
