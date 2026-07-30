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
			<span class="tap-badge">TAP {tap.id}</span>

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
			<div class="gauge-outer">
				<LevelGauge pctRemaining={pct} />
			</div>
			<div class="pct-row">
				<span class="pct-num">{pct.toFixed(0)}</span><span class="pct-sym">%</span>
			</div>
			<p class="remaining">{remainingL}L remaining</p>
			{#if stats}
				<p class="pours">{stats.pour_count} {stats.pour_count !== 1 ? 'pours' : 'pour'}</p>
			{/if}
		</div>
	</div>
</div>

<style>
	@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@1,700&display=swap');

	.card {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #0a0806;
		background-image:
			radial-gradient(ellipse 80% 60% at 75% 50%, rgba(160, 100, 8, 0.08) 0%, transparent 65%),
			radial-gradient(ellipse 40% 40% at 12% 50%, rgba(80, 50, 5, 0.05) 0%, transparent 60%);
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
		border-radius: 16px;
		background: linear-gradient(160deg, rgba(200, 130, 26, 0.09), rgba(160, 100, 8, 0.03));
		border: 1px solid rgba(200, 130, 26, 0.18);
	}

	.plaque-logo {
		width: 120px;
		height: 120px;
		object-fit: contain;
		border-radius: 12px;
		background: rgba(255, 255, 255, 0.04);
		padding: 0.5rem;
	}

	.plaque-name {
		font-size: 0.85rem;
		font-weight: 600;
		letter-spacing: 0.04em;
		color: rgba(210, 175, 110, 0.75);
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
		padding: 0.2rem 0.7rem;
		border-radius: 999px;
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.22em;
		text-transform: uppercase;
		color: rgba(200, 130, 26, 0.9);
		border: 1px solid rgba(200, 130, 26, 0.3);
		background: rgba(200, 130, 26, 0.07);
		width: fit-content;
	}

	.beer-img {
		object-fit: cover;
		border: 2px solid rgba(200, 130, 26, 0.25);
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
		border-radius: 8px;
	}

	.beer-img-can {
		width: 52px;
		height: 90px;
		border-radius: 50% / 12%;
	}

	.name {
		font-family: 'Playfair Display', Georgia, 'Times New Roman', serif;
		font-size: clamp(2.5rem, 5vw, 5rem);
		font-weight: 700;
		font-style: italic;
		color: #f5ead8;
		line-height: 1.1;
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
		font-size: 0.78rem;
		background: rgba(255, 255, 255, 0.05);
		color: rgba(255, 255, 255, 0.38);
		letter-spacing: 0.04em;
	}

	.tag-abv {
		background: rgba(200, 130, 26, 0.08);
		color: rgba(200, 130, 26, 0.65);
	}

	/* Vertical divider */

	.divider {
		width: 1px;
		height: 55%;
		min-height: 180px;
		background: linear-gradient(to bottom, transparent, rgba(200, 130, 26, 0.15), transparent);
		flex-shrink: 0;
	}

	/* RIGHT — stats */

	.stats {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.gauge-outer {
		transform: scale(1.25);
		margin: 2rem 0;
	}

	.pct-row {
		display: flex;
		align-items: flex-start;
		line-height: 1;
	}

	.pct-num {
		font-size: clamp(4rem, 7vw, 6.5rem);
		font-weight: 900;
		color: #f5ead8;
		font-variant-numeric: tabular-nums;
		line-height: 1;
	}

	.pct-sym {
		font-size: clamp(1.4rem, 2.5vw, 2rem);
		font-weight: 700;
		color: rgba(245, 234, 216, 0.38);
		margin-top: 0.4rem;
		margin-left: 0.15rem;
	}

	.remaining {
		font-size: clamp(0.85rem, 1.4vw, 1.05rem);
		color: rgba(210, 175, 110, 0.5);
		margin: 0;
		letter-spacing: 0.04em;
	}

	.pours {
		font-size: clamp(0.7rem, 1.1vw, 0.85rem);
		color: rgba(255, 255, 255, 0.22);
		margin: 0;
		letter-spacing: 0.04em;
	}
</style>
