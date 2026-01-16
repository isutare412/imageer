<script lang="ts">
  import type { ToastLevel } from '$lib/stores/toast.svelte';
  import { toastStore } from '$lib/stores/toast.svelte';

  interface Props {
    id: string;
    message: string;
    level: ToastLevel;
    timeout: number;
    createdAt: number;
  }

  let { id, message, level, timeout, createdAt }: Props = $props();

  let progress = $state(100);
  let animationFrame: number;
  let isHovered = $state(false);
  // svelte-ignore state_referenced_locally
  let remainingTime = $state(timeout);
  // svelte-ignore state_referenced_locally
  let lastUpdateTime = $state(createdAt);

  const alertClasses: Record<ToastLevel, string> = {
    info: 'alert-info',
    success: 'alert-success',
    warning: 'alert-warning',
    error: 'alert-error',
  };

  const progressClasses: Record<ToastLevel, string> = {
    info: 'bg-info-content/50',
    success: 'bg-success-content/50',
    warning: 'bg-warning-content/50',
    error: 'bg-error-content/50',
  };

  const icons: Record<ToastLevel, string> = {
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    warning:
      'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
    error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
  };

  function updateProgress() {
    if (isHovered) {
      animationFrame = requestAnimationFrame(updateProgress);
      return;
    }

    const now = Date.now();
    const elapsed = now - lastUpdateTime;
    lastUpdateTime = now;
    remainingTime = Math.max(0, remainingTime - elapsed);
    progress = (remainingTime / timeout) * 100;

    if (remainingTime > 0) {
      animationFrame = requestAnimationFrame(updateProgress);
    } else {
      toastStore.remove(id);
    }
  }

  function handleClose() {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }
    toastStore.remove(id);
  }

  $effect(() => {
    animationFrame = requestAnimationFrame(updateProgress);

    return () => {
      if (animationFrame) {
        cancelAnimationFrame(animationFrame);
      }
    };
  });
</script>

<div
  class="alert {alertClasses[level]} relative overflow-hidden shadow-lg"
  onmouseenter={() => (isHovered = true)}
  onmouseleave={() => {
    isHovered = false;
    lastUpdateTime = Date.now();
  }}
  role="alert"
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    class="h-6 w-6 shrink-0 stroke-current"
  >
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={icons[level]} />
  </svg>
  <span class="flex-1 whitespace-pre-wrap">{message}</span>
  <button
    type="button"
    class="btn btn-ghost btn-sm btn-square"
    onclick={handleClose}
    aria-label="Close"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      class="h-5 w-5 stroke-current"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  </button>
  <div class="absolute inset-x-0 bottom-0 h-1 bg-black/10">
    <div class="h-full {progressClasses[level]} transition-none" style="width: {progress}%"></div>
  </div>
</div>
