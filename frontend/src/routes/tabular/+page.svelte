<script>		
    import Tabular from './Tabular.svelte';	
</script>

<svelte:head>
	<title>Historical Graph</title>	
</svelte:head>

<section>		
	<!-- load graphical svelte component -->
	<Tabular />	
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}	
</style>