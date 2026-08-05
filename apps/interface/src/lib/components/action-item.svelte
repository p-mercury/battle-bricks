<script lang="ts">
	import type {
		MeleeAction,
		RangeBoltAction,
		RangeRocketAction,
		RangeShellAction,
	} from "$lib/get-actions";
	import type { Action } from "$lib/items";
	import { untrack } from "svelte";
	import NumberInput from "./number-input.svelte";

	let {
		action,
	}: {
		action:
			| RangeBoltAction
			| RangeShellAction
			| RangeRocketAction
			| MeleeAction
			| Action;
	} = $props();

	let ammunition = $state(0);

	$effect(() => {
		if (action.type === "RANGE_BOLT") {
			const newAmmunition = action.ammunition.reduce(
				(state, item) => state + item.capacity,
				0,
			);

			if (newAmmunition != untrack(() => ammunition)) {
				ammunition = newAmmunition;
			}
		} else if (action.type === "RANGE_SHELL") {
			const newAmmunition = action.ammunition.length;
			if (newAmmunition != untrack(() => ammunition)) {
				ammunition = newAmmunition;
			}
		} else if (action.type === "RANGE_ROCKET") {
			const newAmmunition = action.ammunition.length;
			if (newAmmunition != untrack(() => ammunition)) {
				ammunition = newAmmunition;
			}
		}
	});

	$effect(() => {
		if (action.type === "RANGE_BOLT") {
			const newTotal = Math.max(0, Math.trunc(ammunition));
			untrack(() => {
				const count = action.ammunition.length;

				if (count === 0) return;

				const amountPerItem = Math.floor(newTotal / count);
				let remainder = newTotal % count;

				action.ammunition.forEach((item) => {
					item.capacity = amountPerItem + (remainder > 0 ? 1 : 0);

					if (remainder > 0) {
						remainder--;
					}
				});
			});
		}
	});
</script>

<li>
	{#if action.type === "RANGE_BOLT"}
		<b>{action.weapon.name} with {action.ammunition[0].name}</b>
		<div>Fire Rate: {action.weapon.fireRate}</div>
		<div>
			Amunition: <NumberInput bind:value={ammunition} min={0} max={2000} />
		</div>
		<div>Range: {action.weapon.range.min}-{action.weapon.range.max}m</div>
		<div>To hit: ≥{action.b1r}</div>
		{#if action.b2}
			<div>To pierce armor: ≥{action.b2}</div>
		{/if}
		<div>Damage: {action.damage}</div>
	{:else if action.type === "RANGE_SHELL"}
		<b>{action.weapon.name} with {action.ammunition[0].name}</b>
		<div>Fire Rate: {action.weapon.fireRate}</div>
		<div>Amunition: {ammunition}</div>
		<div>Range: {action.weapon.range.min}-{action.weapon.range.max}m</div>
		<div>To hit: ≥{action.b1r}</div>
		{#if action.b2}
			<div>To pierce armor: ≥{action.b2}</div>
		{/if}
		<div>Damage: {action.damage}</div>
	{:else if action.type === "RANGE_ROCKET"}
		<b>{action.weapon.name} with {action.ammunition[0].name}</b>
		<div>Fire Rate: {action.weapon.fireRate}</div>
		<div>Amunition: {ammunition}</div>
		<div>Range: {action.weapon.range.min}-{action.weapon.range.max}m</div>
		<div>To hit: ≥{action.b1r}</div>
		{#if action.b2}
			<div>To pierce armor: ≥{action.b2}</div>
		{/if}
		<div>Damage: {action.damage}</div>
		<div>
			Splash radius: {action.ammunition[0].splashRadius}m
		</div>
	{:else if action.type === "MELEE"}
		<b>{action.weapon.name}</b>
		<div>Attack Speed: {action.weapon.attackSpeed}</div>
		<div>To hit: ≥{action.b1r}</div>
		<div>To pierce armor: ≥{action.b2}</div>
		<div>Damage: {action.damage}</div>
	{:else if action.type === "ACTION"}
		<b>{action.name}</b>
		<div>{action.description}</div>
	{/if}
</li>

<style lang="scss">
	li {
		padding: 0.5rem;
		list-style: none;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightsalmon;
	}
</style>
