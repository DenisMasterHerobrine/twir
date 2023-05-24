import { config } from '@tsuwari/config';
import * as YTSR from '@tsuwari/grpc/generated/ytsr/ytsr';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';
import tlds from 'tlds' assert { type: 'json' };
import ytsrLib, { Video } from 'ytsr';

import { durationToMilliseconds } from './utils/convertDuration.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const linkRegexp = new RegExp(
  `\\S+[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})(?=\\P{L}|$)\\S*`,
  'iu',
);

const youtubeLinkRegexp = /(?:https?:\/\/)?(?:www\.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed)?(?:.*v=|v\/|\/)([\w\-_]+)&?/;

const isYoutubeLink = (l: string) => ['youtube', 'youtu.be'].some(link => l.includes(link));

const ytsrService: YTSR.YtsrServiceImplementation = {
  async search(
    request: YTSR.SearchRequest,
    _context,
  ): Promise<YTSR.DeepPartial<YTSR.SearchResponse>> {
    const videos: Array<YTSR.Song> = [];

    const tracksForSearch: string[] = [];

    const linkMatches = [...request.search.matchAll(linkRegexp)];

    const youtubeMatches = linkMatches.filter(link => isYoutubeLink(link[0]));
    const nonYoutubeMatches = linkMatches.filter(link => !isYoutubeLink(link[0]));

    if (linkMatches.length) {
      if (youtubeMatches.length) {
        tracksForSearch.push(...youtubeMatches.map(m => {
					const fullLink = m[0];
					const match = youtubeLinkRegexp.exec(fullLink)!;

					return `https://www.youtube.com/watch?v=${match[1]!}`;
				}));
      }

      if (nonYoutubeMatches.length) {
        await Promise.all(
          linkMatches.map(async (match) => {
            const request = await fetch(
              `https://api.song.link/v1-alpha.1/links?url=${match[0]}&key=${config.ODESLI_API_KEY}`,
            );
            if (!request.ok) return;

            const data = await request.json();
            const youTube = data.linksByPlatform?.youtube;
            if (!youTube) return;

            tracksForSearch.push(youTube.url);
          }),
        );
      }
    } else {
      tracksForSearch.push(request.search);
    }

    await Promise.all(
      tracksForSearch.map(async (track) => {

        const search = await ytsrLib(track, { limit: 1 });
        const item = search.items.at(0) as Video;
        if (!item) return;

        videos.push({
          title: item.title,
          isLive: item.isLive,
          views: item.views ?? 0,
          id: item.id,
          thumbnailUrl: item.bestThumbnail?.url ?? undefined,
          duration: durationToMilliseconds(item.duration ?? '0:00'),
          author: {
            name: item.author?.name || '',
            channelId: item.author?.channelID || '',
            avatarUrl: item.author?.bestAvatar?.url || '',
          },
        });
      }),
    );

    return {
      songs: videos,
    };
  },
};

grpcServer.add(YTSR.YtsrDefinition, ytsrService);

await grpcServer.listen(`0.0.0.0:${PORTS.YTSR_SERVER_PORT}`);
console.log('YTSR server listening on port', PORTS.YTSR_SERVER_PORT);

process.on('SIGINT', () => grpcServer.shutdown()).on('SIGTERM', () => grpcServer.shutdown());
