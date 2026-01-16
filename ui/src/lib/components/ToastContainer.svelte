<script lang="ts">
  import { fly } from 'svelte/transition';
  import { flip } from 'svelte/animate';
  import { toastStore } from '$lib/stores/toast.svelte';
  import Toast from './Toast.svelte';
</script>

<div class="toast toast-top toast-center z-50 w-full max-w-md gap-2 px-4 md:max-w-xl">
  {#each [...toastStore.toasts.values()].reverse() as toast (toast.id)}
    <div
      out:fly={{ x: 100, duration: 150 }}
      animate:flip={{ duration: 150 }}
      class="animate-slide-in-right w-full"
    >
      <Toast
        id={toast.id}
        message={toast.message}
        level={toast.level}
        timeout={toast.timeout}
        createdAt={toast.createdAt}
      />
    </div>
  {/each}
</div>

<style>
  @keyframes slide-in-right {
    from {
      opacity: 0;
      transform: translateX(100px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .animate-slide-in-right {
    animation: slide-in-right 0.2s ease-out;
  }
</style>
