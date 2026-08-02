<script lang="ts">
	let { pctRemaining }: { pctRemaining: number } = $props();

	let fill = $derived(Math.max(0, Math.min(100, pctRemaining)));
	let ringColor = $derived(
		fill > 50 ? 'var(--color-signal-good)' : fill > 20 ? 'var(--color-accent)' : 'var(--color-signal-critical)'
	);
</script>

<div class="gauge-ring" style="background: conic-gradient({ringColor} {fill}%, var(--color-line) 0)">
	<div class="gauge-hole">
		<span class="gauge-num">{fill.toFixed(0)}</span><span class="gauge-sym">%</span>
	</div>
</div>

<style>
	.gauge-ring {
		--size: 11rem;
		width: var(--size);
		height: var(--size);
		border-radius: 50%;
		padding: 0.65rem;
		display: flex;
		transition: background 1s ease;
	}

	.gauge-hole {
		flex: 1;
		border-radius: 50%;
		background: var(--color-void);
		box-shadow: inset 0 0 0 1px var(--color-line);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.gauge-num {
		font-family: var(--font-mono);
		font-size: 2.6rem;
		font-weight: 700;
		color: var(--color-fg);
		font-variant-numeric: tabular-nums;
		line-height: 1;
	}

	.gauge-sym {
		font-family: var(--font-mono);
		font-size: 0.95rem;
		color: var(--color-fg-muted);
		margin-top: 0.3rem;
	}
</style>
