<script lang="ts">
  import { invalidateAll } from '$app/navigation';
  import { getApiClient, unwrapEmpty, type Image, type Project } from '$lib/api';
  import {
    toastStore,
    EmptyState,
    Pagination,
    Badge,
    ConfirmModal,
    imagePreferencesStore,
  } from '$lib';

  let { data } = $props();

  // State
  let selectedProjectId = $state<string | null>(null);
  let images = $state<Image[]>([]);
  let total = $state(0);
  let offset = $state(0);
  let loading = $state(false);
  let expandedImageId = $state<string | null>(null);

  // Delete modal state
  let deleteModal = $state({ open: false, image: null as Image | null, loading: false });

  // Tab state for expanded images
  let activeTabByImage = $state<Record<string, string>>({});

  function getActiveTab(imageId: string): string {
    return activeTabByImage[imageId] ?? 'original';
  }

  function setActiveTab(imageId: string, tab: string) {
    activeTabByImage = { ...activeTabByImage, [imageId]: tab };
  }

  // Derived
  let selectedProject = $derived(
    data.projects.find((p: Project) => p.id === selectedProjectId) ?? null
  );

  // Fetch images when project or pagination/sort changes
  async function fetchImages() {
    if (!selectedProjectId) {
      images = [];
      total = 0;
      return;
    }

    loading = true;
    const client = getApiClient();
    const result = await client.GET('/api/v1/admin/projects/{projectId}/images', {
      params: {
        path: { projectId: selectedProjectId },
        query: {
          offset,
          limit: imagePreferencesStore.perPage,
          sortBy: imagePreferencesStore.sortBy,
          sortOrder: imagePreferencesStore.sortOrder,
        },
      },
    });

    if (!result.error) {
      images = result.data.items;
      total = result.data.total;
    }
    loading = false;
  }

  // Handle project selection
  function handleProjectChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    selectedProjectId = target.value || null;
    offset = 0;
    expandedImageId = null;
    fetchImages();
  }

  // Handle sort/perPage changes
  function handleSortByChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    imagePreferencesStore.sortBy = target.value as 'createdAt' | 'updatedAt';
    offset = 0;
    fetchImages();
  }

  function handleSortOrderChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    imagePreferencesStore.sortOrder = target.value as 'ASC' | 'DESC';
    offset = 0;
    fetchImages();
  }

  function handlePerPageChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    imagePreferencesStore.perPage = parseInt(target.value, 10) as 10 | 20 | 50;
    offset = 0;
    fetchImages();
  }

  function handlePageChange(newOffset: number) {
    offset = newOffset;
    fetchImages();
  }

  // Helper functions
  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  function openDeleteModal(image: Image) {
    deleteModal = { open: true, image, loading: false };
  }

  async function confirmDelete() {
    if (!deleteModal.image || !selectedProjectId) return;

    deleteModal.loading = true;
    const client = getApiClient();
    const result = await client.DELETE('/api/v1/admin/projects/{projectId}/images/{imageId}', {
      params: {
        path: { projectId: selectedProjectId, imageId: deleteModal.image.id },
      },
    });

    if (!result.error) {
      unwrapEmpty(result);
      toastStore.success('Image deleted successfully');
      deleteModal = { open: false, image: null, loading: false };
      await invalidateAll(); // Refresh project list with updated image counts
      fetchImages();
    } else {
      deleteModal.loading = false;
    }
  }
</script>

<svelte:head>
  <title>Images | Imageer</title>
</svelte:head>

