<script lang="ts">
  import { goto, invalidateAll } from '$app/navigation';
  import { Badge, ConfirmModal, FormField, MultiSelect, Select, toastStore } from '$lib';
  import { getApiClient, type UpdateServiceAccountRequest } from '$lib/api';

  let { data } = $props();

  // Form state (synced via $effect below)
  // svelte-ignore state_referenced_locally
  let name = $state(data.serviceAccount.name);
  // svelte-ignore state_referenced_locally
  let accessScope = $state<'FULL' | 'PROJECT'>(data.serviceAccount.accessScope);
  // svelte-ignore state_referenced_locally
  let projectIds = $state<string[]>(data.serviceAccount.projects.map((p) => p.id));
  // svelte-ignore state_referenced_locally
  let expireAt = $state(formatExpireAt(data.serviceAccount.expireAt));

  let saving = $state(false);
  let errors = $state<{ name?: string; projectIds?: string }>({});
  let deleteModal = $state({ open: false, loading: false });

  function formatExpireAt(expireAt: string | undefined): string {
    return expireAt ? new Date(expireAt).toISOString().slice(0, 16) : '';
  }

  // Sync form state when page data changes (e.g., on page refresh or navigation)
  $effect(() => {
    name = data.serviceAccount.name;
    accessScope = data.serviceAccount.accessScope;
    projectIds = data.serviceAccount.projects.map((p) => p.id);
    expireAt = formatExpireAt(data.serviceAccount.expireAt);
  });

  let projectOptions = $derived(
    data.projects.items.map((p) => ({
      value: p.id,
      label: p.name,
    }))
  );

  const accessScopeOptions = [
    { value: 'PROJECT', label: 'Project - Access to specific projects only' },
    { value: 'FULL', label: 'Full - Access to all projects (admin)' },
  ];

  function isExpired(): boolean {
    if (!data.serviceAccount.expireAt) return false;
    return new Date(data.serviceAccount.expireAt) < new Date();
  }

  function validateForm(): boolean {
    errors = {};

    if (!name.trim()) {
      errors.name = 'Name is required';
    }

    if (accessScope === 'PROJECT' && projectIds.length === 0) {
      errors.projectIds = 'Select at least one project';
    }

    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();

    if (!validateForm()) return;

    saving = true;

    const client = getApiClient();
    const body: UpdateServiceAccountRequest = {
      name: name.trim(),
      accessScope,
      projectIds: accessScope === 'PROJECT' ? projectIds : [],
    };

    if (expireAt) {
      body.expireAt = new Date(expireAt).toISOString();
    }

    const result = await client.PUT('/api/v1/admin/service-accounts/{serviceAccountId}', {
      params: { path: { serviceAccountId: data.serviceAccount.id } },
      body,
    });

    if (!result.error) {
      toastStore.success('Service account updated successfully');
      // Refresh the page data - the $effect will sync form state from updated data
      await invalidateAll();
    }

    saving = false;
  }

  async function confirmDelete() {
    deleteModal.loading = true;
    const client = getApiClient();
    const result = await client.DELETE('/api/v1/admin/service-accounts/{serviceAccountId}', {
      params: { path: { serviceAccountId: data.serviceAccount.id } },
    });

    if (!result.error) {
      toastStore.success(`Service account "${data.serviceAccount.name}" deleted successfully`);
      goto('/service-accounts');
    } else {
      deleteModal.loading = false;
    }
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
</script>

<svelte:head>
  <title>{data.serviceAccount.name} | Imageer</title>
</svelte:head>

<div class="space-y-6">
  <!-- Page header -->
  <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
    <div class="min-w-0">
      <nav class="breadcrumbs text-sm">
        <ul>
          <li><a href="/service-accounts" class="link link-hover">Service Accounts</a></li>
          <li>{data.serviceAccount.name}</li>
        </ul>
      </nav>
      <div class="mt-2 flex items-center gap-2">
        <h1 class="text-2xl font-bold">{data.serviceAccount.name}</h1>
        {#if isExpired()}
          <Badge variant="error">Expired</Badge>
        {/if}
      </div>
      <p class="text-base-content/60 mt-1 text-sm">
        Created {formatDate(data.serviceAccount.createdAt)} · Updated {formatDate(
          data.serviceAccount.updatedAt
        )}
      </p>
    </div>
    <div class="flex shrink-0 gap-2 self-end sm:self-auto">
      <button
        type="button"
        class="btn btn-error btn-outline btn-sm"
        onclick={() => (deleteModal.open = true)}
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
        <span class="hidden sm:inline">Delete</span>
      </button>
    </div>
  </div>

  <!-- Form -->
  <form onsubmit={handleSubmit} class="space-y-6">
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-lg">Account Details</h2>

        <div class="mt-4 grid max-w-2xl min-w-0 gap-4">
          <FormField label="Service Account ID" name="id">
            <input
              type="text"
              class="input input-bordered w-full font-mono text-sm"
              value={data.serviceAccount.id}
              readonly
            />
          </FormField>

          <FormField label="Name" name="name" required error={errors.name}>
            <input
              type="text"
              id="name"
              name="name"
              class="input input-bordered w-full"
              class:input-error={errors.name}
              placeholder="e.g., production-api"
              bind:value={name}
              required
            />
          </FormField>

          <FormField
            label="Access Scope"
            name="accessScope"
            required
            hint="PROJECT scope limits access to specific projects. FULL scope grants admin-level access."
          >
            <Select name="accessScope" options={accessScopeOptions} bind:value={accessScope} />
          </FormField>

          {#if accessScope === 'PROJECT'}
            <FormField
              label="Projects"
              name="projectIds"
              required
              error={errors.projectIds}
              hint="Select which projects this service account can access"
            >
              <MultiSelect
                name="projectIds"
                options={projectOptions}
                bind:values={projectIds}
                placeholder="Select projects..."
              />
            </FormField>
          {/if}

          <FormField label="Expiration Date" name="expireAt" hint="Leave empty for no expiration">
            <input
              type="datetime-local"
              id="expireAt"
              name="expireAt"
              class="input input-bordered w-full"
              bind:value={expireAt}
            />
          </FormField>
        </div>
      </div>
    </div>

    <!-- API Key notice -->
    <div class="alert alert-info">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6 shrink-0 stroke-current"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <div>
        <p class="font-semibold">API Key Not Shown</p>
        <p class="text-sm">
          For security reasons, the API key cannot be retrieved after creation. If you've lost the
          API key, delete this service account and create a new one.
        </p>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-2">
      <a href="/service-accounts" class="btn">Back to Service Accounts</a>
      <button type="submit" class="btn btn-primary" disabled={saving}>
        {#if saving}
          <span class="loading loading-spinner loading-sm"></span>
        {/if}
        Save Changes
      </button>
    </div>
  </form>
</div>

<!-- Delete confirmation modal -->
<ConfirmModal
  bind:open={deleteModal.open}
  title="Delete Service Account"
  message="Are you sure you want to delete the service account '{data.serviceAccount
    .name}'? This will immediately revoke all API access using this account."
  confirmText="Delete"
  confirmVariant="error"
  loading={deleteModal.loading}
  onconfirm={confirmDelete}
  oncancel={() => (deleteModal = { open: false, loading: false })}
/>
