<script lang="ts">
  import { goto, invalidateAll } from '$app/navigation';
  import { getApiClient, unwrapEmpty, type Project } from '$lib/api';
  import { toastStore, EmptyState, Pagination, Badge, ConfirmModal } from '$lib';

  let { data } = $props();

  let deleteModal = $state({ open: false, project: null as Project | null, loading: false });

  function handlePageChange(newOffset: number) {
    const url = new URL(window.location.href);
    url.searchParams.set('offset', newOffset.toString());
    goto(url.toString());
  }

  function openDeleteModal(project: Project) {
    deleteModal = { open: true, project, loading: false };
  }

  async function confirmDelete() {
    if (!deleteModal.project) return;

    deleteModal.loading = true;
    const client = getApiClient();
    const result = await client.DELETE('/api/v1/admin/projects/{projectId}', {
      params: { path: { projectId: deleteModal.project.id } },
    });

    if (!result.error) {
      unwrapEmpty(result);
      toastStore.success(`Project "${deleteModal.project.name}" deleted successfully`);
      deleteModal = { open: false, project: null, loading: false };
      invalidateAll();
    } else {
      deleteModal.loading = false;
    }
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }
</script>

<div class="space-y-6">
  <!-- Page header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Projects</h1>
      <p class="text-base-content/60 mt-1">Manage your image processing projects</p>
    </div>
    <a href="/projects/new" class="btn btn-primary">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      New Project
    </a>
  </div>

  <!-- Projects table -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body p-0">
      {#if data.projects.items.length === 0}
        <EmptyState
          title="No projects yet"
          description="Create your first project to start processing images"
        >
          {#snippet actions()}
            <a href="/projects/new" class="btn btn-primary">Create Project</a>
          {/snippet}
        </EmptyState>
      {:else}
        <div class="overflow-x-auto">
          <table class="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Presets</th>
                <th>Created</th>
                <th>Updated</th>
                <th class="w-24">Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each data.projects.items as project}
                <tr class="hover">
                  <td>
                    <a href="/projects/{project.id}" class="link link-hover font-medium">
                      {project.name}
                    </a>
                  </td>
                  <td>
                    <div class="flex flex-wrap gap-1">
                      {#each project.presets.slice(0, 3) as preset}
                        <Badge variant={preset.default ? 'primary' : 'neutral'} size="sm" truncate>
                          {preset.name}
                        </Badge>
                      {/each}
                      {#if project.presets.length > 3}
                        <Badge variant="neutral" size="sm" outline>
                          +{project.presets.length - 3} more
                        </Badge>
                      {/if}
                      {#if project.presets.length === 0}
                        <span class="text-base-content/40 text-sm">No presets</span>
                      {/if}
                    </div>
                  </td>
                  <td class="text-base-content/60">{formatDate(project.createdAt)}</td>
                  <td class="text-base-content/60">{formatDate(project.updatedAt)}</td>
                  <td>
                    <div class="flex gap-1">
                      <a
                        href="/projects/{project.id}"
                        class="btn btn-ghost btn-sm btn-square"
                        aria-label="Edit project"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                      </a>
                      <button
                        type="button"
                        class="btn btn-ghost btn-sm btn-square text-error"
                        onclick={() => openDeleteModal(project)}
                        aria-label="Delete project"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="border-base-200 border-t p-4">
          <Pagination
            total={data.projects.total}
            limit={data.limit}
            offset={data.offset}
            onPageChange={handlePageChange}
          />
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Delete confirmation modal -->
<ConfirmModal
  bind:open={deleteModal.open}
  title="Delete Project"
  message="Are you sure you want to delete the project '{deleteModal.project
    ?.name}'? This action cannot be undone and will delete all associated images."
  confirmText="Delete"
  confirmVariant="error"
  loading={deleteModal.loading}
  onconfirm={confirmDelete}
  oncancel={() => (deleteModal = { open: false, project: null, loading: false })}
/>
