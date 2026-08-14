import type { LayoutLoad } from "./$types";
import { loadLocale } from "wuchale/load-utils";
import { getClient } from "$lib/clients/univeral-client.svelte";

export const load: LayoutLoad = async (event) => {
	loadLocale(event.data.locale);
	const client = getClient(event);

	const [{ user }, { units }, { loadouts }] = await Promise.all([
		client.policy.getUser({}),
		client.catalogue.listUnits({}),
		client.catalogue.listLoadouts({}),
	]);

	return {
		policyUser: user,
		units: Object.fromEntries(units.map((i) => [i.id, i])),
		loadouts: Object.fromEntries(loadouts.map((i) => [i.id, i])),
	};
};
