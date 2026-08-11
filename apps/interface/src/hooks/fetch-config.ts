import { type Handle } from "@sveltejs/kit";

export const fetchConfig: Handle = async ({ event, resolve }) => {
	return await resolve(event, {
		filterSerializedResponseHeaders: () => true,
	});
};
