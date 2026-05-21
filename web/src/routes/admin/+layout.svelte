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
		{ href: '/admin/pours', label: 'Pour History' },
		{ href: '/admin/features', label: 'Features' }
	];
</script>

{#if isLogin}
	{@render children()}
{:else if token}
	<div class="flex h-screen text-gray-900">
		<aside class="sidebar flex w-52 flex-col">
			<div class="brand-header">
				<div class="beer-icon" aria-hidden="true">
					<svg width="22" height="22" viewBox="0 0 24 24" fill="none">
						<rect x="9" y="2" width="6" height="3" rx="1" fill="currentColor" />
						<rect x="8" y="5" width="8" height="13" rx="1" fill="currentColor" opacity="0.85" />
						<path
							d="M16 8h3a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1h-3"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
						<rect x="9" y="18" width="6" height="2" rx="0.5" fill="currentColor" opacity="0.45" />
						<line
							x1="10.5"
							y1="8"
							x2="10.5"
							y2="15"
							stroke="rgba(255,255,255,0.15)"
							stroke-width="1"
						/>
					</svg>
				</div>
				<div>
					<div class="brand-name">GoPints</div>
					<div class="brand-sub">Admin</div>
				</div>
			</div>

			<nav class="flex flex-1 flex-col p-3">
				<div class="flex flex-col gap-0.5">
					{#each navLinks as link}
						<a
							href={link.href}
							class="nav-link {page.url.pathname.startsWith(link.href) ? 'active' : ''}"
						>
							{link.label}
						</a>
					{/each}
				</div>

				<div class="mt-auto pt-6 pb-1">
					<a href="/" class="view-site-link">
						<svg
							width="12"
							height="12"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
						>
							<path d="M19 12H5M12 5l-7 7 7 7" />
						</svg>
						View Site
					</a>
				</div>
			</nav>

			<div class="sidebar-footer">
				<button onclick={logout} class="logout-btn">Log out</button>
			</div>
		</aside>

		<main class="flex-1 overflow-auto bg-stone-50 p-8">
			{@render children()}
		</main>
	</div>
{/if}

<style>
	@import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&display=swap');

	.sidebar {
		background-color: #16100a;
		background-image:
			radial-gradient(ellipse 80% 40% at 20% 0%, rgba(200, 130, 20, 0.09) 0%, transparent 100%),
			radial-gradient(ellipse 60% 40% at 80% 100%, rgba(160, 80, 10, 0.06) 0%, transparent 100%);
		border-right: 1px solid rgba(200, 140, 30, 0.12);
		color: #e8d5b0;
	}

	.brand-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1.1rem 1.25rem 1rem;
		border-bottom: 1px solid rgba(200, 140, 30, 0.15);
	}

	.beer-icon {
		color: #c8821a;
		flex-shrink: 0;
	}

	.brand-name {
		font-family: 'Bebas Neue', Georgia, serif;
		font-size: 1.35rem;
		letter-spacing: 0.1em;
		color: #c8821a;
		line-height: 1.1;
	}

	.brand-sub {
		font-size: 0.6rem;
		letter-spacing: 0.25em;
		text-transform: uppercase;
		color: rgba(200, 130, 26, 0.4);
		margin-top: 1px;
	}

	.nav-link {
		display: flex;
		align-items: center;
		padding: 0.5rem 0.75rem;
		border-radius: 5px;
		font-size: 0.85rem;
		font-weight: 500;
		color: rgba(232, 213, 176, 0.55);
		text-decoration: none;
		transition:
			color 0.15s,
			background 0.15s;
		border-left: 2px solid transparent;
	}

	.nav-link:hover {
		color: rgba(232, 213, 176, 0.88);
		background: rgba(200, 130, 26, 0.08);
	}

	.nav-link.active {
		color: #c8821a;
		background: rgba(200, 130, 26, 0.13);
		border-left-color: #c8821a;
	}

	.view-site-link {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.45rem 0.75rem;
		border-radius: 5px;
		font-size: 0.78rem;
		font-weight: 500;
		color: rgba(200, 130, 26, 0.5);
		text-decoration: none;
		border: 1px solid rgba(200, 130, 26, 0.18);
		transition:
			color 0.15s,
			background 0.15s,
			border-color 0.15s;
	}

	.view-site-link:hover {
		color: #c8821a;
		background: rgba(200, 130, 26, 0.08);
		border-color: rgba(200, 130, 26, 0.4);
	}

	.sidebar-footer {
		border-top: 1px solid rgba(200, 140, 30, 0.12);
		padding: 0.625rem;
	}

	.logout-btn {
		width: 100%;
		padding: 0.45rem 0.75rem;
		border-radius: 5px;
		font-size: 0.78rem;
		text-align: left;
		color: rgba(232, 213, 176, 0.3);
		background: none;
		border: none;
		cursor: pointer;
		transition:
			color 0.15s,
			background 0.15s;
	}

	.logout-btn:hover {
		color: rgba(232, 213, 176, 0.7);
		background: rgba(255, 255, 255, 0.04);
	}
</style>
