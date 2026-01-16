<script lang="ts">
  import {
    getApiClient,
    type Image,
    type ImageFormat,
    type Project,
    type UploadUrl,
  } from '$lib/api';
  import { toastStore, FormField, Select, Badge } from '$lib';

  let { data } = $props();

  // Form state
  let selectedProjectId = $state('');
  let isDragging = $state(false);
  let selectedFile: File | null = $state(null);
  let previewUrl = $state('');

  // Upload state
  type UploadStep = 'idle' | 'creating_url' | 'uploading' | 'processing' | 'done' | 'error';
  let uploadStep: UploadStep = $state('idle');
  let uploadUrl: UploadUrl | null = $state(null);
  let uploadedImage: Image | null = $state(null);
  let errorMessage = $state('');

  // Derived state
  let selectedProject = $derived(
    data.projects.items.find((p: Project) => p.id === selectedProjectId)
  );

  let projectOptions = $derived(
    data.projects.items.map((p: Project) => ({
      value: p.id,
      label: p.name,
    }))
  );

  const formatMap: Record<string, ImageFormat> = {
    'image/jpeg': 'JPEG',
    'image/jpg': 'JPEG',
    'image/png': 'PNG',
    'image/webp': 'WEBP',
    'image/avif': 'AVIF',
    'image/heic': 'HEIC',
  };

  function getImageFormat(mimeType: string): ImageFormat | null {
    return formatMap[mimeType] ?? null;
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    isDragging = true;
  }

  function handleDragLeave(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;

    const file = e.dataTransfer?.files[0];
    if (file) {
      handleFile(file);
    }
  }

  function handleFileInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
      handleFile(file);
    }
  }

  function handleFile(file: File) {
    const format = getImageFormat(file.type);
    if (!format) {
      toastStore.warning(`Unsupported file type: ${file.type}`);
      return;
    }

    selectedFile = file;
    previewUrl = URL.createObjectURL(file);
    resetUploadState();
  }

  function resetUploadState() {
    uploadStep = 'idle';
    uploadUrl = null;
    uploadedImage = null;
    errorMessage = '';
  }

  function clearFile() {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    selectedFile = null;
    previewUrl = '';
    resetUploadState();
  }

  async function startUpload() {
    if (!selectedFile || !selectedProjectId) {
      toastStore.warning('Please select a project and an image');
      return;
    }

    const format = getImageFormat(selectedFile.type);
    if (!format) {
      toastStore.error('Unsupported file format');
      return;
    }

    const client = getApiClient();

    try {
      // Step 1: Create upload URL
      uploadStep = 'creating_url';
      const urlResult = await client.POST('/api/v1/projects/{projectId}/images/upload-url', {
        params: { path: { projectId: selectedProjectId } },
        body: {
          fileName: selectedFile.name,
          format,
        },
      });

      if (urlResult.error) {
        throw new Error('Failed to create upload URL');
      }

      uploadUrl = urlResult.data;

      // Step 2: Upload file to presigned URL
      uploadStep = 'uploading';
      const uploadResponse = await fetch(uploadUrl.url, {
        method: 'PUT',
        headers: uploadUrl.header,
        body: selectedFile,
      });

      if (!uploadResponse.ok) {
        throw new Error(`Upload failed: ${uploadResponse.statusText}`);
      }

      // Step 3: Wait for processing
      uploadStep = 'processing';
      const imageResult = await client.GET('/api/v1/projects/{projectId}/images/{imageId}', {
        params: {
          path: { projectId: selectedProjectId, imageId: uploadUrl.imageId },
          query: { waitUntilProcessed: true },
        },
      });

      if (imageResult.error) {
        throw new Error('Failed to get image status');
      }

      uploadedImage = imageResult.data;

      if (uploadedImage.state === 'READY') {
        uploadStep = 'done';
        toastStore.success('Image uploaded and processed successfully');
      } else {
        uploadStep = 'error';
        errorMessage = `Image processing failed with state: ${uploadedImage.state}`;
        toastStore.error(errorMessage);
      }
    } catch (err) {
      uploadStep = 'error';
      errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
      toastStore.error(errorMessage);
    }
  }

  function getStepStatus(
    step: UploadStep,
    currentStep: UploadStep
  ): 'pending' | 'active' | 'done' | 'error' {
    const steps: UploadStep[] = ['creating_url', 'uploading', 'processing', 'done'];
    const stepIndex = steps.indexOf(step);
    const currentIndex = steps.indexOf(currentStep);

    if (currentStep === 'error') {
      if (stepIndex < currentIndex) return 'done';
      if (stepIndex === currentIndex) return 'error';
      return 'pending';
    }

    // When upload is complete, mark all steps including 'done' as done
    if (currentStep === 'done') return 'done';

    if (stepIndex < currentIndex) return 'done';
    if (stepIndex === currentIndex) return 'active';
    return 'pending';
  }

  function getStateBadgeVariant(state: string): 'success' | 'warning' | 'error' | 'neutral' {
    switch (state) {
      case 'READY':
        return 'success';
      case 'UPLOAD_PENDING':
      case 'PROCESSING':
        return 'warning';
      case 'FAILED':
      case 'UPLOAD_EXPIRED':
        return 'error';
      default:
        return 'neutral';
    }
  }
