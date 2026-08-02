<script lang="ts">
	let {
		label = "",
		value = $bindable(false),
		disabled = false,
		successMessage,
		errorMessage,
		role = "switch",

		onclick = () => {},
	}: {
		label?: string;

		value?: boolean;

		disabled?: boolean | null;

		successMessage?: string | string[] | null;
		errorMessage?: string | string[] | null;

		role?: string;

		onclick?: ((e: MouseEvent) => void) | null;
	} = $props();

	const ariaChecked = $derived(value ? "true" : "false");

	let inputEl: HTMLInputElement | null = null;

	const uid = `switch-${Math.random().toString(36).slice(2)}`;

	export function focus() {
		inputEl?.focus();
	}
</script>

<div id="wrapper">
	{#if label}
		<div id="label"><label for={uid}>{label}</label></div>
	{/if}

	<label class="switch" aria-live="polite" {onclick}>
		<input
			id={uid}
			bind:this={inputEl}
			type="checkbox"
			bind:checked={value}
			{disabled}
			{role}
			aria-checked={ariaChecked}
		/>
		<span class="slider round {value ? 'on' : 'off'}" aria-hidden="true"></span>
	</label>

	{#if errorMessage}
		<div id="errorMessage">
			{Array.isArray(errorMessage) ? errorMessage.join("\n") : errorMessage}
		</div>
	{/if}

	{#if successMessage}
		<div id="successMessage">
			{Array.isArray(successMessage)
				? successMessage.join("\n")
				: successMessage}
		</div>
	{/if}
</div>

<style lang="scss">
	#wrapper {
		display: grid;
		gap: 0.4rem;
		align-items: center;
	}

	#label {
		color: black;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.25;

		label {
			cursor: pointer;
		}
	}

	.switch {
		position: relative;
		display: inline-block;
		width: 44px;
		height: 24px;
	}

	.switch input {
		position: absolute;
		inset: 0;
		margin: 0;
		opacity: 0;
		width: 100%;
		height: 100%;
	}

	.slider {
		position: absolute;
		inset: 0;
		cursor: pointer;
		background: lightgray;
		transition:
			background 0.2s ease,
			box-shadow 0.2s ease;
		border-radius: 999px;
	}

	.slider::before {
		content: "";
		position: absolute;
		height: 18px;
		width: 18px;
		left: 3px;
		top: 3px;
		background: white;
		border-radius: 50%;
		transition: transform 0.2s ease;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
	}

	.switch input:checked + .slider {
		background: green;
	}
	.switch input:checked + .slider::before {
		transform: translateX(20px);
	}

	.switch input:focus + .slider {
		outline: 2px solid transparent;
		box-shadow: 0 0 0 3px rgba(79, 140, 255, 0.35);
	}

	.switch input:disabled + .slider {
		opacity: 0.6;
		cursor: not-allowed;
	}

	/* Messages */
	#successMessage,
	#errorMessage {
		font-size: 0.875rem;
		white-space: pre-wrap;
		margin-top: 2px;
	}

	#successMessage {
		color: green;
	}
	#errorMessage {
		color: red;
	}
</style>
