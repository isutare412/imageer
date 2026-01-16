<script lang="ts">
  import { goto } from '$app/navigation';
  import {
    getApiClient,
    unwrap,
    type CreateServiceAccountRequest,
    type ServiceAccountWithApiKey,
  } from '$lib/api';
  import { toastStore, FormField, Select, MultiSelect, ApiKeyDisplay } from '$lib';

  let { data } = $props();

  let name = $state('');
  let accessScope = $state<'FULL' | 'PROJECT'>('PROJECT');
  let projectIds = $state<string[]>([]);
  let expireAt = $state('');

  let loading = $state(false);
  let errors = $state<{ name?: string; projectIds?: string }>({});

  // After creation, store the result with API key
  let createdAccount = $state<ServiceAccountWithApiKey | null>(null);

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

    loading = true;

    const client = getApiClient();
    const body: CreateServiceAccountRequest = {
      name: name.trim(),
      accessScope,
      projectIds: accessScope === 'PROJECT' ? projectIds : [],
    };

    if (expireAt) {
      body.expireAt = new Date(expireAt).toISOString();
    }

    const result = await client.POST('/api/v1/admin/service-accounts', { body });

    if (!result.error) {
      const account = unwrap(result);
      createdAccount = account;
      toastStore.success(`Service account "${account.name}" created successfully`);
    }

    loading = false;
  }

  function handleDone() {
    if (createdAccount) {
      goto(`/service-accounts/${createdAccount.id}`);
    } else {
      goto('/service-accounts');
    }
  }
</script>

<div class="space-y-6">
  <!-- Page header -->
  <div>
    <nav class="breadcrumbs text-sm">
      <ul>
        <li><a href="/service-accounts" class="link link-hover">Service Accounts</a></li>
        <li>New Service Account</li>
      </ul>
    </nav>
    <h1 class="mt-2 text-2xl font-bold">Create Service Account</h1>
    <p class="text-base-content/60 mt-1">Set up API access credentials for your applications</p>
  </div>

  {#if createdAccount}
    <!-- Success state - show API key -->
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <div class="mb-4 flex items-center gap-3">
          <div
            class="bg-success/20 text-success flex h-12 w-12 items-center justify-center rounded-full"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-bold">Service Account Created</h2>
            <p class="text-base-content/60">{createdAccount.name}</p>
          </div>
        </div>

        <ApiKeyDisplay apiKey={createdAccount.apiKey} />

        <div class="mt-6 flex justify-end">
          <button type="button" class="btn btn-primary" onclick={handleDone}> Done </button>
        </div>
      </div>
    </div>
  {:else}
    <!-- Form -->
    <form onsubmit={handleSubmit} class="space-y-6">
      <div class="card bg-base-100 shadow-sm">
        <div class="card-body">
          <h2 class="card-title text-lg">Account Details</h2>

          <div class="mt-4 grid max-w-2xl min-w-0 gap-4">
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

      <!-- Actions -->
      <div class="flex justify-end gap-2">
        <a href="/service-accounts" class="btn">Cancel</a>
        <button type="submit" class="btn btn-primary" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
          {/if}
          Create Service Account
        </button>
      </div>
    </form>
  {/if}
</div>
