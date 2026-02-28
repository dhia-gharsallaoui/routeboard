import { Pin } from "lucide-react";
import type { Route } from "../types";
import { RouteCard } from "./RouteCard";

interface RouteGridProps {
	groupedRoutes: Record<string, Route[]>;
	view: "grid" | "list";
	allRoutes: Route[];
	isFavorite: (id: string) => boolean;
	onToggleFavorite: (id: string) => void;
	favorites: Set<string>;
	onSelectRoute?: (id: string) => void;
	focusedRouteId?: string | null;
}

export function RouteGrid({
	groupedRoutes,
	view,
	allRoutes,
	isFavorite,
	onToggleFavorite,
	favorites,
	onSelectRoute,
	focusedRouteId,
}: RouteGridProps) {
	const groups = Object.entries(groupedRoutes);
	let globalIndex = 0;

	// Collect pinned routes
	const pinnedRoutes = allRoutes.filter((r) => favorites.has(r.id) && r.url);

	return (
		<div className="space-y-8">
			{/* Pinned section */}
			{pinnedRoutes.length > 0 && (
				<section className="animate-fade-in">
					<div className="flex items-center gap-3 mb-4">
						<Pin className="w-3.5 h-3.5 text-amber-400" />
						<h2 className="font-display font-semibold text-sm uppercase tracking-wider text-amber-400">Pinned</h2>
						<div className="flex-1 h-px bg-amber-400/20" />
						<span className="text-xs font-mono text-amber-400/60">{pinnedRoutes.length}</span>
					</div>
					<div className={view === "grid" ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3" : "flex flex-col gap-2"}>
						{pinnedRoutes.map((route) => {
							const idx = globalIndex++;
							return (
								<RouteCard
									key={`pinned-${route.id}`}
									route={route}
									index={idx}
									view={view}
									isFavorite={true}
									onToggleFavorite={onToggleFavorite}
									onSelect={onSelectRoute}
									isFocused={focusedRouteId === route.id}
								/>
							);
						})}
					</div>
				</section>
			)}

			{/* Regular groups */}
			{groups.map(([group, routes]) => (
				<section key={group} className="animate-fade-in">
					<div className="flex items-center gap-3 mb-4">
						<h2 className="font-display font-semibold text-sm uppercase tracking-wider text-tx3">{group}</h2>
						<div className="flex-1 h-px bg-line" />
						<span className="text-xs font-mono text-tx3">{routes.length}</span>
					</div>

					<div className={view === "grid" ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3" : "flex flex-col gap-2"}>
						{routes.map((route) => {
							const idx = globalIndex++;
							return (
								<RouteCard
									key={route.id}
									route={route}
									index={idx}
									view={view}
									isFavorite={isFavorite(route.id)}
									onToggleFavorite={onToggleFavorite}
									onSelect={onSelectRoute}
									isFocused={focusedRouteId === route.id}
								/>
							);
						})}
					</div>
				</section>
			))}
		</div>
	);
}
