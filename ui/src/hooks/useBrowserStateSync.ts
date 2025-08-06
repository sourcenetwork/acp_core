import {
  PersistedSandboxData,
  usePlaygroundStore,
} from "@/stores/playgroundStore";
import { useEffect, useRef } from "react";

interface useBrowserSandboxSyncOptions {
  channelName: string;
}

export function useBrowserSandboxSync(options: useBrowserSandboxSyncOptions) {
  const { channelName } = options;
  const channelRef = useRef<BroadcastChannel | null>(null);
  const isReceiver = useRef(false);

  const { handleSandboxSyncChangeReceived } = usePlaygroundStore((state) => ({
    handleSandboxSyncChangeReceived: state.handleSandboxSyncChangeReceived,
  }));

  useEffect(() => {
    // Don't run if BroadcastChannel is not supported
    if (!globalThis.BroadcastChannel) return;

    channelRef.current = new BroadcastChannel(channelName);

    // Listen for messages from other windows
    channelRef.current.onmessage = async (event: {
      data: { sandboxes: PersistedSandboxData[] };
    }) => {
      const { data } = event;
      if (!data.sandboxes) return;

      isReceiver.current = true;

      try {
        // Set the sandboxes state to the received sandbox data
        usePlaygroundStore.setState(data);
        await handleSandboxSyncChangeReceived(data.sandboxes);
      } finally {
        isReceiver.current = false;
      }
    };

    // Cleanup on unmount
    return () => {
      channelRef.current?.close();
      channelRef.current = null;
    };
  }, [channelName]);

  // Subscribe to sandbox data state changes
  useEffect(() => {
    const unsubscribe = usePlaygroundStore.subscribe(
      (state) => state.sandboxes,
      (sandboxes) => {
        // Don't broadcast if we're receiving
        if (isReceiver.current || !channelRef.current) return;
        // Broadcast the change to sandbox data
        channelRef.current.postMessage({ sandboxes });
      }
    );

    return unsubscribe;
  }, []);

  return;
}
