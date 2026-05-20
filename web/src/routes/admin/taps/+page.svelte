<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import type { Tap, Keg } from '$lib/api.js';

	let taps = $state<Tap[]>([]);
	let kegs = $state<Keg[]>([]);
	let loading = $state(true);
	let error = $state('');
	let selections = $state<Record<number, string>>({});

	onMount(async () => {
		try {
			const [tapsRes, kegsRes] = await Promise.all([fetch('/api/taps'), fetch('/api/kegs')]);
			if (!tapsRes.ok || !kegsRes.ok) throw new Error('fetch failed');
			taps = (await tapsRes.json()) ?? [];
			kegs = (await kegsRes.json()) ?? [];
			for (const tap of taps) {
				selections[tap.id] = tap.keg_id ? String(tap.keg_id) : '';
			}
		} catch {
			error = 'Failed to load taps. Check that the server is running.';
		} finally {
			loading = false;
		}
	});

	async function assignKeg(tapId: number) {
		const kegId = selections[tapId];
		if (!kegId) {
			await adminFetch(`/api/taps/${tapId}/keg`, { method: 'DELETE' });
		} else {
			await adminFetch(`/api/taps/${tapId}/keg`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ keg_id: parseInt(kegId) })
			});
		}
		const res = await fetch('/api/taps');
		taps = (await res.json()) ?? [];
		for (const tap of taps) {
			selections[tap.id] = tap.keg_id ? String(tap.keg_id) : '';
		}
	}
</script>

<div class="max-w-2xl">
	<h1 class="mb-6 text-2xl font-bold">Taps</h1>

	{#if error}
		<div class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading}
		<p class="text-gray-500">Loading…</p>
	{:else if taps.length === 0}
		<p class="text-gray-500">No taps configured. Taps are created automatically when the agent starts.</p>
	{:else}
		<div class="space-y-3">
			{#each taps as tap (tap.id)}
				<div class="flex items-center gap-4 rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
					<div class="min-w-[64px] text-center">
						<div class="text-2xl font-black text-gray-800">#{tap.id}</div>
						<div class="text-xs text-gray-400">tap</div>
					</div>
					<div class="flex-1">
						{#if tap.keg}
							<div class="text-sm font-medium text-gray-900">{tap.keg.beer_name}</div>
							<div class="text-xs text-gray-500">{tap.keg.brewery || ''} · {(tap.keg.capacity_ml / 1000).toFixed(1)}L</div>
						{:else}
							<div class="text-sm text-gray-400 italic">No keg assigned</div>
						{/if}
					</div>
					<div class="flex items-center gap-2">
						<select
							bind:value={selections[tap.id]}
							class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
						>
							<option value="">— Remove keg —</option>
							{#each kegs as keg (keg.id)}
								<option value={String(keg.id)}>{keg.beer_name}</option>
							{/each}
						</select>
						<button
							onclick={() => assignKeg(tap.id)}
							class="rounded-lg bg-blue-600 px-3 py-2 text-sm font-semibold text-white hover:bg-blue-700"
						>
							Apply
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
