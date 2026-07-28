<script lang="ts">
	import { getFields } from "./get-fields";

	let bingo = $state(false);

	const size = 5;
	const storageKey = "bingo-grid";

	function generateGrid(): { label: string; selected: boolean }[][] {
		return getFields(size * size).reduce<
			{ label: string; selected: boolean }[][]
		>((rows, item, i) => {
			if (i % size === 0) rows.push([]);
			rows[rows.length - 1].push({ label: item, selected: false });
			return rows;
		}, []);
	}

	function loadGrid(): { label: string; selected: boolean }[][] {
		if (typeof localStorage !== "undefined") {
			const stored = localStorage.getItem(storageKey);
			if (stored) {
				try {
					const parsed = JSON.parse(stored);
					// Basic sanity check that it matches the expected size
					if (
						Array.isArray(parsed) &&
						parsed.length === size &&
						parsed.every(
							(row: unknown) => Array.isArray(row) && row.length === size,
						)
					) {
						return parsed;
					}
				} catch {
					// ignore parse errors and fall back to generating a new grid
				}
			}
		}

		return generateGrid();
	}

	let grid = $state(loadGrid());

	$effect(() => {
		if (typeof localStorage !== "undefined") {
			localStorage.setItem(storageKey, JSON.stringify(grid));
		}
	});

	$effect(() => {
		let newBingo = false;
		let diagonal1 = 0;
		let diagonal2 = 0;
		let columns: boolean[][] = Array.from({ length: grid[0].length }, () =>
			Array(grid.length),
		);
		grid.forEach((row, x) => {
			let selected = 0;
			row.forEach((node, y) => {
				columns[y][x] = node.selected;
				if (node.selected) {
					selected++;
					if (x === y) {
						diagonal1++;
					}
					if (x + y === row.length - 1) {
						diagonal2++;
					}
				}
			});
			if (selected === row.length) {
				newBingo = true;
			}
		});

		if (diagonal1 === grid.length || diagonal2 === grid.length) {
			newBingo = true;
		}

		columns.forEach((column) => {
			let selected = 0;
			column.forEach((node) => {
				if (node) {
					selected++;
				}
			});
			if (selected === column.length) {
				newBingo = true;
			}
		});

		bingo = newBingo;
	});
</script>

<div class="grid-wrapper">
	<div class="grid" style="--size: {size};">
		{#each grid as row}
			{#each row as node}
				<button
					class:selected={node.selected}
					onclick={() => {
						node.selected = !node.selected;
					}}
				>
					<div>
						{node.label}
					</div>
				</button>
			{/each}
		{/each}
	</div>
</div>

<button class="generate" onclick={() => (grid = generateGrid())}>
	Generate new
</button>

{#if bingo}
	<h2>Bingo</h2>
{/if}

<style lang="scss">
	.grid-wrapper {
		width: 100dvw;
		height: 100dvh;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(
			var(--size),
			calc(
				(min(100dvw, 100dvh) - (0.4rem * (var(--size) - 1)) - 0.8rem) /
					var(--size)
			)
		);
		grid-template-rows: repeat(
			var(--size),
			calc(
				(min(100dvw, 100dvh) - (0.4rem * (var(--size) - 1)) - 0.8rem) /
					var(--size)
			)
		);
		gap: 0.4rem;
	}

	.grid button {
		border: none;
		border-radius: 1rem;
		padding: 0;
		margin: 0;

		div {
			font-size: 0.9rem;
			box-sizing: border-box;
			height: 100%;
			width: 100%;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			padding: 0.4rem;
		}

		&.selected {
			div {
				background-color: mediumseagreen;
			}
		}
	}

	.generate {
		position: fixed;
		top: 1rem;
		right: 1rem;
		padding: 0.6rem 1.2rem;
		font-size: 1rem;
		border: 2px solid black;
		border-radius: 0.5rem;
		background-color: white;
		cursor: pointer;

		&:hover {
			background-color: #f0f0f0;
		}
	}

	h2 {
		margin: 0;
		position: fixed;
		top: 50%;
		left: 0;
		width: 100%;
		height: 4rem;
		transform: translateY(-50%);
		background-color: mediumseagreen;
		border: 2px solid black;
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
