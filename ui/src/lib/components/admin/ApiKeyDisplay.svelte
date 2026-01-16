<script lang="ts">
  import { toastStore } from '$lib/stores/toast.svelte';

  interface Props {
    apiKey: string;
  }

  let { apiKey }: Props = $props();

  let copied = $state(false);

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(apiKey);
      copied = true;
      toastStore.success('API key copied to clipboard');
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch {
      toastStore.error('Failed to copy to clipboard');
    }
  }
</script>

<div class="alert alert-warning">
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
      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
    />
  </svg>
  <div class="flex-1">
    <h3 class="font-bold">Save your API key now!</h3>
    <p class="text-sm">This key will only be shown once. Store it securely.</p>
  </div>
</div>

<div class="bg-base-200 mt-4 rounded-lg p-4">
  <label class="label" for="api-key-display">
    <span class="label-text font-medium">API Key</span>
  </label>
  <div class="flex min-w-0 gap-2">
    <input
      id="api-key-display"
      type="text"
      class="input input-bordered min-w-0 flex-1 font-mono text-sm"
      value={apiKey}
      readonly
    />
    <button
      type="button"
      class="btn {copied ? 'btn-success' : 'btn-primary'}"
      onclick={copyToClipboard}
    >
      {#if copied}
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
            d="M5 13l4 4L19 7"
          />
        </svg>
        Copied!
      {:else}
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
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
          />
        </svg>
        Copy
      {/if}
    </button>
  </div>
</div>
