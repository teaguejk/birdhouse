import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8090";

interface DeviceLastStatus {
  detecting?: boolean;
  uptime_seconds?: number;
  captures?: number;
  uploads?: { success: number; failed: number };
}

interface DeviceStatus {
  id: string;
  name: string;
  location: string;
  active: boolean;
  online: boolean;
  last_seen_at: string | null;
  last_status: DeviceLastStatus | null;
}

function timeAgo(dateStr: string): string {
  const seconds = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}

interface DeviceStatusPanelProps {
  selectedDeviceId: string | null;
  onSelectDevice: (id: string | null) => void;
}

export function DeviceStatusPanel({ selectedDeviceId, onSelectDevice }: DeviceStatusPanelProps) {
  const [devices, setDevices] = useState<DeviceStatus[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE_URL}/devices/status`)
      .then((res) => res.json())
      .then((data) => setDevices(data || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle className="text-lg">Devices</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (devices.length === 0) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle className="text-lg">Devices</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No devices registered</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="text-lg">Devices</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {devices.map((device) => (
          <button
            key={device.id}
            onClick={() => onSelectDevice(selectedDeviceId === device.id ? null : device.id)}
            className={cn(
              "flex w-full items-center justify-between rounded-md border px-4 py-3 text-left transition-colors",
              selectedDeviceId === device.id
                ? "border-primary bg-accent"
                : "hover:bg-accent/50"
            )}
          >
            <div>
              <p className="font-medium">{device.name}</p>
              <div className="flex items-center gap-2">
                {device.location && (
                  <span className="text-xs text-muted-foreground">{device.location}</span>
                )}
                {device.online && device.last_status?.detecting != null && (
                  <span className="text-xs text-muted-foreground">
                    {device.last_status.detecting ? "Detecting" : "Paused"}
                  </span>
                )}
                {device.online && device.last_status?.captures != null && (
                  <span className="text-xs text-muted-foreground">
                    {device.last_status.captures} captures
                  </span>
                )}
              </div>
            </div>
            <div className="flex items-center gap-2">
              {device.last_seen_at && (
                <span className="text-xs text-muted-foreground">
                  {timeAgo(device.last_seen_at)}
                </span>
              )}
              <Badge variant={device.online ? "default" : "secondary"}>
                {device.online ? "Online" : device.active ? "Offline" : "Inactive"}
              </Badge>
            </div>
          </button>
        ))}
      </CardContent>
    </Card>
  );
}
