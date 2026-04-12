"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { createDevice } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

export default function CreateDevicePage() {
  const { token } = useAuth();
  const router = useRouter();
  const [name, setName] = useState("");
  const [location, setLocation] = useState("");
  const [loading, setLoading] = useState(false);
  const [apiKey, setApiKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !name) return;

    setLoading(true);
    try {
      const result = await createDevice(token, { name, location });
      setApiKey(result.api_key);
    } catch (err) {
      console.error("Failed to create device:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    if (apiKey) {
      await navigator.clipboard.writeText(apiKey);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="mx-auto max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>Add Device</CardTitle>
          <CardDescription>Register a new birdhouse device</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Backyard Birdhouse"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Input
                id="location"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="Backyard oak tree"
              />
            </div>
            <div className="flex gap-3">
              <Button type="submit" disabled={loading || !name}>
                {loading ? "Creating..." : "Create Device"}
              </Button>
              <Button type="button" variant="outline" onClick={() => router.push("/admin")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Dialog open={!!apiKey} onOpenChange={() => router.push("/admin")}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Device Created</DialogTitle>
            <DialogDescription>
              Copy the API key below. It will only be shown once.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="rounded-md bg-muted p-4 font-mono text-sm break-all">
              {apiKey}
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
