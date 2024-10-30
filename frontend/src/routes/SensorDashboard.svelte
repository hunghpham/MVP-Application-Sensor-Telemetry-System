<script>  
  import { onMount, onDestroy } from 'svelte';
  import Chart from "chart.js/auto";
  // Import the date adapter
  import 'chartjs-adapter-luxon';  

  let chart;
  let chartCanvas;
  let chart_label = '';

  // chart data for generating series
  let chart_data_temp = [];
  let chart_data_hum = [];
  let chart_data_co2 = [];
 
  let sensors = [];
  let historical_sensor_data = [];
  let sortBy = 'Serial';
  let sortOrder = 'asc';
  let error_fetch_all_sensor = '';
  let error_fetch_sensor_historical_data = '';
  let run_int_fetchAllSensors;  
  let selectedSensorSerial = '';    
  let ws;
  
  // Perform all required task when on mount (page loaded)
  onMount(async () => {    
    fetchAllSensors();  //Fetch sensor data if exist
    generateChart();    //Generate chart for selected sensor if done so
    
    //Run every 2 minutes
    run_int_fetchAllSensors = setInterval(()=> {handleRowClick(selectedSensorSerial)}, 120000);    

    //connect to the server using websocket
    ws = new WebSocket("ws://localhost:8080/ws");

    // Listen for messages from the server over websocket
		ws.onmessage = (event) => {
			
      const jsonObject = JSON.parse(event.data);
      console.log(jsonObject);  //for debugging just print the sensor data coming from the websocket
            
      //Got sensor data, fetch latest sensor reading from database and update the sensor array variable which cause Svelte reactive nature to update the html table automatically with the data 
      fetchAllSensors();

      //Find the sensor from the sensor array variable and check to see make sure the receive sensor data from the websocket is indeed a sensor currently exist   
      //const sensor = sensors.find(sensor => sensor.Serial === jsonObject.Serial);     
      
      //Update the chart with this new data from the websocket if the user choose to monitor this sensor graphically as well
      if (jsonObject.Serial === selectedSensorSerial) {
        //Create the new data point object for temperature
        const new_temp_data_point = {
          x: new Date(jsonObject.Timestamp), 
          y: jsonObject.Reading1
        };

        //Create the new data point object for humidity
        const new_hum_data_point = {
          x: new Date(jsonObject.Timestamp), 
          y: jsonObject.Reading2
        };

        //Create the new data point object for CO2
        const new_CO2_data_point = {
          x: new Date(jsonObject.Timestamp), 
          y: jsonObject.Reading3
        };

        updateChart(new_temp_data_point, new_hum_data_point, new_CO2_data_point);
      }

		};

		// Handle WebSocket errors
		ws.onerror = (event) => {
			console.log("WebSocket error:", event);
		};

  });

  // Clean up the interval on component destroy
  onDestroy(() => {
        clearInterval(run_int_fetchAllSensors);        
        if (chart) {
            chart.destroy();
        }
  });
    
  //Fetch all unique sensors into the array, this will cause a reactive update in svelte table, svelte reactive nature at works here...
  async function fetchAllSensors() {
    try {
      const response = await fetch('http://localhost:8080/api/get_all_sensors');
      if (!response.ok) {
        throw new Error("Failed to load data...");
      }
      sensors = await response.json();              
             
    } catch (err) {
      error_fetch_all_sensor = err.message;
    }
  }  

  //Update Chart with new data points
  function updateChart(new_temp_data_point, new_hum_data_point, new_CO2_data_point) {    
    chart.data.datasets[0].data.push(new_temp_data_point);
    chart.data.datasets[1].data.push(new_hum_data_point);
    chart.data.datasets[2].data.push(new_CO2_data_point);
    chart.update();
  }

  //Fetch historical data for a sensor
  async function fetchSensorHistoricalData(serial, start, end) {
    const base_url = 'http://localhost:8080/api/get_sensor_historical_data';
    const params = {
      serial_number: serial,
      start_dt: start,
      end_dt: end
    } 

    // Create a query string from the parameters
    const queryString = new URLSearchParams(params).toString();

    // Construct the full URL
    const urlWithParams = `${base_url}?${queryString}`;

    try {
      const response = await fetch(urlWithParams);
      if (!response.ok) {
        throw new Error("Failed to load data...");
      }
      historical_sensor_data = await response.json();  
      //console.log(historical_sensor_data);    

      //Got historical data, assign to   
      chart_data_temp = [];
      chart_data_hum = [];
      chart_data_co2 = [];

      //loop through response and assign to chart data
      historical_sensor_data.forEach(item => {
        chart_data_temp.push({time: new Date(item.Timestamp), temperature: item.Reading1})  
        chart_data_hum.push({time: new Date(item.Timestamp), humidity: item.Reading2})
        chart_data_co2.push({time: new Date(item.Timestamp), co2: item.Reading3})
      });

      //now generate the graph for this sensor data
      generateChart();

    } catch (err) {
      error_fetch_sensor_historical_data = err.message;
    }
  }

  //Function event row click
  function handleRowClick(serial) {
    selectedSensorSerial = serial;    
    chart_label = selectedSensorSerial;
    

    // Get current UTC time
    const currentDate = new Date().toISOString();

    // Get the past date by hour
    //const pastDate = new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString();

    // Get the pass date by minute
    const pastDate = new Date(Date.now() - 1 * 60 * 1000).toISOString();

    // console.log("Current UTC time: ", currentDate);
    // console.log("One day before UTC: ", pastDate);
   
    //get historical data based on date range
    fetchSensorHistoricalData(chart_label, currentDate, pastDate);       
  }

  function generateChart() {
    // Check if myChart already exists and destroy it
    if (chart) {
        chart.destroy();
    }

    // Create the chart when the component is mounted
    const ctx = chartCanvas.getContext("2d");   
    
    const series_data = {            
        datasets: [
          {
            label: 'Temperature (°C)',
            data: chart_data_temp.map(d => ({
                x: d.time,  // x-axis is DateTime
                y: d.temperature  // y-axis is temperature
            })),                
            borderColor: 'rgb(255, 99, 132)', // Red line for temperature
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            fill: true,
            tension: 0.1,
            yAxisID: 'y',
          }, 
          {
            label: 'Humidity (%)',
            data: chart_data_hum.map(d => ({
                x: d.time,  // x-axis is DateTime
                y: d.humidity  // y-axis is temperature
            })),                
            borderColor: 'rgb(54, 162, 235)', // Blue line for humidity
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            fill: true,
            tension: 0.1,
            yAxisID: 'y1',
          },
          {
            label: 'CO2 (%)',
            data: chart_data_co2.map(d => ({
                x: d.time,  // x-axis is DateTime
                y: d.co2  // y-axis is temperature
            })),                
            borderColor: 'rgb(255, 206, 86)', // Yellow for CO2
            backgroundColor: 'rgba(255, 206, 86, 0.2)',
            fill: true,
            tension: 0.1,
            yAxisID: 'y2',
          }          
      ]
    };

    const chart_config = {
        type: "line",
        data: series_data,
        options: {
            responsive: true,
            plugins: {
              title: {
                  display: true,  
                  text: chart_label + ' - Real-Time Graph',  // Set the chart title text
                  font: {
                      size: 20  // Adjust the font size of the title
                  },
                  color: 'white'
              },
              legend: {
                labels: {
                    color: 'white' // Legend text color
                },                
              },
              tooltip: {
                bodyColor: 'white', // Tooltip text color
                titleColor: 'white' // Tooltip title color
              }
            },
            scales: {                
                x: {
                    type: 'time',
                    time: {
                        unit: 'minute',  // Set the base unit to minute
                        stepSize: 5,  // Display a tick every 5 minutes
                        displayFormats: {
                          minute: 'MMM D, h:mm a'  // Format for each tick label
                        },                   
                        tooltipFormat: 'MMM D, h:mm a'
                    },
                    title: {
                        display: true,
                        text: 'Time',
                        color: 'white'
                    },
                    ticks: {
                        color: 'white' // X-axis text color
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Temperature (°C)',
                        color: 'white'
                    },
                    type: 'linear',
                    position: 'left',
                    ticks: {
                        color: 'white' // Y-axis text color
                    }
                },
                y1: {
                    title: {
                        display: true,
                        text: 'Humidity (%)',
                        color: 'white'
                    },
                    type: 'linear',
                    position: 'right',
                    ticks: {
                        color: 'white' // Y-axis text color
                    }
                },
                y2: {
                    title: {
                        display: true,
                        text: 'CO2 (%)',
                        color: 'white'
                    },
                    type: 'linear',
                    position: 'left',
                    ticks: {
                        color: 'white' // Y-axis text color
                    }
                }
            }
        }
    };


    //Chart.js setup and configuration for line chart
    chart = new Chart(ctx, chart_config);
  }
 
 

