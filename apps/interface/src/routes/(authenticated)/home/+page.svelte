<script lang="ts">
	import { getFactionName } from "$lib/faction";
	import { goto } from "$app/navigation";
	import { colors } from "$lib/color";
	import SquadItem from "$lib/components/squad-item.svelte";
	import type { GameItem, GameLoadout } from "$lib/game";
	import type { PageData } from "./$types";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";

	let { data }: { data: PageData } = $props();

	let selectedSquad = $state<string>("");
	let selectedDefender = $state<Faction>("" as any);
</script>

<div class="wrapper">
	<section>
		<h2>Squads</h2>
		<button onclick={() => goto("/squads/new")}>New Squad</button>
		<ul>
			{#each Object.values(data.squads) as squad}
				<SquadItem {squad} onclick={() => goto(`/squads/${squad.id}`)} />
			{/each}
		</ul>
	</section>
	<h1>Start a new game</h1>
	<select bind:value={selectedSquad}>
		<option value="" disabled hidden selected>Select Squad</option>
		{#each Object.values(data.squads) as squad}
			<option value={squad!.id}>
				{getFactionName(squad!.faction)}
				{squad!.name}
			</option>
		{/each}
	</select>
	{#if selectedSquad}
		<select bind:value={selectedDefender}>
			<option value="" disabled hidden selected>Select enemy faction</option>
			{#if data.squads[selectedSquad]?.faction !== Faction.GALACTIC_REPUBLIC}
				<option value={Faction.GALACTIC_REPUBLIC}>Galactic Republic</option>
			{/if}
			{#if data.squads[selectedSquad]?.faction !== Faction.SEPARATIST_ALLIANCE}
				<option value={Faction.SEPARATIST_ALLIANCE}>Separatist Alliance</option>
			{/if}
		</select>
		{#if selectedDefender}
			<button
				onclick={() => {
					localStorage.setItem(
						"GAME",
						JSON.stringify(
							{
								attacker: {
									name: data.squads[selectedSquad]!.name,
									faction: data.squads[selectedSquad]!.faction,
									loadouts: Object.fromEntries(
										data.squads[selectedSquad]!.loadouts.map((l, i) => {
											const id = crypto.randomUUID();
											const loadout: GameLoadout = {
												id: id,
												image: l.image,
												name: l.name,
												color: colors[i].hex,
												turnComplete: false,
												unit: l.unit!,
												items: Object.values(
													l.items.reduce(
														(items, item) => {
															if (item.details.value) {
																if ("capacity" in item.details.value) {
																	if (item.id in items) {
																		items[item.id].quantity +=
																			item.details.value.capacity;
																	} else {
																		items[item.id] = {
																			quantity: item.details.value.capacity,
																			item: item,
																		};
																	}
																} else {
																	if (item.id in items) {
																		items[item.id].quantity++;
																	} else {
																		items[item.id] = {
																			quantity: 1,
																			item: item,
																		};
																	}
																}
															}
															return items;
														},
														{} as Record<string, GameItem>,
													),
												),
											};
											if (loadout.unit!.size === 1) {
												loadout.inCover = false;
											}
											return [id, loadout];
										}),
									),
								},
								defender: {
									faction: selectedDefender,
									units: Object.values(data.units).filter((i) =>
										i.factions.includes(selectedDefender),
									),
								},
							},
							(_key, value) =>
								typeof value === "bigint" ? Number(value) : value,
						),
					);
					goto("/game");
				}}
			>
				Start Game
			</button>
		{/if}
	{/if}
</div>

<style lang="scss"></style>
