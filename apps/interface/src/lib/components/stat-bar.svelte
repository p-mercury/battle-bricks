<script lang="ts">
	let { value, size }: { value: number; size: number } = $props();

	let bars = $derived(Array.from({ length: size }, (_, i) => value > i));
</script>

<div>
	<ul>
		{#each bars as bar}
			<li class:active={bar}></li>
		{/each}
	</ul>
</div>

<style lang="scss">
	div {
		display: flex;
	}

	ul {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
		display: flex;
		flex-direction: row;
		height: 1.2rem;
		border: 1px solid black;
	}

	li {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
		list-style: none;
		height: 100%;
		aspect-ratio: 1 / 1;
		background-color: red;
		border: 1px solid black;

		&.active {
			background-color: green;
		}
	}
</style>
