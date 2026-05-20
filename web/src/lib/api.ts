export interface Keg {
	id: number;
	beer_name: string;
	style: string;
	abv: number;
	brewery: string;
	notes: string;
	capacity_ml: number;
	added_at: string;
	image_mime_type: string;
}

export interface Tap {
	id: number;
	keg_id: number | null;
	keg: Keg | null;
}

export interface KegStats {
	keg_id: number;
	pour_count: number;
	poured_ml: number;
	remaining_ml: number;
	pct_remaining: number;
}

export interface Pour {
	id: number;
	tap_id: number;
	volume_ml: number;
	started_at: string;
	ended_at: string;
}

export type PourEventType = 'PourStarted' | 'PourUpdated' | 'PourEnded';

export interface PourEvent {
	type: PourEventType;
	tap_id: number;
	volume_ml: number;
	started_at: string;
	ended_at?: string;
}

export async function fetchTaps(): Promise<Tap[]> {
	const res = await fetch('/api/taps');
	if (!res.ok) throw new Error(`fetchTaps: ${res.status}`);
	return (await res.json()) ?? [];
}

export async function fetchKegStats(kegId: number): Promise<KegStats> {
	const res = await fetch(`/api/kegs/${kegId}/stats`);
	if (!res.ok) throw new Error(`fetchKegStats: ${res.status}`);
	return res.json();
}

export async function fetchBanner(): Promise<string | null> {
	const res = await fetch('/api/banner');
	if (res.status === 404) return null;
	if (!res.ok) return null;
	const data: { message: string } = await res.json();
	return data.message ?? null;
}
