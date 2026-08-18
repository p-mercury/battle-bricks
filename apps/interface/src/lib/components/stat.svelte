<script lang="ts">
	let {
		value,
		label,
		color,
	}: {
		value: number;
		label: string;
		color?: "GRAY" | "RED" | "GREEN" | "BLUE";
	} = $props();
</script>

<div
	class:red={color === "RED"}
	class:green={color === "GREEN"}
	class:blue={color === "BLUE"}
>
	<span class="value">{value}</span>
	<span class="label">{label}</span>
</div>

<style lang="scss">
	div {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		border-radius: 0.4rem;
		border: 0.1rem solid #595d60;
		color: #595d60;
		gap: 0.2rem;

		&.red {
			border-color: #c91a09;
			color: #c91a09;
		}

		&.green {
			border-color: #237841;
			color: #237841;
		}

		&.blue {
			border-color: #0055bf;
			color: #0055bf;
		}
	}

	.value {
		font-size: 1.1rem;
		line-height: 1.1rem;
		font-weight: 700;
	}

	.label {
		font-size: 0.8rem;
		line-height: 0.8rem;
		font-weight: 600;
	}
</style>
