<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { setToken, getToken } from '$lib/admin.js';

	let isSetup = $state(false);
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let loading = $state(false);

	onMount(async () => {
		if (getToken()) {
			goto('/admin/kegs', { replaceState: true });
			return;
		}
		try {
			const res = await fetch('/api/admin/status');
			if (!res.ok) throw new Error(`${res.status}`);
			const data: { password_set: boolean } = await res.json();
			isSetup = !data.password_set;
		} catch {
			error = 'Could not reach the server. Check that it is running.';
		}
	});

	async function submit() {
		error = '';
		if (isSetup && password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}
		if (!password) {
			error = 'Password is required';
			return;
		}
		loading = true;
		try {
			const endpoint = isSetup ? '/api/admin/setup' : '/api/admin/login';
			const res = await fetch(endpoint, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password })
			});
			if (!res.ok) {
				error = isSetup ? 'Setup failed' : 'Incorrect password';
				return;
			}
			const data: { token: string } = await res.json();
			setToken(data.token);
			goto('/admin/kegs', { replaceState: true });
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<div class="w-full max-w-sm rounded-xl bg-white p-8 shadow-md">
		<h1 class="mb-1 text-2xl font-bold text-gray-900">GoPints Admin</h1>
		<p class="mb-6 text-sm text-gray-500">
			{isSetup ? 'Set an admin password to get started.' : 'Sign in to manage your kegerator.'}
		</p>

		{#if error}
			<div class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
		{/if}

		<form onsubmit={(e) => { e.preventDefault(); submit(); }} class="flex flex-col gap-4">
			<div>
				<label for="password" class="mb-1 block text-sm font-medium text-gray-700">
					{isSetup ? 'New Password' : 'Password'}
				</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
					autocomplete={isSetup ? 'new-password' : 'current-password'}
				/>
			</div>

			{#if isSetup}
				<div>
					<label for="confirm" class="mb-1 block text-sm font-medium text-gray-700">
						Confirm Password
					</label>
					<input
						id="confirm"
						type="password"
						bind:value={confirmPassword}
						class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
						autocomplete="new-password"
					/>
				</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
			>
				{loading ? 'Please wait…' : isSetup ? 'Set Password & Continue' : 'Sign In'}
			</button>
		</form>
	</div>
</div>
