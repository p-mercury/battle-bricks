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
		overflow: visible;
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
		filter: drop-shadow(0 0.16rem 0.1rem rgba(0, 0, 0, 0.55))
			drop-shadow(0 0.04rem 0.03rem rgba(0, 0, 0, 0.3))
			drop-shadow(0 0.02rem 0.02rem rgba(255, 255, 255, 0.65));
	}
</style>
