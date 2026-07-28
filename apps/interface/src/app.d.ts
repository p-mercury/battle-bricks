import type { Logger } from "@aws-lambda-powertools/logger";

declare global {
	const SENTRY_ORIGIN: string | null;
	const SENTRY_PROJECT_ID: string | null;
	const SENTRY_PROJECT_SLUG: string | null;
	const SENTRY_ORGANISATION_SLUG: string | null;

	namespace App {
		interface Locals {
			logger: Logger;
			locale: string;
		}
	}
}

declare module "svelte/elements" {
	export interface SVGAttributes {
		"inline-src"?: `${string}` | `./${string}` | `../${string}`;
	}
}

export {};
