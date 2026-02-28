import { useEffect, useState } from "react";
import { Globe, Database, Layers } from "lucide-react";

interface ServiceIconProps {
  serviceName?: string;
  resourceName: string;
  className?: string;
  size?: number;
}

// Module-level cache: slug → { path, color } (persists across renders)
interface CachedIcon {
  path: string;
  color?: string;
}
const svgCache = new Map<string, CachedIcon | null>();

export function ServiceIcon({
  serviceName,
  resourceName,
  className = "",
  size = 28,
}: ServiceIconProps) {
  const name = (serviceName || resourceName).toLowerCase();

  // 1. Try known slug mapping
  const slug = resolveSlug(name);
  if (slug) {
    return (
      <BrandIcon
        slug={slug}
        fallbackPath={fallbackPaths[slug]}
        size={size}
        className={className}
      />
    );
  }

  // 2. Category-based fallback icons
  if (/\b(db|postgres|mysql|mongo|mariadb)\b/.test(name)) {
    return <Database size={size} className={className} />;
  }
  if (/\b(seaweed|storage|s3|ceph)\b/.test(name)) {
    return <Layers size={size} className={className} />;
  }

  // 3. Try the service name itself as a CDN slug (strip common suffixes/prefixes)
  const guessedSlug = guessSlug(name);
  if (guessedSlug) {
    return (
      <BrandIcon
        slug={guessedSlug}
        size={size}
        className={className}
      />
    );
  }

  return <Globe size={size} className={className} />;
}

// Try to derive a Simple Icons slug from the service name.
// Strips common k8s suffixes like -server, -svc, -app, -web, etc.
function guessSlug(name: string): string | undefined {
  const cleaned = name
    .replace(/[-_](server|svc|service|app|web|ui|api|proxy|gateway|controller|operator|backend|frontend|master|main|primary)\b/g, "")
    .replace(/[^a-z0-9]/g, "")
    .trim();
  if (cleaned.length < 2) return undefined;
  return cleaned;
}

function BrandIcon({
  slug,
  fallbackPath,
  size,
  className,
}: {
  slug: string;
  fallbackPath?: string;
  size: number;
  className: string;
}) {
  const cached = svgCache.get(slug);
  const [icon, setIcon] = useState<CachedIcon | null>(
    cached ?? (fallbackPath ? { path: fallbackPath } : null)
  );

  useEffect(() => {
    if (svgCache.has(slug)) {
      const c = svgCache.get(slug);
      if (c) setIcon(c);
      return;
    }

    let cancelled = false;
    fetch(`https://cdn.simpleicons.org/${slug}`)
      .then((res) => {
        if (!res.ok) throw new Error("not found");
        return res.text();
      })
      .then((svgText) => {
        const pathMatch = svgText.match(/<path\s[^>]*d="([^"]+)"/);
        const fillMatch = svgText.match(/fill="(#[0-9A-Fa-f]{3,8})"/);
        if (pathMatch && !cancelled) {
          const entry: CachedIcon = {
            path: pathMatch[1],
            color: fillMatch?.[1],
          };
          svgCache.set(slug, entry);
          setIcon(entry);
        } else {
          svgCache.set(slug, null);
        }
      })
      .catch(() => {
        svgCache.set(slug, null);
      });

    return () => {
      cancelled = true;
    };
  }, [slug]);

  if (!icon) {
    return <Globe size={size} className={className} />;
  }

  return (
    <svg
      viewBox="0 0 24 24"
      width={size}
      height={size}
      className={className}
      fill={icon.color || "currentColor"}
    >
      <path d={icon.path} />
    </svg>
  );
}

// Map service name patterns to Simple Icons slugs
function resolveSlug(name: string): string | undefined {
  for (const [pattern, slug] of slugMap) {
    if (name.includes(pattern)) return slug;
  }
  return undefined;
}

