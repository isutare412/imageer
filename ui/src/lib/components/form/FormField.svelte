<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    label: string;
    name: string;
    required?: boolean;
    error?: string;
    hint?: string;
    children: Snippet;
  }

  let { label, name, required = false, error, hint, children }: Props = $props();
</script>

<div class="form-control w-full min-w-0">
  <label class="label" for={name}>
    <span class="label-text">
      {label}
      {#if required}
        <span class="text-error">*</span>
      {/if}
    </span>
  </label>
  {@render children()}
  {#if error}
    <p class="text-error mt-1 text-xs">{error}</p>
  {:else if hint}
    <p class="text-base-content/60 mt-1 text-xs">{hint}</p>
  {/if}
</div>
