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
		<div class="info">
			<span class="tap-badge">TAP {tap.id}</span>

			{#if hasImage}
				<img src="/api/kegs/{keg.id}/image" alt={keg.beer_name} class="beer-img beer-img-{imageStyle}" />
			{/if}

			<h1 class="name">{keg.beer_name}</h1>

			{#if keg.brewery}
				<div class="brewery-row">
					{#if hasBreweryImage}
						<img src="/api/kegs/{keg.id}/brewery-image" alt="{keg.brewery} logo" class="brewery-logo" />
					{/if}
					<p class="brewery">{keg.brewery}</p>
				</div>
			{/if}

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
		gap: 4rem;
		padding: 3rem 5vw;
		width: 100%;
		max-width: 1300px;
	}

	/* LEFT — beer info */

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

	.brewery-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.brewery-logo {
		width: 28px;
		height: 28px;
		border-radius: 4px;
		object-fit: contain;
		background: rgba(255, 255, 255, 0.05);
		flex-shrink: 0;
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

	.brewery {
		font-size: clamp(1rem, 2vw, 1.4rem);
		color: rgba(210, 175, 110, 0.6);
		margin: 0;
		letter-spacing: 0.03em;
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
