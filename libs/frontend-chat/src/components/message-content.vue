<script setup lang="ts">
import { computed } from 'vue';

import { type MessageChunk, EmoteFlag } from '../types.js';

const props = defineProps<{
	chunks: MessageChunk[]
	isItalic?: boolean
	textShadowColor?: string
	textShadowSize?: number
	userColor: string
	// messageAlign: MessageAlignType
}>();

const getEmoteWidth = (w?: number) => {
	return `${w ? w * 2 : 50}px`;
};

const textShadow = computed(() => {
	if (!props.textShadowColor || !props.textShadowSize) return '';

	const array = Array.from({ length: 5 }).map((_, i) => {
		const n = i + 1;
		return `0px 0px ${props.textShadowSize! + n}px ${props.textShadowColor}`;
	});

	return array.join(', ');
});

// const wordBreak = computed(() => {
// 	return props.messageAlign === 'baseline' ? 'break-all' : 'initial';
// });

const mappedChunks = props.chunks.reduce((acc, chunk) => {
	if (chunk.type === 'text') {
		const lastChunk = acc[acc.length - 1];
		if (lastChunk && lastChunk.type === 'text') {
			lastChunk.value += ' ' + chunk.value;
		} else {
			acc.push(chunk);
		}
	} else {
		acc.push(chunk);
	}

	return acc;
}, [] as MessageChunk[]);
</script>

<template>
	<span
		class="text"
		:style="{
			fontStyle: isItalic ? 'italic' : 'normal',
			color: isItalic ? userColor : 'inherit',
		}"
	>
		<template v-for="(chunk, _) of mappedChunks" :key="_">
			<div
				v-if="['emote', '3rd_party_emote'].includes(chunk.type)"
				class="emote"
			>
				<img
					:src="chunk.type === 'emote'
						? `https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/1.0`
						: chunk.value
					"
					:class="{
						'emote-cursed': chunk.flags?.includes(EmoteFlag.Cursed),
						'flipX': chunk.flags?.includes(EmoteFlag.FlipX),
						'flipY': chunk.flags?.includes(EmoteFlag.FlipY),
					}"
					:style="{
						width: chunk.flags?.includes(EmoteFlag.GrowX) ? getEmoteWidth(chunk.emoteWidth) : undefined,
					}"
				/>

				<img
					v-for="(c, idx) of chunk.zeroWidthModifiers" :key="idx" class="emote-zerowidth"
					:src="c"
				/>
			</div>

			<template v-else-if="['text', 'emoji'].includes(chunk.type)">
				{{ chunk.value }}
			</template>
		</template>
	</span>
</template>

<style scoped>
.text {
	text-shadow: v-bind(textShadow);
	vertical-align: baseline;
}

.text > span {
	font-style: inherit;
	font-weight: inherit;
	font-style: inherit;
}

.emote img {
	height: 1.5em;
}

.text .emote {
	position: relative;
	display: inline-block;
	margin-left: 4px;
	margin-right: 4px;
	vertical-align: middle;
}

.text .emote .emote-zerowidth {
	top: 50%;
	left: 50%;
	bottom: 0;
	transform: translate(-50%, -50%);
	position: absolute;
}

.text .emote .emote-cursed {
	filter: grayscale(1) brightness(0.7) contrast(2.5);
}

.text .emote .flipX {
	transform: scaleX(-1);
}

.text .emote .flipY {
	transform: scaleY(-1);
}
</style>
