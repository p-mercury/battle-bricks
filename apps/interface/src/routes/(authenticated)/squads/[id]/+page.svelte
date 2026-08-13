<script lang="ts">
	import { loadouts as baseLoadouts, type Loadout } from "$lib/loadouts";
	import type { Faction } from "$lib/faction";
	import { untrack } from "svelte";
	import LoadoutsZone from "./loadouts-zone.svelte";
	import type { Squad } from "$lib/squad";
	import { page } from "$app/state";
	import SquadZone from "./squad-zone.svelte";
	import { goto } from "$app/navigation";

	const flipDurationMs = 150;
	let faction = $state<Faction>("GALACTIC_REPUBLIC");
	let loadouts = $state<{ id: string; loadout: Loadout }[]>([]);
	let squad = $state<{ id: string; loadout: Loadout }[]>([]);
	let budget = $derived(
		squad.reduce((sum, item) => sum - item.loadout.price, 1500),
	);
	let name = $state("");

	$effect(() => {
		loadouts = Object.values(baseLoadouts)
			.filter((i) => i.unit.faction.includes(faction))
			.map((i) => ({
				id: crypto.randomUUID(),
				loadout: i,
			}));
		squad = untrack(() => squad).filter((i) =>
			i.loadout.unit.faction.includes(faction),
		);
	});

	$effect(() => {
		if (page.params.id && page.params.id !== "new") {
			let item = localStorage.getItem("SQUADS");
			if (!item) {
				item = "{}";
			}
			const squads = JSON.parse(item) as Record<string, Squad | undefined>;
			const s = squads[page.params.id];
			if (s) {
				faction = s.faction;
				name = s.name;
				squad = s.loadouts.map((l) => ({
					id: crypto.randomUUID(),
					loadout: baseLoadouts[l],
				}));
			}
		}
	});
</script>

<div class="wrapper">
	<header>
		<select bind:value={faction}>
			<option value="GALACTIC_REPUBLIC">Galactic Republic</option>
			<option value="REBEL_ALLIANCE">Rebel Alliance</option>
			<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			<option value="GALACTIC_EMPIRE">Galactic Empire</option>
		</select>
		<input type="text" placeholder="name" bind:value={name} />
		<ln>Budget: {budget}</ln>
		<button
			onclick={() => {
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

				if (page.params.id && page.params.id !== "new") {
					squads[page.params.id] = {
						id: page.params.id,
						name: name,
						faction: faction,
						loadouts: squad.map((l) => l.loadout.id),
					};
				} else {
					const id = crypto.randomUUID();
					squads[id] = {
						id: id,
						name: name,
						faction: faction,
						loadouts: squad.map((l) => l.loadout.id),
					};
				}

				localStorage.setItem("SQUADS", JSON.stringify(squads));

				goto("/squads");
			}}
		>
			Save
		</button>
		<button onclick={() => goto("/squads")}> Cancle </button>
	</header>
	<div class="loadouts">
		<LoadoutsZone bind:items={loadouts} {flipDurationMs} />
	</div>
	<div class="squad">
		<SquadZone
			bind:items={squad}
			{flipDurationMs}
			handleConsider={(event) => {
				if (
					event.detail.items.reduce(
						(sum, item) => sum - item.loadout.price,
						1500,
					) >= 0
				) {
					squad = event.detail.items;
				}
			}}
			handleFinalize={(event) => {
				const items = event.detail.items;
				const next = items.reduce((sum, item) => sum + item.loadout.price, 0);
				const prev = squad.reduce((sum, item) => sum + item.loadout.price, 0);

				if (next <= 1500 || next < prev) {
					squad = items;
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
			1fr 1fr;
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
