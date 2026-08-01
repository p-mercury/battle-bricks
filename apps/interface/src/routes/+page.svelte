<script lang="ts">
	import { getFactionName, type Faction } from "$lib/faction";
	import type { Squad } from "$lib/squad";
	import { browser } from "$app/env";
	import { loadouts } from "$lib/loadouts";
	import { units } from "$lib/units";
	import { goto } from "$app/navigation";

	let selectedSquad = $state<string>("");
	let selectedDefender = $state<Faction>("" as any);

	let squads = $derived.by<Record<string, Squad | undefined>>(() => {
		if (browser) {
			let item = localStorage.getItem("SQUADS");
			if (!item) {
				item = "{}";
			}
			return JSON.parse(item);
		} else {
			return {};
		}
	});
</script>

<div class="wrapper">
	<h1>Start a new game</h1>
	<select bind:value={selectedSquad}>
		<option value="" disabled hidden selected>Select Squad</option>
		{#each Object.values(squads) as squad}
			<option value={squad!.id}>
				{getFactionName(squad!.faction)}
				{squad!.name}
			</option>
		{/each}
	</select>
	{#if selectedSquad}
		<select bind:value={selectedDefender}>
			<option value="" disabled hidden selected>Select enemy faction</option>
			{#if squads[selectedSquad]?.faction !== "GALACTIC_REPUBLIC"}
				<option value="GALACTIC_REPUBLIC">Galactic Republic</option>
			{/if}
			{#if squads[selectedSquad]?.faction !== "SEPARATIST_ALLIANCE"}
				<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			{/if}
		</select>
		{#if selectedDefender}
			<button
				onclick={() => {
					localStorage.setItem(
						"GAME",
						JSON.stringify({
							initiative: 1,
							attacker: {
								name: squads[selectedSquad]!.name,
								faction: squads[selectedSquad]!.faction,
								loadouts: squads[selectedSquad]!.loadouts.map((l, i) => ({
									...loadouts[l],
									initiative: i * 2 + 1,
								})),
							},
							defender: {
								faction: selectedDefender,
								units: Object.values(units).filter((i) =>
									i.faction.includes(selectedDefender),
								),
							},
						}),
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
