import { useState, useEffect, useCallback, useMemo } from "react";
import type { Route, Config } from "../types";
import { useSSE } from "./useSSE";

export function useRoutes() {
  const [allRoutes, setAllRoutes] = useState<Route[]>([]);
  const [config, setConfig] = useState<Config>({
    title: "RouteBoard",
    namespaces: [],
  });
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [namespace, setNamespace] = useState("");
  const [healthFilter, setHealthFilter] = useState("");

  const fetchRoutes = useCallback(async () => {
    try {
      const res = await fetch("/api/routes");
      const data = await res.json();
      setAllRoutes(data || []);
    } catch (err) {
      console.error("Failed to fetch routes:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchConfig = useCallback(async () => {
    try {
      const res = await fetch("/api/config");
      const data = await res.json();
      setConfig(data);
    } catch (err) {
      console.error("Failed to fetch config:", err);
    }
  }, []);

  useEffect(() => {
    fetchRoutes();
    fetchConfig();
  }, [fetchRoutes, fetchConfig]);

  const { connected } = useSSE("/api/events", () => {
    fetchRoutes();
    fetchConfig();
  });

  // Client-side filtering: search + namespace + health + hide no-URL routes
  const filteredRoutes = useMemo(() => {
    let routes = allRoutes.filter((r) => r.url);

    if (namespace) {
      routes = routes.filter((r) => r.namespace === namespace);
    }

    if (healthFilter) {
      routes = routes.filter((r) => r.health === healthFilter);
    }

    if (search) {
      const q = search.toLowerCase();
      routes = routes.filter(
        (r) =>
          r.title.toLowerCase().includes(q) ||
          r.url.toLowerCase().includes(q) ||
          r.description.toLowerCase().includes(q) ||
          r.name.toLowerCase().includes(q) ||
          r.namespace.toLowerCase().includes(q)
      );
    }

    return routes;
  }, [allRoutes, namespace, healthFilter, search]);

  const groupedRoutes = useMemo(() => {
    const groups: Record<string, Route[]> = {};
    for (const route of filteredRoutes) {
      const g = route.group || "default";
      if (!groups[g]) groups[g] = [];
      groups[g].push(route);
    }
    const sortedKeys = Object.keys(groups).sort();
    const result: Record<string, Route[]> = {};
    for (const key of sortedKeys) {
      result[key] = groups[key].sort((a, b) => {
        if (a.order !== b.order) return a.order - b.order;
        return a.title.localeCompare(b.title);
      });
    }
    return result;
  }, [filteredRoutes]);

  return {
    routes: filteredRoutes,
    allRoutes,
    groupedRoutes,
    config,
    loading,
    connected,
    search,
    setSearch,
    namespace,
    setNamespace,
    healthFilter,
    setHealthFilter,
  };
}