// Pattern → Simple Icons slug
const slugMap: [string, string][] = [
  ["grafana", "grafana"],
  ["prometheus", "prometheus"],
  ["alertmanager", "prometheus"],
  ["argocd", "argo"],
  ["argo-cd", "argo"],
  ["jenkins", "jenkins"],
  ["gitlab", "gitlab"],
  ["gitea", "gitea"],
  ["vaultwarden", "vaultwarden"],
  ["vault", "vault"],
  ["traefik", "traefikproxy"],
  ["nginx", "nginx"],
  ["kibana", "kibana"],
  ["jaeger", "jaeger"],
  ["harbor", "harbor"],
  ["redis", "redis"],
  ["rabbitmq", "rabbitmq"],
  ["minio", "minio"],
  ["keycloak", "keycloak"],
  ["sonarqube", "sonarqube"],
  ["pgadmin", "postgresql"],
  ["postgres", "postgresql"],
  ["immich", "immich"],
  ["longhorn", "longhorn"],
  ["rancher", "rancher"],
  ["memos", "memos"],
  ["gatus", "gatus"],
  ["tailscale", "tailscale"],
  ["kubernetes", "kubernetes"],
  ["docker", "docker"],
  ["helm", "helm"],
  ["terraform", "terraform"],
  ["ansible", "ansible"],
  ["nextcloud", "nextcloud"],
  ["plex", "plex"],
  ["jellyfin", "jellyfin"],
  ["homeassistant", "homeassistant"],
  ["home-assistant", "homeassistant"],
  ["pihole", "pihole"],
  ["adguard", "adguard"],
  ["wireguard", "wireguard"],
  ["openvpn", "openvpn"],
  ["syncthing", "syncthing"],
  ["portainer", "portainer"],
  ["drone", "drone"],
  ["gitpod", "gitpod"],
  ["backstage", "backstage"],
  ["mysql", "mysql"],
  ["mariadb", "mariadb"],
  ["mongodb", "mongodb"],
  ["elasticsearch", "elasticsearch"],
  ["opensearch", "opensearch"],
  ["cassandra", "apachecassandra"],
  ["kafka", "apachekafka"],
  ["airflow", "apacheairflow"],
  ["superset", "apachesuperset"],
  ["metabase", "metabase"],
  ["n8n", "n8n"],
  ["authentik", "authentik"],
  ["woodpecker", "woodpeckerci"],
  ["forgejo", "forgejo"],
  ["trivy", "trivy"],
  ["loki", "grafana"],
  ["tempo", "grafana"],
  ["mimir", "grafana"],
];

