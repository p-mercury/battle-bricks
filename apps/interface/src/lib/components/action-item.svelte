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
	import Brick from "$lib/components/brick.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import StatRow from "./stat-row.svelte";

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
		if (
			action.type === "RANGE_BOLT" ||
			action.type === "RANGE_SHELL" ||
			action.type === "RANGE_ROCKET"
		) {
			const newAmmunition = action.ammunition.reduce(
				(state, item) => state + item.capacity,
				0,
			);
			if (newAmmunition != untrack(() => ammunition)) {
				ammunition = newAmmunition;
			}
		}
	});

	$effect(() => {
		if (
			action.type === "RANGE_BOLT" ||
			action.type === "RANGE_SHELL" ||
			action.type === "RANGE_ROCKET"
		) {
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

<section>
	<header>
		<Brick />
	</header>
	<div class="info">
		{#if action.type === "RANGE_BOLT"}
			<h3>{action.weapon.name} with {action.ammunition[0].name}</h3>
			<StatTable>
				<StatRow>
					Fire Rate:
					<span>{action.weapon.fireRate}</span>
				</StatRow>
				<StatRow>
					Amunition:
					<NumberInput bind:value={ammunition} min={0} max={2000} />
				</StatRow>
				<StatRow>
					Range:
					<span>{action.weapon.range.min}-{action.weapon.range.max}m</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.b1r}</span>
				</StatRow>
				{#if action.b2}
					<StatRow>
						To pierce armor:
						<span>≥{action.b2}</span>
					</StatRow>
				{/if}
				<StatRow>
					Damage:
					<span>{action.damage}</span>
				</StatRow>
			</StatTable>
		{:else if action.type === "RANGE_SHELL"}
			<h3>{action.weapon.name} with {action.ammunition[0].name}</h3>
			<StatTable>
				<StatRow>
					Fire Rate:
					<span>{action.weapon.fireRate}</span>
				</StatRow>
				<StatRow>
					Amunition:
					<NumberInput
						bind:value={ammunition}
						min={0}
						max={action.ammunition.length}
					/>
				</StatRow>
				<StatRow>
					Range:
					<span>{action.weapon.range.min}-{action.weapon.range.max}m</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.b1r}</span>
				</StatRow>
				{#if action.b2}
					<StatRow>
						To pierce armor:
						<span>≥{action.b2}</span>
					</StatRow>
				{/if}
				<StatRow>
					Damage:
					<span>{action.damage}</span>
				</StatRow>
			</StatTable>
		{:else if action.type === "RANGE_ROCKET"}
			<h3>{action.weapon.name} with {action.ammunition[0].name}</h3>
			<StatTable>
				<StatRow>
					Fire Rate:
					<span>{action.weapon.fireRate}</span>
				</StatRow>
				<StatRow>
					Amunition:
					<NumberInput
						bind:value={ammunition}
						min={0}
						max={action.ammunition.length}
					/>
				</StatRow>
				<StatRow>
					Range:
					<span>{action.weapon.range.min}-{action.weapon.range.max}m</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.b1r}</span>
				</StatRow>
				{#if action.b2}
					<StatRow>
						To pierce armor:
						<span>≥{action.b2}</span>
					</StatRow>
				{/if}
				<StatRow>
					Damage:
					<span>{action.damage}</span>
				</StatRow>
			</StatTable>
		{:else if action.type === "MELEE"}
			<h3>{action.weapon.name}</h3>
			<StatTable>
				<StatRow>
					Attack Speed:
					<span>{action.weapon.attackSpeed}</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.b1r}</span>
				</StatRow>
				<StatRow>
					To pierce armor:
					<span>≥{action.b2}</span>
				</StatRow>
				<StatRow>
					Damage:
					<span>{action.damage}</span>
				</StatRow>
			</StatTable>
		{:else if action.type === "ACTION"}
			<h3>{action.name}</h3>
			<div>
				<div>{action.description}</div>
			</div>
		{/if}
	</div>
</section>

<style lang="scss">
	@use "sass:color";

	section {
		background: white;
		border-radius: 0.8rem;
		padding: 0.5rem;
		overflow: hidden;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		cursor: pointer;
		display: grid;
		grid-template:
			"info color" auto /
			auto 1.25rem;
		gap: 0.6rem;
	}

	header {
		grid-area: color;
		margin: -0.5rem -0.5rem -0.5rem 0;
		background-color: #595d60;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0 0 0 0.2rem;
			font-size: 1.1rem;
			font-weight: 600;
		}
	}
</style>
