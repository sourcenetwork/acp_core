import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";

const formatFilename = (name: string) => {
  return name
    .replace(/\s+/g, "_") // Replace all whitespace with underscores
    .replace(/[^a-z0-9_\-\.]/gi, "") // Remove any remaining invalid characters
    .replace(/_+/g, "_") // Collapse multiple underscores
    .replace(/^_+|_+$/g, "") // Trim leading/trailing underscores
    .substring(0, 255);
};

export const exportSandboxData = async (
  name?: string,
  activeStateData?: SandboxData
) => {
  if (!activeStateData) return;
  const jsonBlob = new Blob([JSON.stringify(activeStateData, null, 2)], {
    type: "application/json",
  });

  const url = URL.createObjectURL(jsonBlob);
  const link = document.createElement("a");

  const prefix = "acp_export";
  const safeName = name ? formatFilename(name) : null;
  const date = new Date().toISOString().replace(/[-:.]/g, "");
  const segments = [prefix, safeName, date].filter(Boolean);
  const filename = `${segments.join("_")}.json`;

  link.href = url;
  link.download = filename;
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
