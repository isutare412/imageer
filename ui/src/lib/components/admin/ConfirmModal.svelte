<script lang="ts">
  interface Props {
    open: boolean;
    title: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    confirmVariant?: 'primary' | 'error' | 'warning';
    loading?: boolean;
    onconfirm: () => void;
    oncancel: () => void;
  }

  let {
    open = $bindable(),
    title,
    message,
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    confirmVariant = 'primary',
    loading = false,
    onconfirm,
    oncancel,
  }: Props = $props();

  const variantClasses: Record<string, string> = {
    primary: 'btn-primary',
    error: 'btn-error',
    warning: 'btn-warning',
  };

  function handleConfirm() {
    onconfirm();
  }

  function handleCancel() {
    open = false;
    oncancel();
  }

  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      handleCancel();
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape' && !loading) {
      handleCancel();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="modal modal-open" onclick={handleBackdropClick}>
    <div class="modal-box">
      <h3 class="text-lg font-bold">{title}</h3>
      <p class="py-4">{message}</p>
      <div class="modal-action">
        <button type="button" class="btn" disabled={loading} onclick={handleCancel}>
          {cancelText}
        </button>
        <button
          type="button"
          class="btn {variantClasses[confirmVariant]}"
          disabled={loading}
          onclick={handleConfirm}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
          {/if}
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}