</script>
 
<h1 class="text-4xl font-bold text-gray mb-4 p-5">Sensor Dashboard</h1>  

<!-- <table class="w-full border-collapse mt-5 text-center dark:bg-gray-800"> -->
<table class="w-full border-collapse mt-5 text-center rounded-lg overflow-hidden border border-gray-300 dark:border-gray-600 shadow-2xl">
  <thead>
    <tr class="bg-gray-600 text-white cursor-pointer">      
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">Serial Number</th>
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">Timestamp</th>
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">Type</th>
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">Temperature</th>
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">Humidity</th>
      <th class="p-3 border-b border-gray-300 dark:border-gray-600 text-lg">CO2</th>
    </tr>
  </thead>
  <tbody>
    {#if sensors.length > 0}
      {#each sensors as sensor}
        <tr on:click={() => handleRowClick(sensor.Serial)} class="bg-white dark:bg-gray-900 hover:bg-gray-200 dark:hover:bg-gray-700 odd:bg-gray-50 even:bg-gray-200 dark:even:bg-gray-800">          
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Serial}</td>
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{new Date(sensor.Timestamp).toLocaleDateString()} {new Date(sensor.Timestamp).toLocaleTimeString()}</td>
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Type}</td>
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading1.toFixed(2)}°C</td>
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading2.toFixed(2)}%</td>
          <td class="p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading3.toFixed(2)}%</td>
        </tr>
      {/each}
    {:else}
      <p>Loading data...</p>
    {/if}

  </tbody>
</table>

<div class="p-10"></div>

<!-- A canvas so chart.js can draw on -->
<div class="bg-gray-600 text-white p-4 border border-gray-600 rounded-lg shadow-2xl">  
  <canvas bind:this={chartCanvas} class="max-w-full h-[500px]"></canvas> 
</div>


<style>  
  /* table {
    width: 100%;
    border-collapse: collapse;
    margin: 20px 0;    
    text-align: center;
  }

  th, td {
    padding: 12px;
    border-bottom: 1px solid #ddd;    
    width: 20%;    
  }

  th {
    font-size: 18px;    
  }

  td {
    font-size: 16px;    
    font-weight: 500;
  }  

  th {
    cursor: pointer;
    background-color: #6f7072;
    color: white;
  }

  th:hover {
    background-color: #ddd;
  } */

  
  /* tr:nth-child(even) {
    background-color: #f9f9f9;
  }

  tr:hover {
    background-color: #d4d4d4;
  } */

  /* canvas {
      max-width: 100%;
      height: 500px;
  } */
</style>
