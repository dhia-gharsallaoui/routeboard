interface SparklineProps {
	history: string[];
	width?: number;
	height?: number;
}

const colors: Record<string, string> = {
	healthy: "#34d399",
	degraded: "#f59e0b",
	unhealthy: "#ef4444",
	unknown: "#3d3d4e",
};

export function Sparkline({ history, width = 48, height = 12 }: SparklineProps) {
	if (!history || history.length < 2) return null;

	const barWidth = Math.max(1, width / history.length);
	const gap = 0.5;
	const healthyCount = history.filter((h) => h === "healthy").length;
	const uptime = Math.round((healthyCount / history.length) * 100);

	return (
		<span title={`${uptime}% uptime (last ${history.length} checks)`} className="inline-flex items-center">
			<svg width={width} height={height} className="rounded-sm overflow-hidden" role="img" aria-label="Uptime sparkline">
				{history.map((status, i) => (
					<rect
						key={`${i}-${status}`}
						x={i * barWidth + gap / 2}
						y={0}
						width={Math.max(0.5, barWidth - gap)}
						height={height}
						fill={colors[status] || colors.unknown}
						rx={0.5}
					/>
				))}
			</svg>
		</span>
	);
}
