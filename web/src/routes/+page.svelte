<script lang="ts">
	import { onMount } from 'svelte';
	import { pourSocket } from '$lib/ws.svelte.js';
	import KegCard from '$lib/components/KegCard.svelte';
	import BannerStack from '$lib/components/BannerStack.svelte';
	import { fetchKegStats, fetchFeatures, recordManualPour } from '$lib/api.js';
	import type { PageData } from './$types.js';
	import type { KegStats, Tap } from '$lib/api.js';

	let { data }: { data: PageData } = $props();

	let stats = $state<Record<number, KegStats>>({});

	$effect(() => {
		stats = Object.fromEntries(data.stats.map((s: KegStats) => [s.keg_id, s]));
	});

	let activeTaps = $derived((data.taps as Tap[]).filter((t) => t.keg));

	const POUR_PRESETS = [
		{ label: '12 oz', ml: 12 * 29.5735 },
		{ label: '16 oz', ml: 16 * 29.5735 },
		{ label: '24 oz', ml: 24 * 29.5735 }
	] as const;

	let flowEnabled = $state(false);
	let pourOpen = $state<Record<number, boolean>>({});
	let pouring = $state<Record<number, boolean>>({});
	let pourDone = $state<Record<number, boolean>>({});

	onMount(() => {
		pourSocket.connect();
		fetchFeatures().then((f) => (flowEnabled = f.flow_based_pour));
		return () => pourSocket.disconnect();
	});

	async function logPour(tapId: number, ml: number) {
		pourOpen[tapId] = false;
		pouring[tapId] = true;
		try {
			await recordManualPour(tapId, ml);
			pourDone[tapId] = true;
			setTimeout(() => (pourDone[tapId] = false), 2500);
		} finally {
			pouring[tapId] = false;
		}
	}

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

					{#if !flowEnabled}
						<div class="pour-anchor">
							{#if pourDone[tap.id]}
								<div class="pour-done">Pour logged!</div>
							{:else}
								{#if pourOpen[tap.id]}
									<button
										class="pour-backdrop"
										onclick={() => (pourOpen[tap.id] = false)}
										aria-label="Close"
										tabindex="-1"
									></button>
								{/if}
								<div class="pour-trigger-wrap">
									{#if pourOpen[tap.id]}
										<div class="pour-flyout">
											{#each POUR_PRESETS as preset}
												<button
													class="preset-btn"
													disabled={pouring[tap.id]}
													onclick={() => logPour(tap.id, preset.ml)}
												>{preset.label}</button>
											{/each}
										</div>
									{/if}
									<button
										class="pour-fab"
										class:open={pourOpen[tap.id]}
										disabled={pouring[tap.id]}
										onclick={() => (pourOpen[tap.id] = !pourOpen[tap.id])}
									>
										{#if !pouring[tap.id]}
											<svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
												<rect x="9" y="2" width="6" height="3" rx="1" fill="currentColor" />
												<rect x="8" y="5" width="8" height="13" rx="1" fill="currentColor" opacity="0.85" />
												<path d="M16 8h3a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1h-3" stroke="currentColor" stroke-width="1.5" fill="none" />
											</svg>
										{/if}
										{pouring[tap.id] ? 'Logging…' : 'Log Pour'}
									</button>
								</div>
							{/if}
						</div>
					{/if}
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
		position: relative;
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

	/* Pour flyout */

	.pour-anchor {
		position: absolute;
		bottom: 2rem;
		left: 50%;
		transform: translateX(-50%);
	}

	.pour-backdrop {
		position: fixed;
		inset: 0;
		z-index: 10;
		background: transparent;
		border: none;
		cursor: default;
		padding: 0;
	}

	.pour-trigger-wrap {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		z-index: 11;
	}

	.pour-flyout {
		position: absolute;
		bottom: calc(100% + 0.625rem);
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		gap: 0.375rem;
		background: rgba(10, 8, 6, 0.93);
		border: 1px solid rgba(200, 130, 26, 0.3);
		border-radius: 12px;
		padding: 0.5rem;
		backdrop-filter: blur(20px);
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
		white-space: nowrap;
	}

	.preset-btn {
		padding: 0.6rem 1rem;
		border-radius: 8px;
		background: rgba(200, 130, 26, 0.1);
		border: 1px solid rgba(200, 130, 26, 0.2);
		color: rgba(245, 234, 216, 0.85);
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			background 0.15s,
			border-color 0.15s,
			transform 0.1s;
		letter-spacing: 0.01em;
	}

	.preset-btn:not(:disabled):hover {
		background: rgba(200, 130, 26, 0.22);
		border-color: rgba(200, 130, 26, 0.5);
		transform: translateY(-1px);
	}

	.preset-btn:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}

	.pour-fab {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.6rem 1.25rem;
		border-radius: 999px;
		background: rgba(200, 130, 26, 0.15);
		border: 1px solid rgba(200, 130, 26, 0.35);
		color: rgba(200, 130, 26, 0.9);
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		backdrop-filter: blur(8px);
		transition: background 0.2s, border-color 0.2s;
	}

	.pour-fab.open,
	.pour-fab:not(:disabled):hover {
		background: rgba(200, 130, 26, 0.25);
		border-color: rgba(200, 130, 26, 0.6);
	}

	.pour-fab:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.pour-done {
		padding: 0.6rem 1.25rem;
		border-radius: 999px;
		background: rgba(16, 185, 129, 0.15);
		border: 1px solid rgba(16, 185, 129, 0.3);
		color: rgba(16, 185, 129, 0.9);
		font-size: 0.875rem;
		font-weight: 600;
		white-space: nowrap;
	}
</style>
