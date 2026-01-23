import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { env } from '$env/dynamic/public';

export const load: PageLoad = async ({ fetch }) => {
  const client = createApiClient({ fetch, baseUrl: env.PUBLIC_API_BASE_URL });

  // Fetch all projects for the dropdown selector
  const result = await client.GET('/api/v1/admin/projects', {
    params: {
      query: { limit: 1000 }, // Get all projects for dropdown
    },
  });

  const projects = unwrap(result);

  return {
    projects: projects.items,
  };
};
