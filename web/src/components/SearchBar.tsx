import { Search, X } from "lucide-react";
import React, { useState } from "react";

interface SearchBarProps {
	value: string;
	onChange: (value: string) => void;
}

export const SearchBar = React.forwardRef<HTMLInputElement, SearchBarProps>(({ value, onChange }, ref) => {
	const [focused, setFocused] = useState(false);

	return (
		<div className="relative">
			<Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-tx3" />
			<input
				ref={ref}
				type="text"
				value={value}
				onChange={(e) => onChange(e.target.value)}
				onFocus={() => setFocused(true)}
				onBlur={() => setFocused(false)}
				placeholder="Search routes..."
				className="w-full pl-10 pr-9 py-2 bg-surface border border-line rounded-lg text-sm text-tx1 placeholder-tx3 outline-none transition-colors focus:border-accent focus:ring-1 focus:ring-accent/20 font-body"
			/>
			{value ? (
				<button
					type="button"
					onClick={() => onChange("")}
					className="absolute right-3 top-1/2 -translate-y-1/2 text-tx3 hover:text-tx2 transition-colors"
				>
					<X className="w-4 h-4" />
				</button>
			) : (
				!focused && (
					<span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] font-mono text-tx3 border border-line rounded px-1 py-0.5">/</span>
				)
			)}
		</div>
	);
});
