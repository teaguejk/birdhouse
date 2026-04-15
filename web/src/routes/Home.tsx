import { useState } from "react";
import { DeviceStatusPanel } from "@/components/device-status";
import { DeviceImages } from "@/components/device-images";

export default function Home() {
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);

  return (
    <main className="mx-auto flex max-w-3xl flex-col items-center gap-8 p-8 pt-16">
      <DeviceStatusPanel
        selectedDeviceId={selectedDeviceId}
        onSelectDevice={setSelectedDeviceId}
      />
      {selectedDeviceId && (
        <div className="w-full animate-in fade-in slide-in-from-bottom-2 duration-300">
          <DeviceImages deviceId={selectedDeviceId} />
        </div>
      )}
    </main>
  );
}
