import adapterCdk from "@flit/sveltekit-adapter-cdk";
import adapterStatic from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);

export default {
	preprocess: [vitePreprocess()],
	kit: {
		outDir: "dist/.svelte-kit",
		adapter: !!process.env.PUBLIC_NATIVE
			? adapterStatic({ pages: "./dist/static", fallback: "index.html" })
			: adapterCdk({ out: "./dist/cdk" }),
		experimental: {
			tracing: {
				server: true,
			},
			instrumentation: {
				server: true,
			},
		},
	},
};
