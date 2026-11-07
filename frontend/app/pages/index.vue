<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import type { Directive } from 'vue'

const cleanups = new WeakMap<HTMLElement, () => void>()

const vScrollProgress: Directive<HTMLElement, number> = {
  mounted(el, binding) {
    const bar = el.querySelector<HTMLElement>('[data-progress]')
    if (!bar) return

    let frame1 = 0
    let frame2 = 0

    const reset = () => {
      cancelAnimationFrame(frame1)
      cancelAnimationFrame(frame2)

      bar.style.transition = 'none'
      bar.style.width = '0%'
    }

    reset()

    const observer = new IntersectionObserver(([entry]) => {
      reset()

      if (!entry?.isIntersecting) return

      frame1 = requestAnimationFrame(() => {
        frame2 = requestAnimationFrame(() => {
          bar.style.transition = 'width 1000ms ease-out'
          bar.style.width = `${Math.min(100, Math.max(0, Number(binding.value) || 0))}%`
        })
      })
    }, {
      threshold: 0.1
    })

    observer.observe(el)

    cleanups.set(el, () => {
      observer.disconnect()
      cancelAnimationFrame(frame1)
      cancelAnimationFrame(frame2)
    })
  },

  unmounted(el) {
    cleanups.get(el)?.()
    cleanups.delete(el)
  }
}
interface SkillType {
	name: string
	description: string
	score: number
	image_url: string 
  skill_type: number 
}
const {data: skills, pending, error } = await useFetch<SkillType[]>("/server/skill", {lazy: true})
watchEffect(() => {
  if (skills.value) console.log(skills.value)
})

// ===============================
// SKILLS: แท็บแยกประเภท + แสดงเท่าที่จำเป็นก่อน
// skill_type: 1 = ภาษา, 2 = Framework, 3 = อื่น ๆ
// ===============================
const skillCategories = [
  { label: 'ทั้งหมด', value: null },
  { label: 'Languages', value: 1 },
  { label: 'Frameworks', value: 2 },
  { label: 'Others', value: 3 }
] as const

const activeSkillType = ref<number | null>(null)
const showAllSkills = ref(false)

const filteredSkills = computed(() => {
  const list = skills.value ?? []
  return activeSkillType.value === null
    ? list
    : list.filter(skill => skill.skill_type === activeSkillType.value)
})

// แสดงชุดแรก 6 รายการก่อน ไม่ให้หน้าเว็บยาวเกินไปบนมือถือ
const visibleSkills = computed(() =>
  showAllSkills.value ? filteredSkills.value : filteredSkills.value.slice(0, 6)
)

function selectSkillCategory(value: number | null) {
  activeSkillType.value = value
  showAllSkills.value = false
}

// ===============================
// PROJECTS / PORTFOLIO
// แก้ไขข้อมูลผลงานและ URL รูปภาพจริงตรงนี้
// ===============================
interface ProjectType {
  id: number
  title: string
  type: string
  description: string
  features: string[]
  technologies: string[]
  images: string[] // ภาพผลงานกี่รูปก็ได้ รูปแรกใช้เป็นภาพหน้าปก
  demo_url: string
  icon: string
}

