import { Moon, Sun } from "lucide-react";

interface ThemeToggleProps {
	dark: boolean;
	onToggle: () => void;
}

export function ThemeToggle({ dark, onToggle }: ThemeToggleProps) {
	return (
		<button
			type="button"
			onClick={onToggle}
			className="p-2 rounded-lg border border-line bg-surface text-tx2 hover:text-tx1 hover:border-line-hover transition-colors"
			title={dark ? "Switch to light mode" : "Switch to dark mode"}
		>
			{dark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
		</button>
	);
}
