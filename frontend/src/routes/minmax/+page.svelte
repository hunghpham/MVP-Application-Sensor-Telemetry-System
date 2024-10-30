<script>		
    import Tabular from './MinMax.svelte';	
</script>

<svelte:head>
	<title>Min Max Avg Report</title>	
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