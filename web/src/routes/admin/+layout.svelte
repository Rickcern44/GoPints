<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getToken, clearToken, adminFetch } from '$lib/admin.js';
	import { fetchHealth } from '$lib/api.js';

	let { children } = $props();

	let isLogin = $derived(page.url.pathname === '/admin/login');
	// Re-read on every navigation so a fresh login is reflected immediately.
	let token = $derived(page.url.pathname ? getToken() : null);

	let version = $state<string | null>(null);

	$effect(() => {
		if (!token && !isLogin) {
			goto('/admin/login', { replaceState: true });
		}
	});

	$effect(() => {
		if (token && version === null) {
			fetchHealth().then((h) => {
				if (h) version = h.version;
			});
		}
	});

	async function logout() {
		await adminFetch('/api/admin/logout', { method: 'POST' });
		clearToken();
		goto('/admin/login');
	}

	const navLinks = [
		{
			href: '/admin/kegs',
			label: 'Kegs',
			icon: `<path d="M6 2h12l1 4H5L6 2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" fill="none"/><rect x="4" y="6" width="16" height="14" rx="2" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M8 10h8M8 14h5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>`
		},
		{
			href: '/admin/taps',
			label: 'Taps',
			icon: `<path d="M12 3v4M9 7h6M10 7v10M14 7v10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><path d="M8 17h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><circle cx="12" cy="3" r="1" fill="currentColor"/>`
		},
		{
			href: '/admin/banner',
			label: 'Banner',
			icon: `<rect x="3" y="6" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M3 10h18" stroke="currentColor" stroke-width="1.5"/><path d="M7 14h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>`
		},
		{
			href: '/admin/pours',
			label: 'Pours',
			icon: `<path d="M5 3l14 0-2 10H7L5 3Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" fill="none"/><path d="M7 13v5M17 13v5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><path d="M5 18h14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>`
		},
		{
			href: '/admin/features',
			label: 'Features',
			icon: `<circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M12 2v3M12 19v3M2 12h3M19 12h3M4.93 4.93l2.12 2.12M16.95 16.95l2.12 2.12M4.93 19.07l2.12-2.12M16.95 7.05l2.12-2.12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>`
		}
	];
</script>

