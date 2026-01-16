<script lang="ts">
  import type { Snippet } from 'svelte';

  type BadgeVariant =
    | 'neutral'
    | 'primary'
    | 'secondary'
    | 'accent'
    | 'info'
    | 'success'
    | 'warning'
    | 'error';
  type BadgeSize = 'xs' | 'sm' | 'md' | 'lg';

  interface Props {
    variant?: BadgeVariant;
    size?: BadgeSize;
    outline?: boolean;
    truncate?: boolean;
    children: Snippet;
  }

  let {
    variant = 'neutral',
    size = 'md',
    outline = false,
    truncate = false,
    children,
  }: Props = $props();

  const variantClasses: Record<BadgeVariant, string> = {
    neutral: 'badge-neutral',
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    accent: 'badge-accent',
    info: 'badge-info',
    success: 'badge-success',
    warning: 'badge-warning',
    error: 'badge-error',
  };

  const sizeClasses: Record<BadgeSize, string> = {
    xs: 'badge-xs',
    sm: 'badge-sm',
    md: 'badge-md',
    lg: 'badge-lg',
  };
</script>

<span
  class="badge {variantClasses[variant]} {sizeClasses[size]} {outline
    ? 'badge-outline'
    : ''} {truncate ? 'max-w-32 truncate' : ''}"
>
  {@render children()}
</span>
