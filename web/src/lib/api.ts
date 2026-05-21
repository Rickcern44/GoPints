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
	image_style: string;
	brewery_image_mime_type: string;
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

export interface KegSize {
	label: string;
	ml: number;
}

export const KEG_SIZES: KegSize[] = [
	{ label: '1/2 Barrel — Half Keg (15.5 gal / 58.7L)', ml: 58674 },
	{ label: '1/4 Barrel — Pony Keg (7.75 gal / 29.3L)', ml: 29337 },
	{ label: '1/6 Barrel — Torpedo Keg (5.2 gal / 19.5L)', ml: 19550 },
	{ label: 'Cornelius / Corny Keg (5 gal / 18.9L)', ml: 18927 },
	{ label: 'Mini Keg (1.3 gal / 5L)', ml: 5000 }
];

export function kegSizeLabel(ml: number): string {
	return KEG_SIZES.find((s) => s.ml === ml)?.label ?? `Custom (${(ml / 1000).toFixed(1)}L)`;
}

export interface Features {
	flow_based_pour: boolean;
}

export async function fetchFeatures(): Promise<Features> {
	const res = await fetch('/api/features');
	if (!res.ok) return { flow_based_pour: false };
	return res.json();
}

export async function setFeature(name: keyof Features, enabled: boolean): Promise<void> {
	const token = localStorage.getItem('admin_token') ?? '';
	const res = await fetch(`/api/features/${name}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
		body: JSON.stringify({ enabled })
	});
	if (!res.ok) throw new Error(`setFeature: ${res.status}`);
}

export async function recordManualPour(tapId: number, volumeMl: number): Promise<Pour> {
	const res = await fetch(`/api/taps/${tapId}/pour`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ volume_ml: volumeMl })
	});
	if (!res.ok) throw new Error(`recordManualPour: ${res.status}`);
	return res.json();
}

export async function fetchBanner(): Promise<string | null> {
	const res = await fetch('/api/banner');
	if (res.status === 404) return null;
	if (!res.ok) return null;
	const data: { message: string } = await res.json();
	return data.message ?? null;
}