<div class="space-y-6">
  <!-- Page header -->
  <div>
    <h1 class="text-2xl font-bold">Image Management</h1>
    <p class="text-base-content/60 mt-1">Browse and manage images across projects</p>
  </div>

  <!-- Controls bar -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body p-4">
      <div class="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-center">
        <!-- Project selector -->
        <select
          class="select select-bordered w-full sm:w-auto sm:min-w-48"
          value={selectedProjectId ?? ''}
          onchange={handleProjectChange}
        >
          <option value="">Select a project</option>
          {#each data.projects as project (project.id)}
            <option value={project.id}
              >{project.name} ({project.imageCount.toLocaleString()})</option
            >
          {/each}
        </select>

        <!-- Sort controls -->
        <div class="flex flex-wrap gap-2">
          <select
            class="select select-bordered select-sm"
            value={imagePreferencesStore.sortBy}
            onchange={handleSortByChange}
            disabled={!selectedProjectId}
          >
            <option value="createdAt">Created At</option>
            <option value="updatedAt">Updated At</option>
          </select>

          <select
            class="select select-bordered select-sm"
            value={imagePreferencesStore.sortOrder}
            onchange={handleSortOrderChange}
            disabled={!selectedProjectId}
          >
            <option value="DESC">Newest First</option>
            <option value="ASC">Oldest First</option>
          </select>

          <select
            class="select select-bordered select-sm"
            value={imagePreferencesStore.perPage}
            onchange={handlePerPageChange}
            disabled={!selectedProjectId}
          >
            <option value={10}>10 per page</option>
            <option value={20}>20 per page</option>
            <option value={50}>50 per page</option>
          </select>
        </div>

        <!-- Total count -->
        {#if selectedProjectId}
          <div class="text-base-content/60 text-sm sm:ml-auto">
            Total: {total.toLocaleString()} images
          </div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Content area -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body p-0">
      {#if !selectedProjectId}
        <EmptyState
          title="Select a Project"
          description="Choose a project from the dropdown above to view its images."
        />
      {:else if loading}
        <div class="flex items-center justify-center py-12">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else if images.length === 0}
        <EmptyState title="No Images" description="This project doesn't have any images yet." />
      {:else}
        <div class="divide-base-200 divide-y">
          {#each images as image (image.id)}
            <div class="border-base-200 border-b last:border-b-0">
              <!-- Collapsed row -->
              <!-- svelte-ignore a11y_click_events_have_key_events -->
              <!-- svelte-ignore a11y_no_static_element_interactions -->
              <div
                class="hover:bg-base-200/50 flex w-full cursor-pointer items-center gap-2 p-4 text-left transition-colors sm:gap-4"
                onclick={() => (expandedImageId = expandedImageId === image.id ? null : image.id)}
              >
                <!-- Chevron -->
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4 shrink-0 transition-transform {expandedImageId === image.id
                    ? 'rotate-90'
                    : ''}"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 5l7 7-7 7"
                  />
                </svg>

                <!-- Desktop layout -->
                <div class="hidden flex-1 items-center gap-4 sm:flex">
                  <span class="tooltip font-mono text-sm" data-tip={image.id}>
                    {image.id.slice(0, 6)}...{image.id.slice(-6)}
                  </span>
                  <Badge
                    variant={image.state === 'READY'
                      ? 'success'
                      : image.state === 'FAILED'
                        ? 'error'
                        : 'warning'}
                    size="sm"
                  >
                    {image.state}
                  </Badge>
                  <span class="text-base-content/60 text-sm">{image.format}</span>
                  <span class="text-base-content/60 text-sm"
                    >{image.variants?.length ?? 0} variants</span
                  >
                  <span class="text-base-content/60 text-sm">{formatDate(image.createdAt)}</span>
                </div>

                <!-- Mobile layout -->
                <div class="flex flex-1 flex-col gap-1 sm:hidden">
                  <span class="tooltip font-mono text-sm" data-tip={image.id}>
                    {image.id.slice(0, 8)}...{image.id.slice(-8)}
                  </span>
                  <div class="text-base-content/60 flex flex-wrap items-center gap-2 text-xs">
                    <Badge
                      variant={image.state === 'READY'
                        ? 'success'
                        : image.state === 'FAILED'
                          ? 'error'
                          : 'warning'}
                      size="xs"
                    >
                      {image.state}
                    </Badge>
                    <span>{image.format}</span>
                    <span>{image.variants?.length ?? 0} variants</span>
                  </div>
                  <span class="text-base-content/60 text-xs">{formatDate(image.createdAt)}</span>
                </div>

                <!-- Delete button -->
                <button
                  type="button"
                  class="btn btn-ghost btn-sm btn-square text-error shrink-0"
                  onclick={(e) => {
                    e.stopPropagation();
                    openDeleteModal(image);
                  }}
                  aria-label="Delete image"
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

              <!-- Expanded content with tabbed viewer -->
              {#if expandedImageId === image.id}
                <div class="bg-base-200/30 border-base-200 border-t p-4">
                  <!-- Tabs -->
                  <div role="tablist" class="tabs tabs-bordered mb-4">
                    <button
                      type="button"
                      role="tab"
                      class="tab {getActiveTab(image.id) === 'original' ? 'tab-active' : ''}"
                      onclick={() => setActiveTab(image.id, 'original')}
                    >
                      Original
                    </button>
                    {#each image.variants ?? [] as variant (variant.id)}
                      <button
                        type="button"
                        role="tab"
                        class="tab {getActiveTab(image.id) === variant.id ? 'tab-active' : ''}"
                        onclick={() => setActiveTab(image.id, variant.id)}
                      >
                        <span class="flex items-center gap-1">
                          {variant.presetName}
                          {#if variant.state !== 'READY'}
                            <Badge
                              variant={variant.state === 'FAILED' ? 'error' : 'warning'}
                              size="xs"
                            >
                              {variant.state}
                            </Badge>
                          {/if}
                        </span>
                      </button>
                    {/each}
                  </div>

                  <!-- Tab content -->
                  <div class="card bg-base-100">
                    <div class="card-body p-4">
                      {#if getActiveTab(image.id) === 'original'}
                        <!-- Original image tab -->
                        {#if image.state === 'READY' && image.url}
                          <div class="mb-4">
                            <img
                              src={image.url}
                              alt="Original"
                              class="max-h-96 rounded-lg object-contain"
                              loading="lazy"
                            />
                          </div>
                        {:else if image.state === 'FAILED'}
                          <div class="bg-error/10 text-error mb-4 rounded-lg p-4">
                            Image processing failed
                          </div>
                        {:else}
                          <div class="bg-warning/10 text-warning mb-4 rounded-lg p-4">
                            Image is {image.state.toLowerCase().replace('_', ' ')}
                          </div>
                        {/if}
                        <dl class="grid gap-2 text-sm sm:grid-cols-2">
                          <div>
                            <dt class="text-base-content/60">State</dt>
                            <dd>
                              <Badge
                                variant={image.state === 'READY'
                                  ? 'success'
                                  : image.state === 'FAILED'
                                    ? 'error'
                                    : 'warning'}
                                size="sm"
                              >
                                {image.state}
                              </Badge>
                            </dd>
                          </div>
                          <div>
                            <dt class="text-base-content/60">Format</dt>
                            <dd>{image.format}</dd>
                          </div>
                          <div>
                            <dt class="text-base-content/60">Created</dt>
                            <dd>{formatDate(image.createdAt)}</dd>
                          </div>
                          <div>
                            <dt class="text-base-content/60">Updated</dt>
                            <dd>{formatDate(image.updatedAt)}</dd>
                          </div>
                          {#if image.url}
                            <div class="overflow-hidden sm:col-span-2">
                              <dt class="text-base-content/60">URL</dt>
                              <dd class="flex items-center gap-2">
                                <code
                                  class="bg-base-200 min-w-0 flex-1 truncate rounded px-2 py-1 text-xs"
                                >
                                  {image.url}
                                </code>
                                <button
                                  type="button"
                                  class="btn btn-ghost btn-xs"
                                  onclick={() => {
                                    navigator.clipboard.writeText(image.url!);
                                    toastStore.success('URL copied to clipboard');
                                  }}
                                >
                                  Copy
                                </button>
                              </dd>
                            </div>
                          {/if}
                        </dl>
                      {:else}
                        <!-- Variant tab -->
                        {@const variant = image.variants?.find(
                          (v) => v.id === getActiveTab(image.id)
                        )}
                        {#if variant}
                          {#if variant.state === 'READY' && variant.url}
                            <div class="mb-4">
                              <img
                                src={variant.url}
                                alt="Variant: {variant.presetName}"
                                class="max-h-96 rounded-lg object-contain"
                                loading="lazy"
                              />
                            </div>
                          {:else if variant.state === 'FAILED'}
                            <div class="bg-error/10 text-error mb-4 rounded-lg p-4">
                              Variant processing failed
                            </div>
                          {:else}
                            <div class="bg-warning/10 text-warning mb-4 rounded-lg p-4">
                              Variant is {variant.state.toLowerCase().replace('_', ' ')}
                            </div>
                          {/if}
                          <dl class="grid gap-2 text-sm sm:grid-cols-2">
                            <div>
                              <dt class="text-base-content/60">Preset</dt>
                              <dd>{variant.presetName}</dd>
                            </div>
                            <div>
                              <dt class="text-base-content/60">State</dt>
                              <dd>
                                <Badge
                                  variant={variant.state === 'READY'
                                    ? 'success'
                                    : variant.state === 'FAILED'
                                      ? 'error'
                                      : 'warning'}
                                  size="sm"
                                >
                                  {variant.state}
                                </Badge>
                              </dd>
                            </div>
                            <div>
                              <dt class="text-base-content/60">Format</dt>
                              <dd>{variant.format}</dd>
                            </div>
                            <div>
                              <dt class="text-base-content/60">Created</dt>
                              <dd>{formatDate(variant.createdAt)}</dd>
                            </div>
                            <div>
                              <dt class="text-base-content/60">Updated</dt>
                              <dd>{formatDate(variant.updatedAt)}</dd>
                            </div>
                            {#if variant.url}
                              <div class="overflow-hidden sm:col-span-2">
                                <dt class="text-base-content/60">URL</dt>
                                <dd class="flex items-center gap-2">
                                  <code
                                    class="bg-base-200 min-w-0 flex-1 truncate rounded px-2 py-1 text-xs"
                                  >
                                    {variant.url}
                                  </code>
                                  <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => {
                                      navigator.clipboard.writeText(variant.url!);
                                      toastStore.success('URL copied to clipboard');
                                    }}
                                  >
                                    Copy
                                  </button>
                                </dd>
                              </div>
                            {/if}
                          </dl>
                        {/if}
                      {/if}
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          {/each}
        </div>
        <div class="border-base-200 border-t p-4">
          <Pagination
            {total}
            limit={imagePreferencesStore.perPage}
            {offset}
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
  title="Delete Image"
  message="Are you sure you want to delete this image? This will also delete all {deleteModal.image
    ?.variants?.length ?? 0} variants. This action cannot be undone."
  confirmText="Delete"
  confirmVariant="error"
  loading={deleteModal.loading}
  onconfirm={confirmDelete}
  oncancel={() => (deleteModal = { open: false, image: null, loading: false })}
/>
