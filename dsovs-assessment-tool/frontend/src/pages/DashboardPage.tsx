import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getCurrentCatalogue, listAssessments, listProjects, syncCatalogue } from "../api/client";
import type { Assessment, Project, Standard } from "../types";

export default function DashboardPage() {
  const [catalogue, setCatalogue] = useState<Standard | null>(null);
  const [projects, setProjects] = useState<Project[]>([]);
  const [latestAssessments, setLatestAssessments] = useState<Assessment[]>([]);
  const [syncing, setSyncing] = useState(false);
  const [syncMsg, setSyncMsg] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    Promise.all([
      getCurrentCatalogue().catch(() => null),
      listProjects(),
    ])
      .then(async ([cat, projs]) => {
        setCatalogue(cat);
        setProjects(projs);

        // Load latest assessments across all projects
        const all: Assessment[] = [];
        for (const p of projs.slice(0, 5)) {
          const a = await listAssessments(p.id).catch(() => []);
          all.push(...a);
        }
        all.sort(
          (a, b) =>
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
        setLatestAssessments(all.slice(0, 5));
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const handleSync = async () => {
    setSyncing(true);
    setSyncMsg("");
    try {
      const result = await syncCatalogue();
      setSyncMsg(result.message);
      const cat = await getCurrentCatalogue();
      setCatalogue(cat);
    } catch (e: unknown) {
      setSyncMsg("Sync failed: " + (e instanceof Error ? e.message : String(e)));
    } finally {
      setSyncing(false);
    }
  };

  if (loading) return <div className="text-gray-500">Loading…</div>;
  if (error) return <div className="text-red-500">Error: {error}</div>;

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h1>

      {/* Catalogue Card */}
      <div className="bg-white rounded-lg border p-5 mb-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <div className="text-sm text-gray-500 mb-1">DSOVS Catalogue</div>
          {catalogue ? (
            <div>
              <span className="font-semibold text-gray-800">
                {catalogue.name} v{catalogue.version}
              </span>
              <span className="ml-3 text-sm text-gray-500">
                {catalogue.controls.length} controls ·{" "}
                {catalogue.phases.length} phases
              </span>
            </div>
          ) : (
            <span className="text-yellow-600 font-medium">
              No DSOVS catalogue loaded yet. Sync from OWASP to begin.
            </span>
          )}
        </div>
        <div className="flex flex-col items-end gap-1">
          <button
            onClick={handleSync}
            disabled={syncing}
            className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded disabled:opacity-50"
          >
            {syncing ? "Syncing…" : "Sync DSOVS Catalogue"}
          </button>
          {syncMsg && (
            <span className="text-xs text-green-600">{syncMsg}</span>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 mb-6">
        <StatCard label="Total Projects" value={projects.length} />
        <StatCard label="Controls" value={catalogue?.controls.length ?? "—"} />
        <StatCard label="Phases" value={catalogue?.phases.length ?? "—"} />
      </div>

      {/* Latest Assessments */}
      <div className="bg-white rounded-lg border p-5">
        <h2 className="font-semibold text-gray-700 mb-3">Recent Assessments</h2>
        {latestAssessments.length === 0 ? (
          <p className="text-sm text-gray-400">
            No assessments yet.{" "}
            <Link to="/projects" className="text-blue-500 hover:underline">
              Create a project
            </Link>{" "}
            to get started.
          </p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500 border-b">
                <th className="pb-2 pr-4">Name</th>
                <th className="pb-2 pr-4">Status</th>
                <th className="pb-2">Date</th>
              </tr>
            </thead>
            <tbody>
              {latestAssessments.map((a) => (
                <tr key={a.id} className="border-b last:border-0">
                  <td className="py-2 pr-4">
                    <Link
                      to={`/assessments/${a.id}`}
                      className="text-blue-600 hover:underline"
                    >
                      {a.name}
                    </Link>
                  </td>
                  <td className="py-2 pr-4">
                    <StatusBadge status={a.status} />
                  </td>
                  <td className="py-2 text-gray-500">
                    {a.assessment_date
                      ? new Date(a.assessment_date).toLocaleDateString()
                      : "—"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
}: {
  label: string;
  value: string | number;
}) {
  return (
    <div className="bg-white rounded-lg border p-4">
      <div className="text-2xl font-bold text-blue-700">{value}</div>
      <div className="text-sm text-gray-500 mt-1">{label}</div>
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    draft: "bg-gray-100 text-gray-600",
    "in-progress": "bg-yellow-100 text-yellow-700",
    complete: "bg-green-100 text-green-700",
  };
  return (
    <span
      className={`text-xs px-2 py-0.5 rounded-full font-medium ${
        colors[status] ?? "bg-gray-100 text-gray-600"
      }`}
    >
      {status}
    </span>
  );
}
