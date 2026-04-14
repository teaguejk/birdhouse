const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8090";

async function fetchAPI(path: string, options: RequestInit = {}) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `API error: ${res.status}`);
  }

  if (res.status === 204) return null;
  return res.json();
}

function authHeaders(token: string) {
  return { Authorization: `Bearer ${token}` };
}

// auth
export async function getMe(token: string) {
  return fetchAPI("/auth/me", { headers: authHeaders(token) });
}

// devices
export async function getDevices(token: string) {
  return fetchAPI("/devices/", { headers: authHeaders(token) });
}

export async function getDevice(token: string, id: string) {
  return fetchAPI(`/devices/${id}`, { headers: authHeaders(token) });
}

export async function createDevice(token: string, data: { name: string; location: string }) {
  return fetchAPI("/devices/", {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
}

export async function updateDevice(token: string, id: string, data: {
  name?: string;
  location?: string;
  active?: boolean;
  config?: { min_contour_area: number; threshold: number; cooldown_seconds: number };
}) {
  return fetchAPI(`/devices/${id}`, {
    method: "PUT",
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
}

export async function deleteDevice(token: string, id: string) {
  return fetchAPI(`/devices/${id}`, {
    method: "DELETE",
    headers: authHeaders(token),
  });
}

export async function rotateDeviceKey(token: string, id: string) {
  return fetchAPI(`/devices/${id}/rotate-key`, {
    method: "POST",
    headers: authHeaders(token),
  });
}

// device status (public)
export async function getDeviceStatuses() {
  return fetchAPI("/devices/status");
}

// commands
export async function sendCommand(token: string, deviceId: string, action: string, payload?: Record<string, unknown>) {
  return fetchAPI(`/commands/device/${deviceId}`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ action, payload }),
  });
}

// public
export async function getLatestImage() {
  return fetchAPI("/uploads/latest");
}
