<script lang="ts">
	import { untrack } from "svelte";
	import LoadoutsZone from "./loadouts-zone.svelte";
	import type { Squad } from "$lib/squad";
	import { page } from "$app/state";
	import SquadZone from "./squad-zone.svelte";
	import { goto } from "$app/navigation";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import type { PageData } from "./$types";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import { getClient } from "$lib/clients/univeral-client.svelte";
	import { getLoadoutPrice } from "$lib/loadout";

	let { data }: { data: PageData } = $props();

	const flipDurationMs = 150;
	let faction = $state<Faction>(Faction.GALACTIC_REPUBLIC);
	let loadouts = $state<{ id: string; loadout: Loadout }[]>([]);
	let squad = $state<{ id: string; loadout: Loadout }[]>([]);
	let budget = $derived(
		squad.reduce(
			(a, item) =>
				a -
				(item.loadout.unit!.price +
					item.loadout.items.reduce((b, item) => b + item.price, 0)),
			1500,
		),
	);
	let name = $state("");

	$effect(() => {
		loadouts = Object.values(data.loadouts)
			.filter((i) => i.unit!.factions.includes(faction))
			.sort((a, b) => getLoadoutPrice(a) - getLoadoutPrice(b))
			.map((i) => ({
				id: crypto.randomUUID(),
				loadout: i,
			}));
		squad = untrack(() => squad).filter((i) =>
			i.loadout.unit!.factions.includes(faction),
		);
	});

	$effect(() => {
		if (data.squad) {
			faction = data.squad.faction;
			name = data.squad.name;
			squad = data.squad.loadouts.map((loadout) => ({
				id: crypto.randomUUID(),
				loadout,
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

				if (!squad.length) {
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
						loadouts: squad.map((l) => l.loadout.id),
					});
				} else {
					client.catalogue.createSquad({
						name: name,
						faction: faction,
						loadouts: squad.map((l) => l.loadout.id),
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
	<div class="loadouts">
		<LoadoutsZone bind:items={loadouts} {flipDurationMs} />
	</div>
	<div class="squad">
		<SquadZone
			bind:items={squad}
			{flipDurationMs}
			handleConsider={(event) => {
				const nextBudget = event.detail.items.reduce(
					(a, item) => a - getLoadoutPrice(item.loadout),
					1500,
				);
				if (nextBudget >= 0) {
					squad = event.detail.items;
				}
			}}
			handleFinalize={(event) => {
				const items = event.detail.items;

				const nextBudget = items.reduce(
					(remaining, item) => remaining - getLoadoutPrice(item.loadout),
					1500,
				);

				const previousBudget = squad.reduce(
					(remaining, item) => remaining - getLoadoutPrice(item.loadout),
					1500,
				);

				if (nextBudget >= 0 || nextBudget > previousBudget) {
					squad = [...items].sort(
						(a, b) => getLoadoutPrice(a.loadout) - getLoadoutPrice(b.loadout),
					);
				} else {
					squad = [...squad];
				}
			}}
		/>
	</div>
</div>

<style lang="scss">
	.wrapper {
		display: grid;
		grid-template:
			"header header" 4rem
			"loadouts squad" 1fr /
			auto 1fr;
		height: 100dvh;
		width: 100dvw;
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
</style>
