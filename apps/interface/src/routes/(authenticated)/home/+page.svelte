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
	<section class="squads">
		<h2>Your Squads</h2>
		<div>
			<ul>
				<li>
					<button onclick={() => goto("/squads/new")}>New Squad</button>
				</li>
				{#each Object.values(data.squads) as squad}
					<li>
						<SquadItem {squad} onclick={() => goto(`/squads/${squad.id}`)} />
					</li>
				{/each}
			</ul>
		</div>
	</section>
	<section class="game">
		<h2>Start game</h2>
		<div>
			<ul>
				<li>
					<select bind:value={selectedSquad}>
						<option value="" disabled hidden selected>Select Squad</option>
						{#each Object.values(data.squads) as squad}
							<option value={squad!.id}>
								{getFactionName(squad!.faction)}
								{squad!.name}
							</option>
						{/each}
					</select>
				</li>
				{#if selectedSquad}
					<li>
						<select bind:value={selectedDefender}>
							<option value="" disabled hidden selected
								>Select enemy faction</option
							>
							{#if data.squads[selectedSquad]?.faction !== Faction.GALACTIC_REPUBLIC}
								<option value={Faction.GALACTIC_REPUBLIC}
									>Galactic Republic</option
								>
							{/if}
							{#if data.squads[selectedSquad]?.faction !== Faction.SEPARATIST_ALLIANCE}
								<option value={Faction.SEPARATIST_ALLIANCE}
									>Separatist Alliance</option
								>
							{/if}
						</select>
					</li>
					{#if selectedDefender}
						<li>
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
																						items[item.id].maxQuantity +=
																							item.details.value.capacity;
																						items[item.id].quantity +=
																							item.details.value.capacity;
																					} else {
																						items[item.id] = {
																							maxQuantity:
																								item.details.value.capacity,
																							quantity:
																								item.details.value.capacity,
																							item: item,
																						};
																					}
																				} else {
																					if (item.id in items) {
																						items[item.id].maxQuantity++;
																						items[item.id].quantity++;
																					} else {
																						items[item.id] = {
																							maxQuantity: 1,
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
						</li>
					{/if}
				{/if}
			</ul>
		</div>
	</section>
</div>

<style lang="scss">
	.wrapper {
		box-sizing: border-box;
		display: grid;
		grid-template:
			". . . . ." 0.5rem
			". squads . game ." 1fr /
			0.5rem min-content 0.5rem min-content minmax(0.5rem, 1fr);
		height: 100dvh;
		width: 100dvw;
	}

	.squads {
		grid-area: squads;
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

	.game {
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

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		height: max-content;
		width: 100%;
	}

	li {
		all: unset;
	}
</style>
