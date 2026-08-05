<script lang="ts">
	import { scale } from "svelte/transition";

	let {
		value = $bindable(0),
		min = 0,
		max = undefined,
		disabled = false,

		onchange = undefined,
	}: {
		value: number;
		min?: number;
		max?: number;
		disabled?: boolean;

		onchange?: (value: number) => void;
	} = $props();

	let effectiveMax = $derived(max ? max : 100);
	let canIncrement = $derived(!disabled && value < effectiveMax);
	let canDecrement = $derived(!disabled && value > min);
	let inputText = $state(String(value));

	$effect(() => {
		inputText = String(value ?? 0);
	});

	$effect(() => {
		if (onchange) {
			onchange(value);
		}
	});

	function commitInput() {
		if (disabled) return;

		let parsed = Number.parseInt(inputText.trim(), 10);

		if (Number.isNaN(parsed)) {
			inputText = String(value);
			return;
		}

		if (parsed < min) parsed = min;
		if (parsed > effectiveMax) parsed = effectiveMax;

		value = parsed;
		inputText = String(parsed);
	}

	function handleInputKeydown(e: KeyboardEvent) {
		if (e.key === /* @wc-ignore */ "Enter") {
			e.preventDefault();
			(e.target as HTMLInputElement).blur();
			commitInput();
		} else if (e.key === /* @wc-ignore */ "Escape") {
			e.preventDefault();
			(e.target as HTMLInputElement).blur();
			inputText = String(value ?? 0);
		}
	}

	function handleIncrementClick() {
		commitInput();
		if (canIncrement) value++;
	}

	function handleDecrementClick() {
		commitInput();
		if (canDecrement) value--;
	}
</script>

<div class="qty-row">
	<div class="qty">
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
		<input
			in:scale={{ duration: 300 }}
			class="qty-input"
			inputmode="numeric"
			pattern="[0-9]*"
			bind:value={inputText}
			onblur={commitInput}
			onkeydown={handleInputKeydown}
			{disabled}
		/>
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
		gap: 0.1rem;
		background: lightgray;
		border-radius: 0.6rem;
	}

	.micro-btn {
		border: none;
		width: 1.6rem;
		height: 1.6rem;
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
		width: 2.1rem;
		text-align: center;
		border: none;
		background: transparent;
		font-size: 0.8rem;
		outline: none;
		padding: 0;
		color: black;
	}
</style>
