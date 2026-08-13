<script lang="ts">
	import type { DiceRoll } from "@battle-bricks/contracts/common/v1/dice-roll_pb";

	let { roll }: { roll?: DiceRoll } = $props();
</script>

<span>
	{#if roll}
		{roll.count}d{roll.sides}{roll.modifier
			? roll.modifier > 0
				? `+${roll.modifier}`
				: roll.modifier
			: ""}
	{/if}
</span>
