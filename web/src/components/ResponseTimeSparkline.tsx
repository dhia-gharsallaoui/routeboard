interface ResponseTimeSparklineProps {
	history: number[];
	width?: number;
	height?: number;
}

export function ResponseTimeSparkline({ history, width = 200, height = 24 }: ResponseTimeSparklineProps) {
	if (!history || history.length < 2) return null;

	const max = Math.max(...history, 1);
	const min = Math.min(...history);
	const range = max - min || 1;

	const padding = 1;
	const chartH = height - padding * 2;
	const step = width / (history.length - 1);

	const points = history
		.map((v, i) => {
			const x = i * step;
			const y = padding + chartH - ((v - min) / range) * chartH;
			return `${x},${y}`;
		})
		.join(" ");

	const avg = Math.round(history.reduce((a, b) => a + b, 0) / history.length);

	return (
		<span title={`Avg: ${avg}ms / Min: ${min}ms / Max: ${max}ms`} className="inline-flex items-center">
			<svg width={width} height={height} className="overflow-visible" role="img" aria-label="Response time sparkline">
				<polyline points={points} fill="none" stroke="var(--accent)" strokeWidth={1.5} strokeLinejoin="round" strokeLinecap="round" />
			</svg>
		</span>
	);
}
