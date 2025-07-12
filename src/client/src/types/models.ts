export interface Model {
  modelId: string;
  lastModified: string;
  downloads: number;
  tags: string[];
}

export interface ModelStats {
  // placeholder for per-model stats
  [key: string]: unknown;
}

export interface GlobalStats {
  // placeholder for global stats
  [key: string]: unknown;
}
