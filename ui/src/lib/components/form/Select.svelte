<script lang="ts">
  interface Option {
    value: string;
    label: string;
    disabled?: boolean;
  }

  interface Props {
    name: string;
    value?: string;
    options: Option[];
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    onchange?: (value: string) => void;
  }

  let {
    name,
    value = $bindable(''),
    options,
    placeholder,
    disabled = false,
    required = false,
    onchange,
  }: Props = $props();

  function handleChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    value = target.value;
    onchange?.(value);
  }
</script>

<select
  id={name}
  {name}
  class="select select-bordered w-full"
  {disabled}
  {required}
  {value}
  onchange={handleChange}
>
  {#if placeholder}
    <option value="" disabled>{placeholder}</option>
  {/if}
  {#each options as option}
    <option value={option.value} disabled={option.disabled}>
      {option.label}
    </option>
  {/each}
</select>
