<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import Spinner from '$lib/components/Spinner.svelte';

	let current = $state<string | null>(null);
	let message = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let status = $state('');

	onMount(async () => {
		try {
			const res = await fetch('/api/banner');
			if (res.ok) {
				const data: { message: string } = await res.json();
				current = data.message;
				message = data.message;
			}
		} catch { /* banner is optional */ } finally {
			loading = false;
		}
	});

	async function setBanner() {
		if (!message.trim()) return;
		saving = true; status = '';
		try {
			const res = await adminFetch('/api/banner', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ message: message.trim() })
			});
			if (res.ok) { current = message.trim(); status = 'Banner updated'; }
			else { status = 'Failed to update banner'; }
		} finally {
			saving = false;
		}
	}

	async function clearBanner() {
		if (!confirm('Clear the banner?')) return;
		await adminFetch('/api/banner', { method: 'DELETE' });
		current = null;
		message = '';
		status = 'Banner cleared';
	}
</script>

<div class="max-w-xl">
	<div class="console-heading">
		<h1>System Broadcast</h1>
	</div>

	{#if loading}
		<div class="flex items-center gap-2 text-sm text-fg-muted"><Spinner size={16} /> Loading…</div>
	{:else}
		{#if current}
			<!-- Mirrors the live kiosk banner color (see BannerStack.svelte) so this is an accurate preview. -->
			<div class="mb-6 rounded-xl bg-indigo-600 px-5 py-4 text-white shadow-sm">
				<div class="mb-1 text-xs font-medium tracking-wide text-indigo-300 uppercase">
					Active Broadcast
				</div>
				<div class="text-lg font-semibold">{current}</div>
			</div>
		{:else}
			<p class="mb-6 text-sm text-fg-muted">No broadcast is currently active.</p>
		{/if}

		{#if status}
			<div class="mb-4 rounded-lg bg-success-bg px-4 py-2 text-sm text-success">{status}</div>
		{/if}

		<div class="panel">
			<label class="field">
				<span class="cap">Message</span>
				<textarea bind:value={message} rows="3" placeholder="Tap 2 offline for cleaning…"></textarea>
			</label>
			<div class="mt-4 flex gap-3">
				<button onclick={setBanner} disabled={saving || !message.trim()} class="btn-console">
					{#if saving}<Spinner size={14} />{/if}
					{saving ? 'Broadcasting…' : current ? 'Update Broadcast' : 'Send Broadcast'}
				</button>
				{#if current}
					<button onclick={clearBanner} class="btn-console-ghost">Clear</button>
				{/if}
			</div>
		</div>
	{/if}
</div>
