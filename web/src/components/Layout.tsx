import { Command } from "lucide-react";
import type { ReactNode, RefObject } from "react";
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
	dark: boolean;
	onToggleTheme: () => void;
	searchInputRef: RefObject<HTMLInputElement | null>;
	onOpenCommandPalette: () => void;
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
	dark,
	onToggleTheme,
	searchInputRef,
	onOpenCommandPalette,
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
								<SearchBar ref={searchInputRef} value={search} onChange={onSearchChange} />
							</div>
							<button
								type="button"
								onClick={onOpenCommandPalette}
								className="hidden sm:flex items-center gap-1.5 px-2 py-1.5 rounded-lg border border-line bg-surface text-tx3 hover:text-tx2 hover:border-line-hover transition-colors text-[11px] font-mono"
								title="Command palette"
							>
								<Command className="w-3 h-3" />
								<span>K</span>
							</button>
							<NamespaceFilter namespaces={namespaces} value={namespace} onChange={onNamespaceChange} />
							<HealthFilter value={healthFilter} onChange={onHealthFilterChange} routes={allRoutes} />
							<ViewToggle view={view} onChange={onViewChange} />
							<ThemeToggle dark={dark} onToggle={onToggleTheme} />
						</div>
					</div>

					{/* Mobile search */}
					<div className="sm:hidden pb-3">
						<SearchBar ref={searchInputRef} value={search} onChange={onSearchChange} />
					</div>
				</div>
			</header>

			{/* Main content */}
			<main className="max-w-7xl mx-auto px-4 sm:px-6 py-8">{children}</main>
		</div>
	);
}
