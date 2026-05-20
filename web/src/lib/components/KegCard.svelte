<script lang="ts">
	import LevelGauge from './LevelGauge.svelte';
	import type { Tap, KegStats } from '$lib/api.js';

	let { tap, stats }: { tap: Tap; stats: KegStats | undefined } = $props();

	let keg = $derived(tap.keg!);
	let pct = $derived(stats?.pct_remaining ?? 100);
	let remainingL = $derived(((stats?.remaining_ml ?? keg.capacity_ml) / 1000).toFixed(1));
	let hasImage = $derived(!!keg.image_mime_type);
</script>

<div class="flex h-full w-full flex-col items-center justify-center gap-8 bg-gray-900 p-8 text-white">
	{#if hasImage}
		<img
			src="/api/kegs/{keg.id}/image"
			alt={keg.beer_name}
			class="h-24 w-24 rounded-full object-cover ring-4 ring-gray-700"
		/>
	{/if}

	<div class="text-center">
		<h1 class="text-5xl font-bold tracking-tight">{keg.beer_name}</h1>
		{#if keg.brewery}
			<p class="mt-1 text-2xl text-gray-400">{keg.brewery}</p>
		{/if}
		<div class="mt-2 flex justify-center gap-4 text-lg text-gray-300">
			{#if keg.style}<span>{keg.style}</span>{/if}
			{#if keg.style && keg.abv > 0}<span class="text-gray-600">·</span>{/if}
			{#if keg.abv > 0}<span>{keg.abv}% ABV</span>{/if}
		</div>
	</div>

	<div class="flex items-center gap-10">
		<LevelGauge pctRemaining={pct} />
		<div class="text-center">
			<div class="text-6xl font-black tabular-nums">{pct.toFixed(0)}%</div>
			<div class="mt-1 text-xl text-gray-400">{remainingL}L remaining</div>
			{#if stats}
				<div class="mt-2 text-sm text-gray-500">
					{stats.pour_count}
					{stats.pour_count !== 1 ? 'pours' : 'pour'}
				</div>
			{/if}
		</div>
	</div>
</div>
