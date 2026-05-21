<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchFeatures, setFeature, type Features } from '$lib/api.js';

	let features = $state<Features>({ flow_based_pour: false });
	let saving = $state<Partial<Record<keyof Features, boolean>>>({});
	let error = $state('');

	onMount(async () => {
		features = await fetchFeatures();
	});

	async function toggle(name: keyof Features) {
		saving[name] = true;
		error = '';
		try {
			await setFeature(name, !features[name]);
			features[name] = !features[name];
		} catch {
			error = 'Failed to save feature. Are you still logged in?';
		} finally {
			saving[name] = false;
		}
	}
</script>

<div class="page">
	<h1 class="title">Features</h1>
	<p class="subtitle">Toggle optional capabilities. Changes take effect immediately.</p>

	{#if error}
		<p class="error">{error}</p>
	{/if}

	<div class="feature-list">
		<div class="feature-card">
			<div class="feature-info">
				<p class="feature-name">Flow-Based Pour</p>
				<p class="feature-desc">
					Enable hardware flow meter pour tracking via GPIO. When on, manual pour buttons are hidden
					and pours are recorded automatically by the flow sensor.
				</p>
			</div>
			<button
				class="toggle"
				class:on={features.flow_based_pour}
				disabled={saving.flow_based_pour}
				onclick={() => toggle('flow_based_pour')}
				aria-pressed={features.flow_based_pour}
				aria-label="Toggle Flow-Based Pour"
			>
				<span class="thumb"></span>
			</button>
		</div>
	</div>
</div>

<style>
	.page {
		max-width: 680px;
	}

	.title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #1c1917;
		margin: 0 0 0.25rem;
	}

	.subtitle {
		color: #78716c;
		font-size: 0.9rem;
		margin: 0 0 2rem;
	}

	.error {
		background: #fef2f2;
		color: #b91c1c;
		border: 1px solid #fecaca;
		border-radius: 6px;
		padding: 0.6rem 0.9rem;
		font-size: 0.875rem;
		margin-bottom: 1.25rem;
	}

	.feature-list {
		display: flex;
		flex-direction: column;
		gap: 1px;
		border: 1px solid #e7e5e4;
		border-radius: 10px;
		overflow: hidden;
		background: #e7e5e4;
	}

	.feature-card {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 2rem;
		padding: 1.25rem 1.5rem;
		background: #fff;
	}

	.feature-info {
		flex: 1;
		min-width: 0;
	}

	.feature-name {
		font-size: 0.95rem;
		font-weight: 600;
		color: #1c1917;
		margin: 0 0 0.3rem;
	}

	.feature-desc {
		font-size: 0.825rem;
		color: #78716c;
		line-height: 1.5;
		margin: 0;
	}

	/* Toggle switch */
	.toggle {
		position: relative;
		width: 44px;
		height: 24px;
		border-radius: 999px;
		border: none;
		background: #d6d3d1;
		cursor: pointer;
		transition: background 0.2s;
		flex-shrink: 0;
		padding: 0;
	}

	.toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.toggle.on {
		background: #c8821a;
	}

	.thumb {
		position: absolute;
		top: 3px;
		left: 3px;
		width: 18px;
		height: 18px;
		border-radius: 50%;
		background: #fff;
		transition: transform 0.2s;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
	}

	.toggle.on .thumb {
		transform: translateX(20px);
	}
</style>
