<script setup lang="ts">
import { Chart, ChartConfiguration, CategoryScale, LinearScale, LineController, PointElement, LineElement, Tooltip } from "chart.js";
import { onMounted, useTemplateRef } from "vue";



const canvas = useTemplateRef<HTMLCanvasElement>("canvas")

onMounted(() => {
    const color = '#948aa3'

    Chart.register(CategoryScale, LinearScale, LineController, PointElement, LineElement, Tooltip)
    Chart.defaults.font = {
        size: 18,
        family: '"Montserrat", sans-serif',
    }

    if (canvas.value == null) return

    const options: ChartConfiguration = {
        type: 'line',
        data: {
            labels: ['10.11.', '11.11.', '12.11.', '13.11.', '14.11.', '15.11.', '16.11.', '17.11.', '18.11.', '19.11.', '20.11.', '21.11.', 'Dnes'],
            datasets: [{
                label: 'Rychlost (CPM)',
                yAxisID: 'rychlost',
                data: [176, 192, 199, 202, 212, 232, 235, 220, 254, 199, 202, 212, 232],
                borderWidth: 6,
                tension: 0,
                borderColor: '#FFF6',
                pointBorderColor: 'white',
                pointBackgroundColor: 'white',
                pointRadius: 4,
            },
            {
                label: 'Přesnost',
                yAxisID: 'presnost',
                data: [84, 84, 84, 84, 100, 84, 78, 84, 80, 84, 84, 99, 98],
                borderWidth: 6,
                tension: 0,
                borderColor: '#86487999',
                pointBorderColor: '#864879',
                pointBackgroundColor: '#864879',
                pointRadius: 4,
            }]
        },
        options: {
            onHover: (_, activeElements) => {
                if (activeElements?.length > 0) {
                    canvas.value!.style.cursor = 'pointer'
                } else {
                    canvas.value!.style.cursor = 'auto'
                }
            },
            maintainAspectRatio: false,
            layout: {
                padding: 10
            },
            scales: {
                x: {
                    grid: {
                        display: false,
                        color: 'transparent'
                    },
                    ticks: {
                        color: color,
                        maxTicksLimit: 7,
                    }
                },
                rychlost: {
                    position: 'left',
                    grid: {
                        color: color,
                        lineWidth: 2,
                    },
                    ticks: {
                        count: 4,
                        precision: -1,
                        color: color,
                        padding: 6,
                    },
                    border: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: 'Rychlost (CPM)',
                        color: '#FFFD',
                        padding: 4,
                        font: {
                            weight: 500
                        }
                    }
                },
                presnost: {
                    position: 'right',
                    grid: {
                        color: 'transparent',
                        lineWidth: 2,
                        drawOnChartArea: false,
                    },
                    ticks: {
                        count: 4,
                        precision: 0,
                        color: color,
                        padding: 4,
                    },
                    border: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: 'Přesnost (%)',
                        color: '#b657a3',
                        padding: 4,
                        font: {
                            size: 21,
                            weight: 500
                        }
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    position: 'nearest',
                    displayColors: false,
                    titleAlign: 'center',
                    bodyAlign: 'center',
                    backgroundColor: '#000',
                    caretSize: 0,
                    caretPadding: 10,
                    titleFont: {
                        size: 16,
                    },
                    callbacks: {
                        label: function (tooltipItem) {
                            if (tooltipItem.datasetIndex == 0) return tooltipItem.formattedValue + ' CPM'
                            else return tooltipItem.formattedValue + ' %'
                        }
                    },
                    
                },
            },
        },
    }
    new Chart(canvas.value, options)
})

</script>
<template>
    <div>
        <canvas ref="canvas" />
    </div>
</template>
<style scoped>
div {
    background-color: #3f3351;
    border-radius: 10px;
    padding: 6px 2px 0px 2px;
}
</style>