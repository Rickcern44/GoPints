<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import type { Keg } from '$lib/api.js';

	let keg = $state<Keg | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');
	let success = $state('');

	let form = $state({ beer_name: '', brewery: '', style: '', abv: '', capacity_ml: '' });
	let imageFile = $state<File | null>(null);
	let hasImage = $state(false);

	onMount(async () => {
		try {
			const res = await fetch(`/api/kegs/${page.params.id}`);
			if (!res.ok) { error = 'Keg not found'; return; }
			keg = await res.json();
			form = {
				beer_name: keg!.beer_name,
				brewery: keg!.brewery,
				style: keg!.style,
				abv: keg!.abv ? String(keg!.abv) : '',
				capacity_ml: String(keg!.capacity_ml)
			};
			hasImage = !!keg!.image_mime_type;
		} catch {
			error = 'Failed to load keg. Check that the server is running.';
		} finally {
			loading = false;
		}
	});

	async function save() {
		error = ''; success = '';
		saving = true;
		try {
			const res = await adminFetch(`/api/kegs/${page.params.id}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					beer_name: form.beer_name,
					brewery: form.brewery,
					style: form.style,
					abv: parseFloat(form.abv) || 0,
					capacity_ml: parseFloat(form.capacity_ml)
				})
			});
			if (!res.ok) { error = 'Failed to save'; return; }
			success = 'Saved';
			keg = await res.json();
		} finally {
			saving = false;
		}
	}

	async function uploadImage() {
		if (!imageFile) return;
		error = ''; success = '';
		saving = true;
		try {
			const res = await adminFetch(`/api/kegs/${page.params.id}/image`, {
				method: 'PUT',
				headers: { 'Content-Type': imageFile.type },
				body: imageFile
			});
			if (!res.ok) { error = 'Image upload failed'; return; }
			hasImage = true;
			success = 'Image uploaded';
			imageFile = null;
		} finally {
			saving = false;
		}
	}

	async function deleteImage() {
		if (!confirm('Remove this image?')) return;
		await adminFetch(`/api/kegs/${page.params.id}/image`, { method: 'DELETE' });
		hasImage = false;
	}
</script>

<div class="max-w-xl">
	<a href="/admin/kegs" class="mb-6 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900">
		← Back to Kegs
	</a>

	{#if loading}
		<p class="text-gray-500">Loading…</p>
	{:else if error && !keg}
		<p class="text-red-600">{error}</p>
	{:else if keg}
		<h1 class="mb-6 text-2xl font-bold">{keg.beer_name}</h1>

		{#if error}
			<div class="mb-4 rounded-lg bg-red-50 px-4 py-2 text-sm text-red-700">{error}</div>
		{/if}
		{#if success}
			<div class="mb-4 rounded-lg bg-green-50 px-4 py-2 text-sm text-green-700">{success}</div>
		{/if}

		<div class="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-lg font-semibold">Details</h2>
			<div class="grid grid-cols-2 gap-4">
				<label class="col-span-2 flex flex-col gap-1 text-sm font-medium text-gray-700">
					Beer Name
					<input bind:value={form.beer_name} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Brewery
					<input bind:value={form.brewery} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Style
					<input bind:value={form.style} class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					ABV %
					<input bind:value={form.abv} type="number" step="0.1" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" />
				</label>
				<label class="flex flex-col gap-1 text-sm font-medium text-gray-700">
					Capacity (mL)
					<input bind:value={form.capacity_ml} type="number" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none" />
				</label>
			</div>
			<div class="mt-4 flex justify-end">
				<button
					onclick={save}
					disabled={saving}
					class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
				>
					{saving ? 'Saving…' : 'Save Changes'}
				</button>
			</div>
		</div>

		<div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-lg font-semibold">Image</h2>
			{#if hasImage}
				<img
					src="/api/kegs/{keg.id}/image"
					alt={keg.beer_name}
					class="mb-4 h-32 w-32 rounded-lg object-cover"
				/>
				<button onclick={deleteImage} class="text-sm text-red-600 hover:underline">
					Remove image
				</button>
			{:else}
				<p class="mb-3 text-sm text-gray-500">No image uploaded.</p>
				<input
					type="file"
					accept="image/*"
					onchange={(e) => { imageFile = (e.currentTarget as HTMLInputElement).files?.[0] ?? null; }}
					class="mb-3 block text-sm"
				/>
				<button
					onclick={uploadImage}
					disabled={!imageFile || saving}
					class="rounded-lg bg-gray-800 px-4 py-2 text-sm font-semibold text-white hover:bg-gray-700 disabled:opacity-40"
				>
					Upload
				</button>
			{/if}
		</div>
	{/if}
</div>

