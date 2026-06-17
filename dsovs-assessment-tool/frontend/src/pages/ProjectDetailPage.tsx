import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  createAssessment,
  deleteAssessment,
  getProject,
  getTrends,
  listAssessments,
} from "../api/client";
import TrendLineChart from "../components/TrendLineChart";
import type { Assessment, AssessmentCreate, Project, Trends } from "../types";

export default function ProjectDetailPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const id = Number(projectId);

  const [project, setProject] = useState<Project | null>(null);
  const [assessments, setAssessments] = useState<Assessment[]>([]);
  const [trends, setTrends] = useState<Trends | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<AssessmentCreate>({
    name: "",
    assessor: "",
    scope: "",
  });
  const [creating, setCreating] = useState(false);

  const load = async () => {
    setLoading(true);
    try {
      const [proj, asmts, tr] = await Promise.all([
        getProject(id),
        listAssessments(id),
        getTrends(id),
      ]);
      setProject(proj);
      setAssessments(asmts);
      setTrends(tr);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [id]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    try {
      const a = await createAssessment(id, form);
      navigate(`/assessments/${a.id}`);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (aId: number) => {
    if (!confirm("Delete this assessment?")) return;
    await deleteAssessment(aId);
    load();
  };

  if (loading) return <div className="text-gray-500">Loading…</div>;
  if (error) return <div className="text-red-500">Error: {error}</div>;
  if (!project) return <div className="text-gray-400">Project not found.</div>;

  return (
    <div>
      <div className="flex items-center gap-2 text-sm text-gray-500 mb-4">
        <Link to="/projects" className="hover:underline text-blue-600">
          Projects
        </Link>
        <span>/</span>
        <span className="text-gray-700">{project.name}</span>
      </div>

      <div className="bg-white border rounded-lg p-5 mb-6">
        <h1 className="text-xl font-bold text-gray-800">{project.name}</h1>
        {project.client_name && (
          <p className="text-sm text-gray-500 mt-1">
            Client: {project.client_name}
          </p>
        )}
        {project.owner && (
          <p className="text-sm text-gray-500">Owner: {project.owner}</p>
        )}
        {project.description && (
          <p className="text-sm text-gray-600 mt-2">{project.description}</p>
        )}
      </div>

      {/* Trend Chart */}
      {trends && trends.trend_points.length > 0 && (
        <div className="bg-white border rounded-lg p-5 mb-6">
          <h2 className="font-semibold text-gray-700 mb-3">
            Maturity Trend
          </h2>
          <TrendLineChart trendPoints={trends.trend_points} />
        </div>
      )}

      {/* Assessments */}
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-800">Assessments</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded"
        >
          + New Assessment
        </button>
      </div>

      {showForm && (
        <form
          onSubmit={handleCreate}
          className="bg-white border rounded-lg p-5 mb-4"
        >
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm text-gray-600 mb-1">
                Assessment Name *
              </label>
              <input
                required
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="Q3 2024 Assessment"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-600 mb-1">
                Assessor
              </label>
              <input
                value={form.assessor ?? ""}
                onChange={(e) =>
                  setForm({ ...form, assessor: e.target.value })
                }
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="Jane Smith"
              />
            </div>
            <div className="sm:col-span-2">
              <label className="block text-sm text-gray-600 mb-1">Scope</label>
              <input
                value={form.scope ?? ""}
                onChange={(e) => setForm({ ...form, scope: e.target.value })}
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="What systems / services are in scope?"
              />
            </div>
          </div>
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={creating}
              className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded disabled:opacity-50"
            >
              {creating ? "Creating…" : "Create & Start"}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="border text-sm px-4 py-2 rounded"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {assessments.length === 0 ? (
        <div className="text-sm text-gray-400">
          No assessments yet. Create one above.
        </div>
      ) : (
        <div className="grid gap-3">
          {assessments.map((a) => (
            <div
              key={a.id}
              className="bg-white border rounded-lg p-4 flex items-center justify-between"
            >
              <div>
                <Link
                  to={`/assessments/${a.id}`}
                  className="font-medium text-blue-700 hover:underline"
                >
                  {a.name}
                </Link>
                <div className="text-xs text-gray-400 mt-0.5">
                  {a.assessor && `Assessor: ${a.assessor} · `}
                  {a.assessment_date
                    ? new Date(a.assessment_date).toLocaleDateString()
                    : "No date"}
                  {" · "}
                  <span className="capitalize">{a.status}</span>
                </div>
              </div>
              <div className="flex gap-2">
                <Link
                  to={`/assessments/${a.id}/results`}
                  className="text-sm text-green-600 border border-green-200 px-3 py-1 rounded hover:bg-green-50"
                >
                  Results
                </Link>
                <Link
                  to={`/assessments/${a.id}`}
                  className="text-sm text-blue-600 border border-blue-200 px-3 py-1 rounded hover:bg-blue-50"
                >
                  Edit
                </Link>
                <button
                  onClick={() => handleDelete(a.id)}
                  className="text-sm text-red-500 border border-red-200 px-3 py-1 rounded hover:bg-red-50"
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
