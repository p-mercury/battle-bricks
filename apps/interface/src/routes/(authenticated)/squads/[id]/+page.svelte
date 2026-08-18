<script lang="ts">
	import { untrack } from "svelte";
	import UnitsZone from "./units-zone.svelte";
	import type { Squad } from "$lib/squad";
	import { page } from "$app/state";
	import SquadZone from "./squad-zone.svelte";
	import { goto } from "$app/navigation";
	import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
	import type { PageData } from "./$types";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import { getClient } from "$lib/clients/univeral-client.svelte";
	import { getUnitPrice } from "$lib/unit";

	let { data }: { data: PageData } = $props();

	const flipDurationMs = 150;
	let faction = $state<Faction>(Faction.GALACTIC_REPUBLIC);
	let units = $state<{ id: string; unit: Unit }[]>([]);
	let squad = $state<{ id: string; unit: Unit }[]>([]);
	let budget = $derived(
		squad.reduce(
			(a, item) =>
				a -
				(item.unit!.price +
					item.unit.items.reduce((b, item) => b + item.price, 0)),
			1500,
		),
	);
	let name = $state("");

	$effect(() => {
		units = Object.values(data.units)
			.filter((i) => i.factions.includes(faction))
			.sort((a, b) => getUnitPrice(a) - getUnitPrice(b))
			.map((i) => ({
				id: crypto.randomUUID(),
				unit: i,
			}));
		squad = untrack(() => squad).filter((i) =>
			i.unit!.factions.includes(faction),
		);
	});

	$effect(() => {
		if (data.squad) {
			faction = data.squad.faction;
			name = data.squad.name;
			squad = data.squad.units.map((unit) => ({
				id: crypto.randomUUID(),
				unit,
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
						updateMask: ["name", "faction", "units"],
						name: name,
						faction: faction,
						units: squad.map((l) => l.unit.id),
					});
				} else {
					client.catalogue.createSquad({
						name: name,
						faction: faction,
						units: squad.map((l) => l.unit.id),
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
		<UnitsZone bind:items={units} {flipDurationMs} />
	</div>
	<div class="squad">
		<SquadZone
			bind:items={squad}
			{flipDurationMs}
			handleConsider={(event) => {
				const nextBudget = event.detail.items.reduce(
					(a, item) => a - getUnitPrice(item.unit),
					1500,
				);
				if (nextBudget >= 0) {
					squad = event.detail.items;
				}
			}}
			handleFinalize={(event) => {
				const items = event.detail.items;

				const nextBudget = items.reduce(
					(remaining, item) => remaining - getUnitPrice(item.unit),
					1500,
				);

				const previousBudget = squad.reduce(
					(remaining, item) => remaining - getUnitPrice(item.unit),
					1500,
				);

				if (nextBudget >= 0 || nextBudget > previousBudget) {
					squad = [...items].sort(
						(a, b) => getUnitPrice(a.unit) - getUnitPrice(b.unit),
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
