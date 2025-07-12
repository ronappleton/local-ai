export interface User {
  id: number;
  username: string;
  email?: string;
  admin?: boolean | number;
  verified?: boolean | number;
}

export async function fetchUsers(): Promise<User[]> {
  const res = await fetch("/api/users");
  if (!res.ok) {
    throw new Error("Failed to fetch users");
  }
  return res.json();
}

export interface ModelInfo {
  modelId: string;
  lastModified: string;
  downloads: number;
  tags: string[];
}

export async function fetchModels(pipeline: string): Promise<ModelInfo[]> {
  const res = await fetch(
    "/api/models?pipeline=" + encodeURIComponent(pipeline),
  );
  if (!res.ok) {
    throw new Error("Failed to fetch models");
  }
  return res.json();
}

export interface Project {
  id: number;
  name: string;
  created: string;
}

export interface Message {
  id: number;
  projectId: number;
  role: string;
  content: string;
  timestamp: string;
}

export async function fetchProjects(): Promise<Project[]> {
  const res = await fetch("/api/projects");
  if (!res.ok) {
    throw new Error("Failed to fetch projects");
  }
  return res.json();
}

export async function createProject(name: string): Promise<Project> {
  const res = await fetch("/api/projects", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({name}),
  });
  if (!res.ok) {
    throw new Error("Failed to create project");
  }
  return res.json();
}

export async function addMessage(projectId: number, role: string, content: string): Promise<Message> {
  const res = await fetch("/api/messages", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({projectId, role, content}),
  });
  if (!res.ok) {
    throw new Error("Failed to add message");
  }
  return res.json();
}

export async function createUser(username: string, email: string, password: string): Promise<User> {
  const res = await fetch("/api/users", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({username, email, password}),
  });
  if (!res.ok) {
    throw new Error("Failed to create user");
  }
  return res.json();
}

export async function promoteUser(username: string): Promise<void> {
  const res = await fetch("/api/users/" + encodeURIComponent(username) + "/promote", {
    method: "POST"
  });
  if (!res.ok) {
    throw new Error("Failed to promote user");
  }
}

export async function downloadModel(modelId: string): Promise<void> {
  const res = await fetch("/api/models/" + encodeURIComponent(modelId) + "/download", {
    method: "POST"
  });
  if (!res.ok) {
    throw new Error("Failed to download model");
  }
}

export async function setActiveModel(modelId: string): Promise<void> {
  const res = await fetch("/api/models/" + encodeURIComponent(modelId) + "/activate", {
    method: "POST"
  });
  if (!res.ok) {
    throw new Error("Failed to set active model");
  }
}

export async function getModelStatus(): Promise<{ activeModel: string }> {
  const res = await fetch("/api/models/status");
  if (!res.ok) {
    throw new Error("Failed to get model status");
  }
  return res.json();
}
