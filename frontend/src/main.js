import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from "vue-router";
import App from './App.vue'

import Home from "./views/Home.vue";
import TeacherDetail from "./views/TeacherDetail.vue";

import { DynamicScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {path: '/', component: Home},
        {path: '/teacher/:id', component: TeacherDetail, props: true}
    ]
})

const app = createApp(App)
app.use(router)
app.component('DynamicScroller', DynamicScroller)
app.mount('#app')
