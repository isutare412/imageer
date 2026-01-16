// Admin-related types

export interface PresetData {
  id?: string;
  name: string;
  default: boolean;
  format?: string;
  quality?: number;
  fit?: string;
  anchor?: string;
  width?: number;
  height?: number;
}
