"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090";

interface Upload {
  id: string;
  filename: string;
  original_name: string;
  mime_type: string;
  url: string;
  status: string;
  created_at: string;
}

interface PaginatedUploads {
  data: Upload[] | null;
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

interface DeviceImagesProps {
  deviceId: string;
}

export function DeviceImages({ deviceId }: DeviceImagesProps) {
  const [uploads, setUploads] = useState<PaginatedUploads | null>(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  useEffect(() => {
    setPage(1);
  }, [deviceId]);

  useEffect(() => {
    setLoading(true);
    fetch(`${API_BASE_URL}/uploads/device/${deviceId}?page=${page}&page_size=12`)
      .then((res) => res.json())
      .then((data) => setUploads(data))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [deviceId, page]);

  if (loading) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle className="text-lg">Images</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} className="aspect-square w-full rounded-md" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const images = uploads?.data || [];
  const pagination = uploads?.pagination;

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Images</CardTitle>
          {pagination && pagination.total > 0 && (
            <span className="text-sm text-muted-foreground">
              {pagination.total} total
            </span>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {images.length === 0 ? (
          <p className="text-sm text-muted-foreground">No images from this device yet</p>
        ) : (
          <>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
              {images.map((upload, index) => (
                <a
                  key={upload.id}
                  href={upload.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="group relative overflow-hidden rounded-md border bg-muted"
                >
                  {index === 0 && page === 1 && (
                    <Badge className="absolute top-2 left-2 z-10">Latest</Badge>
                  )}
                  <img
                    src={upload.url}
                    alt={upload.original_name}
                    className="aspect-square w-full object-cover transition-transform group-hover:scale-105"
                  />
                  <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/60 to-transparent p-2">
                    <p className="truncate text-xs text-white">
                      {new Date(upload.created_at).toLocaleString()}
                    </p>
                  </div>
                </a>
              ))}
            </div>
            {pagination && pagination.total_pages > 1 && (
              <div className="flex items-center justify-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={!pagination.has_prev}
                  onClick={() => setPage((p) => p - 1)}
                >
                  Previous
                </Button>
                <span className="text-sm text-muted-foreground">
                  {pagination.page} / {pagination.total_pages}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={!pagination.has_next}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
