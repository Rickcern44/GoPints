import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		VitePWA({
			registerType: 'autoUpdate',
			strategies: 'generateSW',
			injectRegister: 'auto',
			manifest: {
				name: 'GoPints',
				short_name: 'GoPints',
				description: 'Kegerator monitoring and management',
				theme_color: '#16100a',
				background_color: '#16100a',
				display: 'standalone',
				orientation: 'portrait-primary',
				scope: '/',
				start_url: '/',
				icons: [
					{ src: 'pwa-64x64.png', sizes: '64x64', type: 'image/png' },
					{ src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png' },
					{ src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png' },
					{
						src: 'maskable-icon-512x512.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'maskable'
					}
				],
				shortcuts: [
					{
						name: 'Manage Kegs',
						short_name: 'Kegs',
						url: '/admin/kegs',
						icons: [{ src: 'pwa-192x192.png', sizes: '192x192' }]
					},
					{
						name: 'Manage Taps',
						short_name: 'Taps',
						url: '/admin/taps',
						icons: [{ src: 'pwa-192x192.png', sizes: '192x192' }]
					}
				]
			},
			workbox: {
				// Network-first for API and WebSocket — real-time data must be fresh
				runtimeCaching: [
					{
						urlPattern: /^\/api\//,
						handler: 'NetworkFirst',
						options: {
							cacheName: 'api-cache',
							networkTimeoutSeconds: 5,
							cacheableResponse: { statuses: [0, 200] }
						}
					}
				],
				globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}']
			}
		})
	],
	server: {
		proxy: {
			'/api': { target: 'http://localhost:8080', ws: true }
		}
	}
});
