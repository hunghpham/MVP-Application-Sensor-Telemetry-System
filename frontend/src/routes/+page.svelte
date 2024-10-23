<script>		
    import SensorDashboard from './SensorDashboard.svelte';	
</script>

<svelte:head>
	<title>Dashboard</title>	
</svelte:head>

<section>		
	<!-- load dashboard svelte component -->
	<SensorDashboard />	
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
