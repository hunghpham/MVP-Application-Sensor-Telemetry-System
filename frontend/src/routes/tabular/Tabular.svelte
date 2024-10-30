<script>
  import { onMount, onDestroy } from 'svelte';       
  
  let sensors = [];
  let historical_sensor_data = [];
  let sortBy = 'Serial';
  let sortOrder = 'asc';
  let selectedSensor = '';  // For dropdown selection
  let error_fetch_sensor_historical_data = '';

    
  // Fetch data when the component is mounted
  onMount(async () => {
    fetchAllSensors();     
  });

  // Clean up the interval on component destroy
  onDestroy(() => {                        
  });

  // Provide simple column sorting for the table
  function sort(column) {
    if (sortBy === column) {
      sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    } else {
      sortBy = column;
      sortOrder = 'asc';
    }

    historical_sensor_data = historical_sensor_data.slice().sort((a, b) => {
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
    } catch (err) {
      error_fetch_sensor_historical_data = err.message;
    }
  }

  // Handle time range button clicks
  function setTimeRange(range) {
    console.log(`Time range set to: ${range}`);
    console.log(selectedSensor);

    // Get current UTC time
    const currentDate = new Date().toISOString();
    let pastDate;

    if (range == "1 hour") {            
      // Get the past date by hour
      pastDate = new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString();                          
    }
    else if (range == "1 day") {            
      // Get the past date by day
      pastDate = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString();                          
    }
    else if (range == "1 week") {
      // Get the past date by week
      pastDate = new Date(Date.now() - 168 * 60 * 60 * 1000).toISOString();        
    }
    else {
      // Get the past date by month
      pastDate = new Date(Date.now() - 744 * 60 * 60 * 1000).toISOString();        
    }

    //get historical data based on date range
    fetchSensorHistoricalData(selectedSensor, currentDate, pastDate);    
  }
            
</script>
    
<h1 class="text-4xl font-bold text-gray mb-4 p-5">Historical Data</h1>  

<!-- Dropdown and Buttons -->
<div>
  <label for="sensor-select">Select Sensor: </label>
  <select id="sensor-select" bind:value={selectedSensor} class="bg-gray-600 text-white border border-gray-600 rounded-md cursor-pointer transition-colors duration-300 hover:bg-gray-500 dark:bg-gray-900 dark:border-gray-700 m-2 p-2 text-lg shadow-2xl">
    <option value="" disabled selected class="bg-gray-600 text-white">Select a sensor</option>
    {#each sensors as sensor}
      <option value={sensor.Serial} class="bg-gray-600 text-white">{sensor.Serial}</option>
    {/each}
  </select>

  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 hour')}>1 hour</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 day')}>1 day</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 week')}>1 week</button>
  <button class="bg-gray-600 text-white border-none cursor-pointer rounded-md transition-colors duration-300 hover:bg-gray-500 m-2 p-2 text-lg shadow-2xl" on:click={() => setTimeRange('1 month')}>1 month</button>
</div>


<!-- Table -->
<table class="w-full border-collapse mt-5 text-center rounded-lg overflow-hidden border border-gray-300 dark:border-gray-600 shadow-2xl">
  <thead>
    <tr class="bg-gray-600 text-white cursor-pointer">
      <th class="{sortBy === 'Serial' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Serial')}  >Serial</th>
      <th class="{sortBy === 'Timestamp' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Timestamp')}>Timestamp</th>
      <th class="{sortBy === 'Type' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Type')}>Type</th>
      <th class="{sortBy === 'Reading1' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Reading1')}>Temperature</th>
      <th class="{sortBy === 'Reading2' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Reading2')}>Humidity</th>
      <th class="{sortBy === 'Reading3' ? 'sort-' + sortOrder : ''} w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-lg" on:click={() => sort('Reading3')}>CO2</th>
    </tr>
  </thead>
  <tbody>
    {#if historical_sensor_data.length > 0}
      {#each historical_sensor_data as sensor}
        <tr class="bg-white dark:bg-gray-900 hover:bg-gray-100 dark:hover:bg-gray-700 odd:bg-gray-50 even:bg-gray-200 dark:even:bg-gray-800">
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Serial}</td>
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{new Date(sensor.Timestamp).toLocaleDateString()} {new Date(sensor.Timestamp).toLocaleTimeString()}</td>
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Type}</td>
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading1.toFixed(2)}°C</td>
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading2.toFixed(2)}%</td>
          <td class="w-1/6 p-3 border-b border-gray-300 dark:border-gray-600 text-base font-medium">{sensor.Reading3.toFixed(2)}%</td>
        </tr>
      {/each}
    {:else}
      <p>No data...</p>
    {/if}

  </tbody>
</table>
  

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
  }   */

  /* th {
    cursor: pointer;
    background-color: #6f7072;
    color: white;
  } */

  /* th:hover {
    background-color: #ddd;
  } */

  th.sort-asc::after {
    content: ' ▲';
  }

  th.sort-desc::after {
    content: ' ▼';
  }

  /* tr:nth-child(even) {
    background-color: #f9f9f9;
  }

  tr:hover {
    background-color: #d4d4d4;
  }       */

  /* select, button {
    margin: 10px;
    padding: 10px;
    font-size: 16px;
  }  */

  /* button {
    background-color: #6f7072;
    color: white;
    border: none;
    cursor: pointer;
  }

  button:hover {
    background-color: #5a5b5d;
  }
</style>
