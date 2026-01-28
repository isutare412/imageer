<script lang="ts">
  import './layout.css';
  import { ToastContainer, themeStore } from '$lib';
  import { navigating } from '$app/state';
  import { onMount } from 'svelte';

  let { children } = $props();

  onMount(() => {
    themeStore.init();
  });
</script>

{#if navigating.from}
  <progress class="progress progress-primary fixed top-0 right-0 left-0 z-50 h-1 rounded-none"
  ></progress>
{/if}

<svelte:head>
  <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
  <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
  <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
  <link rel="manifest" href="/site.webmanifest" />
</svelte:head>
<ToastContainer />
{@render children()}
