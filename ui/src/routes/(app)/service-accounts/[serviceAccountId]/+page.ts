import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch, params }) => {
  const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

  // Load service account and projects in parallel
  const [saResult, projectsResult] = await Promise.all([
    client.GET('/api/v1/admin/service-accounts/{serviceAccountId}', {
      params: { path: { serviceAccountId: params.serviceAccountId } },
    }),
    client.GET('/api/v1/admin/projects', {
      params: { query: { offset: 0, limit: 100 } },
    }),
  ]);

  const serviceAccount = unwrap(saResult);
  const projects = unwrap(projectsResult);

  return {
    serviceAccount,
    projects,
  };
};
