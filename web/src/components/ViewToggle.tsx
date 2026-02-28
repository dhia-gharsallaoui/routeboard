import { LayoutGrid, List } from "lucide-react";

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
        <LayoutGrid className="w-4 h-4" />
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
        <List className="w-4 h-4" />
      </button>
    </div>
  );
}
