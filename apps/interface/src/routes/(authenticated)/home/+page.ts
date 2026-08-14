import type { PageLoad } from "./$types";
import { getClient } from "$lib/clients/univeral-client.svelte";

export const load: PageLoad = async (event) => {
	const client = getClient(event);

	const [{ squads }] = await Promise.all([client.catalogue.listSquads({})]);

	return {
		squads: Object.fromEntries(squads.map((i) => [i.id, i])),
	};
};
