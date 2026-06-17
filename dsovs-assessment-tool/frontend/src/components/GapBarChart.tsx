import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import type { ControlGap } from "../types";

interface Props {
  gaps: ControlGap[];
  limit?: number;
}

const colorForGap = (gap: number) => {
  if (gap >= 3) return "#ef4444";
  if (gap === 2) return "#f59e0b";
  if (gap === 1) return "#3b82f6";
  return "#10b981";
};

export default function GapBarChart({ gaps, limit = 10 }: Props) {
  const data = gaps
    .filter((g) => g.gap > 0)
    .slice(0, limit)
    .map((g) => ({
      name: g.code || g.title.slice(0, 20),
      Gap: g.gap,
      gap: g.gap,
    }));

  if (data.length === 0) {
    return <p className="text-gray-500 text-sm">No gaps found.</p>;
  }

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} layout="vertical">
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis type="number" domain={[0, 3]} tickCount={4} />
        <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 11 }} />
        <Tooltip />
        <Bar dataKey="Gap">
          {data.map((entry, i) => (
            <Cell key={i} fill={colorForGap(entry.gap)} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
