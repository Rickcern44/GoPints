<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { setToken, getToken } from '$lib/admin.js';

	let isSetup = $state(false);
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let loading = $state(false);

	onMount(async () => {
		if (getToken()) {
			goto('/admin/kegs', { replaceState: true });
			return;
		}
		try {
			const res = await fetch('/api/admin/status');
			if (!res.ok) throw new Error(`${res.status}`);
			const data: { password_set: boolean } = await res.json();
			isSetup = !data.password_set;
		} catch {
			error = 'Could not reach the server. Check that it is running.';
		}
	});

	async function submit() {
		error = '';
		if (isSetup && password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}
		if (!password) {
			error = 'Password is required';
			return;
		}
		loading = true;
		try {
			const endpoint = isSetup ? '/api/admin/setup' : '/api/admin/login';
			const res = await fetch(endpoint, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password })
			});
			if (!res.ok) {
				error = isSetup ? 'Setup failed' : 'Incorrect password';
				return;
			}
			const data: { token: string } = await res.json();
			setToken(data.token);
			goto('/admin/kegs', { replaceState: true });
		} finally {
			loading = false;
		}
	}
</script>

<div class="cover-stage">
	<div class="cover-card">
		<div class="logo-chip">
			<svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true">
				<rect x="9" y="2" width="6" height="3" rx="1" fill="currentColor" />
				<rect x="8" y="5" width="8" height="13" rx="1" fill="currentColor" opacity="0.85" />
				<path d="M16 8h3a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1h-3" stroke="currentColor" stroke-width="1.5" fill="none" />
			</svg>
		</div>
		<p class="eyebrow">System Access</p>
		<h1>GoPints Console</h1>
		<p class="byline">
			{isSetup ? 'Set an admin password to configure the system.' : 'Authenticate to access the console.'}
		</p>

		{#if error}
			<div class="mb-4 rounded-lg bg-error-bg px-4 py-3 text-sm text-error">{error}</div>
		{/if}

		<form onsubmit={(e) => { e.preventDefault(); submit(); }} class="flex flex-col gap-4">
			<label class="field">
				<span class="cap">{isSetup ? 'New Password' : 'Password'}</span>
				<input
					id="password"
					type="password"
					bind:value={password}
					autocomplete={isSetup ? 'new-password' : 'current-password'}
				/>
			</label>

			{#if isSetup}
				<label class="field">
					<span class="cap">Confirm Password</span>
					<input id="confirm" type="password" bind:value={confirmPassword} autocomplete="new-password" />
				</label>
			{/if}

			<button type="submit" disabled={loading} class="btn-console w-full justify-center mt-2">
				{loading ? 'Please wait…' : isSetup ? 'Set Password & Continue' : 'Sign In'}
			</button>
		</form>
	</div>
</div>

<style>
	.cover-stage {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-void);
		background-image: radial-gradient(
			ellipse 70% 50% at 50% 0%,
			color-mix(in srgb, var(--color-accent) 6%, transparent) 0%,
			transparent 70%
		);
	}

	.cover-card {
		width: 100%;
		max-width: 24rem;
		background: var(--color-panel-raised);
		border: 1px solid var(--color-line);
		border-top: 2px solid var(--color-accent);
		border-radius: 3px;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.35);
		padding: 2.5rem 2.25rem 2.25rem;
		text-align: center;
	}

	.logo-chip {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 3.25rem;
		height: 3.25rem;
		background: var(--color-void);
		border: 1px solid var(--color-line);
		clip-path: polygon(20% 0, 100% 0, 100% 80%, 80% 100%, 0 100%, 0 20%);
		color: var(--color-accent);
		margin-bottom: 1.1rem;
	}

	.eyebrow {
		font-family: var(--font-mono);
		font-size: 0.68rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--color-fg-muted);
		margin: 0 0 0.35rem;
	}

	.cover-card h1 {
		font-family: var(--font-brand);
		font-size: 1.7rem;
		letter-spacing: 0.03em;
		text-transform: uppercase;
		color: var(--color-fg);
		margin: 0 0 0.4rem;
	}

	.byline {
		font-family: var(--font-mono);
		font-size: 0.82rem;
		color: var(--color-fg-muted);
		margin: 0 0 1.75rem;
	}

	.cover-card form {
		text-align: left;
	}
</style>
