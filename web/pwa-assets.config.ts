import { defineConfig, minimal2023Preset } from '@vite-pwa/assets-generator/config';

export default defineConfig({
	headLinkOptions: { preset: 'resvg' },
	preset: {
		...minimal2023Preset,
		maskable: {
			sizes: [512],
			padding: 0.1,
			resizeOptions: { background: '#16100a' }
		},
		apple: {
			sizes: [180],
			padding: 0.1,
			resizeOptions: { background: '#16100a' }
		},
		transparent: {
			sizes: [64, 192, 512],
			padding: 0,
			favicons: [[64, 'favicon.ico']]
		}
	},
	images: ['static/icon.svg']
});
