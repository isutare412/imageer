<script lang="ts">
  import FormField from '../form/FormField.svelte';
  import Select from '../form/Select.svelte';
  import type { PresetData } from '$lib/types/admin';

  interface Props {
    preset: PresetData;
    index: number;
    onremove: () => void;
  }

  let { preset = $bindable(), index, onremove }: Props = $props();

  const formatOptions = [
    { value: '', label: 'None (keep original)' },
    { value: 'JPEG', label: 'JPEG' },
    { value: 'PNG', label: 'PNG' },
    { value: 'WEBP', label: 'WebP' },
    { value: 'AVIF', label: 'AVIF' },
    { value: 'HEIC', label: 'HEIC' },
  ];

  const fitOptions = [
    { value: '', label: 'None' },
    { value: 'COVER', label: 'Cover - Scale to cover, may crop' },
    { value: 'CONTAIN', label: 'Contain - Fit within bounds' },
    { value: 'FILL', label: 'Fill - Stretch to fill' },
  ];

  const anchorOptions = [
    { value: '', label: 'None' },
    { value: 'SMART', label: 'Smart - Auto-detect focus' },
    { value: 'CENTER', label: 'Center' },
    { value: 'NORTH', label: 'North (Top)' },
    { value: 'EAST', label: 'East (Right)' },
    { value: 'SOUTH', label: 'South (Bottom)' },
    { value: 'WEST', label: 'West (Left)' },
  ];
</script>

<div class="card bg-base-200 p-4">
  <div class="mb-3 flex items-center justify-between">
    <h4 class="font-medium">Preset #{index + 1}</h4>
    <button
      type="button"
      class="btn btn-ghost btn-sm btn-square text-error"
      onclick={onremove}
      aria-label="Remove preset"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5"
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

  <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
    <FormField label="Name" name="preset-name-{index}" required>
      <input
        type="text"
        id="preset-name-{index}"
        class="input input-bordered w-full"
        placeholder="e.g., w600h800"
        bind:value={preset.name}
        required
      />
    </FormField>

    <FormField label="Format" name="preset-format-{index}">
      <Select name="preset-format-{index}" options={formatOptions} bind:value={preset.format} />
    </FormField>

    <FormField label="Quality (1-100)" name="preset-quality-{index}" hint="Default: 80">
      <input
        type="number"
        id="preset-quality-{index}"
        class="input input-bordered w-full"
        min="1"
        max="100"
        placeholder="80"
        bind:value={preset.quality}
      />
    </FormField>

    <FormField label="Fit Mode" name="preset-fit-{index}">
      <Select name="preset-fit-{index}" options={fitOptions} bind:value={preset.fit} />
    </FormField>

    <FormField label="Anchor" name="preset-anchor-{index}">
      <Select name="preset-anchor-{index}" options={anchorOptions} bind:value={preset.anchor} />
    </FormField>

    <div class="flex items-end gap-2">
      <FormField label="Width (px)" name="preset-width-{index}">
        <input
          type="number"
          id="preset-width-{index}"
          class="input input-bordered w-full"
          min="1"
          placeholder="600"
          bind:value={preset.width}
        />
      </FormField>
      <span class="text-base-content/50 mb-3">×</span>
      <FormField label="Height (px)" name="preset-height-{index}">
        <input
          type="number"
          id="preset-height-{index}"
          class="input input-bordered w-full"
          min="1"
          placeholder="800"
          bind:value={preset.height}
        />
      </FormField>
    </div>
  </div>

  <div class="mt-4">
    <label class="label cursor-pointer justify-start gap-2">
      <input type="checkbox" class="checkbox shrink-0" bind:checked={preset.default} />
      <span class="label-text">Default preset</span>
    </label>
  </div>
</div>
