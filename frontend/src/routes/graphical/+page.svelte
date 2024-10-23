<script>		
    import Graphical from './Graphical.svelte';	
</script>

<svelte:head>
	<title>Historical Graph</title>	
</svelte:head>

<section>	
	<!-- load graphical svelte component -->
	<Graphical />	
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