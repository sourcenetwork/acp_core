import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";

export const computeShortHash = async (data: string) => {
  const encodedData = new TextEncoder().encode(data);
  const buffer = await crypto.subtle.digest("SHA-256", encodedData);
  return Array.from(new Uint8Array(buffer))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("")
    .substring(0, 8);
};

export const exportSandboxData = async (activeStateData?: SandboxData) => {
  if (!activeStateData) return;

  const jsonBlob = new Blob([JSON.stringify(activeStateData, null, 2)], {
    type: "application/json",
  });

  const url = URL.createObjectURL(jsonBlob);
  const link = document.createElement("a");
  const filename = await computeShortHash(JSON.stringify(activeStateData));
  link.href = url;
  link.download = `acp-playground-export-${filename}.json`;
  link.click();
  URL.revokeObjectURL(url);
};

function removeElement(elem: Element) {
  document.body.removeChild(elem);
}

export const importSandboxData = async (): Promise<
  Partial<SandboxData> | undefined | false
> => {
  return new Promise((resolve, reject) => {
    const fileInput = document.createElement("input");
    fileInput.type = "file";
    fileInput.accept = ".json";
    fileInput.style.display = "none";
    document.body.appendChild(fileInput);

    const cleanUp = (error?: string) => {
      removeElement(fileInput);
      if (error) reject(new Error(error));
    };

    const onChangeHandler = (event: Event) => {
      const file = (event.target as HTMLInputElement).files?.[0];

      if (!file || file?.type !== "application/json") {
        cleanUp("No file selected or invalid");
        return;
      }

      const reader = new FileReader();

      reader.onload = () => {
        try {
          if (typeof reader.result !== "string") {
            throw new Error("Invalid file format");
          }

          const parsedData = JSON.parse(reader.result) as Partial<SandboxData>;
          resolve(parsedData);
        } catch (error) {
          reject(error as Error);
        } finally {
          cleanUp();
        }
      };

      reader.onerror = () => cleanUp("Failed to read file");
      reader.readAsText(file);
    };

    const onCloseHandler = () => {
      cleanUp();
      resolve(false);
    };

    fileInput.addEventListener("change", onChangeHandler);
    fileInput.addEventListener("cancel", onCloseHandler);

    fileInput.click();
  });
};
