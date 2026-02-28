export function timeAgo(dateStr: string): string {
	const seconds = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
	if (seconds < 10) return "just now";
	if (seconds < 60) return `${seconds}s ago`;
	const minutes = Math.floor(seconds / 60);
	if (minutes < 60) return `${minutes}m ago`;
	const hours = Math.floor(minutes / 60);
	return `${hours}h ago`;
}

export function formatDate(dateStr: string): string {
	return new Date(dateStr).toLocaleString(undefined, {
		dateStyle: "medium",
		timeStyle: "short",
	});
}
