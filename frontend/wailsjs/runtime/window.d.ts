interface Window {
  go: Record<string, Record<string, Record<string, (...args: any[]) => Promise<any>>>>;
  runtime: Record<string, (...args: any[]) => any>;
}
