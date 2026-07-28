<script lang="ts">
	import { getLoadoutStats } from "$lib/get-loadout-stats";
	import type { Loadout } from "$lib/loadouts";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: { loadout: Loadout; selected?: Boolean; onclick?: () => void } = $props();

	let loadoutStats = $derived(getLoadoutStats(loadout));
</script>

<section {onclick} class:selected>
	<b>{loadout.name}</b>
	<div>Price: {loadoutStats.price}</div>
	<div>Health: {loadoutStats.unit.health}</div>
	<div>
		Carry Weight: {loadoutStats.carryWeight}/{loadoutStats.unit.carryCapacity}
	</div>
</section>

<style lang="scss">
	section {
		box-sizing: border-box;
		width: 100%;
		height: 100%;
		padding: 0.5rem;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightcyan;

		&.selected {
			background-color: lightcoral;
		}
	}
</style>
