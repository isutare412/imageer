<script lang="ts">
  // Dashboard page - shown after login
</script>

<div class="space-y-6">
  <!-- Page header -->
  <div>
    <h1 class="text-2xl font-bold">Dashboard</h1>
    <p class="text-base-content/60 mt-1">Welcome to Imageer administration panel</p>
  </div>

  <!-- Stats cards -->
  <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-base-content/70 text-base font-medium">Projects</h2>
        <p class="text-3xl font-bold">-</p>
      </div>
    </div>

    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-base-content/70 text-base font-medium">Service Accounts</h2>
        <p class="text-3xl font-bold">-</p>
      </div>
    </div>

    <div class="card bg-base-100 shadow-sm">
      <div class="card-body">
        <h2 class="card-title text-base-content/70 text-base font-medium">Total Images</h2>
        <p class="text-3xl font-bold">-</p>
      </div>
    </div>
  </div>

  <!-- Quick actions -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body">
      <h2 class="card-title text-lg">Quick Actions</h2>
      <div class="mt-2 flex flex-wrap gap-2">
        <a href="/projects" class="btn btn-primary btn-sm">
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
              d="M12 4v16m8-8H4"
            />
          </svg>
          New Project
        </a>
        <a href="/service-accounts" class="btn btn-outline btn-sm">
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
              d="M12 4v16m8-8H4"
            />
          </svg>
          New Service Account
        </a>
      </div>
    </div>
  </div>
</div>
