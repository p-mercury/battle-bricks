import type { LayoutLoad } from "./$types";
import { loadLocale } from "wuchale/load-utils";

import "./../locales/js.loader";
import "./../locales/libJs.loader";
import "./../locales/libSvelte.loader.svelte";
import "./../locales/svelte.loader.svelte";

export const load: LayoutLoad = () => loadLocale("en");
