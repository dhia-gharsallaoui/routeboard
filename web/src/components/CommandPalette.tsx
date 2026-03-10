import { Command, Grid3X3, List, Moon, Search, Sun, X } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import type { Route } from "../types";
import { ServiceIcon } from "./ServiceIcon";

interface CommandPaletteProps {
	routes: Route[];
	dark: boolean;
	view: "grid" | "list";
	onSelectRoute: (id: string) => void;
	onToggleTheme: () => void;
	onChangeView: (view: "grid" | "list") => void;
	onFilterNamespace: (ns: string) => void;
	onFilterHealth: (h: string) => void;
	onClose: () => void;
}

interface PaletteItem {
	id: string;
	label: string;
	section: string;
	icon?: React.ReactNode;
	action: () => void;
}

export function CommandPalette({
	routes,
	dark,
	view,
	onSelectRoute,
	onToggleTheme,
	onChangeView,
	onFilterNamespace,
	onFilterHealth,
	onClose,
}: CommandPaletteProps) {
	const [query, setQuery] = useState("");
	const [selectedIndex, setSelectedIndex] = useState(0);
	const inputRef = useRef<HTMLInputElement>(null);
	const listRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		inputRef.current?.focus();
	}, []);

	const items = useMemo<PaletteItem[]>(() => {
		const all: PaletteItem[] = [];

		// Routes
		for (const r of routes) {
			all.push({
				id: `route:${r.id}`,
				label: r.title,
				section: "Routes",
				icon: <ServiceIcon serviceName={r.serviceName} resourceName={r.name} size={16} />,
				action: () => {
					onSelectRoute(r.id);
					onClose();
				},
			});
		}

		// Actions
		all.push({
			id: "action:theme",
			label: dark ? "Switch to light mode" : "Switch to dark mode",
			section: "Actions",
			icon: dark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />,
			action: () => {
				onToggleTheme();
				onClose();
			},
		});
		all.push({
			id: "action:view",
			label: view === "grid" ? "Switch to list view" : "Switch to grid view",
			section: "Actions",
			icon: view === "grid" ? <List className="w-4 h-4" /> : <Grid3X3 className="w-4 h-4" />,
			action: () => {
				onChangeView(view === "grid" ? "list" : "grid");
				onClose();
			},
		});
		all.push({
			id: "action:clear-filters",
			label: "Clear all filters",
			section: "Actions",
			icon: <X className="w-4 h-4" />,
			action: () => {
				onFilterNamespace("");
				onFilterHealth("");
				onClose();
			},
		});

		// Filters
		const namespaces = [...new Set(routes.map((r) => r.namespace))].sort();
		for (const ns of namespaces) {
			all.push({
				id: `filter:ns:${ns}`,
				label: `Namespace: ${ns}`,
				section: "Filters",
				action: () => {
					onFilterNamespace(ns);
					onClose();
				},
			});
		}
		for (const h of ["healthy", "degraded", "unhealthy"]) {
			all.push({
				id: `filter:health:${h}`,
				label: `Health: ${h}`,
				section: "Filters",
				action: () => {
					onFilterHealth(h);
					onClose();
				},
			});
		}

		return all;
	}, [routes, dark, view, onSelectRoute, onToggleTheme, onChangeView, onFilterNamespace, onFilterHealth, onClose]);

	const filtered = useMemo(() => {
		if (!query) return items.slice(0, 15);
		const q = query.toLowerCase();
		return items.filter((item) => item.label.toLowerCase().includes(q)).slice(0, 15);
	}, [items, query]);

	// Reset selection on filter change
	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally reset on filtered.length change
	useEffect(() => {
		setSelectedIndex(0);
	}, [filtered.length]);

	// Scroll selected into view
	useEffect(() => {
		const el = listRef.current?.children[selectedIndex] as HTMLElement | undefined;
		el?.scrollIntoView({ block: "nearest" });
	}, [selectedIndex]);

	useEffect(() => {
		const onKeyDown = (e: KeyboardEvent) => {
			switch (e.key) {
				case "ArrowDown":
					e.preventDefault();
					setSelectedIndex((i) => Math.min(i + 1, filtered.length - 1));
					break;
				case "ArrowUp":
					e.preventDefault();
					setSelectedIndex((i) => Math.max(i - 1, 0));
					break;
				case "Enter":
					e.preventDefault();
					filtered[selectedIndex]?.action();
					break;
				case "Escape":
					e.preventDefault();
					onClose();
					break;
			}
		};
		window.addEventListener("keydown", onKeyDown);
		return () => window.removeEventListener("keydown", onKeyDown);
	}, [filtered, selectedIndex, onClose]);

	// Group by section
	const sections: { name: string; items: (PaletteItem & { globalIdx: number })[] }[] = [];
	let idx = 0;
	for (const item of filtered) {
		const last = sections[sections.length - 1];
		if (last && last.name === item.section) {
			last.items.push({ ...item, globalIdx: idx });
		} else {
			sections.push({ name: item.section, items: [{ ...item, globalIdx: idx }] });
		}
		idx++;
	}

	return (
		<>
			{/* biome-ignore lint/a11y/noStaticElementInteractions: backdrop overlay dismiss */}
			<div role="presentation" className="fixed inset-0 z-[60] bg-deep/60 backdrop-blur-sm animate-overlay-fade-in" onClick={onClose} />
			<div className="fixed inset-0 z-[60] flex items-start justify-center pt-[20vh] pointer-events-none">
				<div className="pointer-events-auto w-[520px] max-w-[calc(100vw-2rem)] bg-card border border-line rounded-xl shadow-2xl overflow-hidden animate-card-enter">
					{/* Search input */}
					<div className="flex items-center gap-3 px-4 py-3 border-b border-line">
						<Command className="w-4 h-4 text-tx3 flex-shrink-0" />
						<input
							ref={inputRef}
							type="text"
							value={query}
							onChange={(e) => setQuery(e.target.value)}
							placeholder="Type a command or search..."
							className="flex-1 bg-transparent text-sm text-tx1 placeholder-tx3 outline-none font-body"
						/>
						{query && (
							<button type="button" onClick={() => setQuery("")} className="text-tx3 hover:text-tx2 transition-colors">
								<X className="w-4 h-4" />
							</button>
						)}
					</div>

					{/* Results */}
					<div ref={listRef} className="max-h-[320px] overflow-y-auto py-2">
						{filtered.length === 0 && <p className="px-4 py-6 text-center text-sm text-tx3">No results found</p>}
						{sections.map((section) => (
							<div key={section.name}>
								<div className="px-4 py-1.5">
									<span className="text-[10px] font-semibold uppercase tracking-wider text-tx3">{section.name}</span>
								</div>
								{section.items.map((item) => (
									<button
										key={item.id}
										type="button"
										onClick={item.action}
										className={`w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors ${
											item.globalIdx === selectedIndex ? "bg-accent/10 text-accent" : "text-tx1 hover:bg-elevated"
										}`}
									>
										{item.icon && <span className="flex-shrink-0 text-tx2 w-5 h-5 flex items-center justify-center">{item.icon}</span>}
										{!item.icon && <Search className="w-4 h-4 text-tx3 flex-shrink-0" />}
										<span className="truncate">{item.label}</span>
									</button>
								))}
							</div>
						))}
					</div>

					{/* Footer hints */}
					<div className="flex items-center gap-4 px-4 py-2 border-t border-line text-[10px] font-mono text-tx3">
						<span>
							<kbd className="px-1 py-0.5 rounded border border-line bg-elevated">↑↓</kbd> navigate
						</span>
						<span>
							<kbd className="px-1 py-0.5 rounded border border-line bg-elevated">↵</kbd> select
						</span>
						<span>
							<kbd className="px-1 py-0.5 rounded border border-line bg-elevated">esc</kbd> close
						</span>
					</div>
				</div>
			</div>
		</>
	);
}
