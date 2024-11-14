export async function loadGoWASM(path: string) {
  const go = new (window as any).Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(path),
    go.importObject
  );
  go.run(result.instance);
  return result;
}