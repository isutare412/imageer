<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import { goto } from '$app/navigation';
	import { getApiClient } from '$lib/api';

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

	<div class="drawer-content flex flex-col bg-base-200">
		<!-- Mobile navbar -->
		<div class="navbar bg-base-100 lg:hidden border-b border-base-300">
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
			<div class="max-w-5xl mx-auto p-6">
				{@render children()}
			</div>
		</main>
	</div>

	<!-- Sidebar -->
	<div class="drawer-side z-40">
		<label for="app-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
		<aside class="w-64 min-h-full bg-base-100 border-r border-base-300">
			<div class="flex flex-col h-full">
				<!-- Logo -->
				<div class="p-4 border-b border-base-300">
					<h1 class="text-xl font-bold">Imageer</h1>
				</div>

				<!-- Navigation -->
				<nav class="flex-1 p-4">
					<ul class="menu gap-1">
						<li>
							<a
								href="/"
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
				<div class="p-4 border-t border-base-300">
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
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium truncate">{data.user.nickname}</p>
							<p class="text-xs text-base-content/60 truncate">{data.user.email}</p>
						</div>
					</div>
					<button class="btn btn-ghost btn-sm w-full mt-3 justify-start" onclick={handleSignOut}>
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
		</aside>
	</div>
</div>
