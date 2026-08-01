<script lang="ts">
	import type { Loadout } from "$lib/loadouts";
	import Item from "$lib/components/item.svelte";
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import NumberInput from "$lib/components/number-input.svelte";

	let {
		loadout,
		initiative = $bindable(),
		selected = false,
		onclick = () => {},
	}: {
		loadout: Loadout;
		initiative: number;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();

	let liElement: HTMLLIElement;

	$effect(() => {
		if (selected) {
			liElement.scrollIntoView({ behavior: "smooth", block: "nearest" });
		}
	});
</script>

<li {onclick} bind:this={liElement}>
	<section class:selected>
		<header>
			<h3>{loadout.name}</h3>
			<NumberInput min={1} max={100} bind:value={initiative} />
		</header>
		<div class="content">
			<StatTable>
				<StatRow>
					Health:
					<NumberInput bind:value={loadout.unit.health} />
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
			<div class="items-wrapper">
				<div class="items">
					{#each loadout.items as item}
						<Item {item} />
					{/each}
				</div>
			</div>
		</div>
	</section>
</li>

<style lang="scss">
	@use "sass:color";

	li {
		all: unset;
	}

	section {
		box-sizing: border-box;
		width: 100%;
		height: 16rem;
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

			header {
				background-color: color.adjust(#636669, $lightness: -5%);
			}
		}

		&.selected {
			box-shadow: inset 0 0 0 4px #2452b6;

			header {
				background-color: #2452b6;
			}
		}
	}

	header {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		box-sizing: border-box;
		margin: 0;
		background-color: #636669;
		color: white;
		width: 100%;
		padding: 0.4rem 0.5rem 0.4rem 0.5rem;
	}

	h3 {
		font-size: 1.1rem;
		font-weight: 600;
		margin: 0;
		padding: 0;
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
