import { useQuery } from '@tanstack/vue-query';
import {
	GetUsersRequest_Order,
	GetUsersRequest_SortBy,
} from '@twir/grpc/generated/api/api/community';
import { ComputedRef, Ref, unref } from 'vue';

import { unprotectedClient } from './twirp.js';

const sortBy = {
	'watched': GetUsersRequest_SortBy.Watched,
	'messages': GetUsersRequest_SortBy.Messages,
	'emotes': GetUsersRequest_SortBy.Emotes,
	'usedChannelPoints': GetUsersRequest_SortBy.UsedChannelPoints,
};

export type SortKey = keyof typeof sortBy

export type GetCommunityUsersOpts = {
	limit: number;
	page: number;
	desc: boolean;
	sortBy: SortKey;
	channelId?: string
}

export const useCommunityUsers = (opts: Ref<GetCommunityUsersOpts> | ComputedRef<GetCommunityUsersOpts>) => {
	return useQuery({
		queryKey: ['communityUsers', opts],
		queryFn: async () => {
			const rawOpts = unref(opts);
			if (!rawOpts.channelId) return;
			console.log(rawOpts);

			const order = rawOpts.desc ? GetUsersRequest_Order.Desc: GetUsersRequest_Order.Asc;
			const call = await unprotectedClient.communityGetUsers({
				limit: rawOpts.limit,
				page: rawOpts.page,
				order,
				sortBy: sortBy[rawOpts.sortBy],
				channelId: rawOpts.channelId,
			}, { timeout: 5000 });
			return call.response;
		},
	});
};
