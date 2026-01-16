import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: LayoutLoad = async ({ fetch, url }) => {
	const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

	const result = await client.GET('/api/v1/users/me');

	if (result.error) {
		const loginUrl = '/login?redirect=' + encodeURIComponent(url.pathname);
		redirect(307, loginUrl);
	}

	const user = unwrap(result);

	return {
		user
	};
};
