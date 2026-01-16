<script lang="ts">
  interface Props {
    total: number;
    limit: number;
    offset: number;
    onPageChange: (offset: number) => void;
  }

  let { total, limit, offset, onPageChange }: Props = $props();

  let currentPage = $derived(Math.floor(offset / limit) + 1);
  let totalPages = $derived(Math.ceil(total / limit));

  function goToPage(page: number) {
    if (page < 1 || page > totalPages) return;
    onPageChange((page - 1) * limit);
  }

  function getVisiblePages(): (number | 'ellipsis')[] {
    const pages: (number | 'ellipsis')[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (currentPage > 3) {
        pages.push('ellipsis');
      }

      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push('ellipsis');
      }

      pages.push(totalPages);
    }

    return pages;
  }
</script>

{#if totalPages > 1}
  <div class="flex items-center justify-between">
    <div class="text-base-content/60 text-sm">
      Showing {offset + 1} - {Math.min(offset + limit, total)} of {total}
    </div>
    <div class="join">
      <button
        type="button"
        class="join-item btn btn-sm"
        disabled={currentPage === 1}
        onclick={() => goToPage(currentPage - 1)}
        aria-label="Previous page"
      >
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
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>
      {#each getVisiblePages() as page}
        {#if page === 'ellipsis'}
          <button type="button" class="join-item btn btn-sm btn-disabled">...</button>
        {:else}
          <button
            type="button"
            class="join-item btn btn-sm {currentPage === page ? 'btn-active' : ''}"
            onclick={() => goToPage(page)}
          >
            {page}
          </button>
        {/if}
      {/each}
      <button
        type="button"
        class="join-item btn btn-sm"
        disabled={currentPage === totalPages}
        onclick={() => goToPage(currentPage + 1)}
        aria-label="Next page"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
{/if}
