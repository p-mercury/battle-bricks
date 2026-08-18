<script lang="ts">
	import type { Loadout } from "$lib/game";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import BrickCard from "$lib/components/brick-card.svelte";
	import Stat from "$lib/components/stat.svelte";
	import { getUnitPrice } from "$lib/unit";

	let {
		faction,
		loadout,
		selected = false,
		onclick = () => {},
	}: {
		faction: Faction;
		loadout: Loadout;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<BrickCard {selected} {onclick} bind:color={loadout.color}>
	<section>
		{#if loadout.image}
			<img
				alt={loadout.name}
				src={loadout.image}
				class:republic={faction === Faction.GALACTIC_REPUBLIC}
				class:separatists={faction === Faction.SEPARATIST_ALLIANCE}
			/>
		{/if}
		<div class="info">
			<h3>{loadout.name}</h3>
		</div>
		<div class="price">
			{getUnitPrice(loadout.unit)}c
		</div>
		<div class="stats">
			<Stat label="HP" color="GREEN" value={loadout.unit.health} />
			<Stat label="SZ" color="BLUE" value={loadout.unit.size} />
			<Stat label="GS" color="RED" value={loadout.unit.marksmanship || 0} />
			<Stat label="AC" color="GREEN" value={loadout.unit.armorClass} />
			<Stat label="SP" color="BLUE" value={loadout.unit.speed} />
			<Stat label="MS" color="RED" value={loadout.unit.meleeAbility || 0} />
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
			"info info price" 1.3rem
			"image stats stats" 7.5rem /
			7.5rem auto min-content;
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

		&.republic {
			background: linear-gradient(135deg, #1c1c1c, #4f424f);
		}

		&.separatists {
			background: linear-gradient(135deg, #2b1515, #2d3f54);
		}
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

	.price {
		grid-area: price;
		color: black;
		font-weight: 500;
		background-color: orange;
		line-height: 1rem;
		padding: 0.2rem;
		border-radius: 0.4rem;
	}
</style>
