export async function loadPlaygroundWasm(path: string) {
  const go = new window.Go();

  const result = await WebAssembly.instantiateStreaming(
    fetch(path),
    go.importObject
  );

  void go.run(result.instance);

  return result;
}