// ตัวอย่างข้อมูล: เปลี่ยนเป็นผลงานจริงได้ตามต้องการ
const projects: ProjectType[] = [
  {
    id: 1,
    title: 'ระบบบริหารงานขายรถยนต์',
    type: 'Sales Management System',
    description: 'ระบบจัดการข้อมูลลูกค้า การจองรถ สถานะการขาย และรายงานสรุปสำหรับทีมขาย',
    features: ['จัดการข้อมูลลูกค้าและรายการจอง', 'ติดตามสถานะการขาย', 'รายงานสรุปยอดขาย'],
    technologies: ['Go', 'Gin', 'PostgreSQL', 'Tailwind CSS'],
    // ตัวอย่าง: images: ['/images/projects/car-sales/1.webp', '/images/projects/car-sales/2.webp']
    images: [],
    demo_url: '', // ลิงก์เว็บไซต์จริง (ถ้ามี)
    icon: 'heroicons:truck-20-solid'
  },
  {
    id: 2,
    title: 'ระบบบัญชีและการเงิน',
    type: 'Accounting & Finance System',
    description: 'ระบบบันทึกรายรับรายจ่าย ออกเอกสารทางการเงิน และแสดงภาพรวมรายการบัญชี',
    features: ['บันทึกรายรับและรายจ่าย', 'ค้นหาและจัดการเอกสาร', 'ดูรายงานทางการเงิน'],
    technologies: ['PHP', 'MySQL', 'JavaScript'],
    images: [],
    demo_url: '',
    icon: 'heroicons:banknotes-20-solid'
  },
  {
    id: 3,
    title: 'ระบบจัดการงานและทีม',
    type: 'Task Management System',
    description: 'ระบบจัดการงานที่ช่วยให้ทีมติดตามสถานะงาน ผู้รับผิดชอบ และกำหนดส่งได้สะดวก',
    features: ['เพิ่มและมอบหมายงาน', 'กำหนดสถานะและวันครบกำหนด', 'ติดตามงานผ่านแดชบอร์ด'],
    technologies: ['Vue', 'Nuxt', 'Tailwind CSS'],
    images: [],
    demo_url: '',
    icon: 'heroicons:clipboard-document-list-20-solid'
  },
  {
    id: 4,
    title: 'ระบบจัดการคลังสินค้า',
    type: 'Inventory Management System',
    description: 'ระบบบันทึกสินค้า ควบคุมจำนวนคงเหลือ และตรวจสอบการรับเข้า-จ่ายออก',
    features: ['จัดการข้อมูลสินค้า', 'บันทึกสินค้าเข้าและออก', 'ตรวจสอบยอดคงเหลือ'],
    technologies: ['Go', 'HTMX', 'PostgreSQL'],
    images: [],
    demo_url: '',
    icon: 'heroicons:cube-20-solid'
  }
]

const projectCarousel = ref<HTMLElement | null>(null)
const selectedProject = ref<ProjectType | null>(null)
const activeProjectImage = ref(0)
const projectThumbnails = ref<HTMLElement | null>(null)
const touchStart = ref<{ x: number; y: number } | null>(null)

function slideProjects(direction: -1 | 1) {
  const carousel = projectCarousel.value
  if (!carousel) return

  const card = carousel.querySelector<HTMLElement>('[data-project-card]')
  const distance = card ? card.getBoundingClientRect().width + 16 : carousel.clientWidth
  carousel.scrollBy({ left: distance * direction, behavior: 'smooth' })
}

function openProject(project: ProjectType) {
  activeProjectImage.value = 0 // เปิดผลงานใหม่ให้เริ่มจากรูปแรกเสมอ
  selectedProject.value = project
}

function closeProject() {
  selectedProject.value = null
  touchStart.value = null
}

async function selectProjectImage(index: number) {
  const count = selectedProject.value?.images.length ?? 0
  if (!count) return

  // วนกลับไปต้นทาง/ท้ายสุด เมื่อกดเลยรูปแรกหรือรูปสุดท้าย
  activeProjectImage.value = ((index % count) + count) % count
  await nextTick()

  // เลื่อนแถบภาพตัวอย่างตามรูปที่กำลังดู
  const button = projectThumbnails.value?.children[activeProjectImage.value] as HTMLElement | undefined
  button?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
}

function slideProjectImages(direction: -1 | 1) {
  void selectProjectImage(activeProjectImage.value + direction)
}

function onProjectImageTouchStart(event: TouchEvent) {
  const touch = event.changedTouches[0]
  if (touch) touchStart.value = { x: touch.clientX, y: touch.clientY }
}

function onProjectImageTouchEnd(event: TouchEvent) {
  const start = touchStart.value
  const touch = event.changedTouches[0]
  touchStart.value = null
  if (!start || !touch) return

  const dx = touch.clientX - start.x
  const dy = touch.clientY - start.y
  // จับเฉพาะการปัดแนวนอน ไม่รบกวนการเลื่อนหน้าแนวตั้ง
  if (Math.abs(dx) > 45 && Math.abs(dx) > Math.abs(dy)) {
    slideProjectImages(dx < 0 ? 1 : -1)
  }
}

function handleProjectKeydown(event: KeyboardEvent) {
  if (!selectedProject.value) return
  if (event.key === 'Escape') closeProject()
  else if (event.key === 'ArrowRight') slideProjectImages(1)
  else if (event.key === 'ArrowLeft') slideProjectImages(-1)
}

