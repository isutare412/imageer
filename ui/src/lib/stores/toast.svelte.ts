import { SvelteMap } from 'svelte/reactivity';

export type ToastLevel = 'info' | 'success' | 'warning' | 'error';

export interface Toast {
  id: string;
  message: string;
  level: ToastLevel;
  timeout: number;
  createdAt: number;
}

export interface ToastOptions {
  message: string;
  level?: ToastLevel;
  timeout?: number;
}

const DEFAULT_TIMEOUT = 5000;

function createToastStore() {
  const toasts = new SvelteMap<string, Toast>();

  function add(options: ToastOptions): string {
    const id = crypto.randomUUID();
    const toast: Toast = {
      id,
      message: options.message,
      level: options.level ?? 'info',
      timeout: options.timeout ?? DEFAULT_TIMEOUT,
      createdAt: Date.now(),
    };

    toasts.set(id, toast);

    return id;
  }

  function remove(id: string) {
    toasts.delete(id);
  }

  function clear() {
    toasts.clear();
  }

  // Convenience methods
  function info(message: string, timeout?: number) {
    return add({ message, level: 'info', timeout });
  }

  function success(message: string, timeout?: number) {
    return add({ message, level: 'success', timeout });
  }

  function warning(message: string, timeout?: number) {
    return add({ message, level: 'warning', timeout });
  }

  function error(message: string, timeout?: number) {
    return add({ message, level: 'error', timeout });
  }

  return {
    get toasts() {
      return toasts;
    },
    add,
    remove,
    clear,
    info,
    success,
    warning,
    error,
  };
}

export const toastStore = createToastStore();
