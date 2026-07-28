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
			en: [["en"], ["de"], ["bg"]],
		},
		translate: async (messages, instruction) => {
			const { text } = await generateText({
				model: mammouth("gpt-5.5"),
				instructions: PROJECT_CONTEXT + instruction,
				prompt: messages,
			});
			return text
				.replace(/^```[\w]*\n?/, "")
				.replace(/\n?```$/, "")
				.trim();
		},
	},
});

const PROJECT_CONTEXT =
	"This is a marketplace enabling organisations to sell and purchase spare parts for industrial machinery. The tone should be professional but approachable. The platfrom is called jumper.de so do not translate jumper or jumper.de unless its talking about something else. If technical terms such as MT-9510, Conversion Kit, Lead Backer or Device Carrier appear, do not translate them, just copy them over. Respond with raw text content only, do not wrap your response in markdown code blocks. ";