onMounted(() => window.addEventListener('keydown', handleProjectKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', handleProjectKeydown))
</script>

<template>
  <div class="min-h-dvh bg-gray-50 text-gray-700">



    <main
      class="w-full max-w-[1000px] mx-auto px-4 py-6
             flex flex-col gap-6"
    >

      <!-- =================================== -->
      <!-- PROFILE -->
      <!-- =================================== -->
      <section
        id="home"
        class="bg-white border border-gray-300 rounded-md
               p-5 sm:p-8 scroll-mt-24"
      >
        <div
          class="flex flex-col md:flex-row
                 items-center gap-6 md:gap-10"
        >

          <!-- PROFILE IMAGE -->
          <div
            class="relative w-44 sm:w-52 md:w-56
                   aspect-square shrink-0"
          >
            <div
              class="w-full h-full rounded-full
                     bg-gray-100 border-4 border-white
                     shadow-md flex items-center justify-center
                     overflow-hidden"
            >
              <Icon
                name="heroicons:user-20-solid"
                class="w-[65%] h-[65%] text-blue-300"
              />
            </div>

            <!-- STATUS -->
            <div
              class="absolute bottom-2 right-2
                     w-6 h-6 bg-green-400 rounded-full
                     border-4 border-white"
            ></div>
          </div>


          <!-- INFORMATION -->
          <div
            class="flex flex-col gap-4
                   text-center md:text-left min-w-0"
          >

            <div>
              <h1
                class="text-2xl sm:text-3xl md:text-4xl
                       font-bold text-gray-800"
              >
                Hello, I'm a
              </h1>

              <h2
                class="text-2xl sm:text-3xl md:text-4xl
                       font-bold text-blue-400 mt-1"
              >
                Full Stack Developer
              </h2>
            </div>

            <p class="text-gray-500 text-sm leading-7">
              นักพัฒนาเว็บไซต์ที่สนใจการสร้าง Web Application
              ทั้ง Frontend และ Backend ชื่นชอบการออกแบบ UI
              ที่เรียบง่าย สวยงาม และรองรับการใช้งานบนทุกอุปกรณ์
            </p>

            <!-- TAGS -->
            <div
              class="flex flex-wrap justify-center
                     md:justify-start gap-2"
            >
              <span
                class="bg-gray-100 border border-gray-200
                       rounded-md px-3 py-1 text-xs"
              >
                Frontend
              </span>

              <span
                class="bg-gray-100 border border-gray-200
                       rounded-md px-3 py-1 text-xs"
              >
                Backend
              </span>

              <span
                class="bg-gray-100 border border-gray-200
                       rounded-md px-3 py-1 text-xs"
              >
                Database
              </span>

              <span
                class="bg-gray-100 border border-gray-200
                       rounded-md px-3 py-1 text-xs"
              >
                UI Design
              </span>
            </div>

            <!-- BUTTONS -->
            <div
              class="flex flex-wrap justify-center
                     md:justify-start gap-3 mt-2"
            >
              <a
                href="#skills"
                class="flex items-center gap-2 bg-blue-400
                       text-white px-5 py-2.5 rounded-md
                       text-sm hover:bg-blue-500
                       transition-all duration-300"
              >
                <Icon
                  name="heroicons:code-bracket-20-solid"
                  class="w-5 h-5"
                />
                My Skills
              </a>

              <a
                href="#projects"
                class="flex items-center gap-2 bg-white
                       border border-gray-300
                       px-5 py-2.5 rounded-md text-sm
                       hover:bg-blue-50 hover:border-blue-200
                       transition-all duration-300"
              >
                <Icon name="heroicons:squares-2x2-20-solid" class="w-5 h-5" />
                My Projects
              </a>

              <a
                href="#contact"
                class="flex items-center gap-2 bg-white
                       border border-gray-300
                       px-5 py-2.5 rounded-md text-sm
                       hover:bg-gray-100
                       transition-all duration-300"
              >
                <Icon
                  name="heroicons:chat-bubble-left-right-20-solid"
                  class="w-5 h-5"
                />
                Contact Me
              </a>
            </div>

          </div>
        </div>
      </section>


      <!-- =================================== -->
      <!-- SKILLS: 2 columns on mobile + filters -->
      <!-- =================================== -->
      <section
        id="skills"
        class="bg-white border border-gray-300 rounded-md p-4 sm:p-6 scroll-mt-24"
      >
        <div class="flex items-center gap-3 border-b border-gray-200 pb-4 mb-4">
          <div class="w-10 h-10 rounded-md bg-blue-50 flex items-center justify-center shrink-0">
            <Icon name="heroicons:code-bracket-square-20-solid" class="w-6 h-6 text-blue-400" />
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="font-bold text-gray-800">My Skills</h2>
            <p class="text-xs text-gray-400">ทักษะและเทคโนโลยีที่ใช้งาน</p>
          </div>
          <span class="text-xs text-gray-400 whitespace-nowrap">{{ filteredSkills.length }} ทักษะ</span>
        </div>

        <!-- บนมือถือเลื่อนแท็บแนวนอน แทนการขยายหน้าในแนวตั้ง -->
        <div class="skill-tabs flex gap-2 overflow-x-auto pb-1 mb-4" aria-label="เลือกประเภททักษะ">
          <button
            v-for="category in skillCategories"
            :key="category.label"
            type="button"
            :aria-pressed="activeSkillType === category.value"
            class="shrink-0 text-xs sm:text-sm px-3.5 py-2 rounded-md border font-medium transition-colors"
            :class="activeSkillType === category.value
              ? 'bg-blue-400 border-blue-400 text-white'
              : 'bg-white border-gray-200 text-gray-500 hover:border-blue-200 hover:bg-blue-50'"
            @click="selectSkillCategory(category.value)"
          >
            {{ category.label }}
          </button>
        </div>

        <div v-if="pending && !skills" class="py-8 text-center text-sm text-gray-400">
          กำลังโหลดทักษะ...
        </div>
        <div v-else-if="error" class="py-8 text-center text-sm text-red-500">
          ไม่สามารถโหลดข้อมูลทักษะได้
        </div>
        <div v-else-if="filteredSkills.length === 0" class="py-8 text-center text-sm text-gray-400">
          ยังไม่มีทักษะในหมวดนี้
        </div>

        <!-- มือถือ 2 คอลัมน์ / เดสก์ท็อป 3 คอลัมน์ -->
        <div v-else class="grid grid-cols-2 lg:grid-cols-3 gap-2.5 sm:gap-4">
          <div
            v-for="skill in visibleSkills"
            :key="skill.name"
            class="group min-w-0 border border-gray-200 rounded-md p-3 sm:p-4 hover:border-blue-300 hover:shadow-sm transition-all duration-300"
          >
            <div class="flex items-start justify-between gap-2 mb-3">
              <div class="w-9 h-9 sm:w-11 sm:h-11 shrink-0 bg-blue-50 rounded-md flex items-center justify-center overflow-hidden">
                <img
                  v-if="skill.image_url"
                  :src="skill.image_url"
                  :alt="skill.name"
                  class="w-7 h-7 sm:w-9 sm:h-9 object-contain"
                  loading="lazy"
                />
                <Icon v-else name="heroicons:code-bracket-20-solid" class="w-5 h-5 text-blue-400" />
              </div>
              <span class="text-xs sm:text-sm font-bold text-blue-500 whitespace-nowrap">
                {{ skill.score }}%
              </span>
            </div>

            <h3 class="text-sm font-semibold text-gray-700 truncate" :title="skill.name">
              {{ skill.name }}
            </h3>
            <p class="text-[11px] text-gray-400 mt-1 line-clamp-1" :title="skill.description">
              {{ skill.description }}
            </p>

            <!-- ยังคงใช้ directive เดิม เพื่อให้แถบคะแนน animate เมื่อเลื่อนมาถึง -->
            <div
              v-scroll-progress="skill.score"
              class="w-full h-1.5 sm:h-2 bg-gray-100 rounded-full overflow-hidden mt-3"
            >
              <div
                data-progress
                class="h-full w-0 rounded-full"
                :class="skill.skill_type === 1 ? 'bg-blue-400' : skill.skill_type === 2 ? 'bg-green-400' : 'bg-orange-400'"
              ></div>
            </div>
          </div>
        </div>

        <button
          v-if="filteredSkills.length > 6"
          type="button"
          class="w-full mt-4 py-2.5 flex items-center justify-center gap-2 text-sm font-medium rounded-md border border-gray-200 text-gray-600 hover:bg-blue-50 hover:border-blue-200 transition-colors"
          :aria-expanded="showAllSkills"
          @click="showAllSkills = !showAllSkills"
        >
          {{ showAllSkills ? 'แสดงน้อยลง' : `ดูทักษะทั้งหมด (${filteredSkills.length})` }}
          <Icon
            :name="showAllSkills ? 'heroicons:chevron-up-20-solid' : 'heroicons:chevron-down-20-solid'"
            class="w-4 h-4"
          />
        </button>
      </section>

      <!-- =================================== -->
      <!-- PROJECTS - เลื่อนผลงานซ้าย / ขวา -->
      <!-- =================================== -->
      <section
        id="projects"
        class="bg-white border border-gray-300 rounded-md p-5 sm:p-6 scroll-mt-24"
      >
        <div class="flex flex-wrap items-center justify-between gap-4 border-b border-gray-200 pb-4 mb-5">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-50 rounded-md flex items-center justify-center">
              <Icon name="heroicons:squares-2x2-20-solid" class="w-6 h-6 text-blue-400" />
            </div>
            <div>
              <h2 class="text-xl font-bold text-gray-800">My Projects</h2>
              <p class="text-xs text-gray-400 mt-1">ผลงานเว็บไซต์และระบบที่พัฒนา</p>
            </div>
          </div>

          <!-- ปุ่มเลื่อนซ้าย-ขวา -->
          <div class="flex items-center gap-2">
            <button
              type="button"
              aria-label="เลื่อนไปยังผลงานก่อนหน้า"
              class="w-9 h-9 border border-gray-200 rounded-md flex items-center justify-center hover:bg-blue-50 hover:border-blue-200 transition-colors"
              @click="slideProjects(-1)"
            >
              <Icon name="heroicons:chevron-left-20-solid" class="w-5 h-5" />
            </button>
            <button
              type="button"
              aria-label="เลื่อนไปยังผลงานถัดไป"
              class="w-9 h-9 border border-gray-200 rounded-md flex items-center justify-center hover:bg-blue-50 hover:border-blue-200 transition-colors"
              @click="slideProjects(1)"
            >
              <Icon name="heroicons:chevron-right-20-solid" class="w-5 h-5" />
            </button>
          </div>
        </div>

        <p class="text-sm text-gray-500 mb-4">
          เลื่อนเพื่อดูผลงานเพิ่มเติม หรือกดที่ผลงานเพื่ออ่านรายละเอียดของระบบ
        </p>

        <!-- รองรับการปัดบนมือถือ และปุ่มเลื่อนบนเดสก์ท็อป -->
        <div
          ref="projectCarousel"
          class="portfolio-track flex gap-4 overflow-x-auto snap-x snap-mandatory scroll-smooth pb-2"
          aria-label="รายการผลงาน"
        >
          <button
            v-for="project in projects"
            :key="project.id"
            type="button"
            data-project-card
            class="group w-[85%] sm:w-[48%] lg:w-[340px] shrink-0 snap-start overflow-hidden border border-gray-200 rounded-lg bg-white text-left hover:border-blue-300 hover:shadow-md transition-all duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-300"
            :aria-label="`ดูรายละเอียด ${project.title}`"
            @click="openProject(project)"
          >
            <!-- พื้นที่สำหรับรูปภาพผลงานจริง -->
            <div class="relative h-44 sm:h-48 bg-gray-50 overflow-hidden border-b border-gray-100">
              <img
                v-if="project.images.length > 0"
                :src="project.images[0]"
                :alt="`ภาพหน้าจอ ${project.title}`"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                loading="lazy"
              />
              <!-- ภาพตัวอย่างแทนรูป หากยังไม่ได้เพิ่ม images -->
              <div v-else class="w-full h-full bg-gradient-to-br from-gray-50 to-blue-50 flex items-center justify-center p-5">
                <div class="w-full max-w-[250px] rounded-md border border-gray-200 bg-white shadow-sm overflow-hidden">
                  <div class="h-7 bg-white border-b border-gray-100 px-3 flex items-center gap-1.5">
                    <span class="w-2 h-2 rounded-full bg-red-200"></span>
                    <span class="w-2 h-2 rounded-full bg-yellow-200"></span>
                    <span class="w-2 h-2 rounded-full bg-green-200"></span>
                    <span class="ml-auto text-[9px] text-gray-400">Preview</span>
                  </div>
                  <div class="flex h-[108px]">
                    <div class="w-12 bg-gray-50 border-r border-gray-100 p-2 space-y-2">
                      <div class="h-1.5 bg-blue-200 rounded-full"></div>
                      <div class="h-1.5 bg-gray-200 rounded-full"></div>
                      <div class="h-1.5 bg-gray-200 rounded-full"></div>
                    </div>
                    <div class="flex-1 p-3">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-7 h-7 rounded-md bg-blue-50 flex items-center justify-center">
                          <Icon :name="project.icon" class="w-4 h-4 text-blue-400" />
                        </div>
                        <div class="flex-1 space-y-1.5">
                          <div class="h-1.5 w-2/3 rounded-full bg-gray-300"></div>
                          <div class="h-1.5 w-1/2 rounded-full bg-gray-100"></div>
                        </div>
                      </div>
                      <div class="grid grid-cols-3 gap-2 mb-2">
                        <div class="h-7 rounded bg-blue-50"></div>
                        <div class="h-7 rounded bg-green-50"></div>
                        <div class="h-7 rounded bg-orange-50"></div>
                      </div>
                      <div class="h-1.5 rounded-full bg-gray-100"></div>
                    </div>
                  </div>
                </div>
              </div>
              <span class="absolute top-3 left-3 rounded-md bg-white/95 border border-gray-100 px-2.5 py-1 text-[11px] font-medium text-gray-500 shadow-sm">
                {{ project.type }}
              </span>
            </div>

            <div class="p-4 flex flex-col gap-2">
              <h3 class="font-bold text-gray-800 text-base line-clamp-1">{{ project.title }}</h3>
              <p class="text-xs text-gray-500 leading-6 line-clamp-2 min-h-12">
                {{ project.description }}
              </p>
              <div class="flex items-center justify-between gap-2 pt-2 border-t border-gray-100 mt-1">
                <span class="text-xs text-blue-500 font-semibold">ดูรายละเอียดผลงาน</span>
                <Icon name="heroicons:arrow-up-right-20-solid" class="w-4 h-4 text-blue-400 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform" />
              </div>
            </div>
          </button>
        </div>
        <p class="text-center text-xs text-gray-400 mt-3 sm:hidden">ปัดซ้าย–ขวาเพื่อดูผลงานเพิ่มเติม</p>
      </section>

      <!-- =================================== -->
      <!-- MODAL: รายละเอียดผลงานที่คลิก -->
      <!-- =================================== -->
      <Teleport to="body">
        <div
          v-if="selectedProject"
          class="fixed inset-0 z-[9999] bg-black/50 flex items-center justify-center p-4"
          @click.self="closeProject"
        >
          <div
            role="dialog"
            aria-modal="true"
            aria-labelledby="project-detail-title"
            class="w-full max-w-[720px] max-h-[92dvh] overflow-y-auto bg-white rounded-lg shadow-xl"
          >
            <div class="sticky top-0 z-20 bg-white flex items-start justify-between gap-4 p-4 sm:p-6 border-b border-gray-100">
              <div class="min-w-0">
                <p class="text-xs font-semibold text-blue-500 mb-2">{{ selectedProject.type }}</p>
                <h2 id="project-detail-title" class="text-xl sm:text-2xl font-bold text-gray-800">
                  {{ selectedProject.title }}
                </h2>
              </div>
              <button
                type="button"
                class="shrink-0 w-9 h-9 rounded-md bg-gray-100 hover:bg-gray-200 flex items-center justify-center"
                aria-label="ปิดรายละเอียดผลงาน"
                @click="closeProject"
              >
                <Icon name="heroicons:x-mark-20-solid" class="w-5 h-5" />
              </button>
            </div>

            <div class="p-4 sm:p-6 space-y-5">
              <!-- แกลเลอรีภาพ: รองรับจำนวนภาพไม่จำกัด -->
              <div v-if="selectedProject.images.length" class="space-y-3">
                <div
                  class="relative w-full aspect-video rounded-lg overflow-hidden bg-gray-50 border border-gray-200"
                  role="region"
                  aria-roledescription="carousel"
                  :aria-label="`รูปภาพผลงาน ${selectedProject.title}`"
                  @touchstart.passive="onProjectImageTouchStart"
                  @touchend.passive="onProjectImageTouchEnd"
                >
                  <!-- เลื่อนสไลด์ด้วย translateX: รูปเปลี่ยนอย่างนุ่มนวล -->
                  <div
                    class="flex h-full transition-transform duration-300 ease-out"
                    :style="{ transform: `translateX(-${activeProjectImage * 100}%)` }"
                  >
                    <div
                      v-for="(imageUrl, index) in selectedProject.images"
                      :key="`${selectedProject.id}-${index}`"
                      class="w-full h-full shrink-0 flex items-center justify-center bg-gray-50"
                    >
                      <img
                        :src="imageUrl"
                        :alt="`${selectedProject.title} รูปที่ ${index + 1}`"
                        class="w-full h-full object-contain select-none"
                        loading="lazy"
                        draggable="false"
                      />
                    </div>
                  </div>

                  <template v-if="selectedProject.images.length > 1">
                    <button
                      type="button"
                      aria-label="รูปก่อนหน้า"
                      class="absolute left-2 sm:left-3 top-1/2 -translate-y-1/2 w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-white/90 text-gray-700 border border-gray-200 shadow-sm hover:bg-blue-50 flex items-center justify-center"
                      @click="slideProjectImages(-1)"
                    >
                      <Icon name="heroicons:chevron-left-20-solid" class="w-5 h-5" />
                    </button>
                    <button
                      type="button"
                      aria-label="รูปถัดไป"
                      class="absolute right-2 sm:right-3 top-1/2 -translate-y-1/2 w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-white/90 text-gray-700 border border-gray-200 shadow-sm hover:bg-blue-50 flex items-center justify-center"
                      @click="slideProjectImages(1)"
                    >
                      <Icon name="heroicons:chevron-right-20-solid" class="w-5 h-5" />
                    </button>
                  </template>

                  <span class="absolute bottom-2 right-2 px-2.5 py-1 rounded-md bg-black/60 text-white text-xs">
                    {{ activeProjectImage + 1 }} / {{ selectedProject.images.length }}
                  </span>
                </div>

                <!-- คลิกภาพเล็กเพื่อข้ามไปยังรูปที่ต้องการได้ทันที -->
                <div
                  v-if="selectedProject.images.length > 1"
                  ref="projectThumbnails"
                  class="portfolio-track flex gap-2 overflow-x-auto py-1"
                  aria-label="ภาพตัวอย่างผลงาน"
                >
                  <button
                    v-for="(imageUrl, index) in selectedProject.images"
                    :key="index"
                    type="button"
                    :aria-label="`ดูรูปที่ ${index + 1}`"
                    :aria-pressed="activeProjectImage === index"
                    class="w-16 h-12 sm:w-20 sm:h-14 shrink-0 border-2 rounded-md overflow-hidden bg-gray-50 transition-colors"
                    :class="activeProjectImage === index ? 'border-blue-400' : 'border-gray-200 hover:border-blue-200'"
                    @click="selectProjectImage(index)"
                  >
                    <img :src="imageUrl" :alt="`ภาพตัวอย่าง ${index + 1}`" class="w-full h-full object-cover" loading="lazy" />
                  </button>
                </div>
                <p v-if="selectedProject.images.length > 1" class="text-center text-xs text-gray-400 sm:hidden">
                  ปัดรูปซ้าย–ขวาเพื่อดูเพิ่มเติม
                </p>
              </div>

              <div v-else class="w-full aspect-video rounded-lg border border-gray-100 bg-blue-50 flex flex-col items-center justify-center text-gray-400">
                <Icon :name="selectedProject.icon" class="w-12 h-12 text-blue-300 mb-2" />
                <p class="text-sm">ยังไม่มีรูปภาพผลงาน</p>
              </div>

              <div>
                <h3 class="text-sm font-bold text-gray-800 mb-2">เกี่ยวกับระบบ</h3>
                <p class="text-sm text-gray-600 leading-7">{{ selectedProject.description }}</p>
              </div>

              <div>
                <h3 class="text-sm font-bold text-gray-800 mb-2">ฟังก์ชันหลัก</h3>
                <ul class="space-y-2">
                  <li v-for="feature in selectedProject.features" :key="feature" class="flex items-start gap-2 text-sm text-gray-600">
                    <Icon name="heroicons:check-circle-20-solid" class="w-5 h-5 shrink-0 text-green-500" />
                    <span>{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <div>
                <h3 class="text-sm font-bold text-gray-800 mb-2">เทคโนโลยีที่ใช้</h3>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="tech in selectedProject.technologies"
                    :key="tech"
                    class="text-xs bg-blue-50 border border-blue-100 text-blue-600 px-3 py-1.5 rounded-md"
                  >
                    {{ tech }}
                  </span>
                </div>
              </div>

              <a
                v-if="selectedProject.demo_url"
                :href="selectedProject.demo_url"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 px-4 py-2.5 rounded-md bg-blue-400 hover:bg-blue-500 text-white text-sm transition-colors"
              >
                เปิดเว็บไซต์ผลงาน
                <Icon name="heroicons:arrow-top-right-on-square-20-solid" class="w-4 h-4" />
              </a>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- =================================== -->
      <!-- CONTACT -->
      <!-- =================================== -->
      <section
        id="contact"
        class="bg-white border border-gray-300
               rounded-md p-5 sm:p-6 scroll-mt-24"
      >

        <!-- HEADER -->
        <div
          class="flex items-center gap-3
                 border-b border-gray-200 pb-4 mb-5"
        >
          <div
            class="w-10 h-10 bg-blue-50 rounded-md
                   flex items-center justify-center"
          >
            <Icon
              name="heroicons:chat-bubble-left-right-20-solid"
              class="w-6 h-6 text-blue-400"
            />
          </div>

          <div>
            <h2 class="text-xl font-bold text-gray-800">
              Contact Me
            </h2>

            <p class="text-xs text-gray-400 mt-1">
              ช่องทางการติดต่อ
            </p>
          </div>
        </div>


        <p class="text-sm text-gray-500 mb-6 leading-7">
          หากสนใจร่วมงาน หรือต้องการพูดคุยเกี่ยวกับการพัฒนาเว็บไซต์
          สามารถติดต่อผ่านช่องทางด้านล่างได้เลยครับ
        </p>


        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">

          <!-- EMAIL -->
          <a
            href="mailto:you@example.com"
            class="group flex items-center gap-4
                   border border-gray-200 rounded-md p-4
                   hover:border-blue-300
                   hover:bg-blue-50/50
                   transition-all duration-300"
          >
            

            <div class="flex-1 min-w-0">
              <h3 class="text-sm font-semibold">
                Email
              </h3>

              <p class="text-xs text-gray-400 mt-1 break-all">
                you@example.com
              </p>
            </div>
            <div
              class="w-12 h-12 shrink-0
                     bg-blue-50 rounded-md
                     flex items-center justify-center
                     group-hover:bg-blue-100"
            >
              <Icon
                name="heroicons:envelope-20-solid"
                class="w-7 h-7 text-blue-400"
              />
            </div>
          </a>
          
        </div>

      </section>

    </main>


    <!-- FOOTER -->
    <footer class="border-t border-gray-200 bg-white mt-6">
      <div
        class="max-w-[1000px] mx-auto px-4 py-6
               flex flex-col sm:flex-row
               justify-between items-center gap-3"
      >

        <p class="text-xs text-gray-400">
          © 2026 My Portfolio
        </p>
      </div>
    </footer>

  </div>
</template>
<style scoped>
/* ซ่อน scrollbar แต่ยังเลื่อนด้วยนิ้ว/ทัชแพดได้ */
.skill-tabs { scrollbar-width: none; }
.skill-tabs::-webkit-scrollbar { display: none; }

.portfolio-track {
  scrollbar-width: none;
}
.portfolio-track::-webkit-scrollbar {
  display: none;
}
</style>
