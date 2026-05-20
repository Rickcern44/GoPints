<script lang="ts">
	import { onMount } from 'svelte';
	import { pourSocket } from '$lib/ws.svelte.js';
	import KegCard from '$lib/components/KegCard.svelte';
	import BannerStack from '$lib/components/BannerStack.svelte';
	import { fetchKegStats } from '$lib/api.js';
	import type { PageData } from './$types.js';
	import type { KegStats, Tap } from '$lib/api.js';

	let { data }: { data: PageData } = $props();

	let stats = $state<Record<number, KegStats>>({});

	$effect(() => {
		stats = Object.fromEntries(data.stats.map((s: KegStats) => [s.keg_id, s]));
	});

	let activeTaps = $derived((data.taps as Tap[]).filter((t) => t.keg));

	onMount(() => {
		pourSocket.connect();
		return () => pourSocket.disconnect();
	});

	$effect(() => {
		const ended = pourSocket.lastEnded;
		if (!ended) return;
		const tap = activeTaps.find((t) => t.id === ended.tap_id);
		if (tap?.keg) {
			fetchKegStats(tap.keg.id).then((s) => {
				stats = { ...stats, [s.keg_id]: s };
			});
		}
	});
</script>

<div class="relative h-screen w-screen overflow-hidden bg-gray-900">
	<BannerStack taps={activeTaps} {stats} customMessage={data.banner} />

	{#if activeTaps.length === 0}
		<div class="flex h-full items-center justify-center text-2xl text-white">No kegs on tap</div>
	{:else}
		<div
			class="flex h-full snap-x snap-mandatory overflow-x-auto [&::-webkit-scrollbar]:hidden"
			style="scrollbar-width: none;"
		>
			{#each activeTaps as tap (tap.id)}
				<div class="h-full w-screen shrink-0 snap-center">
					<KegCard {tap} stats={stats[tap.keg!.id]} />
				</div>
			{/each}
		</div>
	{/if}
</div>
