export interface Model {
  modelId: string;
  lastModified: string;
  downloads: number;
  tags: string[];
}

export interface ModelDetail extends Model {
  sha: string;
  files: string[];
}

export interface ModelStats {
  // placeholder for per-model stats
  [key: string]: unknown;
}

export interface GlobalStats {
  total_models: number;
}
