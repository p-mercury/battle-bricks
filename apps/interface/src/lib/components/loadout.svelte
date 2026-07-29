<script lang="ts">
	import type { Loadout } from "$lib/loadouts";
	import Item from "./item.svelte";
	import StatBar from "./stat-bar.svelte";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: { loadout: Loadout; selected?: Boolean; onclick?: () => void } = $props();
</script>

<li {onclick} class:selected>
	<h3>{loadout.name}</h3>
	<div class="content">
		<div class="stats">
			<div class="row">
				<div>Price:</div>
				<div>{loadout.price}</div>
			</div>
			<div class="row">
				<div>Health:</div>
				<div>{loadout.unit.health}</div>
			</div>
			<div class="row">
				<div>Carry Weight:</div>
				<div>{loadout.carryWeight}/{loadout.unit.carryCapacity}</div>
			</div>
			<div class="row">
				<div>Size:</div>
				<StatBar value={loadout.unit.size} size={3} />
			</div>
			<div class="row">
				<div>Armor Class:</div>
				<StatBar value={loadout.unit.armorClass} size={4} />
			</div>
			<div class="row">
				<div>Marksmanship:</div>
				<StatBar value={loadout.unit.marksmanship} size={4} />
			</div>
			<div class="row">
				<div>Melee Ability:</div>
				<StatBar value={loadout.unit.meleeAbility} size={4} />
			</div>
		</div>
		<div class="items-wrapper">
			<div class="items">
				{#each loadout.items as item}
					<Item {item} />
				{/each}
			</div>
		</div>
	</div>
</li>

<style lang="scss">
	li {
		box-sizing: border-box;
		width: 100%;
		height: 16rem;
		min-height: 16rem;
		display: flex;
		flex-direction: column;
		margin: 0;
		border: 2px solid black;
		border-radius: 0.4rem;
		overflow: hidden;
		cursor: pointer;
		overflow: hidden;

		&.selected {
			background-color: lightcoral;
		}
	}

	h3 {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
		font-size: 1.1rem;
		font-weight: 600;
		background-color: lightgray;
		width: 100%;
		padding: 0.2rem 0.5rem 0.2rem 0.5rem;
	}

	.content {
		display: flex;
		flex: 1;
		min-height: 0;
		gap: 1rem;
	}

	.stats {
		display: grid;
		grid-template-columns: auto auto;
		gap: 0.5rem;
		padding: 0.5rem;
		flex-shrink: 0;
	}

	.row {
		display: grid;
		list-style: none;
		grid-column: 1 / -1;
		grid-template-columns: subgrid;
	}

	.items-wrapper {
		box-sizing: border-box;
		height: 100%;
		overflow: scroll;
	}

	.items {
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
		padding: 0.5rem;
	}
</style>
