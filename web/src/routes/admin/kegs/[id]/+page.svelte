<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { adminFetch } from '$lib/admin.js';
	import Spinner from '$lib/components/Spinner.svelte';
	import { KEG_SIZES, fetchFeatures } from '$lib/api.js';
	import type { Keg } from '$lib/api.js';

	let keg = $state<Keg | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');
	let success = $state('');

	let form = $state({ beer_name: '', brewery: '', style: '', abv: '', capacity: '', custom_ml: '', image_style: 'circle' });
	let imageFile = $state<File | null>(null);
	let imagePreview = $state<string | null>(null);
	let hasImage = $state(false);
	let dragOver = $state(false);
	let fileInput = $state<HTMLInputElement | undefined>(undefined);

	let breweryFile = $state<File | null>(null);
	let breweryPreview = $state<string | null>(null);
	let hasBreweryImage = $state(false);
	let breweryDragOver = $state(false);
	let breweryFileInput = $state<HTMLInputElement | undefined>(undefined);

	let remoteImageUrlsEnabled = $state(false);
	let imageMode = $state<'upload' | 'url'>('upload');
	let breweryMode = $state<'upload' | 'url'>('upload');
	let imageUrlValue = $state('');
	let breweryUrlValue = $state('');
	let importingImage = $state(false);
	let importingBreweryImage = $state(false);

	onMount(async () => {
		try {
			const res = await fetch(`/api/kegs/${page.params.id}`);
			if (!res.ok) { error = 'Keg not found'; return; }
			keg = await res.json();
			const matchedSize = KEG_SIZES.find((s) => s.ml === keg!.capacity_ml);
			form = {
				beer_name: keg!.beer_name,
				brewery: keg!.brewery,
				style: keg!.style,
				abv: keg!.abv ? String(keg!.abv) : '',
				capacity: matchedSize ? String(matchedSize.ml) : 'custom',
				custom_ml: matchedSize ? '' : String(keg!.capacity_ml),
				image_style: keg!.image_style || 'circle'
			};
			hasImage = !!keg!.image_mime_type;
			hasBreweryImage = !!keg!.brewery_image_mime_type;
			const features = await fetchFeatures();
			remoteImageUrlsEnabled = features.remote_image_urls;
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
					capacity_ml:
						form.capacity === 'custom'
							? parseFloat(form.custom_ml)
							: parseInt(form.capacity),
					image_style: form.image_style
				})
			});
			if (!res.ok) { error = 'Failed to save'; return; }
			success = 'Saved';
			keg = await res.json();
		} finally {
			saving = false;
		}
	}

	function handleFileSelect(file: File | null | undefined) {
		if (!file) return;
		if (imagePreview) URL.revokeObjectURL(imagePreview);
		imageFile = file;
		imagePreview = URL.createObjectURL(file);
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragOver = false;
		const file = e.dataTransfer?.files[0];
		if (file?.type.startsWith('image/')) handleFileSelect(file);
	}

	function clearSelection() {
		imageFile = null;
		if (imagePreview) {
			URL.revokeObjectURL(imagePreview);
			imagePreview = null;
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
			clearSelection();
		} finally {
			saving = false;
		}
	}

	async function deleteImage() {
		if (!confirm('Remove this image?')) return;
		await adminFetch(`/api/kegs/${page.params.id}/image`, { method: 'DELETE' });
		hasImage = false;
	}

	async function importImageFromUrl() {
		if (!imageUrlValue.trim()) return;
		error = ''; success = '';
		importingImage = true;
		try {
			const res = await adminFetch(`/api/kegs/${page.params.id}/image/from-url`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ url: imageUrlValue.trim() })
			});
			if (!res.ok) {
				const body = await res.json().catch(() => ({}));
				error = body.error || 'Could not import image from that URL';
				return;
			}
			hasImage = true;
			success = 'Image imported';
			imageUrlValue = '';
		} finally {
			importingImage = false;
		}
	}

	async function uploadBreweryImage() {
		if (!breweryFile) return;
		error = ''; success = '';
		saving = true;
		try {
			const res = await adminFetch(`/api/kegs/${page.params.id}/brewery-image`, {
				method: 'PUT',
				headers: { 'Content-Type': breweryFile.type },
				body: breweryFile
			});
			if (!res.ok) { error = 'Brewery image upload failed'; return; }
			hasBreweryImage = true;
			success = 'Brewery image uploaded';
			clearBrewerySelection();
		} finally {
			saving = false;
		}
	}

	async function deleteBreweryImage() {
		if (!confirm('Remove brewery image?')) return;
		await adminFetch(`/api/kegs/${page.params.id}/brewery-image`, { method: 'DELETE' });
		hasBreweryImage = false;
	}

	async function importBreweryImageFromUrl() {
		if (!breweryUrlValue.trim()) return;
		error = ''; success = '';
		importingBreweryImage = true;
		try {
			const res = await adminFetch(`/api/kegs/${page.params.id}/brewery-image/from-url`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ url: breweryUrlValue.trim() })
			});
			if (!res.ok) {
				const body = await res.json().catch(() => ({}));
				error = body.error || 'Could not import image from that URL';
				return;
			}
			hasBreweryImage = true;
			success = 'Brewery image imported';
			breweryUrlValue = '';
		} finally {
			importingBreweryImage = false;
		}
	}

	function handleBreweryFileSelect(file: File | null | undefined) {
		if (!file) return;
		if (breweryPreview) URL.revokeObjectURL(breweryPreview);
		breweryFile = file;
		breweryPreview = URL.createObjectURL(file);
	}

	function handleBreweryDrop(e: DragEvent) {
		e.preventDefault();
		breweryDragOver = false;
		const file = e.dataTransfer?.files[0];
		if (file?.type.startsWith('image/')) handleBreweryFileSelect(file);
	}

	function clearBrewerySelection() {
		breweryFile = null;
		if (breweryPreview) {
			URL.revokeObjectURL(breweryPreview);
			breweryPreview = null;
		}
	}
