<script setup lang="ts">
import { Chart, ChartConfiguration, CategoryScale, LinearScale, LineController, PointElement, LineElement, Tooltip } from "chart.js";
import { computed, onMounted, ref, useTemplateRef, watch } from "vue";

const props = defineProps({
    presnosti: {
        type: Array<number>,
        default: [NaN, NaN, NaN, NaN, NaN]
    },
    rychlosti: {
        type: Array<number>,
        default: [NaN, NaN, NaN, NaN, NaN]
    }
})

let chart: Chart
const canvas = useTemplateRef<HTMLCanvasElement>("canvas")
const neniCo = ref(false)

const dny = computed(() => {
    let delka = props.rychlosti.length == 0 ? 13 : props.rychlosti.length

    const arr: string[] = []
    for (let i = delka - 1; i > 0; i--) {
        const d = new Date()
        d.setDate(d.getDate() - i)
        arr.push(`${d.getDate()}.${d.getMonth() + 1}.`)
    }
    arr.push("Dnes")
    return arr
})

watch(props, function () {
    if (props.presnosti.every(v => isNaN(v))) {
        neniCo.value = true
        return
    }

    chart.data.labels = dny.value
    chart.data.datasets[0].data = props.rychlosti
    chart.data.datasets[1].data = props.presnosti

    chart.update()
})

onMounted(() => {
    const color = "#948aa3"

    Chart.register(CategoryScale, LinearScale, LineController, PointElement, LineElement, Tooltip)
    Chart.defaults.font = {
        size: 18,
        family: '"Montserrat", sans-serif',
    }

    if (canvas.value == null) return

    const options: ChartConfiguration = {
        type: "line",
        data: {
            labels: dny.value,
            datasets: [{
                label: "Rychlost (CPM)",
                yAxisID: "rychlost",
                data: props.presnosti,
                borderWidth: 6,
                tension: 0,
                borderColor: "#FFF6",
                pointBorderColor: "white",
                pointBackgroundColor: "white",
                pointRadius: 4,
                spanGaps: true,
            },
            {
                label: "Přesnost",
                yAxisID: "presnost",
                data: props.presnosti,
                borderWidth: 6,
                tension: 0,
                borderColor: "#86487999",
                pointBorderColor: "#864879",
                pointBackgroundColor: "#864879",
                pointRadius: 4,
                spanGaps: true
            }]
        },
        options: {
            onHover: (_, activeElements) => {
                if (activeElements?.length > 0) {
                    canvas.value!.style.cursor = "pointer"
                } else {
                    canvas.value!.style.cursor = "auto"
                }
            },
            locale: "cs-CZ",
            maintainAspectRatio: false,
            layout: {
                padding: {
                    left: 10,
                    right: 10,
                    top: 5,
                }
            },
            scales: {
                x: {
                    grid: {
                        display: false,
                        color: "transparent"
                    },
                    ticks: {
                        color: color,
                    }
                },
                rychlost: {
                    position: "left",
                    grid: {
                        color: color,
                        lineWidth: 2,
                    },
                    ticks: {
                        count: 4,
                        color: color,
                        padding: 6,
                        callback: function (value) {
                            if (typeof value == 'string') return value
                            return value.toFixed(0)
                        }
                    },
                    border: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: "Rychlost (CPM)",
                        color: "#FFFD",
                        padding: 4,
                        font: {
                            weight: 500
                        }
                    }
                },
                presnost: {
                    position: "right",
                    grid: {
                        color: "transparent",
                        lineWidth: 2,
                        drawOnChartArea: false,
                    },
                    ticks: {
                        count: 4,
                        color: color,
                        padding: 4,
                        callback: function (value) {
                            if (typeof value == 'string') return value
                            return value.toFixed(1)
                        }
                    },
                    border: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: "Přesnost (%)",
                        color: "#b657a3",
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
                    position: "nearest",
                    displayColors: false,
                    titleAlign: "center",
                    bodyAlign: "center",
                    backgroundColor: "#000",
                    caretSize: 0,
                    caretPadding: 10,
                    titleFont: {
                        size: 16,
                    },
                    callbacks: {
                        label: function (tooltipItem) {
                            if (tooltipItem.datasetIndex == 0) return tooltipItem.formattedValue + " CPM"
                            else return tooltipItem.formattedValue + " %"
                        }
                    },

                },
            },
        },
    }
    chart = new Chart(canvas.value, options)
})

</script>
<template>
    <div>
        <canvas ref="canvas" />
        <span v-if="neniCo">Zatím tu není co zobrazit.</span>
    </div>
</template>
<style scoped>
div {
    background-color: #3f3351;
    border-radius: 10px;
    padding: 12px 2px 0px 2px;
    position: relative;
}

span {
    position: absolute;
    top: 34%;
    left: 31%;
    font-size: 20px;
    background-color: var(--tmave-fialova);
    padding: 4px;
}
</style>