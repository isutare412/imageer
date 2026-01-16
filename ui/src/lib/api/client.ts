import createClient, { type Middleware } from 'openapi-fetch';
import type { paths, components } from './schema';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { toastStore } from '$lib/stores/toast.svelte';

export type ApiClient = ReturnType<typeof createClient<paths>>;

export type Schemas = components['schemas'];

// Re-export commonly used types
export type Project = Schemas['Project'];
export type Projects = Schemas['Projects'];
export type Preset = Schemas['Preset'];
export type Image = Schemas['Image'];
export type ImageVariant = Schemas['ImageVariant'];
export type UploadUrl = Schemas['UploadUrl'];
export type User = Schemas['User'];
export type ServiceAccount = Schemas['ServiceAccount'];
export type ServiceAccountWithApiKey = Schemas['ServiceAccountWithApiKey'];
export type ServiceAccounts = Schemas['ServiceAccounts'];
export type AppError = Schemas['AppError'];

// Enums
export type ImageState = Schemas['ImageState'];
export type ImageVariantState = Schemas['ImageVariantState'];
export type ImageFormat = Schemas['ImageFormat'];
export type ImageFit = Schemas['ImageFit'];
export type ImageAnchor = Schemas['ImageAnchor'];
export type UserRole = Schemas['UserRole'];
export type ServiceAccountAccessScope = Schemas['ServiceAccountAccessScope'];

// Request types
export type CreateProjectRequest = Schemas['CreateProjectAdminRequest'];
export type UpdateProjectRequest = Schemas['UpdateProjectAdminRequest'];
export type CreatePresetRequest = Schemas['CreatePresetRequest'];
export type UpsertPresetRequest = Schemas['UpsertPresetRequest'];
export type CreateUploadUrlRequest = Schemas['CreateUploadUrlRequest'];
export type CreateServiceAccountRequest = Schemas['CreateServiceAccountAdminRequest'];
export type UpdateServiceAccountRequest = Schemas['UpdateServiceAccountAdminRequest'];
export type ReprocessImagesRequest = Schemas['ReprocessImagesAdminRequest'];

export interface ClientOptions {
  baseUrl?: string;
  fetch?: typeof fetch;
}

/**
 * Creates an API client instance.
 * In SvelteKit, pass the `fetch` from load functions for proper SSR handling.
 */
export function createApiClient(options: ClientOptions = {}): ApiClient {
  const { baseUrl = '', fetch: customFetch } = options;

  const client = createClient<paths>({
    baseUrl,
    fetch: customFetch,
  });

  // Add middleware for handling credentials
  const credentialsMiddleware: Middleware = {
    async onRequest({ request }) {
      return new Request(request, { credentials: 'include' });
    },
  };

  // Add middleware for error toast notifications (browser only)
  const errorToastMiddleware: Middleware = {
    async onResponse({ response }) {
      if (typeof window !== 'undefined' && !response.ok) {
        const status = response.status;
        const body: AppError | null = await response
          .clone()
          .json()
          .catch(() => null);
        const message = body?.message ?? response.statusText;

        if (status >= 500) {
          toastStore.error(message);
        } else if (status >= 400) {
          toastStore.warning(message);
        }
      }
      return response;
    },
  };

  client.use(credentialsMiddleware);
  client.use(errorToastMiddleware);

  return client;
}

// Default client for browser usage
let browserClient: ApiClient | null = null;

/**
 * Gets the browser API client singleton.
 * Only use this in browser context, not in SSR.
 */
export function getApiClient(): ApiClient {
  if (typeof window === 'undefined') {
    throw new Error(
      'getApiClient() should only be called in browser context. Use createApiClient() with fetch in SSR.'
    );
  }

  if (!browserClient) {
    browserClient = createApiClient({ baseUrl: PUBLIC_API_BASE_URL });
  }

  return browserClient;
}
