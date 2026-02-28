export interface Route {
  id: string;
  name: string;
  namespace: string;
  source: "Ingress" | "HTTPRoute";
  url: string;
  hosts: string[];
  paths: string[];
  tls: boolean;
  serviceName?: string;
  servicePort?: string;
  title: string;
  description: string;
  icon: string;
  group: string;
  order: number;
  hidden: boolean;
  labels?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface ChangeEvent {
  type: "added" | "updated" | "deleted";
  route: Route;
}

export interface Config {
  title: string;
  namespaces: string[];
}
