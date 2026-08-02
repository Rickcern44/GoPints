<script lang="ts">
	import { pourSocket } from '$lib/ws.svelte.js';
	import type { Tap, KegStats } from '$lib/api.js';

	let {
		taps,
		stats,
		customMessage
	}: {
		taps: Tap[];
		stats: Record<number, KegStats>;
		customMessage: string | null;
	} = $props();

	let activePours = $derived(Object.values(pourSocket.activePours));
	let lowKegs = $derived(taps.filter((t) => t.keg && (stats[t.keg.id]?.pct_remaining ?? 100) < 20));
	let hasBanners = $derived(activePours.length > 0 || lowKegs.length > 0 || !!customMessage);
</script>

{#if hasBanners}
	<div class="fixed top-0 right-0 left-0 z-50 flex flex-col gap-1 p-2 font-mono">
		{#each activePours as pour (pour.tap_id)}
			<div
				class="animate-pulse rounded-[3px] border border-signal-good bg-signal-good/15 px-4 py-3 text-center text-lg font-semibold text-signal-good backdrop-blur"
			>
				TAP {pour.tap_id} — POURING {pour.volume_ml.toFixed(0)} mL
			</div>
		{/each}

		{#each lowKegs as tap (tap.id)}
			<div
				class="rounded-[3px] border border-accent bg-accent/15 px-4 py-3 text-center text-lg font-semibold text-accent backdrop-blur"
			>
				TAP {tap.id} — LOW KEG ({(stats[tap.keg!.id]?.pct_remaining ?? 0).toFixed(0)}% REMAINING)
			</div>
		{/each}

		{#if customMessage}
			<div
				class="rounded-[3px] border-l-2 border-accent bg-panel/95 px-4 py-3 text-center text-lg font-semibold text-fg backdrop-blur"
			>
				{customMessage}
			</div>
		{/if}
	</div>
{/if}