</script>

<div class="max-w-xl">
	<a href="/admin/kegs" class="mb-6 inline-flex items-center gap-1 text-sm text-fg-muted hover:text-fg">
		← Back to Kegs
	</a>

	{#if loading}
		<div class="flex items-center gap-2 text-sm text-fg-muted"><Spinner size={16} /> Loading…</div>
	{:else if error && !keg}
		<p class="text-error">{error}</p>
	{:else if keg}
		<div class="console-heading">
			<h1>{keg.beer_name}</h1>
			<span class="count">KEG-{keg.id.toString().padStart(3, '0')}</span>
		</div>

		{#if error}
			<div class="mb-4 rounded-lg bg-error-bg px-4 py-2 text-sm text-error">{error}</div>
		{/if}
		{#if success}
			<div class="mb-4 rounded-lg bg-success-bg px-4 py-2 text-sm text-success">{success}</div>
		{/if}

		<div class="panel mb-6">
			<h2>Details</h2>
			<div class="grid grid-cols-2 gap-x-6 gap-y-4">
				<label class="field col-span-2">
					<span class="cap">Beer Name</span>
					<input bind:value={form.beer_name} />
				</label>
				<label class="field">
					<span class="cap">Brewery</span>
					<input bind:value={form.brewery} />
				</label>
				<label class="field">
					<span class="cap">Style</span>
					<input bind:value={form.style} />
				</label>
				<label class="field">
					<span class="cap">ABV %</span>
					<input bind:value={form.abv} type="number" step="0.1" />
				</label>
				<div class="field">
					<span class="cap">Keg Size</span>
					<select bind:value={form.capacity}>
						{#each KEG_SIZES as size (size.ml)}
							<option value={String(size.ml)}>{size.label}</option>
						{/each}
						<option value="custom">Custom…</option>
					</select>
					{#if form.capacity === 'custom'}
						<input bind:value={form.custom_ml} type="number" min="1" placeholder="Capacity in mL" />
					{/if}
				</div>

					<div class="col-span-2 flex flex-col gap-1">
						<span class="cap">Beer Image Style</span>
						<div class="flex gap-2 mt-1">
							{#each [
								{ value: 'circle', label: 'Circle' },
								{ value: 'square', label: 'Square' },
								{ value: 'can', label: 'Can' }
							] as opt (opt.value)}
								<button
									type="button"
									onclick={() => (form.image_style = opt.value)}
									class="style-opt {form.image_style === opt.value ? 'active' : ''}"
								>
									<div class="style-preview style-preview-{opt.value}"></div>
									<span>{opt.label}</span>
								</button>
							{/each}
						</div>
					</div>
				</div>

			<div class="mt-5 flex justify-end">
				<button onclick={save} disabled={saving} class="btn-console">
					{#if saving}<Spinner size={14} />{/if}
					{saving ? 'Saving…' : 'Save Changes'}
				</button>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<!-- Beer Image -->
			<div class="panel !p-5">
				<h2 class="!mb-3 !text-base">Beer Image</h2>
				{#if hasImage}
					<div class="flex flex-col items-center gap-3">
						<img
							src="/api/kegs/{keg.id}/image"
							alt={keg.beer_name}
							class="h-20 w-20 border border-line object-cover
								{form.image_style === 'circle' ? 'rounded-full' : form.image_style === 'can' ? 'rounded-[50%_50%_50%_50%/10%_10%_10%_10%]' : 'rounded-lg'}"
							style={form.image_style === 'can' ? 'width:48px;height:96px' : ''}
						/>
						<button
							onclick={deleteImage}
							class="text-xs text-error transition-colors hover:brightness-90 hover:underline"
						>Remove</button>
					</div>
				{:else}
					{#if remoteImageUrlsEnabled}
						<div class="mode-toggle">
							<button type="button" class="mode-btn {imageMode === 'upload' ? 'active' : ''}" onclick={() => (imageMode = 'upload')}>Upload</button>
							<button type="button" class="mode-btn {imageMode === 'url' ? 'active' : ''}" onclick={() => (imageMode = 'url')}>Paste URL</button>
						</div>
					{/if}
					{#if remoteImageUrlsEnabled && imageMode === 'url'}
						<div class="url-import">
							<input type="url" placeholder="https://example.com/beer.png" bind:value={imageUrlValue}
								class="url-input" />
							<button onclick={importImageFromUrl} disabled={importingImage || !imageUrlValue.trim()}
								class="inline-flex items-center gap-1.5 rounded-[3px] bg-accent px-3 py-1.5 text-xs font-mono font-semibold text-void transition-all hover:brightness-110 disabled:opacity-40">
								{#if importingImage}<Spinner size={12} />{/if}
								{importingImage ? 'Importing…' : 'Import'}
							</button>
						</div>
					{:else}
						<input type="file" accept="image/*" bind:this={fileInput}
							onchange={(e) => handleFileSelect((e.currentTarget as HTMLInputElement).files?.[0])}
							class="sr-only" />
						<div class="dropzone" class:drag-over={dragOver} role="button" tabindex="0"
							onclick={() => fileInput?.click()}
							onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && fileInput?.click()}
							ondragover={(e) => { e.preventDefault(); dragOver = true; }}
							ondragleave={(e) => { const r = e.relatedTarget as Node|null; if (!r || !(e.currentTarget as HTMLElement).contains(r)) dragOver = false; }}
							ondrop={handleDrop}>
							{#if imageFile && imagePreview}
								<button class="clear-btn" onclick={(e) => { e.stopPropagation(); clearSelection(); }} aria-label="Clear">✕</button>
								<img src={imagePreview} alt="Preview" class="preview-img" />
								<div class="file-meta">
									<span class="file-name">{imageFile.name}</span>
									<span class="file-size">{(imageFile.size / 1024).toFixed(0)} KB</span>
								</div>
							{:else}
								<div class="drop-icon"><svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg></div>
								<p class="drop-label">Drop or <span class="browse-link">browse</span></p>
								<p class="drop-hint">PNG, JPG, WebP</p>
							{/if}
						</div>
						{#if imageFile}
							<div class="mt-2 flex justify-end">
								<button onclick={uploadImage} disabled={saving}
									class="inline-flex items-center gap-1.5 rounded-[3px] bg-accent px-3 py-1.5 text-xs font-mono font-semibold text-void transition-all hover:brightness-110 disabled:opacity-40">
									{#if saving}<Spinner size={12} />{/if}
									{saving ? 'Uploading…' : 'Upload'}
								</button>
							</div>
						{/if}
					{/if}
				{/if}
			</div>

			<!-- Brewery Image -->
			<div class="panel !p-5">
				<h2 class="!mb-3 !text-base">Brewery Logo</h2>
				{#if hasBreweryImage}
					<div class="flex flex-col items-center gap-3">
						<img
							src="/api/kegs/{keg.id}/brewery-image"
							alt="{keg.brewery} logo"
							class="h-20 w-20 rounded-lg border border-line object-contain p-1"
						/>
						<button
							onclick={deleteBreweryImage}
							class="text-xs text-error transition-colors hover:brightness-90 hover:underline"
						>Remove</button>
					</div>
				{:else}
					{#if remoteImageUrlsEnabled}
						<div class="mode-toggle">
							<button type="button" class="mode-btn {breweryMode === 'upload' ? 'active' : ''}" onclick={() => (breweryMode = 'upload')}>Upload</button>
							<button type="button" class="mode-btn {breweryMode === 'url' ? 'active' : ''}" onclick={() => (breweryMode = 'url')}>Paste URL</button>
						</div>
					{/if}
					{#if remoteImageUrlsEnabled && breweryMode === 'url'}
						<div class="url-import">
							<input type="url" placeholder="https://example.com/logo.png" bind:value={breweryUrlValue}
								class="url-input" />
							<button onclick={importBreweryImageFromUrl} disabled={importingBreweryImage || !breweryUrlValue.trim()}
								class="inline-flex items-center gap-1.5 rounded-[3px] bg-accent px-3 py-1.5 text-xs font-mono font-semibold text-void transition-all hover:brightness-110 disabled:opacity-40">
								{#if importingBreweryImage}<Spinner size={12} />{/if}
								{importingBreweryImage ? 'Importing…' : 'Import'}
							</button>
						</div>
					{:else}
						<input type="file" accept="image/*" bind:this={breweryFileInput}
							onchange={(e) => handleBreweryFileSelect((e.currentTarget as HTMLInputElement).files?.[0])}
							class="sr-only" />
						<div class="dropzone" class:drag-over={breweryDragOver} role="button" tabindex="0"
							onclick={() => breweryFileInput?.click()}
							onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && breweryFileInput?.click()}
							ondragover={(e) => { e.preventDefault(); breweryDragOver = true; }}
							ondragleave={(e) => { const r = e.relatedTarget as Node|null; if (!r || !(e.currentTarget as HTMLElement).contains(r)) breweryDragOver = false; }}
							ondrop={handleBreweryDrop}>
							{#if breweryFile && breweryPreview}
								<button class="clear-btn" onclick={(e) => { e.stopPropagation(); clearBrewerySelection(); }} aria-label="Clear">✕</button>
								<img src={breweryPreview} alt="Preview" class="preview-img" />
								<div class="file-meta">
									<span class="file-name">{breweryFile.name}</span>
									<span class="file-size">{(breweryFile.size / 1024).toFixed(0)} KB</span>
								</div>
							{:else}
								<div class="drop-icon"><svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg></div>
								<p class="drop-label">Drop or <span class="browse-link">browse</span></p>
								<p class="drop-hint">PNG, JPG, WebP</p>
							{/if}
						</div>
						{#if breweryFile}
							<div class="mt-2 flex justify-end">
								<button onclick={uploadBreweryImage} disabled={saving}
									class="inline-flex items-center gap-1.5 rounded-[3px] bg-accent px-3 py-1.5 text-xs font-mono font-semibold text-void transition-all hover:brightness-110 disabled:opacity-40">
									{#if saving}<Spinner size={12} />{/if}
									{saving ? 'Uploading…' : 'Upload'}
								</button>
							</div>
						{/if}
					{/if}
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	/* Style picker */
	.style-opt {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.35rem;
		padding: 0.5rem 0.75rem;
		border-radius: 3px;
		border: 1.5px solid var(--color-line);
		background: var(--color-panel-raised);
		cursor: pointer;
		font-size: 0.72rem;
		font-weight: 500;
		color: var(--color-fg-muted);
		transition: border-color 0.15s, color 0.15s, background 0.15s;
	}

	.style-opt:hover {
		border-color: color-mix(in srgb, var(--color-accent) 50%, transparent);
		color: var(--color-accent);
	}

	.style-opt.active {
		border-color: var(--color-accent);
		background: color-mix(in srgb, var(--color-accent) 10%, transparent);
		color: var(--color-accent);
	}

	.style-preview {
		background: var(--color-fg-muted);
		width: 32px;
		height: 32px;
	}

	.style-preview-circle {
		border-radius: 50%;
	}

	.style-preview-square {
		border-radius: 4px;
	}

	.style-preview-can {
		width: 20px;
		height: 36px;
		border-radius: 50% / 14%;
	}

	/* Upload/URL mode toggle */

	.mode-toggle {
		display: flex;
		gap: 0.25rem;
		margin-bottom: 0.75rem;
	}

	.mode-btn {
		flex: 1;
		padding: 0.35rem 0.5rem;
		border-radius: 3px;
		border: 1.5px solid var(--color-line);
		background: var(--color-panel-raised);
		cursor: pointer;
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-fg-muted);
		transition: border-color 0.15s, color 0.15s, background 0.15s;
	}

	.mode-btn:hover {
		border-color: color-mix(in srgb, var(--color-accent) 50%, transparent);
		color: var(--color-accent);
	}

	.mode-btn.active {
		border-color: var(--color-accent);
		background: color-mix(in srgb, var(--color-accent) 10%, transparent);
		color: var(--color-accent);
	}

	.url-import {
		display: flex;
		gap: 0.5rem;
	}

	.url-input {
		flex: 1;
		min-width: 0;
		border-radius: 3px;
		border: 1px solid var(--color-line);
		padding: 0.5rem 0.75rem;
		font-size: 0.85rem;
	}

	.url-input:focus {
		outline: none;
		border-color: var(--color-accent);
		box-shadow: 0 0 0 1px var(--color-accent);
	}

	.dropzone {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		min-height: 148px;
		padding: 1.75rem 1.5rem;
		border: 2px dashed var(--color-line);
		border-radius: 4px;
		cursor: pointer;
		text-align: center;
		transition: border-color 0.15s, background-color 0.15s;
		outline: none;
		user-select: none;
	}

	.dropzone:hover,
	.dropzone:focus-visible {
		border-color: color-mix(in srgb, var(--color-accent) 50%, transparent);
		background-color: color-mix(in srgb, var(--color-accent) 6%, var(--color-panel-raised));
	}

	.dropzone.drag-over {
		border-color: var(--color-accent);
		background-color: color-mix(in srgb, var(--color-accent) 10%, var(--color-panel-raised));
	}

	.drop-icon {
		color: var(--color-fg-muted);
		transition: color 0.15s;
	}

	.dropzone:hover .drop-icon,
	.dropzone.drag-over .drop-icon {
		color: var(--color-accent);
	}

	.drop-label {
		font-size: 0.875rem;
		color: var(--color-fg-muted);
		margin: 0;
	}

	.browse-link {
		color: var(--color-accent);
		text-decoration: underline;
		text-underline-offset: 2px;
	}

	.drop-hint {
		font-size: 0.72rem;
		color: var(--color-fg-muted);
		margin: 0;
	}

	/* Selected file state */

	.preview-img {
		width: 80px;
		height: 80px;
		border-radius: 3px;
		object-fit: cover;
		border: 1px solid var(--color-line);
	}

	.file-meta {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.15rem;
	}

	.file-name {
		font-size: 0.8rem;
		font-weight: 500;
		color: var(--color-fg);
		max-width: 220px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-size {
		font-size: 0.7rem;
		color: var(--color-fg-muted);
	}

	.clear-btn {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		width: 1.4rem;
		height: 1.4rem;
		border-radius: 50%;
		background: var(--color-fg-muted);
		color: #fff;
		font-size: 0.65rem;
		border: none;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: background-color 0.15s;
		line-height: 1;
	}

	.clear-btn:hover {
		background: var(--color-fg);
	}
</style>
