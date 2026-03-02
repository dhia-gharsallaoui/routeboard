import { useCallback, useState } from "react";

const STORAGE_KEY = "routeboard-favorites";

function loadFavorites(): Set<string> {
	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (stored) return new Set(JSON.parse(stored));
	} catch {
		/* ignore */
	}
	return new Set();
}

function saveFavorites(favorites: Set<string>) {
	localStorage.setItem(STORAGE_KEY, JSON.stringify([...favorites]));
}

export function useFavorites() {
	const [favorites, setFavorites] = useState<Set<string>>(loadFavorites);

	const toggle = useCallback((id: string) => {
		setFavorites((prev) => {
			const next = new Set(prev);
			if (next.has(id)) {
				next.delete(id);
			} else {
				next.add(id);
			}
			saveFavorites(next);
			return next;
		});
	}, []);

	const isFavorite = useCallback((id: string) => favorites.has(id), [favorites]);

	return { favorites, toggle, isFavorite };
}
