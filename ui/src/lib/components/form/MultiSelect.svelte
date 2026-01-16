<script lang="ts">
  interface Option {
    value: string;
    label: string;
  }

  interface Props {
    name: string;
    values: string[];
    options: Option[];
    placeholder?: string;
    disabled?: boolean;
    onchange?: (values: string[]) => void;
  }

  let {
    name,
    values = $bindable([]),
    options,
    placeholder = 'Select items...',
    disabled = false,
    onchange,
  }: Props = $props();

  let isOpen = $state(false);
  let searchQuery = $state('');

  let filteredOptions = $derived(
    options.filter((opt) => opt.label.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  let selectedLabels = $derived(
    values
      .map((v) => options.find((o) => o.value === v)?.label)
      .filter(Boolean)
      .join(', ')
  );

  function toggleOption(optionValue: string) {
    if (values.includes(optionValue)) {
      values = values.filter((v) => v !== optionValue);
    } else {
      values = [...values, optionValue];
    }
    onchange?.(values);
  }

  function removeValue(optionValue: string) {
    values = values.filter((v) => v !== optionValue);
    onchange?.(values);
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      isOpen = false;
    }
  }
</script>

<div class="dropdown w-full" class:dropdown-open={isOpen}>
  <button
    type="button"
    class="input input-bordered flex w-full cursor-pointer items-center justify-between"
    {disabled}
    onclick={() => (isOpen = !isOpen)}
    onkeydown={handleKeydown}
  >
    <span class="truncate {values.length === 0 ? 'text-base-content/50' : ''}">
      {values.length === 0 ? placeholder : selectedLabels}
    </span>
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-4 w-4 shrink-0 transition-transform {isOpen ? 'rotate-180' : ''}"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
    </svg>
  </button>

  {#if isOpen}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="dropdown-content bg-base-100 border-base-300 z-50 mt-1 max-h-60 w-full overflow-auto rounded-lg border shadow-lg"
      onkeydown={handleKeydown}
    >
      <div class="bg-base-100 sticky top-0 p-2">
        <input
          type="text"
          class="input input-bordered input-sm w-full"
          placeholder="Search..."
          bind:value={searchQuery}
        />
      </div>
      <ul class="menu p-2 pt-0">
        {#each filteredOptions as option}
          <li>
            <label class="flex cursor-pointer items-center gap-2">
              <input
                type="checkbox"
                class="checkbox checkbox-sm"
                checked={values.includes(option.value)}
                onchange={() => toggleOption(option.value)}
              />
              <span>{option.label}</span>
            </label>
          </li>
        {:else}
          <li class="text-base-content/50 p-2 text-center text-sm">No options found</li>
        {/each}
      </ul>
    </div>
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-40" onclick={() => (isOpen = false)}></div>
  {/if}
</div>

<!-- Hidden inputs for form submission -->
{#each values as val}
  <input type="hidden" {name} value={val} />
{/each}

<!-- Selected items as chips -->
{#if values.length > 0}
  <div class="mt-2 flex flex-wrap gap-1">
    {#each values as val}
      {@const option = options.find((o) => o.value === val)}
      {#if option}
        <span class="badge badge-neutral gap-1">
          {option.label}
          <button
            type="button"
            class="hover:text-error"
            onclick={() => removeValue(val)}
            aria-label="Remove {option.label}"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-3 w-3"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </span>
      {/if}
    {/each}
  </div>
{/if}
