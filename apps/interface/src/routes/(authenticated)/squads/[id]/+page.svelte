<script lang="ts">
	import { untrack } from "svelte";
	import UnitsZone from "./units-zone.svelte";
	import type { Squad } from "$lib/squad";
	import { page } from "$app/state";
	import LoadoutsZone from "./loadouts-zone.svelte";
	import { goto } from "$app/navigation";
	import type { PageData } from "./$types";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import { getClient } from "$lib/clients/univeral-client.svelte";
	import { getUnitPrice } from "$lib/unit";
	import { getLoadoutPrice } from "$lib/loadout";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";

	let { data }: { data: PageData } = $props();

	const flipDurationMs = 150;
	let faction = $state<Faction>(Faction.GALACTIC_REPUBLIC);
	let units = $state<{ id: string; loadout: Loadout }[]>([]);
	let loadouts = $state<{ id: string; loadout: Loadout }[]>([]);
	let budget = $derived(
		loadouts.reduce((a, i) => a - getLoadoutPrice(i.loadout), 1500),
	);
	let name = $state("");

	$effect(() => {
		units = Object.values(data.units)
			.filter((i) => i.factions.includes(faction))
			.sort((a, b) => getUnitPrice(a) - getUnitPrice(b))
			.map((i) => ({
				id: crypto.randomUUID(),
				loadout: {
					unit: i,
					items: [],
				} as any,
			}));
		loadouts = untrack(() => loadouts).filter((i) =>
			i.loadout!.unit!.factions.includes(faction),
		);
	});

	$effect(() => {
		if (data.squad) {
			faction = data.squad.faction;
			name = data.squad.name;
			loadouts = data.squad.loadouts.map((loadout) => ({
				id: crypto.randomUUID(),
				loadout,
				items: [],
			}));
		}
	});
</script>

<div class="wrapper">
	<header>
		<select bind:value={faction}>
			<option value={Faction.GALACTIC_REPUBLIC}>Galactic Republic</option>
			<option value={Faction.REBEL_ALLIANCE}>Rebel Alliance</option>
			<option value={Faction.SEPARATIST_ALLIANCE}>Separatist Alliance</option>
			<option value={Faction.GALACTIC_EMPIRE}>Galactic Empire</option>
		</select>
		<input type="text" placeholder="name" bind:value={name} />
		<ln>Budget: {budget}</ln>
		<button
			onclick={async () => {
				if (!name) {
					return;
				}

				if (!loadouts.length) {
					return;
				}

				let item = localStorage.getItem("SQUADS");
				if (!item) {
					item = "{}";
				}
				const squads = JSON.parse(item) as Record<string, Squad | undefined>;

				const client = getClient();
				if (page.params.id && page.params.id !== "new") {
					client.catalogue.updateSquad({
						id: page.params.id,
						updateMask: ["name", "faction", "loadouts"],
						name: name,
						faction: faction,
						loadouts: loadouts.map((l) => ({
							unit: l.loadout.unit!.id,
							item: l.loadout.item?.id,
						})),
					});
				} else {
					client.catalogue.createSquad({
						name: name,
						faction: faction,
						loadouts: loadouts.map((l) => ({
							unit: l.loadout.unit!.id,
							item: l.loadout.item?.id,
						})),
					});
				}

				localStorage.setItem("SQUADS", JSON.stringify(squads));
				await new Promise((r) => setTimeout(r, 1000));

				goto("/home");
			}}
		>
			Save
		</button>
		<button onclick={() => goto("/home")}>Cancel</button>
	</header>
	<section class="loadouts">
		<h2>Units</h2>
		<div>
			<UnitsZone bind:items={units} {flipDurationMs} />
		</div>
	</section>
	<section class="squad">
		<h2>Your Squad</h2>
		<div>
			<LoadoutsZone
				bind:items={loadouts}
				{flipDurationMs}
				handleConsider={(event) => {
					const nextBudget = event.detail.items.reduce(
						(a, item) => a - getLoadoutPrice(item.loadout),
						1500,
					);
					if (nextBudget >= 0) {
						loadouts = event.detail.items;
					}
				}}
				handleFinalize={(event) => {
					const items = event.detail.items;

					const nextBudget = items.reduce(
						(remaining, item) => remaining - getLoadoutPrice(item.loadout),
						1500,
					);

					const previousBudget = loadouts.reduce(
						(remaining, item) => remaining - getLoadoutPrice(item.loadout),
						1500,
					);

					if (nextBudget >= 0 || nextBudget > previousBudget) {
						loadouts = [...items].sort(
							(a, b) => getLoadoutPrice(a.loadout) - getLoadoutPrice(b.loadout),
						);
					} else {
						loadouts = [...loadouts];
					}
				}}
			/>
		</div>
	</section>
</div>

<style lang="scss">
	.wrapper {
		display: grid;
		grid-template:
			"header header header header" 4rem
			". loadouts squad ." 1fr
			". . . ." 0 /
			0 auto auto 1fr;
		height: 100dvh;
		width: 100dvw;
		gap: 1rem;
	}

	header {
		all: unset;
		grid-area: header;
		background-color: lightpink;
	}

	.loadouts {
		grid-area: loadouts;
		overflow: hidden scroll;
	}

	.squad {
		grid-area: squad;
		overflow: hidden scroll;
	}

	section {
		grid-area: game;
		display: flex;
		flex-direction: column;
		margin: 0;
		padding: 0;
		border-radius: 0.8rem;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		overflow: hidden;

		h2 {
			margin: 0;
			padding: 0.5rem 1rem 0 1rem;
			font-size: 1.4rem;
			font-weight: 600;
		}

		div {
			max-height: 100%;
			max-width: 100%;
			overflow: hidden scroll;
			padding: 0.8rem;
		}
	}
</style>
