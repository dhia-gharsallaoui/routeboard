import type { Route } from "../types";
import { timeAgo } from "../utils/time";

interface HealthDotProps {
	health: Route["health"];
	checkedAt?: string;
	size?: "sm" | "md";
}

const config = {
	healthy: { color: "bg-emerald-500", label: "Healthy", pulse: true },
	degraded: { color: "bg-amber-500", label: "Degraded", pulse: false },
	unhealthy: { color: "bg-red-500", label: "Unhealthy", pulse: false },
	unknown: { color: "bg-tx3", label: "Unknown", pulse: false },
} as const;

export function HealthDot({ health, checkedAt, size = "sm" }: HealthDotProps) {
	const { color, label, pulse } = config[health] || config.unknown;
	const dotSize = size === "sm" ? "w-2 h-2" : "w-2.5 h-2.5";

	const ago = checkedAt ? timeAgo(checkedAt) : null;
	const tooltip = ago ? `${label} — checked ${ago}` : label;

	return (
		<span className="relative inline-flex" title={tooltip}>
			{pulse && <span className={`absolute inset-0 rounded-full ${color} opacity-40 animate-ping`} />}
			<span className={`relative inline-block rounded-full ${dotSize} ${color}`} />
		</span>
	);
}
