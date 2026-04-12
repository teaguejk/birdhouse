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

export async function updateDevice(token: string, id: string, data: { name?: string; location?: string; active?: boolean }) {
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

// public
export async function getLatestImage() {
  return fetchAPI("/uploads/latest");
}
