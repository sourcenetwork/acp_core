export async function loadPlaygroundWasm(path: string) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const go = new (window as Window).Go();

  const result = await WebAssembly.instantiateStreaming(
    fetch(path),
    go.importObject
  );
  go.run(result.instance);
  return result;
}
