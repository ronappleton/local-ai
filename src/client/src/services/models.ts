// Centralized API calls for model management.
// Moved from api.ts and expanded with clear method names.
import type { Model, ModelDetail, GlobalStats } from "../types/models";

const BASE = "/api/models";

export async function getAllModels(pipeline: string): Promise<Model[]> {
  console.log("[API] GET", `${BASE}?pipeline=${encodeURIComponent(pipeline)}`);
  const res = await fetch(`${BASE}?pipeline=${encodeURIComponent(pipeline)}`);
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to fetch models");
  return res.json();
}

export async function getActiveModels(): Promise<Model[]> {
  // Placeholder: backend endpoint not implemented
  console.log("[API] GET", `${BASE}?status=active`);
  const res = await fetch(`${BASE}?status=active`);
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to fetch active models");
  return res.json();
}

export async function getModelById(id: string): Promise<ModelDetail> {
  console.log("[API] GET", `${BASE}/${encodeURIComponent(id)}`);
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}`);
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to fetch model");
  return res.json();
}

export async function getModelStats(id: string): Promise<ModelDetail> {
  console.log("[API] GET", `${BASE}/${encodeURIComponent(id)}/stats`);
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/stats`);
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to fetch model stats");
  return res.json();
}

export async function activateModel(id: string): Promise<void> {
  console.log("[API] POST", `${BASE}/${encodeURIComponent(id)}/activate`);
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/activate`, {
    method: "POST",
  });
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to activate model");
}

export async function deactivateModel(id: string): Promise<void> {
  console.log("[API] POST", `${BASE}/${encodeURIComponent(id)}/deactivate`);
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/deactivate`, {
    method: "POST",
  });
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to deactivate model");
}

export async function getGlobalStats(): Promise<GlobalStats> {
  console.log("[API] GET", `${BASE}/stats/global`);
  const res = await fetch(`${BASE}/stats/global`);
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to fetch global stats");
  return res.json();
}

export async function syncModels(): Promise<void> {
  console.log("[API] POST", `${BASE}/sync`);
  const res = await fetch(`${BASE}/sync`, { method: "POST" });
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to sync models");
}

export async function refreshModels(pipeline: string): Promise<void> {
  const url = `${BASE}/refresh?pipeline=${encodeURIComponent(pipeline)}`;
  console.log("[API] POST", url);
  const res = await fetch(url, { method: "POST" });
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to refresh models");
}

export async function enableModel(id: string): Promise<void> {
  console.log("[API] POST", `${BASE}/${encodeURIComponent(id)}/enable`);
  const res = await fetch(`${BASE}/${encodeURIComponent(id)}/enable`, {
    method: "POST",
  });
  console.log("[API] response", res.status);
  if (!res.ok) throw new Error("Failed to enable model");
}

export async function downloadModel(
  id: string,
  onProgress?: (pct: number) => void,
): Promise<void> {
  return new Promise((resolve, reject) => {
    console.log("[API] SSE", `${BASE}/${encodeURIComponent(id)}/download`);
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
      console.error("[API] download error");
      es.close();
      reject(new Error("Download failed"));
    });
    es.addEventListener("done", () => {
      es.close();
      resolve();
    });
  });
}
