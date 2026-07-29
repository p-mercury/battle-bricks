<script lang="ts">
	import type { Loadout } from "$lib/loadouts";
	import Item from "./item.svelte";
	import StatBar from "./stat-bar.svelte";
	import StatRow from "./stat-row.svelte";
	import StatTable from "./stat-table.svelte";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: { loadout: Loadout; selected?: Boolean; onclick?: () => void } = $props();
</script>

<li {onclick} class:selected>
	<h3>{loadout.name}</h3>
	<div class="content">
		<StatTable>
			<StatRow>
				Price:
				<div>{loadout.price}c</div>
			</StatRow>
			<StatRow>
				Health:
				<div>{loadout.unit.health}hp</div>
			</StatRow>
			<StatRow>
				Carry Weight:
				<div>{loadout.carryWeight}kg / {loadout.unit.carryCapacity}kg</div>
			</StatRow>
			<StatRow>
				Size:
				<StatBar value={loadout.unit.size} size={3} />
			</StatRow>
			<StatRow>
				Armor Class:
				<StatBar value={loadout.unit.armorClass} size={4} />
			</StatRow>
			<StatRow>
				Marksmanship:
				<StatBar value={loadout.unit.marksmanship} size={4} />
			</StatRow>
			<StatRow>
				Melee Ability:
				<StatBar value={loadout.unit.meleeAbility} size={4} />
			</StatRow>
		</StatTable>
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
		border: 4px solid lightgray;
		border-radius: 0.5rem;
		overflow: hidden;
		cursor: pointer;
		overflow: hidden;

		&.selected {
			border-color: lightcoral;
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
