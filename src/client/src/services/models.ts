// Centralized API calls for model management.
// Moved from api.ts and expanded with clear method names.
import type { Model, ModelDetail, GlobalStats } from "../types/models";

const BASE = "/api/models";

export async function getAllModels(pipeline: string): Promise<Model[]> {
  const res = await fetch(`${BASE}?pipeline=${encodeURIComponent(pipeline)}`);
  if (!res.ok) throw new Error("Failed to fetch models");
  return res.json();
}

export async function getActiveModels(): Promise<Model[]> {
  // Placeholder: backend endpoint not implemented
  const res = await fetch(`${BASE}?status=active`);
  if (!res.ok) throw new Error("Failed to fetch active models");
  return res.json();
}

export async function getModelById(id: string): Promise<ModelDetail> {
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}`);
  if (!res.ok) throw new Error("Failed to fetch model");
  return res.json();
}

export async function getModelStats(id: string): Promise<ModelDetail> {
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/stats`);
  if (!res.ok) throw new Error("Failed to fetch model stats");
  return res.json();
}

export async function activateModel(id: string): Promise<void> {
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/activate`, {
    method: "POST",
  });
  if (!res.ok) throw new Error("Failed to activate model");
}

export async function deactivateModel(id: string): Promise<void> {
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/deactivate`, {
    method: "POST",
  });
  if (!res.ok) throw new Error("Failed to deactivate model");
}

export async function getGlobalStats(): Promise<GlobalStats> {
  const res = await fetch(`${BASE}/stats/global`);
  if (!res.ok) throw new Error("Failed to fetch global stats");
  return res.json();
}

export async function syncModels(): Promise<void> {
  const res = await fetch(`${BASE}/sync`, { method: "POST" });
  if (!res.ok) throw new Error("Failed to sync models");
}

export async function refreshModels(pipeline: string): Promise<void> {
  const res = await fetch(
    `${BASE}/refresh?pipeline=${encodeURIComponent(pipeline)}`,
    { method: "POST" },
  );
  if (!res.ok) throw new Error("Failed to refresh models");
}

export async function enableModel(id: string): Promise<void> {
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/enable`, {
    method: "POST",
  });
  if (!res.ok) throw new Error("Failed to enable model");
}

export async function downloadModel(
  id: string,
  onProgress?: (pct: number) => void,
): Promise<void> {
  return new Promise((resolve, reject) => {
    const es = new EventSource(`${BASE}/${encodeURIComponent(id)}/download`);
    es.onmessage = (e) => {
      const pct = parseInt(e.data);
      if (onProgress) onProgress(pct);
      if (pct >= 100) {
        es.close();
        resolve();
      }
    };
    es.addEventListener("error", () => {
      es.close();
      reject(new Error("Download failed"));
    });
    es.addEventListener("done", () => {
      es.close();
      resolve();
    });
  });
}
