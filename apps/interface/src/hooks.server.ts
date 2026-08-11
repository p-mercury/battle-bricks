import { sequence } from "@sveltejs/kit/hooks";
import { building } from "$app/environment";

import { fetchConfig } from "./hooks/fetch-config";
import { wuchale } from "./hooks/wuchale";

export const handle = building
	? sequence(wuchale)
	: sequence(fetchConfig, wuchale);
