<script>
	import Header from './Header.svelte';
	import '../app.css';
</script>

<!-- <div class="app"> -->
<div class="app">
	<Header />

	<main class="flex-1 flex flex-col p-4 w-full max-w-2xl mx-auto box-border">
		<slot />
	</main>

	<footer>
		<p>At <a href="https://mesalabs.com/">Mesa Labs</a> we measure what matters</p>
	</footer>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 64rem;
		margin: 0 auto;
		box-sizing: border-box;
	} 

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
	}

	footer a {
		font-weight: bold;
	}

	@media (min-width: 480px) {
		footer {
			padding: 12px 0;
		}
	}
</style>
