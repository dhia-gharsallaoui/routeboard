import { ArrowUpRight, Check, Copy, Star, X } from "lucide-react";
import { useEffect, useState } from "react";
import type { Route } from "../types";
import { formatDate, timeAgo } from "../utils/time";
import { HealthDot } from "./HealthDot";
import { ResponseTimeSparkline } from "./ResponseTimeSparkline";
import { ServiceIcon } from "./ServiceIcon";
import { Sparkline } from "./Sparkline";

interface RouteDetailPanelProps {
	route: Route;
	isFavorite: boolean;
	onToggleFavorite: (id: string) => void;
	onClose: () => void;
}

export function RouteDetailPanel({ route, isFavorite, onToggleFavorite, onClose }: RouteDetailPanelProps) {
	const [copied, setCopied] = useState(false);

	useEffect(() => {
		const onKeyDown = (e: KeyboardEvent) => {
			if (e.key === "Escape") onClose();
		};
		window.addEventListener("keydown", onKeyDown);
		return () => window.removeEventListener("keydown", onKeyDown);
	}, [onClose]);

	const handleCopy = () => {
		navigator.clipboard.writeText(route.url).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 1500);
		});
	};

	const uptimePercent =
		route.healthHistory && route.healthHistory.length > 0
			? Math.round((route.healthHistory.filter((h) => h === "healthy").length / route.healthHistory.length) * 100)
			: null;

	const rtHistory = route.responseTimeHistory;
	const rtAvg = rtHistory && rtHistory.length > 0 ? Math.round(rtHistory.reduce((a, b) => a + b, 0) / rtHistory.length) : null;
	const rtMin = rtHistory && rtHistory.length > 0 ? Math.min(...rtHistory) : null;
	const rtMax = rtHistory && rtHistory.length > 0 ? Math.max(...rtHistory) : null;

	return (
		<>
			{/* biome-ignore lint/a11y/noStaticElementInteractions: backdrop overlay dismiss */}
			<div role="presentation" className="fixed inset-0 z-50 bg-deep/60 backdrop-blur-sm animate-overlay-fade-in" onClick={onClose} />

			{/* Panel */}
			<div className="fixed inset-y-0 right-0 z-50 w-[480px] max-w-full bg-card border-l border-line overflow-y-auto animate-slide-in-right">
				{/* Header */}
				<div className="sticky top-0 bg-card/90 backdrop-blur-sm border-b border-line px-6 py-4 flex items-center gap-3">
					<div className="w-10 h-10 rounded-lg bg-elevated border border-line flex items-center justify-center text-accent flex-shrink-0">
						<ServiceIcon serviceName={route.serviceName} resourceName={route.name} size={22} />
					</div>
					<div className="flex-1 min-w-0">
						<h2 className="font-display font-semibold text-tx1 truncate">{route.title}</h2>
						{route.description && <p className="text-xs text-tx2 truncate">{route.description}</p>}
					</div>
					<button type="button" onClick={onClose} className="p-1.5 rounded-lg text-tx3 hover:text-tx1 hover:bg-elevated transition-colors">
						<X className="w-5 h-5" />
					</button>
				</div>

				{/* Actions */}
				<div className="px-6 py-3 flex items-center gap-2 border-b border-line">
					{route.url && (
						<a
							href={route.url}
							target="_blank"
							rel="noopener noreferrer"
							className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-accent text-deep rounded-lg hover:opacity-90 transition-opacity"
						>
							<ArrowUpRight className="w-3.5 h-3.5" />
							Open
						</a>
					)}
					<button
						type="button"
						onClick={handleCopy}
						className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-elevated border border-line text-tx2 rounded-lg hover:text-tx1 hover:border-line-hover transition-colors"
					>
						{copied ? <Check className="w-3.5 h-3.5 text-success" /> : <Copy className="w-3.5 h-3.5" />}
						{copied ? "Copied" : "Copy URL"}
					</button>
					<button
						type="button"
						onClick={() => onToggleFavorite(route.id)}
						className={`inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-elevated border border-line rounded-lg transition-colors ${
							isFavorite ? "text-amber-400 border-amber-400/30" : "text-tx2 hover:text-amber-400"
						}`}
					>
						<Star className="w-3.5 h-3.5" fill={isFavorite ? "currentColor" : "none"} />
						{isFavorite ? "Pinned" : "Pin"}
					</button>
				</div>

				<div className="px-6 py-5 space-y-6">
					{/* Metadata grid */}
					<Section title="Details">
						<MetaRow label="Namespace" value={route.namespace} />
						<MetaRow label="Source" value={route.source} />
						{route.serviceName && <MetaRow label="Service" value={`${route.serviceName}${route.servicePort ? `:${route.servicePort}` : ""}`} />}
						<MetaRow label="Hosts" value={route.hosts?.join(", ") || "—"} />
						<MetaRow label="Paths" value={route.paths?.join(", ") || "/"} />
						<MetaRow label="TLS" value={route.tls ? "Yes" : "No"} />
						<MetaRow label="Created" value={formatDate(route.createdAt)} />
						<MetaRow label="Updated" value={formatDate(route.updatedAt)} />
					</Section>

					{/* Health */}
					<Section title="Health">
						<div className="flex items-center gap-3 mb-3">
							<HealthDot health={route.health} checkedAt={route.healthCheckedAt} size="md" />
							<span className="text-sm text-tx1 font-medium capitalize">{route.health}</span>
							{uptimePercent != null && <span className="text-xs font-mono text-tx3">{uptimePercent}% uptime</span>}
							{route.healthCheckedAt && <span className="text-xs text-tx3 ml-auto">{timeAgo(route.healthCheckedAt)}</span>}
						</div>
						{route.healthHistory && route.healthHistory.length > 1 && <Sparkline history={route.healthHistory} width={320} height={16} />}
					</Section>

					{/* Response time */}
					{route.responseTimeMs != null && route.responseTimeMs > 0 && (
						<Section title="Response Time">
							<div className="flex items-center gap-3 mb-3">
								<ResponseTimeBadgeLg ms={route.responseTimeMs} />
								{rtAvg != null && (
									<div className="flex gap-4 text-xs font-mono text-tx3">
										<span>avg {rtAvg}ms</span>
										<span>min {rtMin}ms</span>
										<span>max {rtMax}ms</span>
									</div>
								)}
							</div>
							{rtHistory && rtHistory.length > 1 && <ResponseTimeSparkline history={rtHistory} width={320} height={24} />}
						</Section>
					)}

					{/* Labels */}
					{route.labels && Object.keys(route.labels).length > 0 && (
						<Section title="Labels">
							<div className="flex flex-wrap gap-1.5">
								{Object.entries(route.labels).map(([k, v]) => (
									<span key={k} className="text-[10px] font-mono px-2 py-1 rounded bg-elevated text-tx2 border border-line">
										{k}={v}
									</span>
								))}
							</div>
						</Section>
					)}
				</div>
			</div>
		</>
	);
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
	return (
		<div>
			<h3 className="text-xs font-semibold uppercase tracking-wider text-tx3 mb-3">{title}</h3>
			{children}
		</div>
	);
}

function MetaRow({ label, value }: { label: string; value: string }) {
	return (
		<div className="flex items-baseline py-1.5 gap-3">
			<span className="text-xs text-tx3 w-20 flex-shrink-0">{label}</span>
			<span className="text-sm text-tx1 font-mono truncate">{value}</span>
		</div>
	);
}

function ResponseTimeBadgeLg({ ms }: { ms: number }) {
	const color = ms < 200 ? "text-success" : ms < 1000 ? "text-amber-500" : "text-danger";
	return <span className={`font-mono text-lg font-semibold ${color}`}>{ms}ms</span>;
}