</script>

<div class="space-y-6">
  <!-- Page header -->
  <div>
    <h1 class="text-2xl font-bold">Upload Test</h1>
    <p class="text-base-content/60 mt-1">Test image upload with presigned URLs</p>
  </div>

  <!-- Project selection -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body">
      <h2 class="card-title text-lg">1. Select Project</h2>
      <FormField label="Project" name="project" required>
        <Select
          name="project"
          bind:value={selectedProjectId}
          options={projectOptions}
          placeholder="Select a project"
          required
        />
      </FormField>
      {#if selectedProject}
        <div class="mt-2">
          <span class="text-base-content/60 text-sm">Presets:</span>
          <div class="mt-1 flex flex-wrap gap-1">
            {#each selectedProject.presets as preset (preset.id)}
              <Badge variant={preset.default ? 'primary' : 'neutral'} size="sm">
                {preset.name}
              </Badge>
            {/each}
            {#if selectedProject.presets.length === 0}
              <span class="text-base-content/40 text-sm">No presets configured</span>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- File drop zone -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body">
      <h2 class="card-title text-lg">2. Select Image</h2>

      {#if !selectedFile}
        <div
          role="button"
          tabindex="0"
          class="cursor-pointer rounded-lg border-2 border-dashed p-8 text-center transition-colors
            {isDragging
            ? 'border-primary bg-primary/10'
            : 'border-base-300 hover:border-primary/50'}"
          ondragover={handleDragOver}
          ondragleave={handleDragLeave}
          ondrop={handleDrop}
          onclick={() => document.getElementById('file-input')?.click()}
          onkeydown={(e) => e.key === 'Enter' && document.getElementById('file-input')?.click()}
        >
          <input
            id="file-input"
            type="file"
            accept="image/jpeg,image/png,image/webp,image/avif,image/heic"
            class="hidden"
            onchange={handleFileInput}
          />
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="text-base-content/40 mx-auto h-12 w-12"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
            />
          </svg>
          <p class="text-base-content/60 mt-4">Drag and drop an image here, or click to select</p>
          <p class="text-base-content/40 mt-1 text-sm">
            Supported formats: JPEG, PNG, WebP, AVIF, HEIC
          </p>
        </div>
      {:else}
        <div class="space-y-4">
          <div class="flex items-start gap-4">
            <div class="bg-base-200 h-32 w-32 flex-shrink-0 overflow-hidden rounded-lg">
              <img src={previewUrl} alt="Preview" class="h-full w-full object-cover" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate font-medium">{selectedFile.name}</p>
              <p class="text-base-content/60 text-sm">
                {(selectedFile.size / 1024).toFixed(1)} KB
              </p>
              <p class="text-base-content/60 text-sm">
                Type: {selectedFile.type}
              </p>
              <button
                type="button"
                class="btn btn-ghost btn-sm mt-2"
                onclick={clearFile}
                disabled={uploadStep !== 'idle' && uploadStep !== 'done' && uploadStep !== 'error'}
              >
                Change file
              </button>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Upload section -->
  <div class="card bg-base-100 shadow-sm">
    <div class="card-body">
      <h2 class="card-title text-lg">3. Upload</h2>

      <div class="space-y-4">
        <!-- Upload button -->
        <button
          type="button"
          class="btn btn-primary"
          onclick={startUpload}
          disabled={!selectedFile ||
            !selectedProjectId ||
            (uploadStep !== 'idle' && uploadStep !== 'done' && uploadStep !== 'error')}
        >
          {#if uploadStep === 'idle' || uploadStep === 'done' || uploadStep === 'error'}
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
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
              />
            </svg>
            {uploadStep === 'done' ? 'Upload Again' : 'Start Upload'}
          {:else}
            <span class="loading loading-spinner loading-sm"></span>
            Uploading...
          {/if}
        </button>

        <!-- Progress steps -->
        {#if uploadStep !== 'idle'}
          <ul class="steps steps-vertical lg:steps-horizontal w-full">
            <li
              class="step"
              class:step-primary={getStepStatus('creating_url', uploadStep) === 'done' ||
                getStepStatus('creating_url', uploadStep) === 'active'}
              class:step-error={getStepStatus('creating_url', uploadStep) === 'error'}
            >
              <span class="text-sm">Create Upload URL</span>
            </li>
            <li
              class="step"
              class:step-primary={getStepStatus('uploading', uploadStep) === 'done' ||
                getStepStatus('uploading', uploadStep) === 'active'}
              class:step-error={getStepStatus('uploading', uploadStep) === 'error'}
            >
              <span class="text-sm">Upload to S3</span>
            </li>
            <li
              class="step"
              class:step-primary={getStepStatus('processing', uploadStep) === 'done' ||
                getStepStatus('processing', uploadStep) === 'active'}
              class:step-error={getStepStatus('processing', uploadStep) === 'error'}
            >
              <span class="text-sm">Processing</span>
            </li>
            <li
              class="step"
              class:step-primary={getStepStatus('done', uploadStep) === 'done'}
              class:step-error={uploadStep === 'error'}
            >
              <span class="text-sm">{uploadStep === 'error' ? 'Failed' : 'Done'}</span>
            </li>
          </ul>
        {/if}

        <!-- Error message -->
        {#if uploadStep === 'error' && errorMessage}
          <div class="alert alert-error">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6 shrink-0 stroke-current"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span>{errorMessage}</span>
          </div>
        {/if}

        <!-- Upload URL info -->
        {#if uploadUrl}
          <div class="collapse-arrow bg-base-200 collapse">
            <input type="checkbox" />
            <div class="collapse-title font-medium">Upload URL Details</div>
            <div class="collapse-content">
              <div class="space-y-2 text-sm">
                <p><span class="font-medium">Image ID:</span> {uploadUrl.imageId}</p>
                <p>
                  <span class="font-medium">Expires At:</span>
                  {new Date(uploadUrl.expiresAt).toLocaleString()}
                </p>
                <p class="break-all"><span class="font-medium">URL:</span> {uploadUrl.url}</p>
              </div>
            </div>
          </div>
        {/if}

        <!-- Result -->
        {#if uploadedImage}
          <div class="card bg-base-200">
            <div class="card-body">
              <h3 class="card-title text-base">Upload Result</h3>
              <div class="space-y-3">
                <div class="flex items-center gap-2">
                  <span class="font-medium">State:</span>
                  <Badge variant={getStateBadgeVariant(uploadedImage.state)}>
                    {uploadedImage.state}
                  </Badge>
                </div>
                <p><span class="font-medium">Image ID:</span> {uploadedImage.id}</p>
                <p><span class="font-medium">Format:</span> {uploadedImage.format}</p>
                <p>
                  <span class="font-medium">Original URL:</span>
                  <a
                    href={uploadedImage.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="link link-primary break-all"
                  >
                    {uploadedImage.url}
                  </a>
                </p>

                {#if uploadedImage.variants && uploadedImage.variants.length > 0}
                  <div>
                    <p class="mb-2 font-medium">Variants:</p>
                    <div class="overflow-x-auto">
                      <table class="table-sm table">
                        <thead>
                          <tr>
                            <th>Preset</th>
                            <th>Format</th>
                            <th>State</th>
                            <th>URL</th>
                          </tr>
                        </thead>
                        <tbody>
                          {#each uploadedImage.variants as variant (variant.id)}
                            <tr>
                              <td>{variant.presetName}</td>
                              <td>{variant.format}</td>
                              <td>
                                <Badge variant={getStateBadgeVariant(variant.state)} size="sm">
                                  {variant.state}
                                </Badge>
                              </td>
                              <td>
                                <a
                                  href={variant.url}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  class="link link-primary text-xs break-all"
                                >
                                  View
                                </a>
                              </td>
                            </tr>
                          {/each}
                        </tbody>
                      </table>
                    </div>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
