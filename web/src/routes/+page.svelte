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

	let currentIndex = $state(0);
	let carouselEl = $state<HTMLDivElement | undefined>(undefined);
	let touchStartX = 0;
	let touchStartY = 0;

	$effect(() => {
		if (currentIndex >= activeTaps.length && activeTaps.length > 0) {
			currentIndex = activeTaps.length - 1;
		}
	});

	function goTo(index: number) {
		const clamped = Math.max(0, Math.min(index, activeTaps.length - 1));
		currentIndex = clamped;
		const slides = carouselEl?.querySelectorAll<HTMLElement>('.carousel-slide');
		slides?.[clamped]?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' });
	}

	function onTouchStart(e: TouchEvent) {
		touchStartX = e.touches[0].clientX;
		touchStartY = e.touches[0].clientY;
	}

	function onTouchEnd(e: TouchEvent) {
		const dx = e.changedTouches[0].clientX - touchStartX;
		const dy = e.changedTouches[0].clientY - touchStartY;
		if (Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > 40) {
			goTo(currentIndex + (dx < 0 ? 1 : -1));
		}
	}

	function onCarouselScroll() {
		if (!carouselEl) return;
		currentIndex = Math.round(carouselEl.scrollLeft / window.innerWidth);
	}

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
		<div
			class="carousel"
			style="scrollbar-width: none;"
			role="region"
			aria-label="Keg carousel"
			bind:this={carouselEl}
			ontouchstart={onTouchStart}
			ontouchend={onTouchEnd}
			onscroll={onCarouselScroll}
		>
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

		{#if activeTaps.length > 1}
			<div class="dot-nav">
				{#each activeTaps as _, i}
					<button
						class="dot"
						class:active={i === currentIndex}
						onclick={() => goTo(i)}
						aria-label="Go to tap {i + 1}"
					></button>
				{/each}
			</div>
		{/if}
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
		background-color: var(--color-void);
		background-image:
			linear-gradient(color-mix(in srgb, var(--color-line) 30%, transparent) 1px, transparent 1px),
			linear-gradient(
				90deg,
				color-mix(in srgb, var(--color-line) 30%, transparent) 1px,
				transparent 1px
			);
		background-size: 48px 48px;
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
		color: color-mix(in srgb, var(--color-accent) 40%, transparent);
		margin-bottom: 0.25rem;
	}

	.empty-title {
		font-family: var(--font-display);
		font-size: 1.6rem;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: var(--color-fg-muted);
		margin: 0;
	}

	.empty-sub {
		font-family: var(--font-mono);
		font-size: 0.85rem;
		color: color-mix(in srgb, var(--color-fg-muted) 60%, transparent);
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

	.dot-nav {
		position: absolute;
		bottom: 1.25rem;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		gap: 0.5rem;
		align-items: center;
		z-index: 10;
	}

	.dot {
		width: 4px;
		height: 14px;
		border-radius: 1px;
		background: var(--color-line);
		border: none;
		padding: 0;
		cursor: pointer;
		transition:
			background 0.2s,
			box-shadow 0.2s,
			transform 0.2s;
	}

	.dot.active {
		background: var(--color-accent);
		box-shadow: 0 0 8px color-mix(in srgb, var(--color-accent) 70%, transparent);
		transform: scaleY(1.3);
	}

	.dot:not(.active):hover {
		background: var(--color-fg-muted);
	}

	.admin-corner {
		position: absolute;
		bottom: 1rem;
		right: 1rem;
		color: color-mix(in srgb, var(--color-fg-muted) 40%, transparent);
		text-decoration: none;
		padding: 0.25rem;
		transition: color 0.2s;
	}

	.admin-corner:hover {
		color: var(--color-accent);
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
		background: var(--color-panel);
		border: 1px solid var(--color-line);
		border-top: 2px solid var(--color-accent);
		border-radius: 3px;
		padding: 0.5rem;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
		white-space: nowrap;
	}

	.preset-btn {
		padding: 0.6rem 1rem;
		border-radius: 3px;
		background: var(--color-void);
		border: 1px solid var(--color-line);
		color: var(--color-fg);
		font-family: var(--font-mono);
		font-size: 0.85rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			border-color 0.15s,
			color 0.15s,
			transform 0.1s;
		letter-spacing: 0.01em;
	}

	.preset-btn:not(:disabled):hover {
		border-color: var(--color-accent);
		color: var(--color-accent);
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
		border-radius: 3px;
		background: var(--color-accent);
		border: 1px solid var(--color-accent);
		color: var(--color-void);
		font-family: var(--font-display);
		font-size: 0.85rem;
		font-weight: 700;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		cursor: pointer;
		transition: box-shadow 0.2s, filter 0.2s;
	}

	.pour-fab.open,
	.pour-fab:not(:disabled):hover {
		filter: brightness(1.1);
		box-shadow: 0 0 0 4px color-mix(in srgb, var(--color-accent) 25%, transparent);
	}

	.pour-fab:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.pour-done {
		padding: 0.6rem 1.25rem;
		border-radius: 3px;
		background: color-mix(in srgb, var(--color-success) 15%, transparent);
		border: 1px solid var(--color-success);
		color: var(--color-success);
		font-family: var(--font-display);
		font-size: 0.85rem;
		font-weight: 700;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		white-space: nowrap;
	}
</style>
