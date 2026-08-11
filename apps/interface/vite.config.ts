import { sveltekit } from "@sveltejs/kit/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";
import { defineConfig } from "vite";
import { wuchale } from "wuchale/vite";

export default defineConfig({
	clearScreen: false,
	plugins: [
		wuchale(),
		sveltekit(),
		basicSsl({
			certDir: "dist",
		}),
	],
	server: {
		strictPort: true,
		host: false,
		port: 5173,
		proxy: {
			"/api": {
				target: "https://battlebricks.games",
				changeOrigin: true,
			},
		},
	},
});
