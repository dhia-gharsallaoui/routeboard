import type { Route } from "../types";

interface RouteCardProps {
  route: Route;
  index: number;
  view: "grid" | "list";
}

export function RouteCard({ route, index, view }: RouteCardProps) {
  const truncatedUrl = route.url
    ? route.url.replace(/^https?:\/\//, "").replace(/\/$/, "")
    : "";

  if (view === "list") {
    return (
      <a
        href={route.url || "#"}
        target="_blank"
        rel="noopener noreferrer"
        className="animate-card-enter group flex items-center gap-4 px-4 py-3 bg-card border border-line rounded-lg hover:border-line-hover transition-all duration-200 hover:shadow-[var(--shadow-card-hover)]"
        style={{ animationDelay: `${index * 30}ms` }}
      >
        <span className="text-2xl flex-shrink-0 w-9 text-center">{route.icon}</span>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-medium text-tx1 truncate">{route.title}</span>
            {route.tls && (
              <svg className="w-3.5 h-3.5 text-success flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                <rect x="3" y="11" width="18" height="11" rx="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            )}
          </div>
          {route.description && (
            <p className="text-xs text-tx3 truncate mt-0.5">{route.description}</p>
          )}
        </div>

        <span className="hidden sm:block font-mono text-xs text-tx3 truncate max-w-[260px]">
          {truncatedUrl}
        </span>

        <div className="flex items-center gap-1.5 flex-shrink-0">
          <SourceBadge source={route.source} />
          <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-elevated text-tx3">
            {route.namespace}
          </span>
        </div>

        <svg className="w-4 h-4 text-tx3 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path d="M7 17L17 7M17 7H7M17 7v10" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </a>
    );
  }

  return (
    <a
      href={route.url || "#"}
      target="_blank"
      rel="noopener noreferrer"
      className="animate-card-enter group relative flex flex-col bg-card border border-line rounded-xl overflow-hidden transition-all duration-300 hover:border-line-hover hover:shadow-[var(--shadow-card-hover)] hover:-translate-y-0.5"
      style={{
        animationDelay: `${index * 50}ms`,
        boxShadow: "var(--shadow-card)",
      }}
    >
      {/* Top accent line */}
      <div className="h-[2px] bg-line transition-all duration-300 group-hover:bg-accent group-hover:shadow-[0_0_12px_var(--accent-glow)]" />

      <div className="p-5 flex flex-col flex-1">
        {/* Header: icon + title + TLS */}
        <div className="flex items-start gap-3.5 mb-3">
          <span className="text-3xl leading-none mt-0.5">{route.icon}</span>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <h3 className="font-display font-semibold text-tx1 truncate text-[15px] leading-tight">
                {route.title}
              </h3>
              {route.tls && (
                <svg className="w-3.5 h-3.5 text-success flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <rect x="3" y="11" width="18" height="11" rx="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
              )}
            </div>
            {route.description && (
              <p className="text-xs text-tx2 mt-1 line-clamp-2 leading-relaxed">
                {route.description}
              </p>
            )}
          </div>
        </div>

        {/* URL */}
        {truncatedUrl && (
          <div className="flex items-center gap-1.5 mb-4 mt-auto">
            <svg className="w-3 h-3 text-tx3 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" strokeLinecap="round" />
              <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" strokeLinecap="round" />
            </svg>
            <span className="font-mono text-xs text-tx3 truncate group-hover:text-accent transition-colors">
              {truncatedUrl}
            </span>
          </div>
        )}

        {/* Badges */}
        <div className="flex items-center gap-1.5 flex-wrap">
          <SourceBadge source={route.source} />
          <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-elevated text-tx3 border border-line">
            {route.namespace}
          </span>
        </div>
      </div>

      {/* Hover arrow */}
      <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
        <svg className="w-4 h-4 text-tx3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path d="M7 17L17 7M17 7H7M17 7v10" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </div>
    </a>
  );
}

function SourceBadge({ source }: { source: "Ingress" | "HTTPRoute" }) {
  const isIngress = source === "Ingress";
  return (
    <span
      className={`text-[10px] font-mono font-medium px-1.5 py-0.5 rounded ${
        isIngress
          ? "bg-accent-soft text-accent"
          : "bg-info-soft text-info"
      }`}
    >
      {source}
    </span>
  );
}
