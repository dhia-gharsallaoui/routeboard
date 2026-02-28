import { Heart } from "lucide-react";
import type { Route } from "../types";

interface HealthFilterProps {
	value: string;
	onChange: (value: string) => void;
	routes: Route[];
}

const statuses = [
	{ value: "healthy", label: "Healthy", color: "bg-emerald-500" },
	{ value: "degraded", label: "Degraded", color: "bg-amber-500" },
	{ value: "unhealthy", label: "Unhealthy", color: "bg-red-500" },
] as const;

export function HealthFilter({ value, onChange, routes }: HealthFilterProps) {
	const counts = {
		healthy: routes.filter((r) => r.health === "healthy").length,
		degraded: routes.filter((r) => r.health === "degraded").length,
		unhealthy: routes.filter((r) => r.health === "unhealthy").length,
	};

	// Only show if there's at least one health check result
	const hasAny = counts.healthy + counts.degraded + counts.unhealthy > 0;
	if (!hasAny) return null;

	return (
		<div className="flex items-center gap-1 bg-surface border border-line rounded-lg overflow-hidden">
			<button
				type="button"
				onClick={() => onChange("")}
				className={`p-2 transition-colors ${value === "" ? "bg-accent-soft text-accent" : "text-tx3 hover:text-tx2"}`}
				title="All statuses"
			>
				<Heart className="w-4 h-4" />
			</button>
			{statuses.map((s) =>
				counts[s.value] > 0 ? (
					<button
						type="button"
						key={s.value}
						onClick={() => onChange(value === s.value ? "" : s.value)}
						className={`flex items-center gap-1.5 px-2 py-2 text-xs font-mono transition-colors ${
							value === s.value ? "bg-accent-soft text-accent" : "text-tx3 hover:text-tx2"
						}`}
						title={`Show ${s.label.toLowerCase()} only`}
					>
						<span className={`w-2 h-2 rounded-full ${s.color}`} />
						<span>{counts[s.value]}</span>
					</button>
				) : null,
			)}
		</div>
	);
}
