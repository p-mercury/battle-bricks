<script lang="ts">
	let {
		value,
		size,
		red = false,
	}: { value: number; size: number; red?: boolean } = $props();

	let bars = $derived(Array.from({ length: size }, (_, i) => value > i));
</script>

<div>
	<ul>
		{#each bars as bar}
			<li>
				{#if bar}
					{#if red}
						<img src="/stud-red.png" alt="stud" />
					{:else}
						<img src="/stud-blue.png" alt="stud" />
					{/if}
				{:else}
					<img src="/stud-gray-light.png" alt="stud" />
				{/if}
			</li>
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
		overflow: hidden;
		gap: 0.2rem;
	}

	li {
		box-sizing: border-box;
		margin: 0;
		padding: 0.1rem 0 0 0;
		list-style: none;
	}

	img {
		height: 1.1rem;
		width: 1.1rem;
	}
</style>
