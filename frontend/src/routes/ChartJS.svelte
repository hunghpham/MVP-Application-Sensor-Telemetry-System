<script>
    import { onMount } from "svelte";
    import Chart from "chart.js/auto";
    // Import the date adapter
    import 'chartjs-adapter-luxon';


    let chart;
    let chartCanvas;

    // Example data: DateTime and Temperature
    let data = [
        { time: new Date('2024-10-10T10:00:00'), temperature: 20 },
        { time: new Date('2024-10-10T11:00:00'), temperature: 22 },
        { time: new Date('2024-10-10T12:00:00'), temperature: 24 },
        { time: new Date('2024-10-10T13:00:00'), temperature: 23 },
        { time: new Date('2024-10-10T14:00:00'), temperature: 25 }
    ];

    onMount(() => {
        // Create the chart when the component is mounted
        const ctx = chartCanvas.getContext("2d");

        //Chart.js setup and configuration for line chart
        chart = new Chart(ctx, {
            type: "line",
            data: {
                //labels: data.map(d => d.time.toLocaleTimeString()), // Convert datetime to readable time
                datasets: [{
                    label: 'Temperature (°C)',
                    data: data.map(d => ({
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
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'hour', // You can change to 'day', 'minute', etc.
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
    });

    // Destroy chart when component is destroyed to prevent memory leaks
    function onDestroy() {
        if (chart) {
            chart.destroy();
        }
    }
</script>

<!-- A canvas so chart.js can draw on -->
<canvas bind:this={chartCanvas}></canvas>

<style>
    canvas {
        max-width: 100%;
        height: 400px;
    }
</style>
