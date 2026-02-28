import { ArrowUpRight, Check, Copy, Link, Lock } from "lucide-react";
import { useState } from "react";
import type { Route } from "../types";
import { HealthDot } from "./HealthDot";
import { ServiceIcon } from "./ServiceIcon";

interface RouteCardProps {
	route: Route;
	index: number;
	view: "grid" | "list";
}

export function RouteCard({ route, index, view }: RouteCardProps) {
	const [copied, setCopied] = useState(false);

	const truncatedUrl = route.url ? route.url.replace(/^https?:\/\//, "").replace(/\/$/, "") : "";

	const handleCopy = (e: React.MouseEvent) => {
		e.preventDefault();
		e.stopPropagation();
		navigator.clipboard.writeText(route.url).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 1500);
		});
	};

	if (view === "list") {
		return (
			<a
				href={route.url || "#"}
				target="_blank"
				rel="noopener noreferrer"
				className="animate-card-enter group flex items-center gap-4 px-4 py-3 bg-card border border-line rounded-lg hover:border-line-hover transition-all duration-200 hover:shadow-[var(--shadow-card-hover)]"
				style={{ animationDelay: `${index * 30}ms` }}
			>
				<div className="flex-shrink-0 w-8 h-8 flex items-center justify-center text-tx2 group-hover:text-accent transition-colors">
					<ServiceIcon serviceName={route.serviceName} resourceName={route.name} size={20} />
				</div>

				<div className="flex-1 min-w-0">
					<div className="flex items-center gap-2">
						<span className="font-medium text-tx1 truncate text-sm">{route.title}</span>
						{route.tls && <Lock className="w-3 h-3 text-success flex-shrink-0" />}
						<HealthDot health={route.health} checkedAt={route.healthCheckedAt} />
					</div>
					{route.description && <p className="text-xs text-tx3 truncate mt-0.5">{route.description}</p>}
				</div>

				<span className="hidden md:block font-mono text-xs text-tx3 truncate max-w-[260px]">{truncatedUrl}</span>

				<div className="flex items-center gap-1.5 flex-shrink-0">
					<SourceBadge source={route.source} />
					<span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-elevated text-tx3 border border-line">{route.namespace}</span>
				</div>

				<button
					type="button"
					onClick={handleCopy}
					className="flex-shrink-0 p-1 text-tx3 opacity-0 group-hover:opacity-100 hover:text-accent transition-all"
					title="Copy URL"
				>
					{copied ? <Check className="w-3.5 h-3.5 text-success" /> : <Copy className="w-3.5 h-3.5" />}
				</button>

				<ArrowUpRight className="w-4 h-4 text-tx3 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" />
			</a>
		);
	}

	return (
		<a
			href={route.url || "#"}
			target="_blank"
			rel="noopener noreferrer"
			className="animate-card-enter group relative flex flex-col bg-card border border-line rounded-xl overflow-hidden transition-all duration-300 hover:border-line-hover hover:shadow-[var(--shadow-card-hover)] hover:-translate-y-1"
			style={{
				animationDelay: `${index * 50}ms`,
				boxShadow: "var(--shadow-card)",
			}}
		>
			{/* Top accent line */}
			<div className="h-[2px] bg-gradient-to-r from-transparent via-line to-transparent transition-all duration-300 group-hover:via-accent group-hover:shadow-[0_0_16px_var(--accent-glow)]" />

			<div className="p-5 flex flex-col flex-1">
				{/* Header: icon + title + TLS */}
				<div className="flex items-start gap-3.5 mb-3">
					<div className="w-10 h-10 rounded-lg bg-elevated border border-line flex items-center justify-center text-tx2 group-hover:text-accent group-hover:border-accent/20 transition-all duration-300 flex-shrink-0">
						<ServiceIcon serviceName={route.serviceName} resourceName={route.name} size={22} />
					</div>
					<div className="flex-1 min-w-0">
						<div className="flex items-center gap-2">
							<h3 className="font-display font-semibold text-tx1 truncate text-[15px] leading-tight">{route.title}</h3>
							{route.tls && <Lock className="w-3.5 h-3.5 text-success flex-shrink-0" />}
							<HealthDot health={route.health} checkedAt={route.healthCheckedAt} size="md" />
						</div>
						{route.description && <p className="text-xs text-tx2 mt-1 line-clamp-2 leading-relaxed">{route.description}</p>}
					</div>
				</div>

				{/* URL */}
				{truncatedUrl && (
					<div className="flex items-center gap-1.5 mb-4 mt-auto">
						<Link className="w-3 h-3 text-tx3 flex-shrink-0" />
						<span className="font-mono text-xs text-tx3 truncate group-hover:text-accent transition-colors">{truncatedUrl}</span>
					</div>
				)}

				{/* Badges */}
				<div className="flex items-center gap-1.5 flex-wrap">
					<SourceBadge source={route.source} />
					<span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-elevated text-tx3 border border-line">{route.namespace}</span>
				</div>
			</div>

			{/* Copy + arrow on hover */}
			<div className="absolute top-4 right-4 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
				<button
					type="button"
					onClick={handleCopy}
					className="p-1.5 rounded-md bg-elevated/80 backdrop-blur-sm border border-line text-tx3 hover:text-accent hover:border-accent/30 transition-colors"
					title="Copy URL"
				>
					{copied ? <Check className="w-3.5 h-3.5 text-success" /> : <Copy className="w-3.5 h-3.5" />}
				</button>
				<div className="p-1.5 rounded-md bg-elevated/80 backdrop-blur-sm border border-line text-tx3">
					<ArrowUpRight className="w-3.5 h-3.5" />
				</div>
			</div>
		</a>
	);
}

function SourceBadge({ source }: { source: "Ingress" | "HTTPRoute" }) {
	const isIngress = source === "Ingress";
	return (
		<span
			className={`text-[10px] font-mono font-medium px-1.5 py-0.5 rounded ${isIngress ? "bg-accent-soft text-accent" : "bg-info-soft text-info"}`}
		>
			{source}
		</span>
	);
}
