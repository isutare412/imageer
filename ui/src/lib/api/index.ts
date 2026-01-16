export { createApiClient, getApiClient, type ApiClient, type ClientOptions } from './client';
export type {
  // Schema types
  Schemas,
  Project,
  Projects,
  Preset,
  Image,
  ImageVariant,
  UploadUrl,
  User,
  ServiceAccount,
  ServiceAccountWithApiKey,
  ServiceAccounts,
  AppError,
  // Enums
  ImageState,
  ImageVariantState,
  ImageFormat,
  ImageFit,
  ImageAnchor,
  UserRole,
  ServiceAccountAccessScope,
  // Request types
  CreateProjectRequest,
  UpdateProjectRequest,
  CreatePresetRequest,
  UpsertPresetRequest,
  CreateUploadUrlRequest,
  CreateServiceAccountRequest,
  UpdateServiceAccountRequest,
  ReprocessImagesRequest,
} from './client';

export { ApiError, isAppError, unwrap, unwrapEmpty } from './error';
