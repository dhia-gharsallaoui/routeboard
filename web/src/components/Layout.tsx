import type { ReactNode } from "react";
import type { Route } from "../types";
import { HealthFilter } from "./HealthFilter";
import { NamespaceFilter } from "./NamespaceFilter";
import { SearchBar } from "./SearchBar";
import { ThemeToggle } from "./ThemeToggle";
import { ViewToggle } from "./ViewToggle";

interface LayoutProps {
	title: string;
	connected: boolean;
	routeCount: number;
	namespaces: string[];
	allRoutes: Route[];
	search: string;
	onSearchChange: (value: string) => void;
	namespace: string;
	onNamespaceChange: (value: string) => void;
	healthFilter: string;
	onHealthFilterChange: (value: string) => void;
	view: "grid" | "list";
	onViewChange: (view: "grid" | "list") => void;
	children: ReactNode;
}

export function Layout({
	title,
	connected,
	routeCount,
	namespaces,
	allRoutes,
	search,
	onSearchChange,
	namespace,
	onNamespaceChange,
	healthFilter,
	onHealthFilterChange,
	view,
	onViewChange,
	children,
}: LayoutProps) {
	return (
		<div className="min-h-screen bg-deep">
			{/* Header */}
			<header className="sticky top-0 z-40 bg-deep/80 backdrop-blur-xl border-b border-line">
				<div className="max-w-7xl mx-auto px-4 sm:px-6">
					{/* Top row: branding + controls */}
					<div className="flex items-center justify-between h-16 gap-4">
						{/* Left: Logo + status */}
						<div className="flex items-center gap-3 flex-shrink-0">
							<h1 className="font-display font-bold text-xl text-tx1 tracking-tight">{title}</h1>
							<div className="flex items-center gap-1.5">
								<div
									className={`w-1.5 h-1.5 rounded-full ${connected ? "bg-success animate-pulse-dot" : "bg-tx3"}`}
									title={connected ? "Connected — live updates active" : "Disconnected"}
								/>
								<span className="text-[11px] font-mono text-tx3 hidden sm:inline">
									{routeCount} {routeCount === 1 ? "route" : "routes"}
								</span>
							</div>
						</div>

						{/* Right: controls */}
						<div className="flex items-center gap-2">
							<div className="hidden sm:block w-64">
								<SearchBar value={search} onChange={onSearchChange} />
							</div>
							<NamespaceFilter namespaces={namespaces} value={namespace} onChange={onNamespaceChange} />
							<HealthFilter value={healthFilter} onChange={onHealthFilterChange} routes={allRoutes} />
							<ViewToggle view={view} onChange={onViewChange} />
							<ThemeToggle />
						</div>
					</div>

					{/* Mobile search */}
					<div className="sm:hidden pb-3">
						<SearchBar value={search} onChange={onSearchChange} />
					</div>
				</div>
			</header>

			{/* Main content */}
			<main className="max-w-7xl mx-auto px-4 sm:px-6 py-8">{children}</main>
		</div>
	);
}
