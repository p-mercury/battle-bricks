<script lang="ts">
	import StatBar from "./stat-bar.svelte";
	import StatTable from "./stat-table.svelte";
	import StatRow from "./stat-row.svelte";
	import type { Item } from "@battle-bricks/contracts/catalogue/v1/item_pb";
	import DiceRoll from "./dice-roll.svelte";

	let {
		item,
	}: {
		item: Item;
	} = $props();

	let isOpen = $state(false);
</script>

<section>
	<button
		class="header"
		onclick={(event) => {
			event.preventDefault();
			event.stopPropagation();
			isOpen = !isOpen;
		}}
		aria-expanded={isOpen}
	>
		<h3>{item.name}</h3>
		<svg
			class:rotated={isOpen}
			width="16"
			height="16"
			viewBox="0 0 16 16"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			<path
				d="M4 6L8 10L12 6"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			/>
		</svg>
	</button>
	<div class:open={isOpen} class="stats" aria-hidden={!isOpen}>
		<StatTable>
			<StatRow>
				Price:
				<div>{item.price}</div>
			</StatRow>
			<StatRow>
				Weight:
				<div>{item.weight}</div>
			</StatRow>
			{#if item.details.case === "blasterBolt"}
				<StatRow>
					Capacity:
					<div>{item.details.value.capacity}</div>
				</StatRow>
				<StatRow>
					Damage:
					<DiceRoll roll={item.details.value.damage} />
				</StatRow>
				<StatRow>
					Armor Piercing:
					<StatBar value={item.details.value.armorPiercing} size={5} red />
				</StatRow>
			{:else if item.details.case === "cannonShell"}
				<StatRow>
					Damage:
					<DiceRoll roll={item.details.value.damage} />
				</StatRow>
				<StatRow>
					Armor Piercing:
					<StatBar value={item.details.value.armorPiercing} size={5} red />
				</StatRow>
			{:else if item.details.case === "launcherRocket"}
				<StatRow>
					Damage:
					<DiceRoll roll={item.details.value.damage} />
				</StatRow>
				<StatRow>
					Splash Radius:
					<StatBar value={item.details.value.splashRadius} size={5} red />
				</StatRow>
			{:else if item.details.case === "blaster"}
				<StatRow>
					Range:
					<div>
						{item.details.value.range!.min}-{item.details.value.range!.max}m
					</div>
				</StatRow>
				<StatRow>
					Fire Rate:
					<div>{item.details.value.fireRate}</div>
				</StatRow>
			{:else if item.details.case === "cannon"}
				<StatRow>
					Range:
					<div>
						{item.details.value.range!.min}-{item.details.value.range!.max}m
					</div>
				</StatRow>
				<StatRow>
					Fire Rate:
					<div>{item.details.value.fireRate}</div>
				</StatRow>
			{:else if item.details.case === "launcher"}
				<StatRow>
					Fire Rate:
					<div>{item.details.value.fireRate}</div>
				</StatRow>
			{:else if item.details.case === "meleeWeapon"}
				<StatRow>
					Attack Speed:
					<StatBar value={item.details.value.attackSpeed} size={5} red />
				</StatRow>
				<StatRow>
					Armor Piercing:
					<StatBar value={item.details.value.armorPiercing} size={5} red />
				</StatRow>
				<StatRow>
					Damage:
					<DiceRoll roll={item.details.value.damage} />
				</StatRow>
			{/if}
		</StatTable>
	</div>
</section>

<style lang="scss">
	@use "sass:color";

	section {
		box-shadow: inset 0 0 0 4px #636669;
		border-radius: 0.4rem;
		background-color: white;
		overflow: hidden;
		height: max-content;
		cursor: default;

		&:hover {
			box-shadow: inset 0 0 0 4px color.adjust(#636669, $lightness: -5%);

			.header {
				background-color: color.adjust(#636669, $lightness: -5%);
			}
		}
	}

	.header {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: space-between;
		background-color: #636669;
		width: 100%;
		padding: 0.2rem 0.5rem;
		border: none;
		cursor: pointer;
		font: inherit;
		color: inherit;
		text-align: left;
	}

	h3 {
		margin: 0;
		padding: 0;
		color: white;
		font-size: 0.9rem;
		font-weight: 600;
	}

	svg {
		flex-shrink: 0;
		transition: transform 0.3s ease;

		&.rotated {
			transform: rotate(180deg);
		}
	}

	.stats {
		box-sizing: border-box;
		padding: 0 0.2rem;
		max-height: 0;
		overflow: hidden;
		visibility: hidden;
		opacity: 0;

		transition:
			max-height 0.3s ease,
			opacity 0.3s ease,
			visibility 0s linear 0.3s;

		&.open {
			max-height: 1000px;
			visibility: visible;
			opacity: 1;

			transition:
				max-height 0.3s ease,
				opacity 0.3s ease,
				visibility 0s linear 0s;
		}
	}
</style>
