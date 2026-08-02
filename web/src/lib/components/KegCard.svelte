<script lang="ts">
	import LevelGauge from './LevelGauge.svelte';
	import type { Tap, KegStats } from '$lib/api.js';

	let { tap, stats }: { tap: Tap; stats: KegStats | undefined } = $props();

	let keg = $derived(tap.keg!);
	let pct = $derived(stats?.pct_remaining ?? 100);
	let remainingL = $derived(((stats?.remaining_ml ?? keg.capacity_ml) / 1000).toFixed(1));
	let hasImage = $derived(!!keg.image_mime_type);
	let hasBreweryImage = $derived(!!keg.brewery_image_mime_type);
	let imageStyle = $derived(keg.image_style || 'circle');
	let statusLed = $derived(pct > 50 ? 'good' : pct > 20 ? 'warn' : 'critical');
</script>

<div class="card">
	<div class="inner">
		{#if keg.brewery}
			<div class="brewery-plaque">
				{#if hasBreweryImage}
					<img src="/api/kegs/{keg.id}/brewery-image" alt="{keg.brewery} logo" class="plaque-logo" />
				{/if}
				<p class="plaque-name">{keg.brewery}</p>
			</div>
		{/if}

		<div class="info">
			<span class="tap-badge"><span class="led {statusLed}"></span>TAP {tap.id}</span>

			{#if hasImage}
				<img src="/api/kegs/{keg.id}/image" alt={keg.beer_name} class="beer-img beer-img-{imageStyle}" />
			{/if}

			<h1 class="name">{keg.beer_name}</h1>

			<div class="tags">
				{#if keg.style}<span class="tag">{keg.style}</span>{/if}
				{#if keg.abv > 0}<span class="tag tag-abv">{keg.abv}% ABV</span>{/if}
			</div>
		</div>

		<div class="divider"></div>

		<div class="stats">
			<LevelGauge pctRemaining={pct} />
			<p class="remaining">{remainingL}L remaining</p>
			{#if stats}
				<p class="pours">{stats.pour_count} {stats.pour_count !== 1 ? 'pours' : 'pour'}</p>
			{/if}
		</div>
	</div>
</div>

<style>
	.card {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-void);
		background-image:
			linear-gradient(
				color-mix(in srgb, var(--color-line) 35%, transparent) 1px,
				transparent 1px
			),
			linear-gradient(
				90deg,
				color-mix(in srgb, var(--color-line) 35%, transparent) 1px,
				transparent 1px
			);
		background-size: 48px 48px;
	}

	.inner {
		display: flex;
		align-items: center;
		gap: 3rem;
		padding: 3rem 5vw;
		width: 100%;
		max-width: 1300px;
	}

	/* LEFT — brewery plaque */

	.brewery-plaque {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.6rem;
		width: 150px;
		flex-shrink: 0;
		padding: 1.5rem 1rem;
		border-radius: 3px;
		background: var(--color-panel);
		border: 1px solid var(--color-line);
		border-top: 2px solid var(--color-accent);
	}

	.plaque-logo {
		width: 120px;
		height: 120px;
		object-fit: contain;
		border-radius: 3px;
		background: rgba(255, 255, 255, 0.04);
		padding: 0.5rem;
	}

	.plaque-name {
		font-family: var(--font-mono);
		font-size: 0.8rem;
		font-weight: 500;
		letter-spacing: 0.02em;
		color: var(--color-fg-muted);
		text-align: center;
		margin: 0;
		display: -webkit-box;
		-webkit-line-clamp: 3;
		line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
		word-break: break-word;
	}

	/* MIDDLE — beer info */

	.info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		min-width: 0;
	}

	.tap-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.25rem 0.7rem;
		border-radius: 3px;
		font-family: var(--font-mono);
		font-size: 0.68rem;
		font-weight: 700;
		letter-spacing: 0.18em;
		text-transform: uppercase;
		color: var(--color-fg-muted);
		border: 1px solid var(--color-line);
		background: var(--color-panel);
		width: fit-content;
	}

	.beer-img {
		object-fit: cover;
		border: 1px solid var(--color-line);
		margin-bottom: 0.15rem;
	}

	.beer-img-circle {
		width: 90px;
		height: 90px;
		border-radius: 50%;
	}

	.beer-img-square {
		width: 90px;
		height: 90px;
		border-radius: 4px;
	}

	.beer-img-can {
		width: 52px;
		height: 90px;
		border-radius: 50% / 12%;
	}

	.name {
		font-family: var(--font-display);
		font-size: clamp(2.5rem, 5vw, 5rem);
		font-weight: 700;
		color: var(--color-fg);
		line-height: 1.05;
		margin: 0;
		word-break: break-word;
	}

	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		margin-top: 0.15rem;
	}

	.tag {
		padding: 0.2rem 0.6rem;
		border-radius: 3px;
		font-family: var(--font-mono);
		font-size: 0.75rem;
		background: var(--color-panel);
		border: 1px solid var(--color-line);
		color: var(--color-fg-muted);
		letter-spacing: 0.02em;
	}

	.tag-abv {
		border-color: color-mix(in srgb, var(--color-accent) 40%, transparent);
		color: var(--color-accent);
	}

	/* Vertical divider */

	.divider {
		width: 1px;
		height: 55%;
		min-height: 180px;
		background: var(--color-line);
		flex-shrink: 0;
	}

	/* RIGHT — stats */

	.stats {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.85rem;
		flex-shrink: 0;
	}

	.remaining {
		font-family: var(--font-mono);
		font-size: clamp(0.85rem, 1.4vw, 1.05rem);
		color: var(--color-fg-muted);
		margin: 0;
		letter-spacing: 0.02em;
	}

	.pours {
		font-family: var(--font-mono);
		font-size: clamp(0.7rem, 1.1vw, 0.85rem);
		color: color-mix(in srgb, var(--color-fg-muted) 65%, transparent);
		margin: 0;
		letter-spacing: 0.02em;
	}
</style>
