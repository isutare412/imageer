import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch, params }) => {
  const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

  const result = await client.GET('/api/v1/admin/projects/{projectId}', {
    params: { path: { projectId: params.projectId } },
  });

  const project = unwrap(result);

  return {
    project,
  };
};
