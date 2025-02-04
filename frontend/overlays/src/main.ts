import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import './style.css';

import App from './App.vue';

const routes = createRouter({
	history: createWebHistory('/overlays'),
	routes: [
		{
			path: '/:apiKey/registry/overlays/:overlayId',
			component: () => import('./pages/overlays.vue'),
		},
		{
			path: '/:apiKey/tts',
			component: () => import('./pages/tts.vue'),
		},
		{
			path: '/:apiKey/obs',
			component: () => import('./pages/obs.vue'),
		},
		{
			path: '/:apiKey/alerts',
			component: () => import('./pages/alerts.vue'),
		},
		{
			path: '/:apiKey/chat',
			component: () => import('./pages/overlays/chat.vue'),
		},
		{
			path: '/:apiKey/kappagen',
			component: () => import('./pages/overlays/kappagen.vue'),
		},
		{
			path: '/:apiKey/brb',
			component: () => import('./pages/overlays/beRightBack.vue'),
		},
	],
});

createApp(App).use(routes).mount('#app');

// refresh the page when new version comes
document.body.addEventListener('plugin_web_update_notice', () => {
  window.location.reload();
});
