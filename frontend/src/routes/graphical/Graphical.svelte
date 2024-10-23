<script>
  import { onMount, onDestroy } from 'svelte';
  //import Chart from "chart.js/auto";
  // Import the date adapter
  import 'chartjs-adapter-luxon';  
  import { Chart, LineController, LineElement, PointElement, LinearScale, Title, Tooltip, Legend, CategoryScale } from 'chart.js';    
  
  let chart;
  let chartCanvas;
  let chart_label = '';

  // chart data for generating series
  let chart_data_temp = [];
  let chart_data_hum = [];
  let chart_data_co2 = [];

  let sensors = [];
    
  let historical_sensor_data = [];    
  let error_fetch_sensor_historical_data = '';
  let selectedSensor = '';  // For dropdown selection
  let error_fetch_all_sensor = '';
     
  
  // Fetch data when the component is mounted
  onMount(async () => {        
    //Load chart js zoom plugins and register when the browser load because we are running on client side
    //because hammer.js expects the window object and this is only in a browser environment, running in server side development like in docker, there is not window
    const zoomPlugin = (await import('chartjs-plugin-zoom')).default;
    const Hammer = (await import('hammerjs')).default;

    Chart.register(LineController, LineElement, PointElement, LinearScale, Title, Tooltip, Legend, CategoryScale, zoomPlugin);
    

    if (typeof window !== 'undefined') {
      // This code runs only in the browser
      console.log('Running on the client side');
    } else {
      // This code runs on the server side (SSR)
      console.log('Running on the server side');
    }
    
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
            borderColor: 'rgb(255, 206, 86)', // Green for CO2
            backgroundColor: 'rgba(255, 206, 86, 0.2)',
            fill: true,
            tension: 0.1,
            yAxisID: 'y2',
          }          
      ]
    };    


    const zoomOptions = {
      pan: {
        enabled: true,
        mode: "x"
      },
      zoom: {
        wheel: {
          enabled: true
        },
        pinch: {
          enabled: true
        },
        mode: "x"
      }
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
                },
                zoom: zoomOptions                
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
            },                     
        }        
    };


    //Chart.js setup and configuration for line chart
    chart = new Chart(ctx, chart_config);
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

<h1 class="text-4xl font-bold text-gray mb-4 p-5">Historical Graph</h1>  

<!-- Dropdown and Buttons -->
<div>
  <label for="sensor-select">Select Sensor: </label>
  <select id="sensor-select" bind:value={selectedSensor}  class="bg-gray-600 text-white border border-gray-600 rounded-md p-2 cursor-pointer transition-colors duration-300 hover:bg-gray-700 dark:bg-gray-900 dark:border-gray-700 m-2 p-2 text-lg shadow-2xl">
    <option value="" disabled selected class="bg-gray-600 text-white">Select a sensor</option>
    {#each sensors as sensor}
      <option class="bg-gray-600 text-white" value={sensor.Serial}>{sensor.Serial}</option>
    {/each}
  </select>

  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 hour')}>1 hour</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 day')}>1 day</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 week')}>1 week</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 month')}>1 month</button>
</div>

<!-- A canvas so chart.js can draw on -->
<div class="bg-gray-600 text-white p-4 border border-gray-600 rounded-lg shadow-2xl">  
  <canvas bind:this={chartCanvas} class="max-w-full h-[500px]"></canvas>
</div>
  
<style>      
  /* canvas {
      max-width: 100%;
      height: 500px;
  } */

  /* select, button {
  margin: 10px;
  padding: 10px;
  font-size: 16px;
  } */

  /* button {
    background-color: #6f7072;
    color: white;
    border: none;
    cursor: pointer;
  }

  button:hover {
    background-color: #5a5b5d;
  } */
</style>
  