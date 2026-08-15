<script lang="ts">
	import type {
		MeleeAction,
		BlasterAction,
		CannonAction,
		LauncherAction,
		Action,
	} from "$lib/get-actions";
	import NumberInput from "./number-input.svelte";
	import Brick from "$lib/components/brick.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import StatRow from "./stat-row.svelte";
	import DiceRoll from "./dice-roll.svelte";

	let {
		action,
	}: {
		action:
			BlasterAction | CannonAction | LauncherAction | MeleeAction | Action;
	} = $props();
</script>

<section>
	<header>
		<Brick />
	</header>
	<div class="info">
		<h3>{action.name}</h3>
		<StatTable>
			{#if action.type === "BLASTER" || action.type === "CANNON"}
				<StatRow>
					Fire Rate:
					<span>{action.weapon.fireRate}</span>
				</StatRow>
				<StatRow>
					Amunition:
					<NumberInput
						bind:value={action.ammunition.quantity}
						min={0}
						max={2000}
					/>
				</StatRow>
				<StatRow>
					Range:
					<span>{action.weapon.range!.min}-{action.weapon.range!.max}m</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.toHit}</span>
				</StatRow>
				<StatRow>
					To pierce:
					<span>≥{action.toPierce}</span>
				</StatRow>
				<StatRow>
					Damage:
					<DiceRoll roll={action.damage} />
				</StatRow>
				<StatRow>
					Damage chance:
					<span>{action.damageChance.toFixed(0)}%</span>
				</StatRow>
			{:else if action.type === "LAUNCHER"}
				<StatRow>
					Fire Rate:
					<span>{action.weapon.fireRate}</span>
				</StatRow>
				<StatRow>
					Amunition:
					<NumberInput
						bind:value={action.ammunition.quantity}
						min={0}
						max={2000}
					/>
				</StatRow>
				<StatRow>
					Range:
					<span>
						{action.ammunition.item.details.value.range?.min}-{action.ammunition
							.item.details.value.range?.max}m
					</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.toHit}</span>
				</StatRow>
				<StatRow>
					To pierce:
					<span>≥{action.toPierce}</span>
				</StatRow>
				<StatRow>
					Damage:
					<DiceRoll roll={action.damage} />
				</StatRow>
				<StatRow>
					Damage chance:
					<span>{action.damageChance.toFixed(0)}%</span>
				</StatRow>
			{:else if action.type === "MELEE"}
				<StatRow>
					Attack Speed:
					<span>{action.weapon.attackSpeed}</span>
				</StatRow>
				<StatRow>
					To hit:
					<span>≥{action.toHit}</span>
				</StatRow>
				<StatRow>
					To pierce:
					<span>≥{action.toPierce}</span>
				</StatRow>
				<StatRow>
					Damage:
					<DiceRoll roll={action.damage} />
				</StatRow>
				<StatRow>
					Damage chance:
					<span>{action.damageChance.toFixed(0)}%</span>
				</StatRow>
			{:else if action.type === "ACTION"}
				<div>
					<div>{action.description}</div>
				</div>
			{/if}
		</StatTable>
	</div>
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
		display: grid;
		grid-template:
			"info color" auto /
			auto 1.25rem;
		gap: 0.6rem;
	}

	header {
		grid-area: color;
		margin: -0.6rem -0.6rem -0.6rem 0;
		background-color: #595d60;
	}

	.info {
		grid-area: info;
		padding: 0.1rem;

		h3 {
			margin: 0;
			padding: 0 0 0.5rem 0;
			font-size: 1.1rem;
			font-weight: 600;
		}
	}
</style>
