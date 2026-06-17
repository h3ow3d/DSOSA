import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Legend,
  ResponsiveContainer,
  Tooltip,
} from "recharts";
import type { PhaseScore } from "../types";

interface Props {
  phaseScores: PhaseScore[];
}

export default function PhaseRadarChart({ phaseScores }: Props) {
  const data = phaseScores.map((ps) => ({
    phase: ps.phase.length > 12 ? ps.phase.slice(0, 12) + "…" : ps.phase,
    Current: ps.current_score,
    Target: ps.target_score,
  }));

  return (
    <ResponsiveContainer width="100%" height={350}>
      <RadarChart data={data}>
        <PolarGrid />
        <PolarAngleAxis dataKey="phase" tick={{ fontSize: 11 }} />
        <PolarRadiusAxis angle={30} domain={[0, 3]} tickCount={4} />
        <Radar
          name="Current"
          dataKey="Current"
          stroke="#3b82f6"
          fill="#3b82f6"
          fillOpacity={0.3}
        />
        <Radar
          name="Target"
          dataKey="Target"
          stroke="#f59e0b"
          fill="#f59e0b"
          fillOpacity={0.15}
        />
        <Legend />
        <Tooltip />
      </RadarChart>
    </ResponsiveContainer>
  );
}
