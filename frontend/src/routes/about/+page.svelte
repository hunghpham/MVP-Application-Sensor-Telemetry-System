<svelte:head>
	<title>Configure Kafka Topic</title>
	<meta name="description" content="About this app" />
</svelte:head>

<div class="text-column">	

	
	<h1 class="text-4xl font-bold text-gray mb-4 p-5">Discovering New Ideas and Innovations</h1>
	<p>This minimal viable product has introduced practical ideas and innovations that enhance and pave the way for new cutting-edge projects.</p>		
	<ul>
		<li><span>Rapid Product Development</span> - Allow developer to rapidly develop and deploy products, hence reduce time to market.</li>
		<li><span>Modern, Flexible and Efficient Architecture</span> - The application architecture illustrated in the diagram below supports a flexible development methodology, eliminates the need for building software from scratch, deploy data compliance with less overhead,  
			      promotes a modular design, remains hardware protocol agnostic, and facilitates the seamless integration of new software products into the ecosystem.
		</li>
		<li><span>Remove Cloud Based System Limitation</span> - Enable us to update our firmware Over the Air (OTA) and via the internet, which is currently not possible for our cloud-based systems. Enable us to reduce on-site calibration duration by an additional 50%. 
			      Furthermore, it opens up opportunities for innovative solutions like "Remote Calibration," eliminating the need for technicians to be physically present. All of this can be achieved using a real-time message broker.	
		</li>
		<li><span>Real-time Data Acquistion</span> - Enable the CM product to manage larger data ingestion, allowing it to process significantly more sensor data without compromising system performance. 
			      This results in faster data retrieval, leading to quicker report generation and graph creation for our clients, particularly those with a high volume of sensors.
		</li>
		<li><span>Enable True Real-time System</span> - Currently, our CM systems do not operate in true real-time. By utilizing WebSocket and JavaScript frameworks like Svelte, we can effortlessly display data immediately as it arrives.</li>
	</ul>
	
</div>

<h1>Application Design Architecture</h1>	
<div class="container">
	<img src="app_flow.png" alt="Placeholder Image">
</div>


<style>
	.container {
		display: flex;
		justify-content: center; /* Center horizontally */
		align-items: center; /* Center vertically */
		
		border: 1px solid #ccc; /* Optional border for visualization */
	}

	img {
		max-width: 100%; /* Ensure image is responsive */
		height: auto; /* Maintain aspect ratio */
	}
</style>