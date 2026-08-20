<script lang="ts">
	import Item from "$lib/components/item.svelte";
	import BrickCard from "./brick-card.svelte";
	import Stat from "$lib/components/stat.svelte";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import { getLoadoutPrice } from "$lib/loadout";

	let {
		loadout = $bindable(),
		drag,
		selected = false,
		onclick = () => {},
	}: {
		loadout: Loadout;
		drag?: Boolean;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();

	let itemId = $state(loadout.item?.id);

	$effect(() => {
		loadout.item = loadout.unit!.loadoutItems.find((i) => i.id === itemId);
	});
</script>

<BrickCard {selected} {drag} {onclick}>
	<section>
		{#if loadout.unit!.image}
			<img alt={loadout.unit!.name} src={loadout.unit!.image} />
		{/if}
		<div class="info">
			<h3>{loadout.unit!.name}</h3>
		</div>
		<div class="price">
			{getLoadoutPrice(loadout)}c
		</div>
		<div class="stats">
			<Stat label="HP" color="GREEN" value={loadout.unit!.health} />
			<Stat label="SZ" color="BLUE" value={loadout.unit!.size} />
			<Stat label="GS" color="RED" value={loadout.unit!.marksmanship || 0} />
			<Stat label="AC" color="GREEN" value={loadout.unit!.armorClass} />
			<Stat label="SP" color="BLUE" value={loadout.unit!.speed} />
			<Stat label="MS" color="RED" value={loadout.unit!.meleeAbility || 0} />
		</div>
		{#if loadout.unit!.loadoutItems.length}
			<select class="loadout" bind:value={itemId}>
				<option value={undefined}>Extra item</option>
				{#each loadout.unit!.loadoutItems as item}
					<option value={item.id}>{item.name} ({item.price}c)</option>
				{/each}
			</select>
		{/if}
		<div class="items-wrapper">
			<div class="items">
				{#each loadout.unit!.items as item}
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
			"info info price" 1.3rem
			"image stats stats" 7.5rem
			"loadout loadout loadout" auto
			"items items items" fit-content(8rem) /
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
	}

	.loadout {
		grid-area: loadout;
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
