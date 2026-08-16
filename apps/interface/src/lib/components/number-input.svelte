<script lang="ts">
	import { scale } from "svelte/transition";

	let {
		value = $bindable(0),
		pre,
		min = 0,
		max = undefined,
		disabled = false,

		onchange = undefined,
	}: {
		value: number;
		pre?: number;
		min?: number;
		max?: number;
		disabled?: boolean;

		onchange?: (value: number) => void;
	} = $props();

	let canIncrement = $derived(!disabled && (max == null || value < max));
	let canDecrement = $derived(!disabled && (min == null || value > min));

	$effect(() => {
		if (onchange) {
			onchange(value);
		}
	});

	$effect(() => {
		if (max != null) value = Math.min(value, max);
		if (min != null) value = Math.max(value, min);
	});

	function handleIncrementClick(event: Event) {
		event.preventDefault();
		event.stopPropagation();
		value++;
	}

	function handleDecrementClick(event: Event) {
		event.preventDefault();
		event.stopPropagation();
		value--;
	}
</script>

<div class="qty-row">
	<div class="qty">
		{#if pre}
			<button
				in:scale={{ duration: 300 }}
				class="micro-btn"
				aria-label="Decrease quantity"
				onclick={(event) => {
					event.preventDefault();
					event.stopPropagation();
					value -= pre;
				}}
				class:disabled={!canDecrement}
			>
				-{pre}
			</button>
		{/if}
		<button
			in:scale={{ duration: 300 }}
			class="micro-btn"
			aria-label="Decrease quantity"
			onclick={handleDecrementClick}
			class:disabled={!canDecrement}
		>
			<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M6 12H18"
					stroke="#000000"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				></path>
			</svg>
		</button>
		<div in:scale={{ duration: 300 }} class="qty-input">
			{value}
		</div>
		<button
			class="micro-btn"
			aria-label="Increase quantity"
			onclick={handleIncrementClick}
			class:disabled={!canIncrement}
		>
			<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M6 12H18"
					stroke="#000000"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				></path>
				<path
					d="M12 6V18"
					stroke="#000000"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				></path>
			</svg>
		</button>
	</div>
</div>

<style lang="scss">
	.qty-row {
		display: flex;
		align-items: center;
	}

	.qty {
		display: flex;
		align-items: center;
		gap: 0.2rem;
		background: lightgray;
		border-radius: 0.6rem;
	}

	.micro-btn {
		border: none;
		width: 1.8rem;
		height: 1.8rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: darkgray;
		border-radius: 0.6rem;
		cursor: pointer;

		svg {
			width: 1rem;
			height: 1rem;

			* {
				stroke: blacks;
			}
		}
	}

	.micro-btn.disabled {
		opacity: 0.4;
		pointer-events: none;
	}

	.qty-input {
		width: 1.5rem;
		text-align: center;
		border: none;
		background: transparent;
		font-size: 0.8rem;
		outline: none;
		padding: 0;
		color: black;
	}
</style>
