<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { LayoutData } from './$types';
  import { goto } from '$app/navigation';
  import { getApiClient } from '$lib/api';
  import { themeStore } from '$lib';

  let { children, data }: { children: Snippet; data: LayoutData } = $props();

  let drawerOpen = $state(false);

  async function handleSignOut() {
    const client = getApiClient();
    await client.POST('/api/v1/auth/sign-out');
    goto('/login');
  }
</script>

<div class="drawer lg:drawer-open">
  <input id="app-drawer" type="checkbox" class="drawer-toggle" bind:checked={drawerOpen} />

  <div class="drawer-content bg-base-200 flex flex-col">
    <!-- Mobile navbar -->
    <div class="navbar bg-base-100 border-base-300 border-b lg:hidden">
      <div class="flex-none">
        <label for="app-drawer" class="btn btn-square btn-ghost">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            class="inline-block h-5 w-5 stroke-current"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            ></path>
          </svg>
        </label>
      </div>
      <div class="flex-1">
        <span class="text-xl font-bold">Imageer</span>
      </div>
    </div>

    <!-- Main content -->
    <main class="flex-1 overflow-auto">
      <div class="mx-auto max-w-5xl p-6">
        {@render children()}
      </div>
    </main>
  </div>

  <!-- Sidebar -->
  <div class="drawer-side z-40">
    <label for="app-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
    <aside class="bg-base-100 border-base-300 min-h-full w-64 border-r">
      <div class="flex h-full flex-col">
        <!-- Logo -->
        <div class="border-base-300 border-b p-4">
          <h1 class="text-xl font-bold">Imageer</h1>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 p-4">
          <ul class="menu gap-1">
            <li>
              <a href="/" class="flex items-center gap-3" onclick={() => (drawerOpen = false)}>
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
                    d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                  />
                </svg>
                Dashboard
              </a>
            </li>
            <li>
              <a
                href="/projects"
                class="flex items-center gap-3"
                onclick={() => (drawerOpen = false)}
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
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                  />
                </svg>
                Projects
              </a>
            </li>
            <li>
              <a
                href="/service-accounts"
                class="flex items-center gap-3"
                onclick={() => (drawerOpen = false)}
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
                    d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                  />
                </svg>
                Service Accounts
              </a>
            </li>
          </ul>
        </nav>

        <!-- User section -->
        <div class="border-base-300 border-t p-4">
          <div class="flex items-center gap-3">
            {#if data.user.photoUrl}
              <div class="avatar">
                <div class="w-10 rounded-full">
                  <img src={data.user.photoUrl} alt={data.user.nickname} />
                </div>
              </div>
            {:else}
              <div class="avatar placeholder">
                <div class="bg-neutral text-neutral-content w-10 rounded-full">
                  <span class="text-sm">{data.user.nickname.charAt(0).toUpperCase()}</span>
                </div>
              </div>
            {/if}
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium">{data.user.nickname}</p>
              <p class="text-base-content/60 truncate text-xs">{data.user.email}</p>
            </div>
          </div>
          <div class="mt-3 flex gap-2">
            <button
              class="btn btn-ghost btn-sm flex-1 justify-start"
              onclick={() => themeStore.toggle()}
              aria-label="Toggle theme"
            >
              {#if themeStore.current === 'dark'}
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
                    d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
                  />
                </svg>
                Light
              {:else}
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
                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
                  />
                </svg>
                Dark
              {/if}
            </button>
            <button class="btn btn-ghost btn-sm flex-1 justify-start" onclick={handleSignOut}>
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
                  d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                />
              </svg>
              Sign out
            </button>
          </div>
        </div>
      </div>
    </aside>
  </div>
</div>
