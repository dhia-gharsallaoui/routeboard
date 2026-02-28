import type { Route } from "../types";
import { RouteCard } from "./RouteCard";

interface RouteGridProps {
  groupedRoutes: Record<string, Route[]>;
  view: "grid" | "list";
}

export function RouteGrid({ groupedRoutes, view }: RouteGridProps) {
  const groups = Object.entries(groupedRoutes);
  let globalIndex = 0;

  return (
    <div className="space-y-8">
      {groups.map(([group, routes]) => (
        <section key={group} className="animate-fade-in">
          <div className="flex items-center gap-3 mb-4">
            <h2 className="font-display font-semibold text-sm uppercase tracking-wider text-tx3">
              {group}
            </h2>
            <div className="flex-1 h-px bg-line" />
            <span className="text-xs font-mono text-tx3">
              {routes.length}
            </span>
          </div>

          <div
            className={
              view === "grid"
                ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3"
                : "flex flex-col gap-2"
            }
          >
            {routes.map((route) => {
              const idx = globalIndex++;
              return (
                <RouteCard
                  key={route.id}
                  route={route}
                  index={idx}
                  view={view}
                />
              );
            })}
          </div>
        </section>
      ))}
    </div>
  );
}
