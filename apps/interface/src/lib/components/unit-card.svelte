<script lang="ts">
	import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
	import Item from "$lib/components/item.svelte";
	import BrickCard from "./brick-card.svelte";
	import Stat from "$lib/components/stat.svelte";

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
			<Stat label="HP" color="GREEN" value={unit.health} />
			<Stat label="SZ" color="BLUE" value={unit.size} />
			<Stat label="GS" color="RED" value={unit.marksmanship || 0} />
			<Stat label="AC" color="GREEN" value={unit.armorClass} />
			<Stat label="SP" color="BLUE" value={unit.speed} />
			<Stat label="MS" color="RED" value={unit.meleeAbility || 0} />
		</div>
		<div class="items-wrapper">
			<div class="items">
				{#each unit.items as item}
					<Item {item} />
				{/each}
			</div>
		</div>
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		padding: 0.6rem;
		display: grid;
		width: 18.6rem;
		grid-template:
			"info info" 1.3rem
			"image stats" 7.5rem
			"items items" fit-content(8rem) /
			7.5rem auto;
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

	.items-wrapper {
		grid-area: items;
		overflow-y: scroll;
		padding-left: 0.1rem;
	}

	.items {
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
		padding: 0;
	}

	.stats {
		grid-area: stats;
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		grid-auto-rows: 1fr;
		gap: 0.4rem;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0;
			font-size: 1.2rem;
			line-height: 1.2rem;
			font-weight: 600;
		}
	}
</style>
