import { ArrowUpRight, Check, Copy, Link, Lock, Star } from "lucide-react";
import React, { useState } from "react";
import type { Route } from "../types";
import { HealthDot } from "./HealthDot";
import { ServiceIcon } from "./ServiceIcon";
import { Sparkline } from "./Sparkline";

interface RouteCardProps {
	route: Route;
	index: number;
	view: "grid" | "list";
	isFavorite: boolean;
	onToggleFavorite: (id: string) => void;
	onSelect?: (id: string) => void;
	isFocused?: boolean;
}

export function RouteCard({ route, index, view, isFavorite, onToggleFavorite, onSelect, isFocused }: RouteCardProps) {
	const [copied, setCopied] = useState(false);
	const cardRef = React.useRef<HTMLButtonElement>(null);

	React.useEffect(() => {
		if (isFocused && cardRef.current) {
			cardRef.current.scrollIntoView({ block: "nearest", behavior: "smooth" });
		}
	}, [isFocused]);

	const truncatedUrl = route.url ? route.url.replace(/^https?:\/\//, "").replace(/\/$/, "") : "";

	const handleCopy = (e: React.MouseEvent) => {
		e.preventDefault();
		e.stopPropagation();
		navigator.clipboard.writeText(route.url).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 1500);
		});
	};

	const handleStar = (e: React.MouseEvent) => {
		e.preventDefault();
		e.stopPropagation();
		onToggleFavorite(route.id);
	};

	const focusRing = isFocused ? "ring-2 ring-accent ring-offset-2 ring-offset-deep" : "";

	if (view === "list") {
		return (
			<button
				ref={cardRef}
				type="button"
				onClick={() => onSelect?.(route.id)}
				className={`animate-card-enter group flex items-center gap-4 px-4 py-3 bg-card border border-line rounded-lg hover:border-line-hover transition-all duration-200 hover:shadow-[var(--shadow-card-hover)] cursor-pointer text-left w-full ${focusRing}`}
				style={{ animationDelay: `${index * 30}ms` }}
			>
				<button
					type="button"
					onClick={handleStar}
					className={`flex-shrink-0 transition-colors ${isFavorite ? "text-amber-400" : "text-tx3 opacity-0 group-hover:opacity-100 hover:text-amber-400"}`}
					title={isFavorite ? "Unpin" : "Pin to top"}
				>
					<Star className="w-3.5 h-3.5" fill={isFavorite ? "currentColor" : "none"} />
				</button>

				<div className="flex-shrink-0 w-8 h-8 flex items-center justify-center text-tx2 group-hover:text-accent transition-colors">
					<ServiceIcon serviceName={route.serviceName} resourceName={route.name} size={20} />
				</div>

				<div className="flex-1 min-w-0">
					<div className="flex items-center gap-2">
						<span className="font-medium text-tx1 truncate text-sm">{route.title}</span>
						{route.tls && <Lock className="w-3 h-3 text-success flex-shrink-0" />}
						<HealthDot health={route.health} checkedAt={route.healthCheckedAt} />
						{route.healthHistory && <Sparkline history={route.healthHistory} />}
						<ResponseTimeBadge ms={route.responseTimeMs} />
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

				{route.url && (
					<a
						href={route.url}
						target="_blank"
						rel="noopener noreferrer"
						onClick={(e) => e.stopPropagation()}
						className="flex-shrink-0 p-1 text-tx3 opacity-0 group-hover:opacity-100 hover:text-accent transition-all"
						title="Open URL"
					>
						<ArrowUpRight className="w-4 h-4" />
					</a>
				)}
			</button>
		);
	}

	return (
		<button
			ref={cardRef}
			type="button"
			onClick={() => onSelect?.(route.id)}
			className={`animate-card-enter group relative flex flex-col bg-card border border-line rounded-xl overflow-hidden transition-all duration-300 hover:border-line-hover hover:shadow-[var(--shadow-card-hover)] hover:-translate-y-1 cursor-pointer text-left w-full ${focusRing}`}
			style={{
				animationDelay: `${index * 50}ms`,
				boxShadow: "var(--shadow-card)",
			}}
		>
			{/* Top accent line */}
			<div className="h-[2px] bg-gradient-to-r from-transparent via-line to-transparent transition-all duration-300 group-hover:via-accent group-hover:shadow-[0_0_16px_var(--accent-glow)]" />

			<div className="p-5 flex flex-col flex-1">
				{/* Header: icon + title + TLS + health */}
				<div className="flex items-start gap-3.5 mb-3">
					<div className="w-10 h-10 rounded-lg bg-elevated border border-line flex items-center justify-center text-tx2 group-hover:text-accent group-hover:border-accent/20 transition-all duration-300 flex-shrink-0">
						<ServiceIcon serviceName={route.serviceName} resourceName={route.name} size={22} />
					</div>
					<div className="flex-1 min-w-0">
						<div className="flex items-center gap-2">
							<h3 className="font-display font-semibold text-tx1 truncate text-[15px] leading-tight">{route.title}</h3>
							{route.tls && <Lock className="w-3.5 h-3.5 text-success flex-shrink-0" />}
							<HealthDot health={route.health} checkedAt={route.healthCheckedAt} size="md" />
							<ResponseTimeBadge ms={route.responseTimeMs} />
						</div>
						{route.description && <p className="text-xs text-tx2 mt-1 line-clamp-2 leading-relaxed">{route.description}</p>}
					</div>
				</div>

				{/* URL + sparkline */}
				<div className="flex items-center gap-2 mb-4 mt-auto">
					{truncatedUrl && (
						<div className="flex items-center gap-1.5 flex-1 min-w-0">
							<Link className="w-3 h-3 text-tx3 flex-shrink-0" />
							<span className="font-mono text-xs text-tx3 truncate group-hover:text-accent transition-colors">{truncatedUrl}</span>
						</div>
					)}
					{route.healthHistory && route.healthHistory.length > 1 && <Sparkline history={route.healthHistory} />}
				</div>

				{/* Badges */}
				<div className="flex items-center gap-1.5 flex-wrap">
					<SourceBadge source={route.source} />
					<span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-elevated text-tx3 border border-line">{route.namespace}</span>
				</div>
			</div>

			{/* Star + Copy + arrow on hover */}
			<div className="absolute top-4 right-4 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
				<button
					type="button"
					onClick={handleStar}
					className={`p-1.5 rounded-md bg-elevated/80 backdrop-blur-sm border border-line transition-colors ${isFavorite ? "text-amber-400" : "text-tx3 hover:text-amber-400"}`}
					title={isFavorite ? "Unpin" : "Pin to top"}
				>
					<Star className="w-3.5 h-3.5" fill={isFavorite ? "currentColor" : "none"} />
				</button>
				<button
					type="button"
					onClick={handleCopy}
					className="p-1.5 rounded-md bg-elevated/80 backdrop-blur-sm border border-line text-tx3 hover:text-accent hover:border-accent/30 transition-colors"
					title="Copy URL"
				>
					{copied ? <Check className="w-3.5 h-3.5 text-success" /> : <Copy className="w-3.5 h-3.5" />}
				</button>
				{route.url && (
					<a
						href={route.url}
						target="_blank"
						rel="noopener noreferrer"
						onClick={(e) => e.stopPropagation()}
						className="p-1.5 rounded-md bg-elevated/80 backdrop-blur-sm border border-line text-tx3 hover:text-accent hover:border-accent/30 transition-colors"
						title="Open URL"
					>
						<ArrowUpRight className="w-3.5 h-3.5" />
					</a>
				)}
			</div>
		</button>
	);
}

function ResponseTimeBadge({ ms }: { ms?: number }) {
	if (ms == null || ms <= 0) return null;
	const color = ms < 200 ? "text-success" : ms < 1000 ? "text-amber-500" : "text-danger";
	return (
		<span className={`font-mono text-[10px] ${color} flex-shrink-0`} title={`Response time: ${ms}ms`}>
			{ms}ms
		</span>
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
