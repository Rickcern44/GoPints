<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';

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
	<h1 class="mb-6 text-2xl font-bold">Custom Banner</h1>

	{#if loading}
		<p class="text-gray-500">Loading…</p>
	{:else}
		{#if current}
			<div class="mb-6 rounded-xl bg-indigo-600 px-5 py-4 text-white shadow-sm">
				<div class="mb-1 text-xs font-medium tracking-wide text-indigo-300 uppercase">
					Current Banner
				</div>
				<div class="text-lg font-semibold">{current}</div>
			</div>
		{:else}
			<p class="mb-6 text-sm text-gray-500">No custom banner is currently set.</p>
		{/if}

		{#if status}
			<div class="mb-4 rounded-lg bg-green-50 px-4 py-2 text-sm text-green-700">{status}</div>
		{/if}

		<div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
			<label class="mb-1 flex flex-col gap-1 text-sm font-medium text-gray-700">
				Message
				<textarea
					bind:value={message}
					rows="3"
					placeholder="Tap 2 offline for cleaning…"
					class="mb-4 w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
				></textarea>
			</label>
			<div class="flex gap-3">
				<button
					onclick={setBanner}
					disabled={saving || !message.trim()}
					class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-700 disabled:opacity-50"
				>
					{saving ? 'Saving…' : current ? 'Update Banner' : 'Set Banner'}
				</button>
				{#if current}
					<button
						onclick={clearBanner}
						class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
					>
						Clear Banner
					</button>
				{/if}
			</div>
		</div>
	{/if}
</div>
