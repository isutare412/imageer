import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch, url }) => {
  const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

  const offset = parseInt(url.searchParams.get('offset') ?? '0', 10);
  const limit = parseInt(url.searchParams.get('limit') ?? '10', 10);

  const result = await client.GET('/api/v1/admin/projects', {
    params: {
      query: { offset, limit },
    },
  });

  const projects = unwrap(result);

  return {
    projects,
    offset,
    limit,
  };
};
