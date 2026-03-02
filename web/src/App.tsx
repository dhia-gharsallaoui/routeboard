import { useState } from "react";
import { EmptyState } from "./components/EmptyState";
import { Layout } from "./components/Layout";
import { RouteGrid } from "./components/RouteGrid";
import { SkeletonGrid } from "./components/Skeleton";
import { useFavorites } from "./hooks/useFavorites";
import { useRoutes } from "./hooks/useRoutes";

function App() {
	const { routes, allRoutes, groupedRoutes, config, loading, connected, search, setSearch, namespace, setNamespace, healthFilter, setHealthFilter } =
		useRoutes();

	const { favorites, toggle, isFavorite } = useFavorites();
	const [view, setView] = useState<"grid" | "list">("grid");

	const hasRoutes = Object.keys(groupedRoutes).length > 0;
	const isSearching = search !== "" || namespace !== "" || healthFilter !== "";

	return (
		<Layout
			title={config.title}
			connected={connected}
			routeCount={routes.length}
			namespaces={config.namespaces}
			allRoutes={allRoutes}
			search={search}
			onSearchChange={setSearch}
			namespace={namespace}
			onNamespaceChange={setNamespace}
			healthFilter={healthFilter}
			onHealthFilterChange={setHealthFilter}
			view={view}
			onViewChange={setView}
		>
			{loading ? (
				<SkeletonGrid />
			) : hasRoutes ? (
				<RouteGrid
					groupedRoutes={groupedRoutes}
					view={view}
					allRoutes={allRoutes}
					isFavorite={isFavorite}
					onToggleFavorite={toggle}
					favorites={favorites}
				/>
			) : (
				<EmptyState searching={isSearching} />
			)}
		</Layout>
	);
}

export default App;
