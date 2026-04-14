import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { getDevice, getDeviceStatuses, updateDevice, deleteDevice, rotateDeviceKey, sendCommand } from "@/lib/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface DeviceLastStatus {
  detecting?: boolean;
  uptime_seconds?: number;
  captures?: number;
  uploads?: { success: number; failed: number };
}

interface DeviceConfig {
  min_contour_area: number;
  threshold: number;
  cooldown_seconds: number;
}

interface Device {
  id: string;
  name: string;
  location: string;
  active: boolean;
  config: DeviceConfig;
  created_at: string;
  updated_at: string;
}

export default function DeviceDetail() {
  const { id } = useParams<{ id: string }>();
  const { token } = useAuth();
  const navigate = useNavigate();

  const [device, setDevice] = useState<Device | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [name, setName] = useState("");
  const [location, setLocation] = useState("");
  const [active, setActive] = useState(true);

  const [deviceStatus, setDeviceStatus] = useState<DeviceLastStatus | null>(null);
  const [online, setOnline] = useState(false);

  const [newKey, setNewKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [sending, setSending] = useState(false);

  useEffect(() => {
    if (!token || !id) return;
    getDevice(token, id)
      .then((d) => {
        setDevice(d);
        setName(d.name);
        setLocation(d.location);
        setActive(d.active);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [token, id]);

  useEffect(() => {
    if (!id) return;
    const fetchStatus = () =>
      getDeviceStatuses()
        .then((statuses: { id: string; online: boolean; last_status: DeviceLastStatus | null }[]) => {
          const match = statuses?.find((s: { id: string }) => s.id === id);
          if (match) {
            setOnline(match.online);
            setDeviceStatus(match.last_status);
          }
        })
        .catch(console.error);

    fetchStatus();
    const interval = setInterval(fetchStatus, 15000);
    return () => clearInterval(interval);
  }, [id]);

  const handleSave = async () => {
    if (!token || !id) return;
    setSaving(true);
    try {
      const updated = await updateDevice(token, id, { name, location, active });
      setDevice(updated);
    } catch (err) {
      console.error("Failed to update device:", err);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!token || !id) return;
    try {
      await deleteDevice(token, id);
      navigate("/admin");
    } catch (err) {
      console.error("Failed to delete device:", err);
    }
  };

  const handleRotateKey = async () => {
    if (!token || !id) return;
    try {
      const result = await rotateDeviceKey(token, id);
      setNewKey(result.api_key);
    } catch (err) {
      console.error("Failed to rotate key:", err);
    }
  };

  const handleCommand = async (action: string, payload?: Record<string, unknown>) => {
    if (!token || !id) return;
    setSending(true);
    try {
      await sendCommand(token, id, action, payload);
    } catch (err) {
      console.error("Failed to send command:", err);
    } finally {
      setSending(false);
    }
  };

  const handleCopy = async () => {
    if (newKey) {
      await navigator.clipboard.writeText(newKey);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  if (loading) {
    return (
      <div className="mx-auto max-w-lg space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!device) {
    return <p className="text-muted-foreground">Device not found.</p>;
  }

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{device.name}</h1>
        <Button variant="outline" size="sm" onClick={() => navigate("/admin")}>
          Back
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Device Settings</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input id="name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="location">Location</Label>
            <Input id="location" value={location} onChange={(e) => setLocation(e.target.value)} />
          </div>
          <div className="flex items-center justify-between">
            <Label htmlFor="active">Active</Label>
            <Switch id="active" checked={active} onCheckedChange={setActive} />
          </div>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save Changes"}
          </Button>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Controls</CardTitle>
            <div className="flex items-center gap-2">
              <Badge variant={online ? "default" : "secondary"}>
                {online ? "Online" : "Offline"}
              </Badge>
              {online && deviceStatus?.detecting != null && (
                <Badge variant={deviceStatus.detecting ? "default" : "outline"}>
                  {deviceStatus.detecting ? "Detecting" : "Paused"}
                </Badge>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {online && deviceStatus && (
            <div className="grid grid-cols-3 gap-4 rounded-md bg-muted p-3 text-center text-sm">
              <div>
                <p className="text-muted-foreground">Captures</p>
                <p className="font-medium">{deviceStatus.captures ?? 0}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Uploads</p>
                <p className="font-medium">{deviceStatus.uploads?.success ?? 0}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Failed</p>
                <p className="font-medium">{deviceStatus.uploads?.failed ?? 0}</p>
              </div>
            </div>
          )}
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => handleCommand("start_detection")}
              disabled={sending}
            >
              Start Detection
            </Button>
            <Button
              variant="outline"
              onClick={() => handleCommand("stop_detection")}
              disabled={sending}
            >
              Stop Detection
            </Button>
            <Button
              onClick={() => handleCommand("capture")}
              disabled={sending}
            >
              Capture Now
            </Button>
          </div>
          <Separator />
          <div className="space-y-3">
            <p className="text-sm font-medium">Motion Settings</p>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="contour">Min Contour Area</Label>
                <Input
                  id="contour"
                  type="number"
                  key={`contour-${device.config?.min_contour_area}`}
                  defaultValue={device.config?.min_contour_area ?? 500}
                  min={100}
                  step={100}
                  onBlur={async (e) => {
                    const val = Number(e.target.value);
                    try {
                      const updated = await updateDevice(token!, id!, {
                        config: { ...device.config, min_contour_area: val },
                      } as any);
                      setDevice(updated);
                    } catch (err) { console.error(err); }
                  }}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="threshold">Motion Threshold</Label>
                <Input
                  id="threshold"
                  type="number"
                  key={`threshold-${device.config?.threshold}`}
                  defaultValue={device.config?.threshold ?? 25}
                  min={1}
                  max={255}
                  onBlur={async (e) => {
                    const val = Number(e.target.value);
                    try {
                      const updated = await updateDevice(token!, id!, {
                        config: { ...device.config, threshold: val },
                      } as any);
                      setDevice(updated);
                    } catch (err) { console.error(err); }
                  }}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="cooldown">Cooldown (seconds)</Label>
              <Input
                id="cooldown"
                type="number"
                key={`cooldown-${device.config?.cooldown_seconds}`}
                defaultValue={device.config?.cooldown_seconds ?? 2}
                min={0}
                step={0.5}
                className="w-1/2"
                onBlur={async (e) => {
                  const val = Number(e.target.value);
                  try {
                    const updated = await updateDevice(token!, id!, {
                      config: { ...device.config, cooldown_seconds: val },
                    } as any);
                    setDevice(updated);
                  } catch (err) { console.error(err); }
                }}
              />
            </div>
            <p className="text-xs text-muted-foreground">
              Changes are saved to the device and pushed immediately.
            </p>
          </div>
        </CardContent>
      </Card>

      <Separator />

      <Card>
        <CardHeader>
          <CardTitle>API Key</CardTitle>
        </CardHeader>
        <CardContent>
          <AlertDialog>
            <AlertDialogTrigger render={<Button variant="outline" />}>
              Rotate API Key
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Rotate API Key?</AlertDialogTitle>
                <AlertDialogDescription>
                  This will invalidate the current key. The device will need to be updated with the new key.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction onClick={handleRotateKey}>Rotate Key</AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardContent>
      </Card>

      <Card className="border-destructive">
        <CardHeader>
          <CardTitle className="text-destructive">Danger Zone</CardTitle>
        </CardHeader>
        <CardContent>
          <AlertDialog>
            <AlertDialogTrigger render={<Button variant="destructive" />}>
              Delete Device
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete {device.name}?</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. All uploads associated with this device will remain but the device will be removed.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction onClick={handleDelete}>Delete</AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardContent>
      </Card>

      <Dialog open={!!newKey} onOpenChange={() => setNewKey(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>New API Key</DialogTitle>
            <DialogDescription>
              Copy the new API key below. It will only be shown once.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="rounded-md bg-muted p-4 font-mono text-sm break-all">
              {newKey}
            </div>
            <Button onClick={handleCopy} className="w-full">
              {copied ? "Copied!" : "Copy API Key"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
