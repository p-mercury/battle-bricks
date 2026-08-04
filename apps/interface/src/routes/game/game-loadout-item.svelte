<script lang="ts">
	import Item from "$lib/components/item.svelte";
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import NumberInput from "$lib/components/number-input.svelte";
	import type { GameLoadout } from "$lib/game";
	import Switcher from "$lib/components/switcher.svelte";
	import SelectColor from "$lib/components/select-color.svelte";

	let {
		loadout,
		selected = false,
		onclick = () => {},
	}: {
		loadout: GameLoadout;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<li {onclick} class:selected>
	<section>
		<header onclick={(e) => e.stopPropagation()}>
			<SelectColor bind:value={loadout.color} />
		</header>
		{#if loadout.image}
			<img alt={loadout.name} src={loadout.image} />
		{/if}
		<div class="info">
			<h3>{loadout.name}</h3>
			<StatTable>
				<StatRow>
					Health:
					<NumberInput bind:value={loadout.unit.health} />
				</StatRow>
				{#if "inCover" in loadout}
					<StatRow>
						In Cover:
						<Switcher
							bind:value={loadout.inCover}
							onclick={(event) => event.stopPropagation()}
						/>
					</StatRow>
				{/if}
			</StatTable>
		</div>
		<div class="stats">
			<StatTable>
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
	</section>
</li>

<style lang="scss">
	@use "sass:color";

	li {
		all: unset;
		background: white;
		border-radius: 0.8rem;
		padding: 0.5rem;
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
	}

	section {
		display: grid;
		grid-template:
			"color color color" 1.2rem
			"image info items" 8em
			"stats stats stats" 1fr /
			8rem auto auto;
		gap: 0.6rem;
	}

	img {
		max-width: 100%;
		max-height: 100%;
		grid-area: image;
		object-fit: contain;
		border-radius: 0.8rem;
		background: linear-gradient(135deg, #1c1c1c, #4f424f);
		box-sizing: border-box;
	}

	header {
		grid-area: color;
		margin: -0.5rem -0.5rem 0 -0.5rem;
	}

	.stats {
		grid-area: stats;
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
