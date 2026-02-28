import { useState } from "react";
import { useRoutes } from "./hooks/useRoutes";
import { Layout } from "./components/Layout";
import { RouteGrid } from "./components/RouteGrid";
import { EmptyState } from "./components/EmptyState";

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

  if (loading) {
    return (
      <div className="min-h-screen bg-deep flex items-center justify-center">
        <div className="flex flex-col items-center gap-3 animate-fade-in">
          <div className="text-4xl">🧭</div>
          <p className="font-display font-semibold text-tx2 text-sm">
            Discovering routes...
          </p>
        </div>
      </div>
    );
  }

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
      {hasRoutes ? (
        <RouteGrid groupedRoutes={groupedRoutes} view={view} />
      ) : (
        <EmptyState searching={isSearching} />
      )}
    </Layout>
  );
}

export default App;