// Hardcoded fallback paths — shown instantly while CDN fetch is in-flight.
// Only need fallbacks for the most common services to avoid layout shift.
const fallbackPaths: Record<string, string> = {
  grafana:
    "M23.02 10.59a8.578 8.578 0 0 0-.862-3.034 8.911 8.911 0 0 0-1.789-2.445c.337-1.342-.413-2.505-.413-2.505-1.292-.08-2.113.4-2.416.62-.052-.02-.102-.044-.154-.064-.22-.089-.446-.172-.677-.247a9.867 9.867 0 0 0-1.586-.358C14.557.753 12.94 0 12.94 0c-1.804 1.145-2.147 2.744-2.147 2.744l-.018.093a8.96 8.96 0 0 0-1.121.397 8.869 8.869 0 0 0-1.557.87l-.063-.029c-2.497-.955-4.716.195-4.716.195-.203 2.658.996 4.33 1.235 4.636a11.608 11.608 0 0 0-.607 2.635C1.636 12.677.953 15.014.953 15.014c1.926 2.214 4.171 2.351 4.171 2.351l.006-.005c.285.509.615.994.986 1.446.156.19.32.371.488.548-.704 2.009.099 3.68.099 3.68 2.144.08 3.553-.937 3.849-1.173a9.784 9.784 0 0 0 3.164.501h.345l.003.002c1.01 1.44 2.788 1.646 2.788 1.646 1.264-1.332 1.337-2.653 1.337-2.94v-.118c.265-.187.52-.387.758-.6a7.875 7.875 0 0 0 1.415-1.7c1.43.083 2.437-.885 2.437-.885-.236-1.49-1.085-2.216-1.264-2.354l-.065-.046c.008-.092.016-.18.02-.27.011-.162.016-.323.016-.48v-.486a6.215 6.215 0 0 0-2.099-4.103 6.015 6.015 0 0 0-3.222-1.46 6.292 6.292 0 0 0-.85-.048l-.318.013a4.777 4.777 0 0 0-3.335 1.695c-.332.4-.592.84-.768 1.297a4.594 4.594 0 0 0-.312 1.817l.016.255a3.615 3.615 0 0 0 .698 1.82 3.53 3.53 0 0 0 1.827 1.282c.33.098.66.14.971.137l.239-.008a2.634 2.634 0 0 0 1.275-.472.248.248 0 0 0 .039-.35.244.244 0 0 0-.309-.06 2.476 2.476 0 0 1-.836.254l-.166.008a2.59 2.59 0 0 1-.859-.229 2.52 2.52 0 0 1-1.472-1.913 2.306 2.306 0 0 1-.023-.567 3.163 3.163 0 0 1 1.357-2.35 2.946 2.946 0 0 1 1.713-.531h.227a4.041 4.041 0 0 1 1.635.49 3.94 3.94 0 0 1 1.602 1.662 3.77 3.77 0 0 1 .397 1.414l.01.226v.21a6.195 6.195 0 0 1-.088.813 5.31 5.31 0 0 1-.891 2.057 5.052 5.052 0 0 1-3.237 2.014 4.82 4.82 0 0 1-.975.069 6.607 6.607 0 0 1-1.716-.265 6.776 6.776 0 0 1-3.4-2.274 6.616 6.616 0 0 1-1.46-3.746l-.01-.348.003-.433a8.707 8.707 0 0 1 .334-2.033c.128-.444.286-.872.473-1.277a7.04 7.04 0 0 1 1.456-2.1c.293-.298.614-.565.953-.763a7.177 7.177 0 0 1 1.649-.77 8 8 0 0 1 2.265-.288 7.917 7.917 0 0 1 2.048.68 8.253 8.253 0 0 1 1.672 1.09l.179.155a8.671 8.671 0 0 1 1.735 2.302l.25.54a8.848 8.848 0 0 1 .45 1.34.186.186 0 0 0 .373-.042c.01-.246.002-.532-.024-.856z",
  vault:
    "M0 0l11.955 24L24 0zm13.366 4.827h1.393v1.38h-1.393zm-2.77 5.569H9.22V8.993h1.389zm0-2.087H9.22V6.906h1.389zm0-2.086H9.22V4.819h1.389zm2.087 6.263h-1.377V11.08h1.388zm0-2.09h-1.377V8.993h1.388zm0-2.087h-1.377V6.906h1.388zm0-2.086h-1.377V4.819h1.388zm.683.683h1.393v1.389h-1.393zm0 3.475V8.993h1.389v1.388Z",
  vaultwarden:
    "M0 0l11.955 24L24 0zm13.366 4.827h1.393v1.38h-1.393zm-2.77 5.569H9.22V8.993h1.389zm0-2.087H9.22V6.906h1.389zm0-2.086H9.22V4.819h1.389zm2.087 6.263h-1.377V11.08h1.388zm0-2.09h-1.377V8.993h1.388zm0-2.087h-1.377V6.906h1.388zm0-2.086h-1.377V4.819h1.388zm.683.683h1.393v1.389h-1.393zm0 3.475V8.993h1.389v1.388Z",
  redis:
    "M22.71 13.145c-1.66 2.092-3.452 4.483-7.038 4.483-3.203 0-4.397-2.825-4.48-5.12.701 1.484 2.073 2.685 4.214 2.63 4.117-.133 6.94-3.852 6.94-7.239 0-4.05-3.022-6.972-8.268-6.972-3.752 0-8.4 1.428-11.455 3.685C2.59 6.937 3.885 9.958 4.35 9.626c2.648-1.904 4.748-3.13 6.784-3.744C8.12 9.244.886 17.05 0 18.425c.1 1.261 1.66 4.648 2.424 4.648.232 0 .431-.133.664-.365a100.49 100.49 0 0 0 5.54-6.765c.222 3.104 1.748 6.898 6.014 6.898 3.819 0 7.604-2.756 9.33-8.965.2-.764-.73-1.361-1.261-.73zm-4.349-5.013c0 1.959-1.926 2.922-3.685 2.922-.941 0-1.664-.247-2.235-.568 1.051-1.592 2.092-3.225 3.21-4.973 1.972.334 2.71 1.43 2.71 2.619z",
  nginx:
    "M12 0L1.605 6v12L12 24l10.395-6V6L12 0zm6 16.59c0 .705-.646 1.29-1.529 1.29-.631 0-1.351-.255-1.801-.81l-6-7.141v6.66c0 .721-.57 1.29-1.274 1.29H7.32c-.721 0-1.29-.6-1.29-1.29V7.41c0-.705.63-1.29 1.5-1.29.646 0 1.38.255 1.83.81l5.97 7.141V7.41c0-.721.6-1.29 1.29-1.29h.075c.72 0 1.29.6 1.29 1.29v9.18H18z",
  gitlab:
    "m23.6004 9.5927-.0337-.0862L20.3.9814a.851.851 0 0 0-.3362-.405.8748.8748 0 0 0-.9997.0539.8748.8748 0 0 0-.29.4399l-2.2055 6.748H7.5375l-2.2057-6.748a.8573.8573 0 0 0-.29-.4412.8748.8748 0 0 0-.9997-.0537.8585.8585 0 0 0-.3362.4049L.4332 9.5015l-.0325.0862a6.0657 6.0657 0 0 0 2.0119 7.0105l.0113.0087.03.0213 4.976 3.7264 2.462 1.8633 1.4995 1.1321a1.0085 1.0085 0 0 0 1.2197 0l1.4995-1.1321 2.4619-1.8633 5.006-3.7489.0125-.01a6.0682 6.0682 0 0 0 2.0094-7.003z",
  rabbitmq:
    "M23.035 9.601h-7.677a.956.956 0 01-.962-.962V.962a.956.956 0 00-.962-.956H10.56a.956.956 0 00-.962.956V8.64a.956.956 0 01-.962.962H5.762a.956.956 0 01-.961-.962V.962A.956.956 0 003.839 0H.959a.956.956 0 00-.956.962v22.076A.956.956 0 00.965 24h22.07a.956.956 0 00.962-.962V10.58a.956.956 0 00-.962-.98zm-3.86 8.152a1.437 1.437 0 01-1.437 1.443h-1.924a1.437 1.437 0 01-1.436-1.443v-1.917a1.437 1.437 0 011.436-1.443h1.924a1.437 1.437 0 011.437 1.443z",
};
