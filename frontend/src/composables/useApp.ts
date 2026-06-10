import { SelectFile, SelectDirectory, Extract } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export function useApp() {
  async function selectFile(): Promise<string> {
    return await SelectFile()
  }

  async function selectDirectory(): Promise<string> {
    return await SelectDirectory()
  }

  async function extract(archivePath: string, outputDir: string): Promise<void> {
    return await Extract(archivePath, outputDir)
  }

  function onExtractLog(callback: (line: string) => void) {
    return EventsOn('extract-log', callback)
  }

  return {
    selectFile,
    selectDirectory,
    extract,
    onExtractLog,
  }
}
