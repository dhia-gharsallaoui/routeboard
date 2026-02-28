import { useCallback, useEffect, useState } from "react";
import type { Route } from "../types";

interface UseKeyboardNavOptions {
	routes: Route[];
	disabled: boolean;
	searchInputRef: React.RefObject<HTMLInputElement | null>;
	onSelectRoute: (id: string) => void;
	onToggleFavorite: (id: string) => void;
}

export function useKeyboardNav({ routes, disabled, searchInputRef, onSelectRoute, onToggleFavorite }: UseKeyboardNavOptions) {
	const [focusedIndex, setFocusedIndex] = useState<number>(-1);

	// Reset focus when routes change
	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally reset on routes.length change
	useEffect(() => {
		setFocusedIndex(-1);
	}, [routes.length]);

	const focusedRouteId = focusedIndex >= 0 && focusedIndex < routes.length ? routes[focusedIndex].id : null;

	const handleKeyDown = useCallback(
		(e: KeyboardEvent) => {
			if (disabled) return;

			// Don't capture when in input/textarea
			const tag = (e.target as HTMLElement).tagName;
			if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") {
				if (e.key === "Escape") {
					(e.target as HTMLElement).blur();
					e.preventDefault();
				}
				return;
			}

			switch (e.key) {
				case "/":
					e.preventDefault();
					searchInputRef.current?.focus();
					break;
				case "j":
					e.preventDefault();
					setFocusedIndex((prev) => Math.min(prev + 1, routes.length - 1));
					break;
				case "k":
					e.preventDefault();
					setFocusedIndex((prev) => Math.max(prev - 1, 0));
					break;
				case "Enter":
					if (focusedIndex >= 0 && focusedIndex < routes.length) {
						e.preventDefault();
						onSelectRoute(routes[focusedIndex].id);
					}
					break;
				case "o":
					if (focusedIndex >= 0 && focusedIndex < routes.length) {
						const url = routes[focusedIndex].url;
						if (url) {
							e.preventDefault();
							window.open(url, "_blank", "noopener,noreferrer");
						}
					}
					break;
				case "f":
					if (focusedIndex >= 0 && focusedIndex < routes.length) {
						e.preventDefault();
						onToggleFavorite(routes[focusedIndex].id);
					}
					break;
				case "Escape":
					setFocusedIndex(-1);
					break;
			}
		},
		[disabled, routes, focusedIndex, searchInputRef, onSelectRoute, onToggleFavorite],
	);

	useEffect(() => {
		window.addEventListener("keydown", handleKeyDown);
		return () => window.removeEventListener("keydown", handleKeyDown);
	}, [handleKeyDown]);

	return { focusedRouteId };
}
