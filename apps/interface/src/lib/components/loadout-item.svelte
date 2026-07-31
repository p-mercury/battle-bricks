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
				Size:
				<StatBar value={loadout.unit.size} size={4} />
			</StatRow>
			<StatRow>
				Armor Class:
				<StatBar value={loadout.unit.armorClass} size={4} />
			</StatRow>
			{#if loadout.unit.marksmanship}
				<StatRow>
					Marksmanship:
					<StatBar value={loadout.unit.marksmanship} size={4} red />
				</StatRow>
			{/if}
			{#if loadout.unit.meleeAbility}
				<StatRow>
					Melee Ability:
					<StatBar value={loadout.unit.meleeAbility} size={4} red />
				</StatRow>
			{/if}
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
	@use "sass:color";

	li {
		box-sizing: border-box;
		width: 100%;
		height: 14.2rem;
		display: flex;
		flex-direction: column;
		margin: 0;
		padding: 0;
		box-shadow: inset 0 0 0 4px #636669;
		background-color: white;
		border-radius: 0.5rem;
		overflow: hidden;
		cursor: pointer;

		&:hover {
			box-shadow: inset 0 0 0 4px color.adjust(#636669, $lightness: -5%);

			h3 {
				background-color: color.adjust(#636669, $lightness: -5%);
			}
		}

		&.selected {
			box-shadow: inset 0 0 0 4px #2452b6;

			h3 {
				background-color: #2452b6;
			}
		}
	}

	h3 {
		box-sizing: border-box;
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		background-color: #636669;
		color: white;
		width: 100%;
		padding: 0.4rem 0.5rem 0.2rem 0.5rem;
	}

	.content {
		display: flex;
		flex: 1;
		min-height: 0;
		padding: 0 0.2rem 3px;
		gap: 0.6rem;
	}

	.items-wrapper {
		box-sizing: border-box;
		height: 100%;
		overflow: scroll;
		padding: 0.5rem;
	}

	.items {
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
		padding: 0;
	}
</style>
