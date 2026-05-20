<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import type { Keg } from '$lib/api.js';

	let kegs = $state<Keg[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showCreate = $state(false);
	let form = $state({ beer_name: '', brewery: '', style: '', abv: '', capacity_ml: '' });
	let saving = $state(false);
	let formError = $state('');

	onMount(loadKegs);

	async function loadKegs() {
		loading = true;
		error = '';
		try {
			const res = await fetch('/api/kegs');
			if (!res.ok) throw new Error(`${res.status}`);
			kegs = (await res.json()) ?? [];
		} catch {
			error = 'Failed to load kegs. Check that the server is running.';
		} finally {
			loading = false;
		}
	}

	async function createKeg() {
		formError = '';
		if (!form.beer_name || !form.capacity_ml) {
			formError = 'Beer name and capacity are required';
			return;
		}
		saving = true;
		try {
			const res = await adminFetch('/api/kegs', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					beer_name: form.beer_name,
					brewery: form.brewery,
					style: form.style,
					abv: parseFloat(form.abv) || 0,
					capacity_ml: parseFloat(form.capacity_ml)
				})
			});
			if (!res.ok) { formError = 'Failed to create keg'; return; }
			form = { beer_name: '', brewery: '', style: '', abv: '', capacity_ml: '' };
			showCreate = false;
			await loadKegs();
		} finally {
			saving = false;
		}
	}

	async function deleteKeg(id: number) {
		if (!confirm('Delete this keg?')) return;
		await adminFetch(`/api/kegs/${id}`, { method: 'DELETE' });
		await loadKegs();
	}
</script>

<div class="max-w-4xl">
	<div class="mb-6 flex items-center justify-between">
		<h1 class="text-2xl font-bold">Kegs</h1>
		<button
			onclick={() => (showCreate = !showCreate)}
			class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700"
		>
			{showCreate ? 'Cancel' : '+ New Keg'}
		</button>
	</div>

	{#if showCreate}
		<div class="mb-6 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-lg font-semibold">New Keg</h2>
			{#if formError}
				<div class="mb-4 rounded-lg bg-red-50 px-4 py-2 text-sm text-red-700">{formError}</div>
			{/if}
			<div class="grid grid-cols-2 gap-4">
				<label class="col-span-2 flex flex-col gap-1 text-sm font-medium text-gray-700">
					Beer Name *
					<input bind:value={form.beer_name} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" placeholder="Hazy IPA" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Brewery
					<input bind:value={form.brewery} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" placeholder="Local Brewing Co." />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Style
					<input bind:value={form.style} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" placeholder="IPA" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					ABV %
					<input bind:value={form.abv} type="number" step="0.1" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" placeholder="6.5" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Capacity (mL) *
					<input bind:value={form.capacity_ml} type="number" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" placeholder="19000" />
				</label>
			</div>
			<div class="mt-4 flex justify-end">
				<button
					onclick={createKeg}
					disabled={saving}
					class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
				>
					{saving ? 'Creating…' : 'Create Keg'}
				</button>
			</div>
		</div>
	{/if}

	{#if error}
		<div class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading}
		<p class="text-gray-500">Loading…</p>
	{:else if kegs.length === 0}
		<p class="text-gray-500">No kegs yet. Create one above.</p>
	{:else}
		<div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
			<table class="w-full text-sm">
				<thead class="border-b border-gray-200 bg-gray-50">
					<tr>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Beer</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Brewery</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Style</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">ABV</th>
						<th class="px-4 py-3 text-left font-medium text-gray-600">Capacity</th>
						<th class="px-4 py-3"></th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#each kegs as keg (keg.id)}
						<tr class="hover:bg-gray-50">
							<td class="px-4 py-3 font-medium">{keg.beer_name}</td>
							<td class="px-4 py-3 text-gray-600">{keg.brewery || '—'}</td>
							<td class="px-4 py-3 text-gray-600">{keg.style || '—'}</td>
							<td class="px-4 py-3 text-gray-600">{keg.abv ? `${keg.abv}%` : '—'}</td>
							<td class="px-4 py-3 text-gray-600">{(keg.capacity_ml / 1000).toFixed(1)}L</td>
							<td class="px-4 py-3 text-right">
								<a
									href="/admin/kegs/{keg.id}"
									class="mr-3 text-blue-600 hover:underline"
								>Edit</a>
								<button
									onclick={() => deleteKeg(keg.id)}
									class="text-red-600 hover:underline"
								>Delete</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

