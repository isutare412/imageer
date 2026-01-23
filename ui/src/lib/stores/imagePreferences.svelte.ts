import { browser } from '$app/environment';

type SortBy = 'createdAt' | 'updatedAt';
type SortOrder = 'ASC' | 'DESC';
type PerPage = 10 | 20 | 50;

interface ImagePreferences {
  sortBy: SortBy;
  sortOrder: SortOrder;
  perPage: PerPage;
}

const STORAGE_KEY_PREFIX = 'imageer:images:';

function createImagePreferencesStore() {
  const defaults: ImagePreferences = {
    sortBy: 'createdAt',
    sortOrder: 'DESC',
    perPage: 20,
  };

  function load<K extends keyof ImagePreferences>(key: K): ImagePreferences[K] {
    if (!browser) return defaults[key];
    const stored = localStorage.getItem(`${STORAGE_KEY_PREFIX}${key}`);
    if (!stored) return defaults[key];
    // Parse perPage as number since localStorage stores strings
    if (key === 'perPage') {
      return parseInt(stored, 10) as ImagePreferences[K];
    }
    return stored as ImagePreferences[K];
  }

  function save<K extends keyof ImagePreferences>(key: K, value: ImagePreferences[K]) {
    if (!browser) return;
    localStorage.setItem(`${STORAGE_KEY_PREFIX}${key}`, String(value));
  }

  let sortBy = $state<SortBy>(load('sortBy'));
  let sortOrder = $state<SortOrder>(load('sortOrder'));
  let perPage = $state<PerPage>(load('perPage'));

  return {
    get sortBy() {
      return sortBy;
    },
    set sortBy(value: SortBy) {
      sortBy = value;
      save('sortBy', value);
    },
    get sortOrder() {
      return sortOrder;
    },
    set sortOrder(value: SortOrder) {
      sortOrder = value;
      save('sortOrder', value);
    },
    get perPage() {
      return perPage;
    },
    set perPage(value: PerPage) {
      perPage = value;
      save('perPage', value);
    },
  };
}

export const imagePreferencesStore = createImagePreferencesStore();
