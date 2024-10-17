<script>
  import { onMount, onDestroy } from 'svelte';
  import Chart from "chart.js/auto";
  // Import the date adapter
  import 'chartjs-adapter-luxon';

  let chart;
  let chartCanvas;
  let chart_label = '';

  // Example data: DateTime and Temperature
  let chart_data = [];
  let sensors = [];
    
  let historical_sensor_data = [];    
  let error_fetch_sensor_historical_data = '';
  let selectedSensor = '';  // For dropdown selection
  let error_fetch_all_sensor = '';
      
  
  // Fetch data when the component is mounted
  onMount(async () => {  
    fetchAllSensors();     
    generateChart();           
  });

  // Clean up the interval on component destroy
  onDestroy(() => {          
        if (chart) {
            chart.destroy();
        }
  });


  //Fetch all unique sensors 
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
  
  //Update Chart with new data point
  function updateChart(new_data_point) {
    chart.data.datasets[0].data.push(new_data_point);
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
      chart_data = [];

      //loop through response and assign to chart data
      historical_sensor_data.forEach(item => {
        chart_data.push({time: new Date(item.Timestamp), temperature: item.Reading1})  
      });

      //now generate the graph for this sensor data
      generateChart();

    } catch (err) {
      error_fetch_sensor_historical_data = err.message;
    }
  }

  

  function generateChart() {
    // Check if myChart already exists and destroy it
    if (chart) {
        chart.destroy();
    }

    // Create the chart when the component is mounted
    const ctx = chartCanvas.getContext("2d");    

    //Chart.js setup and configuration for line chart
    chart = new Chart(ctx, {
        type: "line",
        data: {
            //labels: data.map(d => d.time.toLocaleTimeString()), // Convert datetime to readable time
            datasets: [{
                label: chart_label,
                data: chart_data.map(d => ({
                    x: d.time,  // x-axis is DateTime
                    y: d.temperature  // y-axis is temperature
                })),                
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: false,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            plugins: {
              title: {
                  display: true,  
                  text: chart_label,  // Set the chart title text
                  font: {
                      size: 20  // Adjust the font size of the title
                  }
              }
            },
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'hour',  
                        stepSize: 5,  
                        displayFormats: {
                          minute: 'MMM D, h:mm a'  // Format for each tick label
                        },                   
                        tooltipFormat: 'MMM D, h:mm a'
                    },
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Temperature (°C)'
                    }
                }
            }
        }
    });
  }


  // Handle time range button clicks
  function setTimeRange(range) {
    console.log(`Time range set to: ${range}`);
    console.log(selectedSensor);

    chart_label = selectedSensor;

    // Get current UTC time
    const currentDate = new Date().toISOString();
    let pastDate;

    if (range == "1 hour") {            
      // Get the past date by hour
      pastDate = new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString();                          
    }
    else if (range == "1 day") {            
      // Get the past date by hour
      pastDate = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString();                          
    }
    else if (range == "1 week") {
      // Get the past date by hour
      pastDate = new Date(Date.now() - 168 * 60 * 60 * 1000).toISOString();        
    }
    else {
      // Get the past date by hour
      pastDate = new Date(Date.now() - 744 * 60 * 60 * 1000).toISOString();        
    }

    //get historical data based on date range
    fetchSensorHistoricalData(selectedSensor, currentDate, pastDate);    
  }
</script>
        
<!-- Dropdown and Buttons -->
<div>
  <label for="sensor-select">Select Sensor: </label>
  <select id="sensor-select" bind:value={selectedSensor}>
    <option value="" disabled selected>Select a sensor</option>
    {#each sensors as sensor}
      <option value={sensor.Serial}>{sensor.Serial}</option>
    {/each}
  </select>

  <button on:click={() => setTimeRange('1 hour')}>1 hour</button>
  <button on:click={() => setTimeRange('1 day')}>1 day</button>
  <button on:click={() => setTimeRange('1 week')}>1 week</button>
  <button on:click={() => setTimeRange('1 month')}>1 month</button>
</div>

<!-- A canvas so chart.js can draw on -->
<canvas bind:this={chartCanvas}></canvas>
  
<style>      
  canvas {
      max-width: 100%;
      height: 500px;
  }

  select, button {
  margin: 10px;
  padding: 10px;
  font-size: 16px;
  }

  button {
    background-color: #6f7072;
    color: white;
    border: none;
    cursor: pointer;
  }

  button:hover {
    background-color: #5a5b5d;
  }
</style>
  