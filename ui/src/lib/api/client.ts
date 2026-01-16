import createClient, { type Middleware } from 'openapi-fetch';
import type { paths, components } from './schema';

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
		fetch: customFetch
	});

	// Add middleware for handling credentials
	const credentialsMiddleware: Middleware = {
		async onRequest({ request }) {
			return new Request(request, { credentials: 'include' });
		}
	};

	client.use(credentialsMiddleware);

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
		throw new Error('getApiClient() should only be called in browser context. Use createApiClient() with fetch in SSR.');
	}

	if (!browserClient) {
		browserClient = createApiClient();
	}

	return browserClient;
}
