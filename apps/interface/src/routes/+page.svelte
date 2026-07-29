<script lang="ts">
	import Attack from "$lib/components/attack.svelte";
	import Loadout from "$lib/components/loadout.svelte";
	import { getLoadoutStats } from "$lib/get-loadout-stats";
	import { loadouts } from "$lib/loadouts";
	import type { Faction } from "$lib/units";

	let attackerFaction = $state<Faction>("GALACTIC_REPUBLIC");
	let defenderFaction = $state<Faction>("SEPARATIST_ALLIANCE");

	let selectedAttacker = $state<string>();
	let selectedDefender = $state<string>();
</script>

<svelte:head>
	<title>Battle Bricks</title>
</svelte:head>

<div class="wrapper">
	<div>
		<h2>Attacker</h2>
		<select bind:value={attackerFaction}>
			<option value="GALACTIC_REPUBLIC">Galactic Republic</option>
			<option value="REBEL_ALLIANCE">Rebel Alliance</option>
			<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			<option value="GALACTIC_EMPIRE">Galactic Empire</option>
		</select>
		<ul>
			{#each Object.values(loadouts).filter( (l) => getLoadoutStats(l).unit.faction.includes(attackerFaction), ) as loadout}
				<Loadout
					{loadout}
					selected={selectedAttacker === loadout.id}
					onclick={() => {
						if (selectedAttacker !== loadout.id) {
							selectedAttacker = loadout.id;
						} else {
							selectedAttacker = undefined;
						}
					}}
				/>
			{/each}
		</ul>
	</div>
	<div>
		<h2>Defender</h2>
		<select bind:value={defenderFaction}>
			<option value="GALACTIC_REPUBLIC">Galactic Eepublic</option>
			<option value="REBEL_ALLIANCE">Rebel Alliance</option>
			<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			<option value="GALACTIC_EMPIRE">Galactic Empire</option>
		</select>
		<ul>
			{#each Object.values(loadouts).filter( (l) => getLoadoutStats(l).unit.faction.includes(defenderFaction), ) as loadout}
				<Loadout
					{loadout}
					selected={selectedDefender === loadout.id}
					onclick={() => {
						if (selectedDefender !== loadout.id) {
							selectedDefender = loadout.id;
						} else {
							selectedDefender = undefined;
						}
					}}
				/>
			{/each}
		</ul>
	</div>
	{#if selectedAttacker && selectedDefender}
		<Attack
			attacker={loadouts[selectedAttacker]}
			defender={loadouts[selectedDefender]}
		/>
	{/if}
</div>

<style lang="scss">
	.wrapper {
		display: grid;
		grid-template-columns: auto auto auto 1fr;
		gap: 4rem;
		padding: 1rem;
	}

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		height: 100%;
		width: max-content;
	}
</style>
