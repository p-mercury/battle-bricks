<div></div>

<style lang="scss">
	div {
		display: block;
		width: 100%;
		height: 100%;

		&::after {
			content: "";
			display: block;
			width: 100%;
			height: 100%;
			padding: 0.15rem;
			box-sizing: border-box;

			background-image:
				linear-gradient(
					to bottom,
					rgb(255 255 255 / 0.18) 0 2px,
					transparent 2px calc(100% - 2px),
					rgb(0 0 0 / 0.22) calc(100% - 2px) 100%
				),
				url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='34.67' height='34.67'%3E%3Cdefs%3E%3ClinearGradient id='edge' x1='0' y1='0' x2='1' y2='1'%3E%3Cstop offset='0' stop-color='%23fff' stop-opacity='.46'/%3E%3Cstop offset='.46' stop-color='%23fff' stop-opacity='.16'/%3E%3Cstop offset='.54' stop-color='%23000' stop-opacity='.13'/%3E%3Cstop offset='1' stop-color='%23000' stop-opacity='.39'/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle cx='17.335' cy='17.335' r='10.98' fill='none' stroke='url(%23edge)' stroke-width='1.92'/%3E%3C/svg%3E");

			background-size:
				100% 100%,
				1.5rem 1.5rem;

			background-repeat: no-repeat, repeat;

			background-position:
				center,
				left top;

			background-origin: border-box, content-box;
			background-clip: border-box, content-box;
		}
	}
</style>
