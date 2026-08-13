<script lang="ts">
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import Item from "$lib/components/item.svelte";
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import BrickCard from "./brick-card.svelte";

	let {
		loadout,
		drag,
		selected = false,
		onclick = () => {},
	}: {
		loadout: Loadout;
		drag?: Boolean;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<BrickCard {selected} {drag} {onclick}>
	<section>
		{#if loadout.unit}
			{#if loadout.image}
				<img alt={loadout.name} src={loadout.image} />
			{/if}
			<div class="info">
				<h3>{loadout.name}</h3>
			</div>
			<div class="stats">
				<StatTable>
					<StatRow>
						Price:
						<div>
							{loadout.unit.price +
								loadout.items.reduce((total, item) => total + item.price, 0)}c
						</div>
					</StatRow>
					<StatRow>
						Health:
						<div>{loadout.unit!.health}hp</div>
					</StatRow>
					<StatRow>
						Size:
						<StatBar value={loadout.unit.size} size={4} />
					</StatRow>
					<StatRow>
						Speed:
						<StatBar value={loadout.unit.speed} size={4} />
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
			</div>
			<div class="items-wrapper">
				<div class="items">
					{#each loadout.items as item}
						<Item {item} />
					{/each}
				</div>
			</div>
		{/if}
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		padding: 0.6rem;
		display: grid;
		width: 33.7rem;
		height: 13rem;
		grid-template:
			"image stats items" 8rem
			"info stats items" auto /
			8rem min-content 1fr;
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
		padding-left: 0.3rem;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0;
			font-size: 1.2rem;
			font-weight: 600;
		}
	}
</style>
