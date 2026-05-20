<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getToken, clearToken, adminFetch } from '$lib/admin.js';

	let { children } = $props();

	let isLogin = $derived(page.url.pathname === '/admin/login');
	let token = $state(getToken());

	$effect(() => {
		if (!token && !isLogin) {
			goto('/admin/login', { replaceState: true });
		}
	});

	async function logout() {
		await adminFetch('/api/admin/logout', { method: 'POST' });
		clearToken();
		goto('/admin/login');
	}

	const navLinks = [
		{ href: '/admin/kegs', label: 'Kegs' },
		{ href: '/admin/taps', label: 'Taps' },
		{ href: '/admin/banner', label: 'Banner' },
		{ href: '/admin/pours', label: 'Pour History' }
	];
</script>

{#if isLogin}
	{@render children()}
{:else if token}
	<div class="flex h-screen bg-gray-100 text-gray-900">
		<aside class="flex w-52 flex-col bg-gray-900 text-gray-100">
			<div class="border-b border-gray-700 px-5 py-4 text-lg font-bold tracking-tight">
				GoPints Admin
			</div>
			<nav class="flex flex-1 flex-col gap-1 p-3">
				{#each navLinks as link}
					<a
						href={link.href}
						class="rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-gray-700
							{page.url.pathname.startsWith(link.href) ? 'bg-gray-700 text-white' : 'text-gray-300'}"
					>
						{link.label}
					</a>
				{/each}
			</nav>
			<div class="border-t border-gray-700 p-3">
				<button
					onclick={logout}
					class="w-full rounded-md px-3 py-2 text-left text-sm text-gray-400 transition-colors hover:bg-gray-700 hover:text-white"
				>
					Log out
				</button>
			</div>
		</aside>

		<main class="flex-1 overflow-auto p-8">
			{@render children()}
		</main>
	</div>
{/if}
