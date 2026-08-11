import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async (event) => ({
	locale: event.locals.locale,
});
