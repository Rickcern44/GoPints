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

<div class="stage">
	<BannerStack taps={activeTaps} {stats} customMessage={data.banner} />

	{#if activeTaps.length === 0}
		<div class="empty-state">
			<div class="empty-icon" aria-hidden="true">
				<svg width="56" height="56" viewBox="0 0 24 24" fill="none">
					<rect x="9" y="2" width="6" height="3" rx="1" fill="currentColor" />
					<rect x="8" y="5" width="8" height="13" rx="1" fill="currentColor" opacity="0.7" />
					<path
						d="M16 8h3a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1h-3"
						stroke="currentColor"
						stroke-width="1.5"
						fill="none"
					/>
					<rect x="9" y="18" width="6" height="2" rx="0.5" fill="currentColor" opacity="0.4" />
				</svg>
			</div>
			<p class="empty-title">Nothing On Tap</p>
			<p class="empty-sub">Check back once the kegs are loaded up</p>
		</div>
	{:else}
		<div class="carousel" style="scrollbar-width: none;">
			{#each activeTaps as tap (tap.id)}
				<div class="carousel-slide">
					<KegCard {tap} stats={stats[tap.keg!.id]} />
				</div>
			{/each}
		</div>
	{/if}

	<a href="/admin" class="admin-corner" aria-label="Admin">
		<svg
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="1.5"
			stroke-linecap="round"
			stroke-linejoin="round"
			aria-hidden="true"
		>
			<circle cx="12" cy="12" r="3" />
			<path
				d="M12 2v2M12 20v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M2 12h2M20 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"
			/>
		</svg>
	</a>
</div>

<style>
	.stage {
		position: relative;
		height: 100vh;
		width: 100vw;
		overflow: hidden;
		background-color: #111111;
		background-image:
			radial-gradient(ellipse 60% 50% at 50% 0%, rgba(200, 130, 26, 0.05) 0%, transparent 70%),
			radial-gradient(ellipse 80% 40% at 50% 110%, rgba(20, 12, 4, 0.9) 0%, transparent 80%);
	}

	.empty-state {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
	}

	.empty-icon {
		color: rgba(200, 130, 26, 0.3);
		margin-bottom: 0.25rem;
	}

	.empty-title {
		font-size: 1.75rem;
		font-weight: 600;
		letter-spacing: 0.06em;
		color: rgba(255, 255, 255, 0.5);
		margin: 0;
	}

	.empty-sub {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.22);
		margin: 0;
	}

	.carousel {
		display: flex;
		height: 100%;
		overflow-x: auto;
		scroll-snap-type: x mandatory;
	}

	.carousel::-webkit-scrollbar {
		display: none;
	}

	.carousel-slide {
		height: 100%;
		width: 100vw;
		flex-shrink: 0;
		scroll-snap-align: center;
	}

	.admin-corner {
		position: absolute;
		bottom: 1rem;
		right: 1rem;
		color: rgba(255, 255, 255, 0.18);
		text-decoration: none;
		padding: 0.25rem;
		transition: color 0.2s;
	}

	.admin-corner:hover {
		color: rgba(255, 255, 255, 0.55);
	}
</style>
