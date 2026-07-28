<script lang="ts">
	import type { Ammunition, RangeWeapon, MeleeWeapon } from "$lib/items";

	let { item }: { item: Ammunition | RangeWeapon | MeleeWeapon } = $props();
</script>

<section>
	<div>Name: {item.name}</div>
	<div>Price: {item.price}</div>
	<div>Weight: {item.weight}</div>
	{#if item.type === "AMMUNITION"}
		<div>Type: {item.ammunitionType}</div>
		<div>Capacity: {item.capacity}</div>
		<div>Damage: {item.damage}</div>
		<div>Armor Piercing: {item.armorPiercing}</div>
	{:else if item.type === "RANGE_WEAPON"}
		<div>Type: {item.ammunitionType}</div>
		<div>Tange: {item.range.min}-{item.range.max}m</div>
		<div>Fire Rate: {item.fireRate}</div>
	{:else}
		<div>Armor Piercing: {item.armorPiercing}</div>
		<div>Damage: {item.damage}</div>
	{/if}
</section>

<style lang="scss">
	section {
		padding: 0.5rem;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightcoral;
	}
</style>
