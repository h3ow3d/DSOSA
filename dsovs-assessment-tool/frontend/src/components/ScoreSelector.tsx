interface Props {
  value: number | null;
  onChange: (v: number | null) => void;
  disabled?: boolean;
}

const LEVELS = [0, 1, 2, 3];

export default function ScoreSelector({ value, onChange, disabled = false }: Props) {
  return (
    <div className="flex gap-1">
      {LEVELS.map((l) => (
        <button
          key={l}
          type="button"
          disabled={disabled}
          onClick={() => onChange(value === l ? null : l)}
          className={`w-8 h-8 rounded text-sm font-bold border-2 transition-all
            ${
              value === l
                ? "bg-blue-600 border-blue-700 text-white"
                : "bg-white border-gray-300 text-gray-600 hover:border-blue-400"
            }
            ${disabled ? "opacity-40 cursor-not-allowed" : "cursor-pointer"}
          `}
        >
          {l}
        </button>
      ))}
    </div>
  );
}
