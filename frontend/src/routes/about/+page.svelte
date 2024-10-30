<svelte:head>
	<title>Configure Kafka Topic</title>
	<meta name="description" content="About this app" />
</svelte:head>

<div class="text-column">	

	
	<!-- <h1 class="text-4xl font-bold text-gray mb-4 p-5">Discovering New Ideas and Innovations</h1>
	<p>This minimal viable product has introduced practical ideas and innovations that enhance and pave the way for new cutting-edge projects.</p>		
	<ul>
		<li><span>Rapid Product Development</span> - Allow developer to rapidly develop and deploy products, hence reduce time to market.</li>
		<li><span>Modern, Flexible and Efficient Architecture</span> - The application architecture illustrated in the diagram below supports a flexible development methodology, eliminates the need for building software from scratch, deploy data compliance with less overhead,  
			      promotes a modular design, remains hardware protocol agnostic, and facilitates the seamless integration of new software products into the ecosystem.
		</li>
		<li><span>Remove Cloud Based System Limitation</span> - Enable us to update our firmware Over the Air (OTA) and via the internet, which is currently not possible for our cloud-based systems. Enable us to reduce on-site calibration duration by an additional 30%. 
			      Furthermore, it opens up opportunities eliminating Twilio SMS cost by leveraging Svelte with Capacitor.js to develop a mobile application so push notification can use instead of SMS.	
		</li>
		<li><span>Enhance MFA security</span> - In addition to the mobile application free push notification instead of SMS, it can also provide geo-location of the user to enhance MFA security.
		</li>
		<li><span>Real-time Data Acquistion</span> - Enable the CM product to manage larger data ingestion, allowing it to process significantly more sensor data without compromising system performance. 
			      This results in faster data retrieval, leading to quicker report generation and graph creation for our clients, particularly those with a high volume of sensors.
		</li>
		<li><span>Enable True Real-time System</span> - Currently, our CM systems do not operate in true real-time. By utilizing WebSocket and JavaScript frameworks like Svelte, we can effortlessly display data immediately as it arrives.</li>
		</li>
		<li><span>In Field Sensor Status Update</span> - Currently, hardware engineering required faulty sensors to be send back for diagostic, with message broker for IoT devices, sensors can update their statuses live while in the field, allowing
			      software engineer to develop remote application to ingest these data for remote investigation and troubleshooting without the need to send the sensors back to Mesa.
		</li>
	</ul> -->

	<h1 class="text-4xl font-bold text-gray-800 mb-4 p-5">Discovering New Ideas and Innovations</h1>
	<p class="font-semibold text-gray-700">This minimal viable product introduces practical ideas and innovations that enhance and pave the way for cutting-edge projects.</p>		

	<ul class="list-disc pl-8 space-y-4">
		<li class="font-semibold text-gray-700">
			<span>Rapid Product Development:</span> 
			<span class="text-gray-600">Allows developers to rapidly develop and deploy products, significantly reducing time to market.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>Modern, Flexible, and Efficient Architecture:</span> 
			<span class="text-gray-600">Supports a flexible development methodology, eliminates the need for building software from scratch, ensures data compliance with less overhead, promotes a modular design, remains hardware protocol agnostic, and facilitates seamless integration of new software products into the ecosystem.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>Removing Cloud-Based System Limitations:</span> 
			<span class="text-gray-600">Enables firmware updates over the air (OTA) and via the internet, reducing on-site calibration duration by an additional 30%. This also eliminates Twilio SMS costs by leveraging Svelte with Capacitor.js to develop a mobile application that uses push notifications instead of SMS.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>Enhanced MFA Security:</span> 
			<span class="text-gray-600">In addition to free mobile push notifications, the application can provide user geo-location to enhance MFA security.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>Real-time Data Acquisition:</span> 
			<span class="text-gray-600">Enables the CM product to manage larger data ingestion, allowing it to process significantly more sensor data without compromising system performance. This results in faster data retrieval, leading to quicker report generation and graph creation for clients, especially those with a high volume of sensors.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>True Real-time System Capabilities:</span> 
			<span class="text-gray-600">Currently, our CM systems do not operate in true real-time. By utilizing WebSocket and JavaScript frameworks like Svelte, we can effortlessly display data immediately as it arrives.</span>
		</li>
		
		<li class="font-semibold text-gray-700">
			<span>In-Field Sensor Status Updates:</span> 
			<span class="text-gray-600">Instead of sending faulty sensors back for diagnostics, our message broker for IoT devices allows sensors to update their statuses live while in the field. This enables software engineers to develop remote applications that ingest this data for troubleshooting without needing to send sensors back to Mesa.</span>
		</li>
	</ul>

	
</div>

<h1>Application Architecture</h1>	
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