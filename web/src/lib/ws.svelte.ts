import type { PourEvent } from './api.js';

class PourSocket {
	activePours = $state<Record<number, PourEvent>>({});
	lastEnded = $state<PourEvent | null>(null);

	private ws: WebSocket | null = null;

	connect() {
		this.ws?.close();
		const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
		this.ws = new WebSocket(`${protocol}//${location.host}/api/ws`);

		this.ws.onmessage = (e: MessageEvent) => {
			const event: PourEvent = JSON.parse(e.data as string);
			if (event.type === 'PourEnded') {
				const next = { ...this.activePours };
				delete next[event.tap_id];
				this.activePours = next;
				this.lastEnded = event;
			} else {
				this.activePours = { ...this.activePours, [event.tap_id]: event };
			}
		};
	}

	disconnect() {
		this.ws?.close();
		this.ws = null;
	}
}

export const pourSocket = new PourSocket();
