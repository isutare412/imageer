import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { env } from '$env/dynamic/public';

export const load: LayoutLoad = async ({ fetch, url }) => {
  const client = createApiClient({ fetch, baseUrl: env.PUBLIC_API_BASE_URL });

  const result = await client.GET('/api/v1/users/me');

  if (result.error) {
    const loginUrl = '/login?redirect=' + encodeURIComponent(url.pathname);
    redirect(307, loginUrl);
  }

  const user = unwrap(result);

  // Only admins can access the admin panel
  if (user.role !== 'ADMIN') {
    const loginUrl = '/login?error=unauthorized';
    redirect(307, loginUrl);
  }

  return {
    user,
  };
};