{#if isLogin}
	{@render children()}
{:else if token}
	<!-- Desktop: sidebar layout | Mobile: bottom nav layout -->
	<div class="admin-shell">
		<!-- Sidebar (desktop only) -->
		<aside class="sidebar">
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

			<nav class="sidebar-nav">
				<div class="nav-items">
					{#each navLinks as link}
						<a
							href={link.href}
							class="nav-link {page.url.pathname.startsWith(link.href) ? 'active' : ''}"
						>
							{link.label}
						</a>
					{/each}
				</div>

				<div class="sidebar-bottom">
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

		<!-- Main content -->
		<main class="main-content">
			{@render children()}
		</main>

		<!-- Mobile top bar -->
		<header class="mobile-topbar">
			<div class="mobile-brand">
				<div class="beer-icon" aria-hidden="true">
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none">
						<rect x="9" y="2" width="6" height="3" rx="1" fill="currentColor" />
						<rect x="8" y="5" width="8" height="13" rx="1" fill="currentColor" opacity="0.85" />
						<path
							d="M16 8h3a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1h-3"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
					</svg>
				</div>
				<span class="mobile-brand-name">GoPints</span>
				<span class="mobile-brand-sub">Admin</span>
			</div>
			<div class="mobile-topbar-actions">
				<a href="/" class="mobile-view-site" aria-label="View site">
					<svg
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<circle cx="12" cy="12" r="10" />
						<path d="M12 8l4 4-4 4M8 12h8" />
					</svg>
				</a>
				<button onclick={logout} class="mobile-logout" aria-label="Log out">
					<svg
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
						<polyline points="16,17 21,12 16,7" />
						<line x1="21" y1="12" x2="9" y2="12" />
					</svg>
				</button>
			</div>
		</header>

		<!-- Mobile bottom nav -->
		<nav class="mobile-bottomnav" aria-label="Admin navigation">
			{#each navLinks as link}
				<a
					href={link.href}
					class="bottomnav-item {page.url.pathname.startsWith(link.href) ? 'active' : ''}"
					aria-label={link.label}
				>
					<svg
						width="22"
						height="22"
						viewBox="0 0 24 24"
						fill="none"
						aria-hidden="true"
					>
						{@html link.icon}
					</svg>
					<span class="bottomnav-label">{link.label}</span>
				</a>
			{/each}
		</nav>

		{#if version}
			<div class="version-badge" aria-hidden="true">v{version}</div>
		{/if}
	</div>
{/if}

<style>
	/* ─── Shell ──────────────────────────────────────────────────── */
	.admin-shell {
		display: flex;
		height: 100dvh;
		color: var(--color-fg);
		position: relative;
	}

	/* ─── Desktop sidebar ────────────────────────────────────────── */
	.sidebar {
		display: none;
		flex-direction: column;
		width: 13rem;
		flex-shrink: 0;
		background-color: var(--color-panel);
		background-image:
			radial-gradient(
				ellipse 80% 40% at 20% 0%,
				color-mix(in srgb, var(--color-accent) 9%, transparent) 0%,
				transparent 100%
			),
			radial-gradient(
				ellipse 60% 40% at 80% 100%,
				color-mix(in srgb, var(--color-accent) 6%, transparent) 0%,
				transparent 100%
			);
		border-right: 1px solid color-mix(in srgb, var(--color-accent) 12%, transparent);
		color: var(--color-fg-muted);
	}

	.brand-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1.1rem 1.25rem 1rem;
		border-bottom: 1px solid color-mix(in srgb, var(--color-accent) 15%, transparent);
	}

	.beer-icon {
		color: var(--color-accent);
		flex-shrink: 0;
	}

	.brand-name {
		font-family: var(--font-brand);
		font-size: 1.35rem;
		letter-spacing: 0.1em;
		color: var(--color-accent);
		line-height: 1.1;
	}

	.brand-sub {
		font-family: var(--font-mono);
		font-size: 0.6rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: color-mix(in srgb, var(--color-accent) 40%, transparent);
		margin-top: 1px;
	}

	.sidebar-nav {
		display: flex;
		flex: 1;
		flex-direction: column;
		padding: 0.75rem;
	}

	.nav-items {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.nav-link {
		position: relative;
		display: flex;
		align-items: center;
		padding: 0.5rem 0.75rem;
		border-radius: 5px;
		font-size: 0.85rem;
		font-weight: 500;
		color: color-mix(in srgb, var(--color-fg-muted) 55%, transparent);
		text-decoration: none;
		transition:
			color 0.15s,
			background 0.15s;
		border-left: 2px solid transparent;
	}

	.nav-link:hover {
		color: color-mix(in srgb, var(--color-fg-muted) 88%, transparent);
		background: color-mix(in srgb, var(--color-accent) 8%, transparent);
	}

	.nav-link.active {
		color: var(--color-accent);
		background: color-mix(in srgb, var(--color-accent) 13%, transparent);
		border-left-color: var(--color-accent);
	}

	/* Active-tab notch — pokes a small diamond into the border seam, like a connector pin */
	.nav-link.active::after {
		content: '';
		position: absolute;
		right: -1px;
		top: 50%;
		width: 8px;
		height: 8px;
		background: var(--color-void);
		transform: translateY(-50%) rotate(45deg);
		box-shadow: -1px 1px 0 color-mix(in srgb, var(--color-accent) 30%, transparent);
	}

	.sidebar-bottom {
		margin-top: auto;
		padding-top: 1.5rem;
		padding-bottom: 0.25rem;
	}

	.view-site-link {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.45rem 0.75rem;
		border-radius: 5px;
		font-size: 0.78rem;
		font-weight: 500;
		color: color-mix(in srgb, var(--color-accent) 50%, transparent);
		text-decoration: none;
		border: 1px solid color-mix(in srgb, var(--color-accent) 18%, transparent);
		transition:
			color 0.15s,
			background 0.15s,
			border-color 0.15s;
	}

	.view-site-link:hover {
		color: var(--color-accent);
		background: color-mix(in srgb, var(--color-accent) 8%, transparent);
		border-color: color-mix(in srgb, var(--color-accent) 40%, transparent);
	}

	.sidebar-footer {
		border-top: 1px solid color-mix(in srgb, var(--color-accent) 12%, transparent);
		padding: 0.625rem;
	}

	.logout-btn {
		width: 100%;
		padding: 0.45rem 0.75rem;
		border-radius: 5px;
		font-size: 0.78rem;
		text-align: left;
		color: color-mix(in srgb, var(--color-fg-muted) 30%, transparent);
		background: none;
		border: none;
		cursor: pointer;
		transition:
			color 0.15s,
			background 0.15s;
	}

	.logout-btn:hover {
		color: color-mix(in srgb, var(--color-fg-muted) 70%, transparent);
		background: rgba(255, 255, 255, 0.04);
	}

	/* ─── Main content ───────────────────────────────────────────── */
	.main-content {
		flex: 1;
		overflow-y: auto;
		background-color: var(--color-void);
		padding: 2rem;
		/* mobile: top bar + bottom nav */
		padding-top: calc(3.5rem + 2rem);
		padding-bottom: calc(4.5rem + env(safe-area-inset-bottom));
	}

	/* ─── Mobile top bar ─────────────────────────────────────────── */
	.mobile-topbar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		height: 3.5rem;
		padding: 0 1rem;
		padding-top: env(safe-area-inset-top);
		background-color: var(--color-panel);
		background-image: radial-gradient(
			ellipse 80% 100% at 20% 0%,
			color-mix(in srgb, var(--color-accent) 9%, transparent) 0%,
			transparent 100%
		);
		border-bottom: 1px solid color-mix(in srgb, var(--color-accent) 15%, transparent);
		z-index: 50;
	}

	.mobile-brand {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: var(--color-accent);
	}

	.mobile-brand-name {
		font-family: var(--font-brand);
		font-size: 1.25rem;
		letter-spacing: 0.1em;
		color: var(--color-accent);
		line-height: 1;
	}

	.mobile-brand-sub {
		font-size: 0.6rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: color-mix(in srgb, var(--color-accent) 45%, transparent);
		align-self: flex-end;
		padding-bottom: 2px;
	}

	.mobile-topbar-actions {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.mobile-view-site,
	.mobile-logout {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 8px;
		color: color-mix(in srgb, var(--color-fg-muted) 45%, transparent);
		background: none;
		border: none;
		cursor: pointer;
		text-decoration: none;
		transition: color 0.15s, background 0.15s;
	}

	.mobile-view-site:hover,
	.mobile-logout:hover {
		color: color-mix(in srgb, var(--color-fg-muted) 85%, transparent);
		background: rgba(255, 255, 255, 0.06);
	}

	/* ─── Mobile bottom nav ──────────────────────────────────────── */
	.mobile-bottomnav {
		display: flex;
		align-items: stretch;
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		height: calc(4rem + env(safe-area-inset-bottom));
		padding-bottom: env(safe-area-inset-bottom);
		background-color: var(--color-panel);
		border-top: 1px solid color-mix(in srgb, var(--color-accent) 15%, transparent);
		z-index: 50;
	}

	.bottomnav-item {
		display: flex;
		flex: 1;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 3px;
		text-decoration: none;
		color: color-mix(in srgb, var(--color-fg-muted) 40%, transparent);
		font-size: 0.65rem;
		font-weight: 500;
		letter-spacing: 0.03em;
		transition: color 0.15s;
		-webkit-tap-highlight-color: transparent;
	}

	.bottomnav-item:active {
		background: color-mix(in srgb, var(--color-accent) 8%, transparent);
	}

	.bottomnav-item.active {
		color: var(--color-accent);
	}

	.bottomnav-label {
		line-height: 1;
	}

	/* ─── Version badge ──────────────────────────────────────────── */
	.version-badge {
		display: none;
		position: fixed;
		bottom: 0.75rem;
		right: 0.75rem;
		font-size: 0.7rem;
		color: color-mix(in srgb, var(--color-fg-muted) 40%, transparent);
		pointer-events: none;
		z-index: 40;
	}

	/* ─── Responsive breakpoint ──────────────────────────────────── */
	@media (min-width: 768px) {
		.sidebar {
			display: flex;
		}

		.mobile-topbar,
		.mobile-bottomnav {
			display: none;
		}

		.main-content {
			padding: 2rem;
		}

		.version-badge {
			display: block;
		}
	}
</style>
