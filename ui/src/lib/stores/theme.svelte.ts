const STORAGE_KEY = 'theme';
type Theme = 'light' | 'dark';

function createThemeStore() {
  let theme = $state<Theme>('light');

  function getSystemTheme(): Theme {
    if (typeof window === 'undefined') return 'light';
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function applyTheme(t: Theme) {
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-theme', t);
    }
  }

  function init() {
    if (typeof window === 'undefined') return;

    // Try to load from localStorage
    const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
    if (stored === 'light' || stored === 'dark') {
      theme = stored;
    } else {
      // Fall back to OS preference
      theme = getSystemTheme();
    }
    applyTheme(theme);

    // Listen for OS theme changes (only if no stored preference)
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (!localStorage.getItem(STORAGE_KEY)) {
        theme = e.matches ? 'dark' : 'light';
        applyTheme(theme);
      }
    });
  }

  function toggle() {
    theme = theme === 'light' ? 'dark' : 'light';
    localStorage.setItem(STORAGE_KEY, theme);
    applyTheme(theme);
  }

  function set(t: Theme) {
    theme = t;
    localStorage.setItem(STORAGE_KEY, theme);
    applyTheme(theme);
  }

  return {
    get current() {
      return theme;
    },
    init,
    toggle,
    set,
  };
}

export const themeStore = createThemeStore();
