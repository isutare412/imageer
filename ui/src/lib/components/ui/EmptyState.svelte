<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title: string;
    description?: string;
    icon?: Snippet;
    actions?: Snippet;
  }

  let { title, description, icon, actions }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 text-center">
  {#if icon}
    <div class="text-base-content/30 mb-4">
      {@render icon()}
    </div>
  {:else}
    <div class="text-base-content/30 mb-4">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-16 w-16"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
        />
      </svg>
    </div>
  {/if}
  <h3 class="text-lg font-semibold">{title}</h3>
  {#if description}
    <p class="text-base-content/60 mt-1 max-w-sm">{description}</p>
  {/if}
  {#if actions}
    <div class="mt-4">
      {@render actions()}
    </div>
  {/if}
</div>
