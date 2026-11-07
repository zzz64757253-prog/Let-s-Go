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
    { key: 'near1', label: 'รางวัลข้างเคียงรางวัลที่ 1', money: '100,000' },
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
  <div id="loading" class="min-h-dvh bg-gray-50 text-gray-700">
    <main class="w-full max-w-[1000px] mx-auto px-3 py-3 sm:px-4 sm:py-6">
      <div v-if="dataFetch?.response.date" class="flex flex-col gap-3 sm:gap-4">

        <!-- ส่วนหัว + ค้นหา: จบในกล่องเดียว เห็นช่องค้นหาทันทีบนมือถือ -->
        <section class="bg-white border border-gray-200 rounded-md p-3 sm:p-5">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 sm:w-10 sm:h-10 shrink-0 bg-blue-50 rounded-md flex items-center justify-center">
              <Icon name="heroicons:ticket-20-solid" class="w-5 h-5 text-blue-400" />
            </div>
            <div class="min-w-0 flex-1">
              <h1 class="text-base sm:text-xl font-bold text-gray-800 leading-snug">
                ตรวจผลสลากกินแบ่งรัฐบาล
              </h1>
              <p class="text-xs sm:text-sm text-gray-500 mt-0.5">
                งวดวันที่ {{ dataFetch?.response.displayDate?.date }}
                {{ thaiMonths[Number(dataFetch?.response.displayDate?.month)] }}
                {{ Number(dataFetch?.response.displayDate?.year) + 543 }}
              </p>
            </div>
          </div>

          <div class="relative mt-3 sm:mt-4">
            <Icon
              name="heroicons:magnifying-glass-20-solid"
              class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400 pointer-events-none"
            />
            <input
              id="lottery-search"
              v-model="lotteryNumber"
              type="text"
              inputmode="numeric"
              maxlength="6"
              aria-label="ค้นหาหมายเลขสลาก"
              @input="lotteryNumber = lotteryNumber.replace(/\D/g, '')"
              placeholder="กรอกเลขสลากที่ต้องการค้นหา"
              class="w-full h-11 sm:h-12 pl-10 pr-12 bg-gray-50 border border-gray-200 rounded-md text-base text-gray-800 placeholder:text-gray-400 outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100 transition-all"
              :class="lotteryNumber.length > 0 ? 'font-semibold tracking-widest' : ''"
            />
            <button
              v-if="lotteryNumber"
              type="button"
              aria-label="ล้างหมายเลขที่ค้นหา"
              class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-gray-400 hover:text-gray-600"
              @click="lotteryNumber = ''"
            >
              <Icon name="heroicons:x-mark-20-solid" class="w-5 h-5" />
            </button>
          </div>
        </section>

        <!-- หัวข้อผลรางวัลให้กินพื้นที่เพียงบรรทัดเดียว -->
        <div class="flex items-center justify-between gap-3 px-1">
          <h2 class="text-sm sm:text-base font-bold text-gray-800">ผลการออกรางวัล</h2>
          <span class="text-xs text-blue-500">
            {{ lotteryNumber ? 'ผลการค้นหา' : 'ทุกประเภทรางวัล' }}
          </span>
        </div>

        <!-- รายการรางวัล -->
        <section
          v-for="(item, index) in filteredPosts"
          :key="index"
          class="bg-white border border-gray-300 rounded-md overflow-hidden"
        >
          <div class="p-4 sm:p-5 flex flex-wrap items-center justify-between gap-3 border-b border-gray-200">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 shrink-0 bg-blue-50 rounded-md flex items-center justify-center">
                <Icon name="heroicons:trophy-20-solid" class="w-5 h-5 text-blue-400" />
              </div>
              <div>
                <h3 class="text-sm sm:text-base font-semibold text-gray-800">
                  {{ typeReward.find(reward => reward.key === index)?.label }}
                </h3>
                <p class="text-xs text-gray-400 mt-1">ผลรางวัลประจำงวด</p>
              </div>
            </div>
            <div class="rounded-md bg-gray-100 px-3 py-2 text-right">
              <p class="text-[10px] text-gray-400">เงินรางวัล</p>
              <p class="text-sm font-bold text-blue-500">
                {{ typeReward.find(reward => reward.key === index)?.money }} บาท
              </p>
            </div>
          </div>

          <div class="p-4 sm:p-6">
            <!-- รางวัลที่ 1 -->
            <div
              v-if="index === 'first'"
              class="bg-blue-50 border border-blue-100 rounded-md py-7 px-3 text-center"
            >
              <p class="text-xs text-blue-400 mb-3">หมายเลขที่ถูกรางวัล</p>
              <p class="font-bold text-4xl sm:text-5xl tracking-[0.12em] text-blue-500 break-all">
                {{ item.number[0].value }}
              </p>
            </div>

            <!-- รางวัลอื่น ๆ -->
            <div
              v-else
              class="grid gap-2.5 sm:gap-3"
              :class="
                index === 'last2'
                  ? 'grid-cols-1 max-w-[260px] mx-auto'
                  : index === 'last3f' || index === 'last3b'
                    ? 'grid-cols-2 max-w-[440px] mx-auto'
                    : 'grid-cols-2 sm:grid-cols-4'
              "
            >
              <div
                v-for="num in item.number"
                :key="num.round"
                class="min-w-0 bg-gray-50 border border-gray-200 rounded-md px-2 py-3 sm:py-4 text-center font-semibold text-base sm:text-xl text-gray-700 tracking-wider break-all"
                :class="index === 'last2' ? 'text-3xl sm:text-4xl text-blue-500 font-bold py-6' : ''"
              >
                {{ num.value }}
              </div>
            </div>
          </div>
        </section>

        <!-- ไม่พบผลลัพธ์ -->
        <section
          v-if="lotteryNumber && Object.keys(filteredPosts || {}).length === 0"
          class="bg-white border border-gray-300 rounded-md p-10 text-center"
        >
          <div class="w-12 h-12 bg-gray-100 rounded-md flex items-center justify-center mx-auto mb-3">
            <Icon name="heroicons:magnifying-glass-20-solid" class="w-6 h-6 text-gray-400" />
          </div>
          <h3 class="font-semibold text-gray-800">ไม่พบหมายเลขที่ค้นหา</h3>
          <p class="text-sm text-gray-400 mt-2">ลองเปลี่ยนหมายเลข หรือกรอกตัวเลขให้น้อยลง</p>
        </section>
      </div>
    </main>
  </div>
</template>
