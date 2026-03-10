import { useEffect, useState } from "react";

export function useTheme() {
	const [dark, setDark] = useState(() => {
		if (typeof window === "undefined") return true;
		const saved = localStorage.getItem("routeboard-theme");
		if (saved) return saved === "dark";
		return true;
	});

	useEffect(() => {
		const root = document.documentElement;
		if (dark) {
			root.classList.add("dark");
		} else {
			root.classList.remove("dark");
		}
		localStorage.setItem("routeboard-theme", dark ? "dark" : "light");
	}, [dark]);

	const toggle = () => setDark((d) => !d);

	return { dark, toggle };
}
