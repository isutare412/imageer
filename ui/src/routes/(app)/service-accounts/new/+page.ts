import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch }) => {
  const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

  // Load projects for the project selector
  const result = await client.GET('/api/v1/admin/projects', {
    params: {
      query: { offset: 0, limit: 100 }, // Load up to 100 projects
    },
  });

  const projects = unwrap(result);

  return {
    projects,
  };
};
