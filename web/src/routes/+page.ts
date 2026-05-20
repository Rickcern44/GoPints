import type { LoadEvent } from '@sveltejs/kit';
import type { Tap } from '$lib/api.js';

export const ssr = false;

export const load = async ({ fetch }: LoadEvent) => {
	const tapsRes = await fetch('/api/taps');
	if (!tapsRes.ok) return { taps: [], stats: [], banner: null };
	const taps: Tap[] = (await tapsRes.json()) ?? [];

	const tapsWithKegs = taps.filter(
		(t): t is typeof t & { keg: NonNullable<typeof t.keg> } => t.keg != null
	);
	const [stats, bannerRes] = await Promise.all([
		Promise.all(
			tapsWithKegs.map(async (t) => {
				const r = await fetch(`/api/kegs/${t.keg.id}/stats`);
				return r.ok ? r.json() : null;
			})
		),
		fetch('/api/banner')
	]);

	const banner = bannerRes.ok
		? (((await bannerRes.json()) as { message: string }).message ?? null)
		: null;

	return { taps, stats: stats.filter(Boolean), banner };
};
