import { ChevronDown } from "lucide-react";

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
			<ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-tx3 pointer-events-none" />
		</div>
	);
}
