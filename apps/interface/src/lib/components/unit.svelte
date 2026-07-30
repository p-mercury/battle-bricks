<script lang="ts">
	import type { Unit } from "$lib/units";
	import StatBar from "./stat-bar.svelte";
	import StatRow from "./stat-row.svelte";
	import StatTable from "./stat-table.svelte";

	let {
		unit,
		selected = false,
		onclick = () => {},
	}: { unit: Unit; selected?: Boolean; onclick?: () => void } = $props();
</script>

<li {onclick} class:selected>
	<h3>{unit.name}</h3>
	<StatTable>
		<StatRow>
			Price:
			<div>{unit.price}c</div>
		</StatRow>
		<StatRow>
			Health:
			<div>{unit.health}hp</div>
		</StatRow>
		<StatRow>
			Size:
			<StatBar value={unit.size} size={4} />
		</StatRow>
		<StatRow>
			Armor Class:
			<StatBar value={unit.armorClass} size={4} />
		</StatRow>
		{#if unit.marksmanship}
			<StatRow>
				Marksmanship:
				<StatBar value={unit.marksmanship} size={4} />
			</StatRow>
		{/if}
		{#if unit.meleeAbility}
			<StatRow>
				Marksmanship:
				<StatBar value={unit.meleeAbility} size={4} />
			</StatRow>
		{/if}
	</StatTable>
</li>

<style lang="scss">
	li {
		box-sizing: border-box;
		width: 100%;
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
</style>
