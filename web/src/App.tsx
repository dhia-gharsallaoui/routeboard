import { useState } from "react";
import { useRoutes } from "./hooks/useRoutes";
import { Layout } from "./components/Layout";
import { RouteGrid } from "./components/RouteGrid";
import { EmptyState } from "./components/EmptyState";
import { SkeletonGrid } from "./components/Skeleton";

function App() {
  const {
    routes,
    groupedRoutes,
    config,
    loading,
    connected,
    search,
    setSearch,
    namespace,
    setNamespace,
  } = useRoutes();

  const [view, setView] = useState<"grid" | "list">("grid");

  const hasRoutes = Object.keys(groupedRoutes).length > 0;
  const isSearching = search !== "" || namespace !== "";

  return (
    <Layout
      title={config.title}
      connected={connected}
      routeCount={routes.length}
      namespaces={config.namespaces}
      search={search}
      onSearchChange={setSearch}
      namespace={namespace}
      onNamespaceChange={setNamespace}
      view={view}
      onViewChange={setView}
    >
      {loading ? (
        <SkeletonGrid />
      ) : hasRoutes ? (
        <RouteGrid groupedRoutes={groupedRoutes} view={view} />
      ) : (
        <EmptyState searching={isSearching} />
      )}
    </Layout>
  );
}

export default App;
