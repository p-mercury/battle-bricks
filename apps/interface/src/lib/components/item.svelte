<script lang="ts">
	import type { Ammunition, RangeWeapon, MeleeWeapon } from "$lib/items";
	import StatBar from "./stat-bar.svelte";
	import { slide } from "svelte/transition";

	let { item }: { item: Ammunition | RangeWeapon | MeleeWeapon } = $props();

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
	{#if isOpen}
		<div class="stats" transition:slide={{ duration: 200 }}>
			<div class="row">
				<div>Price:</div>
				<div>{item.price}</div>
			</div>
			<div class="row">
				<div>Weight:</div>
				<div>{item.weight}</div>
			</div>
			{#if item.type === "AMMUNITION"}
				<div class="row">
					<div>Type:</div>
					<div>{item.ammunitionType}</div>
				</div>
				<div class="row">
					<div>Capacity:</div>
					<div>{item.capacity}</div>
				</div>
				<div class="row">
					<div>Damage:</div>
					<div>{item.damage}</div>
				</div>
				<div class="row">
					<div>Armor Piercing:</div>
					<StatBar value={item.armorPiercing} size={4} />
				</div>
			{:else if item.type === "RANGE_WEAPON"}
				<div class="row">
					<div>Type:</div>
					<div>{item.ammunitionType}</div>
				</div>
				<div class="row">
					<div>Tange:</div>
					<div>{item.range.min}-{item.range.max}m</div>
				</div>
				<div class="row">
					<div>Fire Rate:</div>
					<div>{item.fireRate}</div>
				</div>
			{:else}
				<div class="row">
					<div>Armor Piercing:</div>
					<StatBar value={item.armorPiercing} size={4} />
				</div>
				<div class="row">
					<div>Damage:</div>
					<div>{item.damage}</div>
				</div>
			{/if}
		</div>
	{/if}
</section>

<style lang="scss">
	section {
		padding: 0;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: white;
		overflow: hidden;
		height: max-content;
	}

	.header {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: space-between;
		background-color: lightgray;
		width: 100%;
		padding: 0.2rem 0.5rem 0.2rem 0.5rem;
		border: none;
		cursor: pointer;
		font: inherit;
		color: inherit;
		text-align: left;

		&:hover {
			background-color: darken(lightgray, 5%);
		}
	}

	h3 {
		margin: 0;
		padding: 0;
		font-size: 1rem;
		font-weight: 600;
	}

	svg {
		flex-shrink: 0;
		transition: transform 0.2s ease;

		&.rotated {
			transform: rotate(180deg);
		}
	}

	.stats {
		display: grid;
		grid-template-columns: auto auto;
		gap: 0.5rem;
		padding: 0.5rem;
		flex-shrink: 0;
	}

	.row {
		display: grid;
		list-style: none;
		grid-column: 1 / -1;
		grid-template-columns: subgrid;
	}
</style>
