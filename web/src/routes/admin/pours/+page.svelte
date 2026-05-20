<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import type { Tap, Pour } from '$lib/api.js';

	const PAGE_SIZE = 20;

	let taps = $state<Tap[]>([]);
	let pours = $state<Pour[]>([]);
	let loading = $state(true);
	let error = $state('');
	let tapFilter = $state('');
	let offset = $state(0);
	let hasMore = $state(false);

	onMount(async () => {
		try {
			const res = await fetch('/api/taps');
			if (res.ok) taps = (await res.json()) ?? [];
		} catch { /* taps filter is optional */ }
		await loadPours();
	});

	async function loadPours() {
		loading = true;
		error = '';
		try {
			const params = new URLSearchParams({
				limit: String(PAGE_SIZE + 1),
				offset: String(offset)
			});
			if (tapFilter) params.set('tap_id', tapFilter);
			const res = await fetch(`/api/pours?${params}`);
			if (!res.ok) throw new Error(`${res.status}`);
			const data: Pour[] = (await res.json()) ?? [];
			hasMore = data.length > PAGE_SIZE;
			pours = data.slice(0, PAGE_SIZE);
		} catch {
			error = 'Failed to load pours. Check that the server is running.';
		} finally {
			loading = false;
		}
	}

	async function applyFilter() {
		offset = 0;
		await loadPours();
	}

	async function deletePour(id: number) {
		if (!confirm('Delete this pour record?')) return;
		await adminFetch(`/api/pours/${id}`, { method: 'DELETE' });
		await loadPours();
	}

	function formatDate(iso: string) {
		return new Date(iso).toLocaleString();
	}
</script>

<div class="max-w-5xl">
	<h1 class="mb-6 text-2xl font-bold">Pour History</h1>

	<div class="mb-4 flex items-center gap-3">
		<select
			bind:value={tapFilter}
			class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
		>
			<option value="">All taps</option>
			{#each taps as tap (tap.id)}
				<option value={String(tap.id)}>Tap {tap.id}{tap.keg ? ` — ${tap.keg.beer_name}` : ''}</option>
			{/each}
		</select>
		<button
			onclick={applyFilter}
			class="rounded-lg bg-gray-800 px-4 py-2 text-sm font-semibold text-white hover:bg-gray-700"
		>
			Filter
		</button>
	</div>

	{#if error}
		<div class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading}
		<p class="text-gray-500">Loading…</p>
	{:else if pours.length === 0}
		<p class="text-gray-500">No pours recorded yet.</p>
	{:else}
		<div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
			<table class="w-full text-sm">
				<thead class="border-b border-gray-200 bg-gray-50">
					<tr>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Tap</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Volume</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Started</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Ended</th>
						<th class="px-4 py-3"></th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#each pours as pour (pour.id)}
						<tr class="hover:bg-gray-50">
							<td class="px-4 py-3 font-medium">#{pour.tap_id}</td>
							<td class="px-4 py-3">{pour.volume_ml.toFixed(0)} mL</td>
							<td class="px-4 py-3 text-gray-600">{formatDate(pour.started_at)}</td>
							<td class="px-4 py-3 text-gray-600">{formatDate(pour.ended_at)}</td>
							<td class="px-4 py-3 text-right">
								<button onclick={() => deletePour(pour.id)} class="text-red-600 hover:underline">
									Delete
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="mt-4 flex gap-3">
			<button
				onclick={() => { offset = Math.max(0, offset - PAGE_SIZE); loadPours(); }}
				disabled={offset === 0}
				class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium hover:bg-gray-50 disabled:opacity-40"
			>
				← Previous
			</button>
			<button
				onclick={() => { offset += PAGE_SIZE; loadPours(); }}
				disabled={!hasMore}
				class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium hover:bg-gray-50 disabled:opacity-40"
			>
				Next →
			</button>
		</div>
	{/if}
</div>
