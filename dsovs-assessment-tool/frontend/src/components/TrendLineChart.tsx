import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { TrendPoint } from "../types";

interface Props {
  trendPoints: TrendPoint[];
  showPhases?: boolean;
}

const PHASE_COLORS = [
  "#3b82f6",
  "#10b981",
  "#f59e0b",
  "#ef4444",
  "#8b5cf6",
  "#ec4899",
  "#06b6d4",
];

export default function TrendLineChart({ trendPoints, showPhases = false }: Props) {
  const data = trendPoints.map((tp) => ({
    name: tp.assessment_name,
    Overall: tp.overall_score,
    ...Object.fromEntries(
      Object.entries(tp.phase_scores).map(([k, v]) => [k, v])
    ),
  }));

  const phaseKeys = trendPoints.length
    ? Object.keys(trendPoints[0].phase_scores)
    : [];

  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="name" tick={{ fontSize: 11 }} />
        <YAxis domain={[0, 3]} tickCount={4} />
        <Tooltip />
        <Legend />
        <Line
          type="monotone"
          dataKey="Overall"
          stroke="#1d4ed8"
          strokeWidth={2}
          dot
        />
        {showPhases &&
          phaseKeys.map((phase, i) => (
            <Line
              key={phase}
              type="monotone"
              dataKey={phase}
              stroke={PHASE_COLORS[i % PHASE_COLORS.length]}
              strokeWidth={1}
              dot={false}
            />
          ))}
      </LineChart>
    </ResponsiveContainer>
  );
}
