<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { adminFetch } from '$lib/admin.js';
	import Spinner from '$lib/components/Spinner.svelte';
	import { KEG_SIZES, kegSizeLabel } from '$lib/api.js';
	import type { Keg } from '$lib/api.js';

	let kegs = $state<Keg[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showCreate = $state(false);
	let form = $state({ beer_name: '', brewery: '', style: '', abv: '', capacity: '', custom_ml: '' });
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
		const capacity_ml =
			form.capacity === 'custom' ? parseFloat(form.custom_ml) : parseInt(form.capacity);
		if (!form.beer_name || !form.capacity || isNaN(capacity_ml) || capacity_ml <= 0) {
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
					capacity_ml
				})
			});
			if (!res.ok) { formError = 'Failed to create keg'; return; }
			form = { beer_name: '', brewery: '', style: '', abv: '', capacity: '', custom_ml: '' };
			showCreate = false;
			await loadKegs();
		} finally {
			saving = false;
		}
	}

	async function deleteKeg(id: number) {
		if (!confirm('Delete this keg?')) return;
		const res = await adminFetch(`/api/kegs/${id}`, { method: 'DELETE' });
		if (!res.ok) {
			error = 'Failed to delete keg';
			return;
		}
		await loadKegs();
	}
</script>

<div class="max-w-4xl">
	<div class="console-heading">
		<h1>Keg Stock</h1>
		<span class="count">{kegs.length} on file</span>
		<button onclick={() => (showCreate = !showCreate)} class="btn-console">
			{showCreate ? 'Cancel' : '+ Register Keg'}
		</button>
	</div>

	{#if showCreate}
		<div class="panel mb-6">
			<h2>Register Keg</h2>
			{#if formError}
				<div class="mb-4 rounded-lg bg-error-bg px-4 py-2 text-sm text-error">{formError}</div>
			{/if}
			<div class="grid grid-cols-1 gap-x-6 gap-y-4 sm:grid-cols-2">
				<label class="field col-span-1 sm:col-span-2">
					<span class="cap">Beer Name *</span>
					<input bind:value={form.beer_name} placeholder="Hazy IPA" />
				</label>
				<label class="field">
					<span class="cap">Brewery</span>
					<input bind:value={form.brewery} placeholder="Local Brewing Co." />
				</label>
				<label class="field">
					<span class="cap">Style</span>
					<input bind:value={form.style} placeholder="IPA" />
				</label>
				<label class="field">
					<span class="cap">ABV %</span>
					<input bind:value={form.abv} type="number" step="0.1" placeholder="6.5" />
				</label>
				<div class="field">
					<span class="cap">Keg Size *</span>
					<select bind:value={form.capacity}>
						<option value="">— Select size —</option>
						{#each KEG_SIZES as size (size.ml)}
							<option value={String(size.ml)}>{size.label}</option>
						{/each}
						<option value="custom">Custom…</option>
					</select>
					{#if form.capacity === 'custom'}
						<input bind:value={form.custom_ml} type="number" min="1" placeholder="Capacity in mL" />
					{/if}
				</div>
			</div>
			<div class="mt-5 flex justify-end">
				<button onclick={createKeg} disabled={saving} class="btn-console">
					{#if saving}<Spinner size={14} />{/if}
					{saving ? 'Registering…' : 'Register'}
				</button>
			</div>
		</div>
	{/if}

	{#if error}
		<div class="mb-4 rounded-lg bg-error-bg px-4 py-3 text-sm text-error">{error}</div>
	{/if}

	{#if loading}
		<div class="flex items-center gap-2 text-sm text-fg-muted"><Spinner size={16} /> Loading…</div>
	{:else if kegs.length === 0}
		<p class="text-fg-muted">No kegs registered yet. Register one above.</p>
	{:else}
		<div class="panel !p-0 overflow-hidden">
			<table class="data-table">
				<thead>
					<tr>
						<th>Beer</th>
						<th class="hidden sm:table-cell">Brewery</th>
						<th class="hidden sm:table-cell">Style</th>
						<th class="hidden sm:table-cell">ABV</th>
						<th class="hidden sm:table-cell">Capacity</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each kegs as keg (keg.id)}
						<tr class="sm:cursor-default cursor-pointer" onclick={() => goto(`/admin/kegs/${keg.id}`)}>
							<td class="font-bold">
								<span>{keg.beer_name}</span>
								{#if keg.brewery || keg.style}
									<div class="sm:hidden mt-0.5 text-xs text-fg-muted">
										{[keg.brewery, keg.style].filter(Boolean).join(' · ')}
									</div>
								{/if}
							</td>
							<td class="hidden sm:table-cell text-fg-muted">{keg.brewery || '—'}</td>
							<td class="hidden sm:table-cell text-fg-muted">{keg.style || '—'}</td>
							<td class="hidden sm:table-cell text-fg-muted">{keg.abv ? `${keg.abv}%` : '—'}</td>
							<td class="hidden sm:table-cell text-fg-muted">{kegSizeLabel(keg.capacity_ml)}</td>
							<td class="text-right whitespace-nowrap">
								<!-- Desktop: Edit + Delete -->
								<a
									href="/admin/kegs/{keg.id}"
									onclick={(e) => e.stopPropagation()}
									class="btn-console-ghost hidden sm:inline-flex mr-2"
								>Edit</a>
								<button
									onclick={(e) => { e.stopPropagation(); deleteKeg(keg.id); }}
									class="btn-console-ghost danger hidden sm:inline-flex"
								>Delete</button>
								<!-- Mobile: chevron -->
								<svg class="sm:hidden inline-block text-fg-muted" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
									<path d="M9 18l6-6-6-6"/>
								</svg>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

