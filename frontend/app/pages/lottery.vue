<script setup lang="ts">
import { Block } from 'notiflix/build/notiflix-aio';
const isMobile = useState<boolean>('isMobile')
const lotteryData = ref<any>(null)
const lotteryNumber = ref('')
const thaiMonths = [
  '',
  'มกราคม',
  'กุมภาพันธ์',
  'มีนาคม',
  'เมษายน',
  'พฤษภาคม',
  'มิถุนายน',
  'กรกฎาคม',
  'สิงหาคม',
  'กันยายน',
  'ตุลาคม',
  'พฤศจิกายน',
  'ธันวาคม'
]
const typeReward = ref<Array<any>>([
    { key: 'first', label: 'รางวัลที่ 1', money: '6,000,000' },
    { key: 'second', label: 'รางวัลที่ 2', money: '200,000' },
    { key: 'third', label: 'รางวัลที่ 3', money: '80,000' },
    { key: 'fourth', label: 'รางวัลที่ 4', money: '40,000' },
    { key: 'fifth', label: 'รางวัลที่ 5', money: '20,000' },
    { key: 'last2', label: 'รางวัลเลขท้าย 2 ตัว', money: '2,000' },
    { key: 'last3f', label: 'รางวัลเลขท้าย 3 ตัวหน้า', money: '4,000' },
    { key: 'last3b', label: 'รางวัลเลขท้าย 3 ตัวหลัง', money: '4,000' }
])
const arrTitle = ref<string[]>([])
const { data: dataFetch, pending, error } = await useFetch<any>('/server/lottery', { lazy: true })
onMounted(() => {
  watch(pending, (isPending) => {
    if (isPending) {
        Block.pulse('#loading', '', {
            backgroundColor: 'transparent',
        })
    } else {
      Block.remove('#loading')
    }
  }, { immediate: true })
})

onMounted(() => {
  watch(dataFetch, (dataFetch) => {
    if (dataFetch?.response.data === undefined) {
        lotteryData.value = []
    } else {
        lotteryData.value = dataFetch?.response.data
        // console.log(lotteryData.value)
    }
  }, { immediate: true })
})

const filteredPosts = computed(() => {
    
    const keyword = lotteryNumber.value.trim()
    if (!keyword) {
        return lotteryData.value
    }
    arrTitle.value = []
    const result: Record<string, any> = {}
    for (const [key, data] of Object.entries<any>(lotteryData.value)) {
        const hasMatch = data.number?.some((num: any) => num.value.includes(keyword))
        if (hasMatch) {
            result[key] = { ...data }
            arrTitle.value.push(key)
        }
    }
    arrTitle.value.filter(num => {
        result[num].number = result[num].number.filter((n: any) => n.value.includes(keyword))
    })
    return result
})
</script>

<template>
  <div id="loading" class="border border-transparent">
    
    <div v-if="dataFetch?.response.date" class="flex justify-center flex-col gap-2 items-center w-full h-full border ">
        <div class=" max-w-[400px] flex justify-center" :class="isMobile ? 'w-full' : 'w-[50%]' "> 
            <input type="text" maxlength="6" @input="lotteryNumber = lotteryNumber.replace(/\D/g, '')" placeholder="ค้นหาเลข ลอตเตอรี่" v-model="lotteryNumber" 
                class="w-[100%] border border-black p-2 rounded-lg text-xl" :class="lotteryNumber.length > 0 ? 'tracking-[2em] indent-[2em] font-bold text-center' : ''"/>
        </div>
        <h1 class="text-center text-4xl font-bold my-4">งวด ณ วันที่  {{dataFetch?.response.displayDate?.date}} {{ thaiMonths[Number(dataFetch?.response.displayDate?.month)] }} {{ Number(dataFetch?.response.displayDate?.year) + 543 }}</h1>
        <div v-for="(item, index) in filteredPosts" :key="index" class="w-full max-w-[800px] p-3 bg-white rounded-lg border border-gray-300">
            {{ typeReward.find(reward => reward.key === index)?.label }} เงินรางวัล {{ typeReward.find(reward => reward.key === index)?.money }}
            <div v-if="index === 'first'" class="flex justify-center font-bold text-2xl w-full p-4 ">
                {{ item.number[0].value }}
            </div>
            <div v-else-if="index === 'second'" class="grid grid-cols-4 text-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'third'" class="grid grid-cols-4 text-center font-bold text-2xl font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'fourth'" class="grid grid-cols-4 text-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'fifth'" class="grid grid-cols-4 text-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'last2'" class="flex justify-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'last3f'" class="flex justify-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else-if="index === 'last3b'" class="flex justify-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
            <div v-else class="flex justify-center font-bold text-2xl">
                <span v-for="num in item.number" :key="num.round" style="margin-right: 8px;">
                    {{ num.value }}
                </span>
            </div>
        </div>
    </div>
  </div>
</template>