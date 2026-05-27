<script lang="ts">
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import Spinner from '$lib/components/Spinner.svelte';
	import { recordManualPour, fetchFeatures } from '$lib/api.js';
	import type { Tap, Keg } from '$lib/api.js';

	let taps = $state<Tap[]>([]);
	let kegs = $state<Keg[]>([]);
	let loading = $state(true);
	let error = $state('');
	let selections = $state<Record<number, string>>({});
	let assigning = $state<Record<number, boolean>>({});

	let showAdd = $state(false);
	let newTapId = $state('');
	let addError = $state('');
	let adding = $state(false);

	let deleting = $state<Record<number, boolean>>({});

	const POUR_PRESETS = [
		{ label: '12 oz', ml: 12 * 29.5735 },
		{ label: '16 oz', ml: 16 * 29.5735 },
		{ label: '24 oz', ml: 24 * 29.5735 }
	] as const;

	let flowEnabled = $state(false);
	let pourOpen = $state<Record<number, boolean>>({});
	let pouring = $state<Record<number, boolean>>({});
	let pourSuccess = $state<Record<number, boolean>>({});

	onMount(async () => {
		await loadAll();
		const features = await fetchFeatures();
		flowEnabled = features.flow_based_pour;
	});

	async function loadAll() {
		loading = true;
		error = '';
		try {
			const [tapsRes, kegsRes] = await Promise.all([fetch('/api/taps'), fetch('/api/kegs')]);
			if (!tapsRes.ok || !kegsRes.ok) throw new Error('fetch failed');
			taps = (await tapsRes.json()) ?? [];
			kegs = (await kegsRes.json()) ?? [];
			syncSelections();
		} catch {
			error = 'Failed to load. Check that the server is running.';
		} finally {
			loading = false;
		}
	}

	function syncSelections() {
		for (const tap of taps) {
			selections[tap.id] = tap.keg_id ? String(tap.keg_id) : '';
		}
	}

	async function assignKeg(tapId: number) {
		assigning[tapId] = true;
		try {
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
			syncSelections();
		} finally {
			assigning[tapId] = false;
		}
	}

	async function addTap() {
		addError = '';
		const id = parseInt(newTapId);
		if (!id || id < 1 || id > 255) {
			addError = 'Tap ID must be 1–255';
			return;
		}
		adding = true;
		try {
			const res = await adminFetch('/api/taps', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ id })
			});
			if (res.status === 409) { addError = 'Tap already exists'; return; }
			if (!res.ok) { addError = 'Failed to add tap'; return; }
			newTapId = '';
			showAdd = false;
			await loadAll();
		} finally {
			adding = false;
		}
	}

	async function deleteTap(id: number) {
		if (!confirm(`Delete tap #${id}? This will un-assign any keg on it.`)) return;
		deleting[id] = true;
		try {
			const res = await adminFetch(`/api/taps/${id}`, { method: 'DELETE' });
			if (!res.ok) { error = 'Failed to delete tap'; return; }
			await loadAll();
		} finally {
			deleting[id] = false;
		}
	}

	async function logPour(tapId: number, ml: number) {
		pourOpen[tapId] = false;
		pouring[tapId] = true;
		try {
			await recordManualPour(tapId, ml);
			pourSuccess[tapId] = true;
			setTimeout(() => (pourSuccess[tapId] = false), 2500);
		} finally {
			pouring[tapId] = false;
		}
	}
</script>

