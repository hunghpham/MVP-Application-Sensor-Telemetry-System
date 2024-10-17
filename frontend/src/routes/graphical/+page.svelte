<script>	
	//import welcome from '$lib/images/svelte-welcome.webp';
	//import welcome_fallback from '$lib/images/svelte-welcome.png';
    import Graphical from './Graphical.svelte';
	//import ChartJS from './ChartJS.svelte';
</script>

<svelte:head>
	<title>Historical Graph</title>	
</svelte:head>

<section>
	<h1>	
		Historical Graph
	</h1>
	
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