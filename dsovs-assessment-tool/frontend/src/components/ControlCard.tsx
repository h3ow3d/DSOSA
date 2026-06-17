import { useCallback, useEffect, useRef, useState } from "react";
import { upsertScore } from "../api/client";
import type { Control, Score, ScoreUpsert } from "../types";
import ScoreSelector from "./ScoreSelector";

interface Props {
  control: Control;
  assessmentId: number;
  initialScore: Score | undefined;
}

const PRIORITIES = ["critical", "high", "medium", "low"];
const CONFIDENCES = ["high", "medium", "low"];

export default function ControlCard({ control, assessmentId, initialScore }: Props) {
  const [score, setScore] = useState<ScoreUpsert>({
    current_level: initialScore?.current_level ?? null,
    target_level: initialScore?.target_level ?? null,
    not_applicable: initialScore?.not_applicable ?? false,
    confidence: initialScore?.confidence ?? null,
    priority: initialScore?.priority ?? null,
    evidence_notes: initialScore?.evidence_notes ?? "",
    action_notes: initialScore?.action_notes ?? "",
  });
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const save = useCallback(
    async (data: ScoreUpsert) => {
      setSaving(true);
      try {
        await upsertScore(assessmentId, control.id, data);
        setSaved(true);
        setTimeout(() => setSaved(false), 1500);
      } finally {
        setSaving(false);
      }
    },
    [assessmentId, control.id]
  );

  const update = (patch: Partial<ScoreUpsert>) => {
    const next = { ...score, ...patch };
    setScore(next);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => save(next), 600);
  };

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  const levels = [...(control.maturity_levels || [])].sort(
    (a, b) => a.level - b.level
  );

  return (
    <div
      className={`border rounded-lg p-4 bg-white shadow-sm ${
        score.not_applicable ? "opacity-60" : ""
      }`}
    >
      <div className="flex items-start justify-between gap-4 mb-2">
        <div>
          <span className="text-xs font-mono text-gray-400 mr-2">
            {control.code}
          </span>
          <span className="font-semibold text-gray-800">{control.title}</span>
          {control.doc_url && (
            <a
              href={control.doc_url}
              target="_blank"
              rel="noopener noreferrer"
              className="ml-2 text-xs text-blue-500 hover:underline"
            >
              ↗ docs
            </a>
          )}
        </div>
        <label className="flex items-center gap-1 text-xs text-gray-500 whitespace-nowrap">
          <input
            type="checkbox"
            checked={score.not_applicable}
            onChange={(e) => update({ not_applicable: e.target.checked })}
          />
          N/A
        </label>
      </div>

      {control.summary && (
        <p className="text-sm text-gray-600 mb-3">{control.summary}</p>
      )}

      {levels.length > 0 && (
        <div className="mb-3 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2">
          {levels.map((lvl) => (
            <div
              key={lvl.id}
              className={`text-xs border rounded p-2 ${
                score.current_level === lvl.level
                  ? "border-blue-400 bg-blue-50"
                  : "border-gray-200 bg-gray-50"
              }`}
            >
              <div className="font-semibold mb-1">
                L{lvl.level}
                {lvl.title ? `: ${lvl.title}` : ""}
              </div>
              {lvl.description && (
                <div className="text-gray-500">{lvl.description}</div>
              )}
            </div>
          ))}
        </div>
      )}

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-3">
        <div>
          <div className="text-xs text-gray-500 mb-1">Current Level</div>
          <ScoreSelector
            value={score.current_level ?? null}
            onChange={(v) => update({ current_level: v })}
            disabled={score.not_applicable}
          />
        </div>
        <div>
          <div className="text-xs text-gray-500 mb-1">Target Level</div>
          <ScoreSelector
            value={score.target_level ?? null}
            onChange={(v) => update({ target_level: v })}
            disabled={score.not_applicable}
          />
        </div>
        <div>
          <div className="text-xs text-gray-500 mb-1">Priority</div>
          <select
            value={score.priority ?? ""}
            onChange={(e) => update({ priority: e.target.value || null })}
            disabled={score.not_applicable}
            className="text-xs border rounded px-2 py-1 w-full"
          >
            <option value="">Auto</option>
            {PRIORITIES.map((p) => (
              <option key={p} value={p}>
                {p}
              </option>
            ))}
          </select>
        </div>
        <div>
          <div className="text-xs text-gray-500 mb-1">Confidence</div>
          <select
            value={score.confidence ?? ""}
            onChange={(e) => update({ confidence: e.target.value || null })}
            disabled={score.not_applicable}
            className="text-xs border rounded px-2 py-1 w-full"
          >
            <option value="">—</option>
            {CONFIDENCES.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <div className="text-xs text-gray-500 mb-1">Evidence Notes</div>
          <textarea
            rows={2}
            value={score.evidence_notes ?? ""}
            onChange={(e) => update({ evidence_notes: e.target.value })}
            disabled={score.not_applicable}
            className="text-xs border rounded px-2 py-1 w-full resize-none"
            placeholder="What evidence supports this score?"
          />
        </div>
        <div>
          <div className="text-xs text-gray-500 mb-1">Action Notes</div>
          <textarea
            rows={2}
            value={score.action_notes ?? ""}
            onChange={(e) => update({ action_notes: e.target.value })}
            disabled={score.not_applicable}
            className="text-xs border rounded px-2 py-1 w-full resize-none"
            placeholder="What actions will improve this?"
          />
        </div>
      </div>

      <div className="text-right text-xs text-gray-400 mt-1 h-4">
        {saving && "Saving…"}
        {saved && !saving && (
          <span className="text-green-500">✓ Saved</span>
        )}
      </div>
    </div>
  );
}
