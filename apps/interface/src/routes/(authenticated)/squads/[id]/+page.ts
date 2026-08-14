import type { PageLoad } from "./$types";
import { getClient } from "$lib/clients/univeral-client.svelte";

export const load: PageLoad = async (event) => {
	if (event.params.id !== "new") {
		const client = getClient(event);

		const [{ squad }] = await Promise.all([
			client.catalogue.getSquad({ id: event.params.id }),
		]);

		return { squad };
	} else {
		return { squad: undefined };
	}
};
