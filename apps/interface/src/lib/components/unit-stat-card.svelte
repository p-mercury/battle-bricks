<script lang="ts">
	import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import BrickCard from "./brick-card.svelte";
	import { getUnitPrice } from "$lib/unit";

	let {
		unit,
		drag,
		selected = false,
		onclick = () => {},
	}: {
		unit: Unit;
		drag?: Boolean;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<BrickCard {selected} {drag} {onclick}>
	<section>
		{#if unit.image}
			<img alt={unit.name} src={unit.image} />
		{/if}
		<div class="info">
			<h3>{unit.name}</h3>
		</div>
		<div class="stats">
			<StatTable>
				<StatRow>
					Price:
					<div>
						{getUnitPrice(unit)}c
					</div>
				</StatRow>
				<StatRow>
					Health:
					<div>{unit.health}hp</div>
				</StatRow>
				<StatRow>
					Size:
					<StatBar value={unit.size} size={5} />
				</StatRow>
				<StatRow>
					Speed:
					<StatBar value={unit.speed} size={5} />
				</StatRow>
				<StatRow>
					Armor Class:
					<StatBar value={unit.armorClass} size={5} />
				</StatRow>
				{#if unit.marksmanship}
					<StatRow>
						Marksmanship:
						<StatBar value={unit.marksmanship} size={5} red />
					</StatRow>
				{/if}
				{#if unit.meleeAbility}
					<StatRow>
						Melee Ability:
						<StatBar value={unit.meleeAbility} size={5} red />
					</StatRow>
				{/if}
			</StatTable>
		</div>
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		padding: 0.6rem;
		display: grid;
		width: 21.65rem;
		grid-template:
			"image stats" 7rem
			"info stats" auto /
			7rem min-content;
		gap: 0.6rem;
	}

	img {
		width: 100%;
		height: 100%;
		grid-area: image;
		object-fit: contain;
		border-radius: 0.8rem;
		padding: 0.4rem;
		background: linear-gradient(135deg, #1c1c1c, #4f424f);
		box-sizing: border-box;
	}

	.stats {
		grid-area: stats;
		padding-left: 0.3rem;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0;
			font-size: 1.1rem;
			font-weight: 600;
		}
	}
</style>
