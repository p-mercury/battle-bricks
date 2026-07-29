<script lang="ts">
	import Attack from "$lib/components/attack.svelte";
	import Loadout from "$lib/components/loadout.svelte";
	import Unit from "$lib/components/unit.svelte";
	import { loadouts } from "$lib/loadouts";
	import { units, type Faction } from "$lib/units";

	let attackerFaction = $state<Faction>("GALACTIC_REPUBLIC");
	let defenderFaction = $state<Faction>("SEPARATIST_ALLIANCE");

	let selectedAttacker = $state<string>();
	let selectedDefender = $state<string>();
</script>

<svelte:head>
	<title>Battle Bricks</title>
</svelte:head>

<div class="wrapper">
	<section>
		<h2>Attacker</h2>
		<select bind:value={attackerFaction}>
			<option value="GALACTIC_REPUBLIC">Galactic Republic</option>
			<option value="REBEL_ALLIANCE">Rebel Alliance</option>
			<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			<option value="GALACTIC_EMPIRE">Galactic Empire</option>
		</select>
		<ul>
			{#each Object.values(loadouts).filter( (i) => i.unit.faction.includes(attackerFaction), ) as loadout}
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
	</section>
	<section>
		<h2>Defender</h2>
		<select bind:value={defenderFaction}>
			<option value="GALACTIC_REPUBLIC">Galactic Eepublic</option>
			<option value="REBEL_ALLIANCE">Rebel Alliance</option>
			<option value="SEPARATIST_ALLIANCE">Separatist Alliance</option>
			<option value="GALACTIC_EMPIRE">Galactic Empire</option>
		</select>
		<ul>
			{#each Object.values(units).filter( (i) => i.faction.includes(defenderFaction), ) as unit}
				<Unit
					{unit}
					selected={selectedDefender === unit.id}
					onclick={() => {
						if (selectedDefender !== unit.id) {
							selectedDefender = unit.id;
						} else {
							selectedDefender = undefined;
						}
					}}
				/>
			{/each}
		</ul>
	</section>
	{#if selectedAttacker && selectedDefender}
		<Attack
			attacker={loadouts[selectedAttacker]}
			defender={units[selectedDefender]}
		/>
	{/if}
</div>

<style lang="scss">
	.wrapper {
		box-sizing: border-box;
		display: grid;
		grid-template-columns: auto auto auto 1fr;
		gap: 2rem;
		height: 100dvh;
	}

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		height: max-content;
		width: max-content;
	}

	section {
		max-height: 100%;
		overflow: scroll;
		padding: 1rem;
	}
</style>
