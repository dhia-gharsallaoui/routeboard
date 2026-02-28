import { useEffect, useRef, useState } from "react";

export function useSSE(url: string, onEvent: (event: string, data: unknown) => void) {
	const [connected, setConnected] = useState(false);
	const onEventRef = useRef(onEvent);
	onEventRef.current = onEvent;

	useEffect(() => {
		let reconnectTimeout: ReturnType<typeof setTimeout>;
		let es: EventSource;

		function connect() {
			es = new EventSource(url);

			es.addEventListener("connected", () => {
				setConnected(true);
			});

			es.addEventListener("route-change", (e) => {
				try {
					const data = JSON.parse(e.data);
					onEventRef.current("route-change", data);
				} catch {
					/* ignore parse errors */
				}
			});

			es.onerror = () => {
				setConnected(false);
				es.close();
				reconnectTimeout = setTimeout(connect, 3000);
			};
		}

		connect();

		return () => {
			clearTimeout(reconnectTimeout);
			es?.close();
			setConnected(false);
		};
	}, [url]);

	return { connected };
}
