import type { AppError } from './client';

/**
 * API Error class that wraps error responses from the server.
 */
export class ApiError extends Error {
	readonly codeId: number;
	readonly codeName: string;
	readonly status: number;

	constructor(error: AppError, status: number) {
		super(error.message);
		this.name = 'ApiError';
		this.codeId = error.codeId;
		this.codeName = error.codeName;
		this.status = status;
	}
}

/**
 * Type guard to check if a value is an AppError.
 */
export function isAppError(value: unknown): value is AppError {
	return (
		typeof value === 'object' &&
		value !== null &&
		'message' in value &&
		'codeId' in value &&
		'codeName' in value
	);
}

/**
 * Unwraps an API response, throwing ApiError if the request failed.
 */
export function unwrap<T>(result: { data?: T; error?: AppError; response: Response }): T {
	if (result.error) {
		throw new ApiError(result.error, result.response.status);
	}

	if (result.data === undefined) {
		throw new Error('Unexpected empty response');
	}

	return result.data;
}

/**
 * Unwraps an API response that may have no content (e.g., DELETE operations).
 */
export function unwrapEmpty(result: { data?: unknown; error?: AppError; response: Response }): void {
	if (result.error) {
		throw new ApiError(result.error, result.response.status);
	}
}
