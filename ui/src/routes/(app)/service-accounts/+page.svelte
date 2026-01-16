<script lang="ts">
  import { goto, invalidateAll } from '$app/navigation';
  import { getApiClient, unwrapEmpty, type ServiceAccount } from '$lib/api';
  import { toastStore, EmptyState, Pagination, Badge, ConfirmModal } from '$lib';

  let { data } = $props();

  let deleteModal = $state({
    open: false,
    serviceAccount: null as ServiceAccount | null,
    loading: false,
  });

  function handlePageChange(newOffset: number) {
    const url = new URL(window.location.href);
    url.searchParams.set('offset', newOffset.toString());
    goto(url.toString());
  }

  function openDeleteModal(serviceAccount: ServiceAccount) {
    deleteModal = { open: true, serviceAccount, loading: false };
  }

  async function confirmDelete() {
    if (!deleteModal.serviceAccount) return;

    deleteModal.loading = true;
    const client = getApiClient();
    const result = await client.DELETE('/api/v1/admin/service-accounts/{serviceAccountId}', {
      params: { path: { serviceAccountId: deleteModal.serviceAccount.id } },
    });

    if (!result.error) {
      unwrapEmpty(result);
      toastStore.success(
        `Service account "${deleteModal.serviceAccount.name}" deleted successfully`
      );
      deleteModal = { open: false, serviceAccount: null, loading: false };
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

  function isExpired(expireAt: string | undefined): boolean {
    if (!expireAt) return false;
    return new Date(expireAt) < new Date();
  }

  function getScopeVariant(scope: string): 'info' | 'warning' {
    return scope === 'FULL' ? 'warning' : 'info';
  }
</script>

<div class="space-y-6">
  <!-- Page header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Service Accounts</h1>
      <p class="text-base-content/60 mt-1">Manage API access credentials</p>
    </div>
    <a href="/service-accounts/new" class="btn btn-primary">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      New Service Account
    </a>
  </div>

  <!-- Service accounts table -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body p-0">
      {#if data.serviceAccounts.items.length === 0}
        <EmptyState
          title="No service accounts yet"
          description="Create a service account to enable API access for your applications"
        >
          {#snippet actions()}
            <a href="/service-accounts/new" class="btn btn-primary">Create Service Account</a>
          {/snippet}
        </EmptyState>
      {:else}
        <div class="overflow-x-auto">
          <table class="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Access Scope</th>
                <th>Projects</th>
                <th>Expires</th>
                <th>Created</th>
                <th class="w-24">Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each data.serviceAccounts.items as sa}
                {@const expired = isExpired(sa.expireAt)}
                <tr class="hover" class:opacity-60={expired}>
                  <td>
                    <div class="flex items-center gap-2">
                      <a href="/service-accounts/{sa.id}" class="link link-hover font-medium">
                        {sa.name}
                      </a>
                      {#if expired}
                        <Badge variant="error" size="xs">Expired</Badge>
                      {/if}
                    </div>
                  </td>
                  <td>
                    <Badge variant={getScopeVariant(sa.accessScope)} size="sm">
                      {sa.accessScope}
                    </Badge>
                  </td>
                  <td>
                    {#if sa.accessScope === 'FULL'}
                      <span class="text-base-content/60 text-sm">All projects</span>
                    {:else if sa.projects.length === 0}
                      <span class="text-base-content/40 text-sm">None</span>
                    {:else}
                      <div class="flex flex-wrap gap-1">
                        {#each sa.projects.slice(0, 2) as project}
                          <Badge variant="neutral" size="sm" truncate>{project.name}</Badge>
                        {/each}
                        {#if sa.projects.length > 2}
                          <Badge variant="neutral" size="sm" outline>
                            +{sa.projects.length - 2} more
                          </Badge>
                        {/if}
                      </div>
                    {/if}
                  </td>
                  <td class="text-base-content/60">
                    {#if sa.expireAt}
                      <span class:text-error={expired}>
                        {formatDate(sa.expireAt)}
                      </span>
                    {:else}
                      <span class="text-base-content/40">Never</span>
                    {/if}
                  </td>
                  <td class="text-base-content/60">{formatDate(sa.createdAt)}</td>
                  <td>
                    <div class="flex gap-1">
                      <a
                        href="/service-accounts/{sa.id}"
                        class="btn btn-ghost btn-sm btn-square"
                        aria-label="Edit service account"
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
                        onclick={() => openDeleteModal(sa)}
                        aria-label="Delete service account"
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
            total={data.serviceAccounts.total}
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
  title="Delete Service Account"
  message="Are you sure you want to delete the service account '{deleteModal.serviceAccount
    ?.name}'? This will immediately revoke all API access using this account."
  confirmText="Delete"
  confirmVariant="error"
  loading={deleteModal.loading}
  onconfirm={confirmDelete}
  oncancel={() => (deleteModal = { open: false, serviceAccount: null, loading: false })}
/>
