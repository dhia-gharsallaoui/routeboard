interface NamespaceFilterProps {
  namespaces: string[];
  value: string;
  onChange: (value: string) => void;
}

export function NamespaceFilter({ namespaces, value, onChange }: NamespaceFilterProps) {
  if (namespaces.length === 0) return null;

  return (
    <div className="relative">
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="appearance-none bg-surface border border-line rounded-lg text-sm text-tx1 py-2 pl-3 pr-9 outline-none transition-colors focus:border-accent focus:ring-1 focus:ring-accent/20 cursor-pointer font-body"
      >
        <option value="">All namespaces</option>
        {namespaces.map((ns) => (
          <option key={ns} value={ns}>
            {ns}
          </option>
        ))}
      </select>
      <svg
        className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-tx3 pointer-events-none"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path d="M6 9l6 6 6-6" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    </div>
  );
}
