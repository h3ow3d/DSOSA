import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { createProject, deleteProject, listProjects } from "../api/client";
import type { Project, ProjectCreate } from "../types";

export default function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<ProjectCreate>({
    name: "",
    client_name: "",
    owner: "",
    description: "",
  });
  const [creating, setCreating] = useState(false);

  const load = () => {
    setLoading(true);
    listProjects()
      .then(setProjects)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    try {
      await createProject(form);
      setForm({ name: "", client_name: "", owner: "", description: "" });
      setShowForm(false);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Delete this project and all its assessments?")) return;
    await deleteProject(id);
    load();
  };

  if (loading) return <div className="text-gray-500">Loading…</div>;
  if (error) return <div className="text-red-500">Error: {error}</div>;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Projects</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded"
        >
          + New Project
        </button>
      </div>

      {showForm && (
        <form
          onSubmit={handleCreate}
          className="bg-white border rounded-lg p-5 mb-6"
        >
          <h2 className="font-semibold text-gray-700 mb-4">Create Project</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm text-gray-600 mb-1">
                Name *
              </label>
              <input
                required
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="My Application"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-600 mb-1">
                Client Name
              </label>
              <input
                value={form.client_name ?? ""}
                onChange={(e) =>
                  setForm({ ...form, client_name: e.target.value })
                }
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="Acme Corp"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-600 mb-1">Owner</label>
              <input
                value={form.owner ?? ""}
                onChange={(e) => setForm({ ...form, owner: e.target.value })}
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="Jane Smith"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-600 mb-1">
                Description
              </label>
              <input
                value={form.description ?? ""}
                onChange={(e) =>
                  setForm({ ...form, description: e.target.value })
                }
                className="w-full border rounded px-3 py-2 text-sm"
                placeholder="Optional description"
              />
            </div>
          </div>
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={creating}
              className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded disabled:opacity-50"
            >
              {creating ? "Creating…" : "Create Project"}
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

      {projects.length === 0 ? (
        <div className="text-gray-500 text-sm">
          No projects yet. Create one to get started.
        </div>
      ) : (
        <div className="grid gap-4">
          {projects.map((p) => (
            <div
              key={p.id}
              className="bg-white border rounded-lg p-4 flex items-center justify-between"
            >
              <div>
                <Link
                  to={`/projects/${p.id}`}
                  className="font-semibold text-blue-700 hover:underline text-lg"
                >
                  {p.name}
                </Link>
                {p.client_name && (
                  <span className="ml-2 text-sm text-gray-500">
                    · {p.client_name}
                  </span>
                )}
                {p.owner && (
                  <span className="ml-2 text-sm text-gray-400">
                    · {p.owner}
                  </span>
                )}
                {p.description && (
                  <p className="text-sm text-gray-500 mt-1">{p.description}</p>
                )}
              </div>
              <div className="flex gap-2">
                <Link
                  to={`/projects/${p.id}`}
                  className="text-sm text-blue-600 border border-blue-200 px-3 py-1 rounded hover:bg-blue-50"
                >
                  View
                </Link>
                <button
                  onClick={() => handleDelete(p.id)}
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
