import * as Js from "./../locales/js.loader.server";
import * as LibJs from "./../locales/libJs.loader.server";
import * as LibSvelte from "./../locales/libSvelte.loader.server.svelte";
import * as Svelte from "./../locales/svelte.loader.server.svelte";
import { locales } from "./../locales/data";
import type { Handle } from "@sveltejs/kit";
import { runWithLocale, loadLocales } from "wuchale/load-utils/server";
import { getAcceptedLanguage } from "@flit/accepted-language";
import { getClient } from "$lib/clients/univeral-client.svelte";
import { Language } from "@battle-bricks/contracts/common/v1/language_pb";
import { localeToLanguage } from "$lib/language";

loadLocales(Js.key, Js.loadCount, Js.loadCatalog, locales);
loadLocales(LibJs.key, LibJs.loadCount, LibJs.loadCatalog, locales);
loadLocales(LibSvelte.key, LibSvelte.loadCount, LibSvelte.loadCatalog, locales);
loadLocales(Svelte.key, Svelte.loadCount, Svelte.loadCatalog, locales);

export const wuchale: Handle = async ({ event, resolve }) => {
	event.locals.locale = getAcceptedLanguage(
		locales,
		event.request.headers.get("accept-language"),
	);

	if (event.route.id) {
		const client = getClient(event);
		try {
			const user = (await client.identity.getUser({})).user!;
			switch (user.language) {
				case Language.EN:
					event.locals.locale = "en";
					break;
				case Language.DE:
					event.locals.locale = "de";
					break;
				case Language.BG:
					event.locals.locale = "bg";
					break;
				default:
					event.locals.locale = getAcceptedLanguage(
						locales,
						event.request.headers.get("accept-language"),
					);
					client.identity.updateUser({
						id: user.id,
						updateMask: ["language"],
						language: localeToLanguage(event.locals.locale),
					});
			}
		} catch {
			event.locals.locale = getAcceptedLanguage(
				locales,
				event.request.headers.get("accept-language"),
			);
		}

		return runWithLocale(event.locals.locale, () => resolve(event));
	} else {
		return resolve(event);
	}
};
