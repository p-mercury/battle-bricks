<script lang="ts">
	import type { Unit } from "$lib/units";

	let { unit }: { unit: Unit } = $props();
</script>

<section>
	<div>{unit.name}</div>
	<div>Price: {unit.price}</div>
	<div>Size: {unit.size}</div>
	<div>Health: {unit.health}</div>
	<div>Armor Class: {unit.armorClass}</div>
	<div>Carry Capacity: {unit.carryCapacity}</div>
	<div>Accuracy: {unit.accuracy}</div>
</section>

<style lang="scss">
	section {
		padding: 0.5rem;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightcyan;
	}
</style>
