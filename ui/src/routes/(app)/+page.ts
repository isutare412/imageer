import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch }) => {
  const client = createApiClient({ fetch, baseUrl: PUBLIC_API_BASE_URL });

  // Fetch stats in parallel
  const [projectsResult, serviceAccountsResult] = await Promise.all([
    client.GET('/api/v1/admin/projects', {
      params: { query: { offset: 0, limit: 1 } },
    }),
    client.GET('/api/v1/admin/service-accounts', {
      params: { query: { offset: 0, limit: 1 } },
    }),
  ]);

  const projects = unwrap(projectsResult);
  const serviceAccounts = unwrap(serviceAccountsResult);

  return {
    stats: {
      projects: projects.total,
      serviceAccounts: serviceAccounts.total,
    },
  };
};
