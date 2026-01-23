<script lang="ts">
  import { goto, invalidateAll } from '$app/navigation';
  import { getApiClient, type UpdateProjectRequest, type UpsertPresetRequest } from '$lib/api';
  import { toastStore, FormField, PresetForm, Badge, ConfirmModal, type PresetData } from '$lib';

  let { data } = $props();

  // Form state (initial values synced via $effect below)
  // svelte-ignore state_referenced_locally
  let name = $state(data.project.name);
  // svelte-ignore state_referenced_locally
  let presets = $state<PresetData[]>(mapPresets(data.project.presets));

  let saving = $state(false);
  let errors = $state<{ name?: string }>({});
  let deleteModal = $state({ open: false, loading: false });

  // Helper to map API presets to local PresetData format
  function mapPresets(apiPresets: typeof data.project.presets): PresetData[] {
    return apiPresets.map((p) => ({
      id: p.id,
      name: p.name,
      default: p.default,
      format: p.format ?? '',
      quality: p.quality,
      fit: p.fit ?? '',
      anchor: p.anchor ?? '',
      width: p.width,
      height: p.height,
    }));
  }

  // Sync form state when page data changes (e.g., on page refresh or navigation)
  $effect(() => {
    name = data.project.name;
    presets = mapPresets(data.project.presets);
  });

  // Track if form has been modified
  let isDirty = $derived.by(() => {
    if (name !== data.project.name) return true;
    if (presets.length !== data.project.presets.length) return true;
    // Could add more detailed comparison here
    return false;
  });

  function addPreset() {
    presets = [
      ...presets,
      {
        name: '',
        default: presets.length === 0,
        format: '',
        quality: undefined,
        fit: '',
        anchor: '',
        width: undefined,
        height: undefined,
      },
    ];
  }

  function removePreset(index: number) {
    presets = presets.filter((_, i) => i !== index);
    const firstPreset = presets[0];
    if (firstPreset && !presets.some((p) => p.default)) {
      firstPreset.default = true;
    }
  }

  function validateForm(): boolean {
    errors = {};

    if (!name.trim()) {
      errors.name = 'Project name is required';
    }

    for (const preset of presets) {
      if (!preset.name.trim()) {
        toastStore.warning('All presets must have a name');
        return false;
      }
    }

    return Object.keys(errors).length === 0;
  }

  function buildPresetRequest(preset: PresetData): UpsertPresetRequest {
    const req: UpsertPresetRequest = {
      name: preset.name,
      default: preset.default,
    };

    if (preset.id) req.id = preset.id;
    if (preset.format) req.format = preset.format as UpsertPresetRequest['format'];
    if (preset.quality) req.quality = preset.quality;
    if (preset.fit) req.fit = preset.fit as UpsertPresetRequest['fit'];
    if (preset.anchor) req.anchor = preset.anchor as UpsertPresetRequest['anchor'];
    if (preset.width) req.width = preset.width;
    if (preset.height) req.height = preset.height;

    return req;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();

    if (!validateForm()) return;

    saving = true;

    const client = getApiClient();
    const body: UpdateProjectRequest = {
      name: name.trim(),
      presets: presets.map(buildPresetRequest),
    };

    const result = await client.PUT('/api/v1/admin/projects/{projectId}', {
      params: { path: { projectId: data.project.id } },
      body,
    });

    if (!result.error) {
      toastStore.success('Project updated successfully');
      // Refresh the page data - the $effect will sync form state from updated data
      await invalidateAll();
    }

    saving = false;
  }

  async function confirmDelete() {
    deleteModal.loading = true;
    const client = getApiClient();
    const result = await client.DELETE('/api/v1/admin/projects/{projectId}', {
      params: { path: { projectId: data.project.id } },
    });

    if (!result.error) {
      toastStore.success(`Project "${data.project.name}" deleted successfully`);
      goto('/projects');
    } else {
      deleteModal.loading = false;
    }
  }

  function handleReprocess() {
    toastStore.warning('Image reprocessing feature is not implemented yet.');
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
  <title>{data.project.name} | Imageer</title>
</svelte:head>

<div class="space-y-6">
  <!-- Page header -->
  <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
    <div class="min-w-0">
      <nav class="breadcrumbs text-sm">
        <ul>
          <li><a href="/projects" class="link link-hover">Projects</a></li>
          <li class="truncate">{data.project.name}</li>
        </ul>
      </nav>
      <h1 class="mt-2 truncate text-2xl font-bold">{data.project.name}</h1>
      <p class="text-base-content/60 mt-1 text-sm">
        {data.project.imageCount.toLocaleString()} images · Created {formatDate(
          data.project.createdAt
        )} · Updated {formatDate(data.project.updatedAt)}
      </p>
    </div>
    <div class="flex shrink-0 gap-2 self-end sm:self-auto">
      <button type="button" class="btn btn-outline btn-sm" onclick={handleReprocess}>
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
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        <span class="hidden sm:inline">Reprocess</span>
      </button>
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
    <!-- Basic info card -->
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-lg">Project Details</h2>
        <div class="mt-4 grid max-w-2xl min-w-0 gap-4">
          <FormField label="Project ID" name="id">
            <input
              type="text"
              class="input input-bordered w-full font-mono text-sm"
              value={data.project.id}
              readonly
            />
          </FormField>
          <FormField label="Project Name" name="name" required error={errors.name}>
            <input
              type="text"
              id="name"
              name="name"
              class="input input-bordered w-full"
              class:input-error={errors.name}
              placeholder="e.g., my-awesome-project"
              bind:value={name}
              required
            />
          </FormField>
        </div>
      </div>
    </div>

    <!-- Presets card -->
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="card-title text-lg">Image Presets</h2>
            <p class="text-base-content/60 mt-1 text-sm">
              Define how images should be processed and converted
            </p>
          </div>
          <button type="button" class="btn btn-outline btn-sm" onclick={addPreset}>
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
                d="M12 4v16m8-8H4"
              />
            </svg>
            Add Preset
          </button>
        </div>

        {#if presets.length === 0}
          <div class="bg-base-200 mt-4 rounded-lg p-8 text-center">
            <p class="text-base-content/60">
              No presets defined. Add presets to specify how images should be processed.
            </p>
            <button type="button" class="btn btn-primary btn-sm mt-4" onclick={addPreset}>
              Add First Preset
            </button>
          </div>
        {:else}
          <div class="mt-4 space-y-4">
            {#each { length: presets.length } as _, index (presets[index]?.id ?? index)}
              {#if presets[index]}
                <PresetForm
                  bind:preset={presets[index]}
                  {index}
                  onremove={() => removePreset(index)}
                />
              {/if}
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-2">
      <a href="/projects" class="btn">Back to Projects</a>
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
  title="Delete Project"
  message="Are you sure you want to delete the project '{data.project
    .name}'? This action cannot be undone and will delete all associated images."
  confirmText="Delete"
  confirmVariant="error"
  loading={deleteModal.loading}
  onconfirm={confirmDelete}
  oncancel={() => (deleteModal = { open: false, loading: false })}
/>
