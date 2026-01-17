import type { PageLoad } from './$types';
import { createApiClient, unwrap } from '$lib/api';
import { env } from '$env/dynamic/public';

export const load: PageLoad = async ({ fetch }) => {
  const client = createApiClient({ fetch, baseUrl: env.PUBLIC_API_BASE_URL });

  // Fetch stats in parallel
  const [projectsResult, serviceAccountsResult] = await Promise.all([
    client.GET('/api/v1/admin/projects', {
      params: { query: { offset: 0, limit: 100 } },
    }),
    client.GET('/api/v1/admin/service-accounts', {
      params: { query: { offset: 0, limit: 1 } },
    }),
  ]);

  const projects = unwrap(projectsResult);
  const serviceAccounts = unwrap(serviceAccountsResult);

  // Calculate total images across all projects
  const totalImages = projects.items.reduce((sum, project) => sum + project.imageCount, 0);

  return {
    stats: {
      projects: projects.total,
      serviceAccounts: serviceAccounts.total,
      totalImages,
    },
  };
};
