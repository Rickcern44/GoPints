<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import Spinner from '$lib/components/Spinner.svelte';
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

<div class="max-w-3xl">
	<div class="console-heading">
		<h1>Pour Log</h1>
		<span class="count">{pours.length} entries</span>
	</div>

	<div class="mb-5 flex items-center gap-3">
		<label class="field flex-1 sm:flex-none">
			<span class="cap">Tap</span>
			<select bind:value={tapFilter}>
				<option value="">All taps</option>
				{#each taps as tap (tap.id)}
					<option value={String(tap.id)}>Tap {tap.id}{tap.keg ? ` — ${tap.keg.beer_name}` : ''}</option>
				{/each}
			</select>
		</label>
		<button onclick={applyFilter} class="btn-console self-end">Filter</button>
	</div>

	{#if error}
		<div class="mb-4 rounded-lg bg-error-bg px-4 py-3 text-sm text-error">{error}</div>
	{/if}

	{#if loading}
		<div class="flex items-center gap-2 text-sm text-fg-muted"><Spinner size={16} /> Loading…</div>
	{:else if pours.length === 0}
		<p class="text-fg-muted">No pours recorded yet.</p>
	{:else}
		<div class="log-panel">
			{#each pours as pour (pour.id)}
				<div class="log-row flex items-center gap-4">
					<div class="readout-badge" style="--rb-size: 2.15rem; font-size: 0.85rem;">{pour.tap_id}</div>
					<div class="flex-1 min-w-0">
						<div class="font-mono text-base font-bold text-fg">{pour.volume_ml.toFixed(0)} mL</div>
						<div class="text-xs font-mono text-fg-muted">
							{formatDate(pour.started_at)} → {formatDate(pour.ended_at)}
						</div>
					</div>
					<button onclick={() => deletePour(pour.id)} class="btn-console-ghost danger">Void</button>
				</div>
			{/each}
		</div>

		<div class="mt-5 flex gap-3">
			<button
				onclick={() => { offset = Math.max(0, offset - PAGE_SIZE); loadPours(); }}
				disabled={offset === 0}
				class="btn-console-ghost"
			>
				← Previous
			</button>
			<button
				onclick={() => { offset += PAGE_SIZE; loadPours(); }}
				disabled={!hasMore}
				class="btn-console-ghost"
			>
				Next →
			</button>
		</div>
	{/if}
</div>

<style>
	.log-panel {
		background: var(--color-panel);
		border: 1px solid var(--color-line);
		border-top: 2px solid var(--color-accent);
		border-radius: 3px;
	}
</style>