<div class="max-w-2xl">
	<div class="mb-6 flex items-center justify-between">
		<h1 class="text-2xl font-bold">Taps</h1>
		<button
			onclick={() => { showAdd = !showAdd; addError = ''; newTapId = ''; }}
			class="btn btn-primary"
		>
			{showAdd ? 'Cancel' : '+ Add Tap'}
		</button>
	</div>

	{#if showAdd}
		<div class="mb-6 rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
			<h2 class="mb-3 text-sm font-semibold text-gray-700">New Tap</h2>
			{#if addError}
				<div class="mb-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{addError}</div>
			{/if}
			<div class="flex items-end gap-3">
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Tap ID (1–255)
					<input
						type="number"
						min="1"
						max="255"
						bind:value={newTapId}
						onkeydown={(e) => e.key === 'Enter' && addTap()}
						placeholder="e.g. 3"
						class="w-28 rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
					/>
				</label>
				<button onclick={addTap} disabled={adding} class="btn btn-primary mb-0.5">
					{#if adding}<Spinner size={14} />{/if}
					{adding ? 'Adding…' : 'Add'}
				</button>
			</div>
		</div>
	{/if}

	{#if error}
		<div class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading}
		<div class="flex items-center gap-2 text-sm text-gray-500">
			<Spinner size={16} />
			Loading…
		</div>
	{:else if taps.length === 0}
		<p class="text-gray-500">No taps yet. Add one above.</p>
	{:else}
		<div class="space-y-3">
			{#each taps as tap (tap.id)}
				<div class="rounded-xl border border-gray-200 bg-white shadow-sm">
					<!-- ── Desktop layout: single flex row ── -->
					<div class="hidden sm:flex items-center gap-4 p-4">
						<div class="min-w-[52px] text-center">
							<div class="text-2xl font-black text-gray-800">#{tap.id}</div>
							<div class="text-xs text-gray-400">tap</div>
						</div>
						<div class="flex-1 min-w-0">
							{#if tap.keg}
								<div class="text-sm font-medium text-gray-900 truncate">{tap.keg.beer_name}</div>
								<div class="text-xs text-gray-500">{tap.keg.brewery || ''} · {(tap.keg.capacity_ml / 1000).toFixed(1)}L</div>
							{:else}
								<div class="text-sm italic text-gray-400">No keg assigned</div>
							{/if}
						</div>
						<div class="flex items-center gap-2 shrink-0">
							{#if !flowEnabled && tap.keg}
								<div class="pour-wrap">
									{#if pourOpen[tap.id]}
										<button class="pour-backdrop" onclick={() => (pourOpen[tap.id] = false)} aria-label="Close" tabindex="-1"></button>
									{/if}
									<div class="pour-trigger">
										{#if pourOpen[tap.id]}
											<div class="pour-flyout">
												{#each POUR_PRESETS as preset}
													<button class="preset-btn" disabled={pouring[tap.id]} onclick={() => logPour(tap.id, preset.ml)}>{preset.label}</button>
												{/each}
											</div>
										{/if}
										<button class="btn btn-pour" class:open={pourOpen[tap.id]} disabled={pouring[tap.id]} onclick={() => (pourOpen[tap.id] = !pourOpen[tap.id])}>
											{#if pouring[tap.id]}<Spinner size={13} />{/if}
											{pouring[tap.id] ? 'Logging…' : 'Log Pour'}
										</button>
									</div>
								</div>
							{/if}
							{#if pourSuccess[tap.id]}
								<span class="text-xs font-medium text-green-600">Logged!</span>
							{/if}
							<select bind:value={selections[tap.id]} class="rounded-lg border border-gray-300 px-2 py-1.5 text-sm focus:border-blue-500 focus:outline-none">
								<option value="">— Remove keg —</option>
								{#each kegs as keg (keg.id)}
									<option value={String(keg.id)}>{keg.beer_name}</option>
								{/each}
							</select>
							<button onclick={() => assignKeg(tap.id)} disabled={assigning[tap.id]} class="btn btn-secondary flex items-center gap-1.5">
								{#if assigning[tap.id]}<Spinner size={13} />{/if}
								Apply
							</button>
							<button onclick={() => deleteTap(tap.id)} disabled={deleting[tap.id]} class="btn btn-danger flex items-center gap-1.5">
								{#if deleting[tap.id]}<Spinner size={13} />{/if}
								Delete
							</button>
						</div>
					</div>

					<!-- ── Mobile layout: stacked ── -->
					<div class="sm:hidden p-4">
						<!-- Header: tap number + keg name -->
						<div class="flex items-center gap-3 mb-3">
							<div class="text-center min-w-[40px]">
								<div class="text-xl font-black text-gray-800">#{tap.id}</div>
								<div class="text-xs text-gray-400">tap</div>
							</div>
							<div class="flex-1 min-w-0">
								{#if tap.keg}
									<div class="text-sm font-medium text-gray-900 truncate">{tap.keg.beer_name}</div>
									<div class="text-xs text-gray-500">{tap.keg.brewery || ''} · {(tap.keg.capacity_ml / 1000).toFixed(1)}L</div>
								{:else}
									<div class="text-sm italic text-gray-400">No keg assigned</div>
								{/if}
							</div>
							{#if pourSuccess[tap.id]}
								<span class="text-xs font-medium text-green-600">Logged!</span>
							{/if}
						</div>

						<!-- Keg select (full width) -->
						<select
							bind:value={selections[tap.id]}
							class="w-full rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:border-blue-500 focus:outline-none mb-2"
						>
							<option value="">— Remove keg —</option>
							{#each kegs as keg (keg.id)}
								<option value={String(keg.id)}>{keg.beer_name}</option>
							{/each}
						</select>

						<!-- Action buttons row -->
						<div class="flex gap-2">
							<button onclick={() => assignKeg(tap.id)} disabled={assigning[tap.id]} class="btn btn-secondary flex-1 flex items-center justify-center gap-1.5">
								{#if assigning[tap.id]}<Spinner size={13} />{/if}
								Apply
							</button>
							{#if !flowEnabled && tap.keg}
								<div class="pour-wrap">
									{#if pourOpen[tap.id]}
										<button class="pour-backdrop" onclick={() => (pourOpen[tap.id] = false)} aria-label="Close" tabindex="-1"></button>
									{/if}
									<div class="pour-trigger">
										{#if pourOpen[tap.id]}
											<div class="pour-flyout">
												{#each POUR_PRESETS as preset}
													<button class="preset-btn" disabled={pouring[tap.id]} onclick={() => logPour(tap.id, preset.ml)}>{preset.label}</button>
												{/each}
											</div>
										{/if}
										<button class="btn btn-pour" class:open={pourOpen[tap.id]} disabled={pouring[tap.id]} onclick={() => (pourOpen[tap.id] = !pourOpen[tap.id])}>
											{#if pouring[tap.id]}<Spinner size={13} />{/if}
											{pouring[tap.id] ? 'Logging…' : 'Log Pour'}
										</button>
									</div>
								</div>
							{/if}
							<button onclick={() => deleteTap(tap.id)} disabled={deleting[tap.id]} class="btn btn-danger flex items-center gap-1.5">
								{#if deleting[tap.id]}<Spinner size={13} />{/if}
								Delete
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.4rem 0.875rem;
		border-radius: 8px;
		font-size: 0.8125rem;
		font-weight: 600;
		cursor: pointer;
		border: none;
		transition: background-color 0.15s, box-shadow 0.15s, opacity 0.15s;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-primary {
		background: #2563eb;
		color: #fff;
	}

	.btn-primary:not(:disabled):hover {
		background: #1d4ed8;
		box-shadow: 0 1px 4px rgba(37, 99, 235, 0.35);
	}

	.btn-secondary {
		background: #f3f4f6;
		color: #374151;
		border: 1px solid #e5e7eb;
	}

	.btn-secondary:not(:disabled):hover {
		background: #e5e7eb;
	}

	.btn-danger {
		background: transparent;
		color: #dc2626;
		border: 1px solid #fecaca;
	}

	.btn-danger:not(:disabled):hover {
		background: #fef2f2;
		border-color: #fca5a5;
	}

	.btn-pour {
		background: rgba(200, 130, 26, 0.1);
		color: #92400e;
		border: 1px solid rgba(200, 130, 26, 0.3);
	}

	.btn-pour:not(:disabled):hover,
	.btn-pour.open {
		background: rgba(200, 130, 26, 0.18);
		border-color: rgba(200, 130, 26, 0.5);
	}

	/* Pour flyout */

	.pour-wrap {
		position: relative;
	}

	.pour-backdrop {
		position: fixed;
		inset: 0;
		z-index: 10;
		background: transparent;
		border: none;
		cursor: default;
		padding: 0;
	}

	.pour-trigger {
		position: relative;
		z-index: 11;
	}

	.pour-flyout {
		position: absolute;
		bottom: calc(100% + 0.5rem);
		right: 0;
		display: flex;
		gap: 0.3rem;
		background: #fff;
		border: 1px solid #e5e7eb;
		border-radius: 10px;
		padding: 0.4rem;
		box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
		white-space: nowrap;
	}

	.preset-btn {
		padding: 0.45rem 0.75rem;
		border-radius: 7px;
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		color: #374151;
		font-size: 0.8125rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s;
	}

	.preset-btn:not(:disabled):hover {
		background: rgba(200, 130, 26, 0.08);
		border-color: rgba(200, 130, 26, 0.35);
		color: #92400e;
	}

	.preset-btn:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}
</style>
