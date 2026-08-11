import { defineConfig } from "wuchale";
import { adapter as svelteAdapter } from "@wuchale/svelte";
import { adapter as jsAdapter } from "wuchale/adapter-vanilla";
import { generateText } from "ai";
import { createOpenAI } from "@ai-sdk/openai";

const mammouth = createOpenAI({
	baseURL: "https://api.mammouth.ai/v1",
	apiKey: process.env.MAMMOUTH_API_KEY,
});

export default defineConfig({
	locales: ["en", "de", "bg"],
	adapters: {
		svelte: svelteAdapter({ loader: "sveltekit" }),
		js: jsAdapter({
			loader: "vite",
			files: [
				"src/**/+{page,layout}.{js,ts}",
				"src/**/+{page,layout}.server.{js,ts}",
			],
		}),
		libSvelte: svelteAdapter({
			loader: "sveltekit",
			files: "../../packages/ui-web/dist/**/*.svelte",
		}),
		libJs: jsAdapter({
			loader: "vite",
			files: "../../packages/ui-web/dist/**/*.js",
		}),
	},
	ai: {
		name: "SONNET-5",
		batchSize: 50,
		parallel: 1,
		group: {
			en: [["en"], ["de"]],
		},
		translate: async (messages, instruction) => {
			const { text } = await generateText({
				model: mammouth("gpt-5.5"),
				instructions: instruction,
				prompt: messages,
			});
			return text
				.replace(/^```[\w]*\n?/, "")
				.replace(/\n?```$/, "")
				.trim();
		},
	},
});
