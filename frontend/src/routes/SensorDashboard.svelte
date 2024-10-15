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
  let sortBy = 'Serial';
  let sortOrder = 'asc';
  let error_fetch_all_sensor = '';
  let error_fetch_sensor_historical_data = '';
  let run_int_fetchAllSensors;  
  let selectedSensorSerial = '';
  let selectedSensorTimeStamp = '';
  
  // Fetch data when the component is mounted
  onMount(async () => {
    fetchAllSensors();
    generateChart();
    
    //Run every 5 seconds
    run_int_fetchAllSensors = setInterval(fetchAllSensors, 2000);    
  });

  // Clean up the interval on component destroy
  onDestroy(() => {
        clearInterval(run_int_fetchAllSensors);        
        if (chart) {
            chart.destroy();
        }
  });

  // Provide simple column sorting for the table
  function sort(column) {
    if (sortBy === column) {
      sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    } else {
      sortBy = column;
      sortOrder = 'asc';
    }

    sensors = sensors.slice().sort((a, b) => {
      let aValue = a[sortBy];
      let bValue = b[sortBy];

      if (typeof aValue === 'string') {
        return sortOrder === 'asc' ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
      }

      return sortOrder === 'asc' ? aValue - bValue : bValue - aValue;
    });
  }

  //Fetch all unique sensors 
  async function fetchAllSensors() {
    try {
      const response = await fetch('http://localhost:8080/api/get_all_sensors');
      if (!response.ok) {
        throw new Error("Failed to load data...");
      }
      sensors = await response.json();  
      
      //Get the sensor currently displaying its data in the chart
      let sensor_in_chart = sensors.find(obj => obj.Serial == selectedSensorSerial);

      //Create the new data point object
      const new_data_point = {
        x: new Date(sensor_in_chart.Timestamp), 
        y: sensor_in_chart.Reading1
      };
      
      const dataLength = chart.data.datasets[0].data.length;
      const latest_date_value_on_chart = new Date(chart.data.datasets[0].data[dataLength - 1].x);      

      console.log(latest_date_value_on_chart);
      console.log(new_data_point.x);

      if (new_data_point.x > latest_date_value_on_chart) {        
        updateChart(new_data_point);
      }
             
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

  //Function event row click
  function handleRowClick(row) {
    selectedSensorSerial = row.Serial;    
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
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            plugins: {
              title: {
                  display: true,  
                  text: chart_label + ' - Real-Time Graph',  // Set the chart title text
                  font: {
                      size: 20  // Adjust the font size of the title
                  }
              }
            },
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'minute',  // Set the base unit to minute
                        //stepSize: 5,  // Display a tick every 5 minutes
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
  
</script>
  
  
<table>
  <thead>
    <tr>
      <th>Serial Number</th>
      <th>Timestamp</th>
      <th>Type</th>
      <th>Reading 1</th>
      <!-- <th>Reading 2</th> -->
    </tr>
  </thead>
  <tbody>
    {#if sensors.length > 0}
      {#each sensors as sensor}
        <tr on:click={() => handleRowClick(sensor)}>
          <td>{sensor.Serial}</td>
          <td>{new Date(sensor.Timestamp).toLocaleDateString()} {new Date(sensor.Timestamp).toLocaleTimeString()}</td>
          <td>{sensor.Type}</td>
          <td>{sensor.Reading1.toFixed(2)}°C</td>
          <!-- <td>{sensor.Reading2.toFixed(2)}%</td> -->
        </tr>
      {/each}
    {:else}
      <p>Loading data...</p>
    {/if}

  </tbody>
</table>

<!-- A canvas so chart.js can draw on -->
<canvas bind:this={chartCanvas}></canvas>

<style>
  table {
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
  }

  th.sort-asc::after {
    content: ' ▲';
  }

  th.sort-desc::after {
    content: ' ▼';
  }

  tr:nth-child(even) {
    background-color: #f9f9f9;
  }

  tr:hover {
    background-color: #d4d4d4;
  }

  canvas {
      max-width: 100%;
      height: 500px;
  }
</style>
