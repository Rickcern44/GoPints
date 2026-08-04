<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchFeatures, setFeature, type Features } from '$lib/api.js';

	let features = $state<Features>({ flow_based_pour: false, remote_image_urls: false });
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
	<div class="console-heading">
		<h1>System Modules</h1>
	</div>
	<p class="subtitle">Optional subsystems. Flip a switch to enable or disable — takes effect immediately.</p>

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
				class="console-switch"
				class:on={features.flow_based_pour}
				disabled={saving.flow_based_pour}
				onclick={() => toggle('flow_based_pour')}
				aria-pressed={features.flow_based_pour}
				aria-label="Toggle Flow-Based Pour"
			>
				<span class="switch-label off">OFF</span>
				<span class="switch-label on">ON</span>
				<span class="switch-thumb"></span>
			</button>
		</div>

		<div class="feature-card">
			<div class="feature-info">
				<p class="feature-name">Remote Image URLs</p>
				<p class="feature-desc">
					Let admins import a keg's beer image or brewery logo by pasting a URL instead of only
					uploading a file. The server fetches and stores the image itself.
				</p>
			</div>
			<button
				class="console-switch"
				class:on={features.remote_image_urls}
				disabled={saving.remote_image_urls}
				onclick={() => toggle('remote_image_urls')}
				aria-pressed={features.remote_image_urls}
				aria-label="Toggle Remote Image URLs"
			>
				<span class="switch-label off">OFF</span>
				<span class="switch-label on">ON</span>
				<span class="switch-thumb"></span>
			</button>
		</div>
	</div>
</div>

<style>
	.page {
		max-width: 680px;
	}

	.subtitle {
		color: var(--color-fg-muted);
		font-size: 0.9rem;
		margin: 0 0 2rem;
	}

	.error {
		background: var(--color-error-bg);
		color: var(--color-error);
		border: 1px solid color-mix(in srgb, var(--color-error) 35%, transparent);
		border-radius: 6px;
		padding: 0.6rem 0.9rem;
		font-size: 0.875rem;
		margin-bottom: 1.25rem;
	}

	.feature-list {
		display: flex;
		flex-direction: column;
		gap: 1px;
		border: 1px solid var(--color-line);
		border-radius: 4px;
		overflow: hidden;
		background: var(--color-line);
	}

	.feature-card {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 2rem;
		padding: 1.25rem 1.5rem;
		background: var(--color-panel-raised);
	}

	.feature-info {
		flex: 1;
		min-width: 0;
	}

	.feature-name {
		font-family: var(--font-brand);
		letter-spacing: 0.05em;
		text-transform: uppercase;
		font-size: 1rem;
		color: var(--color-fg);
		margin: 0 0 0.35rem;
	}

	.feature-desc {
		font-size: 0.825rem;
		color: var(--color-fg-muted);
		line-height: 1.5;
		margin: 0;
	}

	/* Console switch — an illuminated rocker toggle with an OFF/ON readout */
	.console-switch {
		position: relative;
		width: 4.5rem;
		height: 2rem;
		border: 1px solid var(--color-line);
		border-radius: 3px;
		background: var(--color-void);
		cursor: pointer;
		flex-shrink: 0;
		padding: 0;
		overflow: hidden;
	}

	.console-switch:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.switch-label {
		position: absolute;
		top: 50%;
		transform: translateY(-50%);
		font-family: var(--font-mono);
		font-size: 0.6rem;
		font-weight: 700;
		letter-spacing: 0.05em;
		color: var(--color-fg-muted);
		z-index: 1;
	}

	.switch-label.off {
		left: 0.45rem;
	}

	.switch-label.on {
		right: 0.45rem;
	}

	.switch-thumb {
		position: absolute;
		top: 2px;
		bottom: 2px;
		left: 2px;
		width: calc(50% - 4px);
		background: var(--color-line);
		border-radius: 2px;
		transition:
			transform 0.2s ease,
			background 0.2s ease,
			box-shadow 0.2s ease;
	}

	.console-switch.on .switch-thumb {
		transform: translateX(calc(100% + 0px));
		background: var(--color-signal-good);
		box-shadow: 0 0 10px color-mix(in srgb, var(--color-signal-good) 70%, transparent);
	}
</style>
