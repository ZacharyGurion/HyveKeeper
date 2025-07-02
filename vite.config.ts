import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		host: '0.0.0.0',
		port: 3000,
		allowedHosts: [
			//'localhost',
			//'10.0.0.100',
			'dev.zgurion.com',
		],
	}
});
