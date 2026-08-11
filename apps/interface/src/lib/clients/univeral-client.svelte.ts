import { createConnectTransport } from "@connectrpc/connect-web";
import {
	ConnectError,
	type createClient as cCreateClient,
	type CallOptions,
	type Transport,
	makeAnyClient,
} from "@connectrpc/connect";
import { External as FileStagingService } from "@battle-bricks/contracts/filestaging/v1/service_pb";
import { External as IdentityService } from "@battle-bricks/contracts/identity/v1/service_pb";
import { External as PolicyService } from "@battle-bricks/contracts/policy/v1/service_pb";

import {
	isRedirect,
	redirect,
	type LoadEvent,
	type RequestEvent,
} from "@sveltejs/kit";
import wasmUrl from "@battle-bricks/policy-service/policy?url";
import { loadPolicy } from "@open-policy-agent/opa-wasm";
import {
	GetUserResponseSchema,
	type GetUserResponse,
} from "@battle-bricks/contracts/policy/v1/get-user_pb";
import type {
	DescMessage,
	DescMethodUnary,
	MessageInitShape,
	MessageShape,
} from "@bufbuild/protobuf";
import { browser, dev } from "$app/environment";
import { create, toJson } from "@bufbuild/protobuf";
import type { User } from "@battle-bricks/contracts/policy/v1/user_pb";
import { navigating } from "$app/state";

if (dev && !browser) {
	process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
}

const createConnectClient: typeof cCreateClient = (service, transport) => {
	return makeAnyClient(service, (method: any) =>
		createUnaryFn(transport, method),
	) as any;
};

type UnaryFn<I extends DescMessage, O extends DescMessage> = (
	request: MessageInitShape<I>,
	options?: CallOptions,
) => Promise<MessageShape<O>>;
function createUnaryFn<I extends DescMessage, O extends DescMessage>(
	transport: Transport,
	method: DescMethodUnary<I, O>,
): UnaryFn<I, O> {
	return async function (input, options) {
		try {
			var _a, _b;
			const response = await transport.unary(
				method,
				options === null || options === void 0 ? void 0 : options.signal,
				options === null || options === void 0 ? void 0 : options.timeoutMs,
				options === null || options === void 0 ? void 0 : options.headers,
				input,
				options === null || options === void 0 ? void 0 : options.contextValues,
			);
			(_a =
				options === null || options === void 0 ? void 0 : options.onHeader) ===
				null || _a === void 0
				? void 0
				: _a.call(options, response.header);
			(_b =
				options === null || options === void 0 ? void 0 : options.onTrailer) ===
				null || _b === void 0
				? void 0
				: _b.call(options, response.trailer);
			return response.message;
		} catch (err) {
			if (err instanceof ConnectError && isRedirect(err.cause)) {
				throw err.cause;
			}
			throw err;
		}
	};
}

const REFRESH_LOCK_SYMBOL = Symbol();
const OPA_CLIENT_SYMBOL = Symbol();
const POLICY_USER_SYMBOL = Symbol();

class Client {
	private readonly store;
	private readonly fetch;
	private url;

	public readonly fileStaging;
	public readonly identity;
	public readonly policy;

	constructor(event?: LoadEvent | RequestEvent) {
		if (event) {
			this.store = event.params as any;
			this.fetch = event.fetch;
			this.url = event.url;
		} else {
			this.store = {} as any;
			this.fetch = fetch;
		}

		const transport = createConnectTransport({
			baseUrl: `${(this.url || window.location).origin}/api`,
			useBinaryFormat: !dev,
			fetch: async (input: RequestInfo | URL, init?: RequestInit) => {
				await this.store[REFRESH_LOCK_SYMBOL];

				let response = await this.fetch(input, {
					...init,
					credentials: "include",
				});

				if (response.status === 401) {
					if (!this.store[REFRESH_LOCK_SYMBOL]) {
						this.store[REFRESH_LOCK_SYMBOL] = this.fetch(
							new URL(
								`/api/${IdentityService.typeName}/auth/refresh`,
								(this.url || window.location).origin,
							),
							{
								credentials: "include",
							},
						);
					}

					if (!(await this.store[REFRESH_LOCK_SYMBOL]).ok) {
						this.authRedirect();
					}

					response = await this.fetch(input, {
						...init,
						credentials: "include",
					});

					this.store[REFRESH_LOCK_SYMBOL] = undefined;
				}

				return response;
			},
		});

		const policyClient = createConnectClient(PolicyService, transport);

		this.fileStaging = createConnectClient(FileStagingService, transport);
		this.identity = createConnectClient(IdentityService, transport);
		this.policy = {
			getUser: async (_: object) => {
				if (!this.store[POLICY_USER_SYMBOL]) {
					this.store[POLICY_USER_SYMBOL] = policyClient.getUser({});
				}
				return {
					user: (await this.store[POLICY_USER_SYMBOL]).user.value as User,
				};
			},
			evaluate: async <Desc extends DescMessage>(
				service:
					| "catalogue"
					| "identity_customer"
					| "identity_staff"
					| "zoho_crm"
					| "order"
					| "sales"
					| "purchasing",
				method: string,
				contextSchema: Desc,
				context: MessageInitShape<Desc>,
			) => {
				if (!this.store[OPA_CLIENT_SYMBOL]) {
					this.store[OPA_CLIENT_SYMBOL] = (async () => {
						const wasm = await this.fetch(wasmUrl);
						return loadPolicy(await wasm.arrayBuffer());
					})();
				}
				const opaClient: ReturnType<typeof loadPolicy> =
					this.store[OPA_CLIENT_SYMBOL];

				if (!this.store[POLICY_USER_SYMBOL]) {
					this.store[POLICY_USER_SYMBOL] = policyClient.getUser({});
				}
				let policyUser: Promise<GetUserResponse> =
					this.store[POLICY_USER_SYMBOL];

				const [client, user] = await Promise.all([opaClient, policyUser]);

				try {
					return client.evaluate(
						{
							user: toJson(GetUserResponseSchema, user),
							context: toJson(contextSchema, create(contextSchema, context)),
						},
						`services/${service}/${method}/decision`,
					)[0].result as { authz: boolean; output_mask: any };
				} catch {
					return { authz: false, output_mask: {} };
				}
			},
		};
	}

	authRedirect(): never {
		if (browser && navigating.to?.url) {
			this.url = navigating.to?.url;
		}

		const callbackUri = new URL((this.url || window.location).origin);
		callbackUri.pathname = "auth/callback";

		const authUrl = new URL(
			`/api/${IdentityService.typeName}/auth/authorize`,
			(this.url || window.location).origin,
		);
		authUrl.searchParams.set("redirect_uri", callbackUri.href);
		authUrl.searchParams.set("response_type", "code");
		authUrl.searchParams.set("scope", "openid");
		authUrl.searchParams.set("state", (this.url || window.location).href);

		if (browser) {
			window.location.href = authUrl.href;
		}

		redirect(302, authUrl);
	}
}

const CLIENT_SYMBOL = Symbol();

export function getClient(event?: LoadEvent | RequestEvent) {
	if (browser) {
		if (!(window as any)[CLIENT_SYMBOL]) {
			(window as any)[CLIENT_SYMBOL] = new Client();
		}
		return (window as any)[CLIENT_SYMBOL] as Client;
	} else if (event) {
		if (!(event.params as any)[CLIENT_SYMBOL]) {
			(event.params as any)[CLIENT_SYMBOL] = new Client(event);
		}
		return (event.params as any)[CLIENT_SYMBOL] as Client;
	} else {
		throw new Error("Can't create client without event on server");
	}
}
