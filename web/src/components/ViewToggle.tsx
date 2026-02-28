interface ViewToggleProps {
  view: "grid" | "list";
  onChange: (view: "grid" | "list") => void;
}

export function ViewToggle({ view, onChange }: ViewToggleProps) {
  return (
    <div className="flex items-center bg-surface border border-line rounded-lg overflow-hidden">
      <button
        onClick={() => onChange("grid")}
        className={`p-2 transition-colors ${
          view === "grid"
            ? "bg-accent-soft text-accent"
            : "text-tx3 hover:text-tx2"
        }`}
        title="Grid view"
      >
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 16 16">
          <rect x="1" y="1" width="6" height="6" rx="1" />
          <rect x="9" y="1" width="6" height="6" rx="1" />
          <rect x="1" y="9" width="6" height="6" rx="1" />
          <rect x="9" y="9" width="6" height="6" rx="1" />
        </svg>
      </button>
      <div className="w-px h-5 bg-line" />
      <button
        onClick={() => onChange("list")}
        className={`p-2 transition-colors ${
          view === "list"
            ? "bg-accent-soft text-accent"
            : "text-tx3 hover:text-tx2"
        }`}
        title="List view"
      >
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 16 16">
          <rect x="1" y="2" width="14" height="2.5" rx="0.75" />
          <rect x="1" y="6.75" width="14" height="2.5" rx="0.75" />
          <rect x="1" y="11.5" width="14" height="2.5" rx="0.75" />
        </svg>
      </button>
    </div>
  );
}
