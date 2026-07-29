<script lang="ts">
	import { getLoadoutStats } from "$lib/get-loadout-stats";
	import type { Loadout } from "$lib/loadouts";
	import StatBar from "./stat-bar.svelte";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: { loadout: Loadout; selected?: Boolean; onclick?: () => void } = $props();

	let loadoutStats = $derived(getLoadoutStats(loadout));
</script>

<button {onclick} class:selected>
	<b>
		{loadout.name}
	</b>
	<ul>
		<li>
			<div>Price:</div>
			<div>{loadoutStats.price}</div>
		</li>
		<li>
			<div>Health:</div>
			<div>{loadoutStats.unit.health}</div>
		</li>
		<li>
			<div>Carry Weight:</div>
			<div>{loadoutStats.carryWeight}/{loadoutStats.unit.carryCapacity}</div>
		</li>
		<li>
			<div>Size:</div>
			<StatBar value={loadoutStats.unit.size} size={3} />
		</li>
		<li>
			<div>Armor Class:</div>
			<StatBar value={loadoutStats.unit.armorClass} size={4} />
		</li>
		<li>
			<div>Marksmanship:</div>
			<StatBar value={loadoutStats.unit.marksmanship} size={4} />
		</li>
		<li>
			<div>Melee Ability:</div>
			<StatBar value={loadoutStats.unit.meleeAbility} size={4} />
		</li>
	</ul>
</button>

<style lang="scss">
	button {
		all: unset;
		cursor: pointer;
		padding: 0.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightcyan;

		&.selected {
			background-color: lightcoral;
		}
	}

	ul {
		box-sizing: border-box;
		display: grid;
		grid-template-columns: auto auto;
		gap: 0.5rem;
		width: 100%;
		height: 100%;
		padding: 0;
		margin: 0;
	}

	li {
		display: grid;
		list-style: none;
		grid-column: 1 / -1;
		grid-template-columns: subgrid;
	}
</style>
