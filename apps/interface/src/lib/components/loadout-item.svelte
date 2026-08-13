<script lang="ts">
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import Item from "$lib/components/item.svelte";
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import Brick from "$lib/components/brick.svelte";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: { loadout: Loadout; selected?: Boolean; onclick?: () => void } = $props();
</script>

<section {onclick} class:selected>
	<header>
		<Brick />
	</header>
	{#if loadout.unit}
		{#if loadout.unit.image}
			<img alt={loadout.unit.name} src={loadout.unit.image} />
		{/if}
		<div class="info">
			<h3>{loadout.unit.name}</h3>
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
					Speed:
					<StatBar value={loadout.unit.speed} size={4} />
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

<style lang="scss">
	@use "sass:color";

	section {
		background: white;
		border-radius: 0.8rem;
		padding: 0.6rem;
		overflow: hidden;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		cursor: pointer;
		&.selected {
			outline: 4px solid #2388ff;
			outline-offset: 2px;
		}
		display: grid;
		width: 33.7rem;
		height: 15rem;
		grid-template:
			"color color color" 1.25rem
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

		&.republic {
			background: linear-gradient(135deg, #1c1c1c, #4f424f);
		}

		&.separatists {
			background: linear-gradient(135deg, #2b1515, #2d3f54);
		}
	}

	header {
		grid-area: color;
		margin: -0.6rem -0.6rem 0 -0.6rem;
		background-color: #595d60;
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
