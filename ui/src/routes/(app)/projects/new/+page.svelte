<script lang="ts">
  import { goto } from '$app/navigation';
  import {
    getApiClient,
    unwrap,
    type CreateProjectRequest,
    type CreatePresetRequest,
  } from '$lib/api';
  import { toastStore, FormField, PresetForm, type PresetData } from '$lib';

  let name = $state('');
  let presets = $state<PresetData[]>([]);
  let loading = $state(false);
  let errors = $state<{ name?: string }>({});

  function addPreset() {
    presets = [
      ...presets,
      {
        name: '',
        default: presets.length === 0, // First preset is default
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
    // Ensure at least one default if we have presets
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

    // Validate presets
    for (const preset of presets) {
      if (!preset.name.trim()) {
        toastStore.warning('All presets must have a name');
        return false;
      }
    }

    return Object.keys(errors).length === 0;
  }

  function buildPresetRequest(preset: PresetData): CreatePresetRequest {
    const req: CreatePresetRequest = {
      name: preset.name,
      default: preset.default,
    };

    if (preset.format) req.format = preset.format as CreatePresetRequest['format'];
    if (preset.quality) req.quality = preset.quality;
    if (preset.fit) req.fit = preset.fit as CreatePresetRequest['fit'];
    if (preset.anchor) req.anchor = preset.anchor as CreatePresetRequest['anchor'];
    if (preset.width) req.width = preset.width;
    if (preset.height) req.height = preset.height;

    return req;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();

    if (!validateForm()) return;

    loading = true;

    const client = getApiClient();
    const body: CreateProjectRequest = {
      name: name.trim(),
      presets: presets.map(buildPresetRequest),
    };

    const result = await client.POST('/api/v1/admin/projects', { body });

    if (!result.error) {
      const project = unwrap(result);
      toastStore.success(`Project "${project.name}" created successfully`);
      goto(`/projects/${project.id}`);
    } else {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>New Project | Imageer</title>
</svelte:head>

<div class="space-y-6">
  <!-- Page header -->
  <div>
    <nav class="breadcrumbs text-sm">
      <ul>
        <li><a href="/projects" class="link link-hover">Projects</a></li>
        <li>New Project</li>
      </ul>
    </nav>
    <h1 class="mt-2 text-2xl font-bold">Create New Project</h1>
    <p class="text-base-content/60 mt-1">Set up a new image processing project</p>
  </div>

  <!-- Form -->
  <form onsubmit={handleSubmit} class="space-y-6">
    <!-- Basic info card -->
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-lg">Project Details</h2>
        <div class="mt-4 max-w-md">
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
              No presets defined yet. Add presets to specify how images should be processed.
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
      <a href="/projects" class="btn">Cancel</a>
      <button type="submit" class="btn btn-primary" disabled={loading}>
        {#if loading}
          <span class="loading loading-spinner loading-sm"></span>
        {/if}
        Create Project
      </button>
    </div>
  </form>
</div>
