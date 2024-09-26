<script>
  import { onMount } from 'svelte';
  let socket;
  let messages = [];
  let chartData = [];

  onMount(() => {
    socket = new WebSocket('ws://localhost:8080/ws');
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      messages = [...messages, data];

      // Dynamically update chart data
      chartData = [...chartData, data.value];
    };
  });
</script>

<div class="container mx-auto">
  <h1 class="text-3xl font-bold text-center mt-5">Sensor Data Dashboard</h1>

  <!-- Table of readings -->
  <table class="table-auto w-full mt-5">
    <thead>
      <tr>
        <th class="px-4 py-2">Timestamp</th>
        <th class="px-4 py-2">Value</th>
      </tr>
    </thead>
    <tbody>
      {#each messages as message}
        <tr>
          <td class="border px-4 py-2">{message.timestamp}</td>
          <td class="border px-4 py-2">{message.value}</td>
        </tr>
      {/each}
    </tbody>
  </table>

  <!-- Graph -->
  <div class="mt-5">
    <!-- Example: You could integrate a chart library here -->
    <h2 class="text-2xl">Graph of Values</h2>
    <pre>{JSON.stringify(chartData, null, 2)}</pre>
  </div>
</div>

<style>
  /* TailwindCSS is included via postcss and purgecss */
</style>
