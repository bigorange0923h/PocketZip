import { SelectFile, SelectDirectory, Extract, ExtractWithPassword, GetPasswordCandidates, SavePassword } from '../../wailsjs/go/main/App'
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

  async function extractWithPassword(archivePath: string, outputDir: string, password: string): Promise<void> {
    return await ExtractWithPassword(archivePath, outputDir, password)
  }

  async function getPasswordCandidates(archivePath: string): Promise<string[]> {
    return await GetPasswordCandidates(archivePath)
  }

  async function savePassword(archivePath: string, password: string): Promise<void> {
    return await SavePassword(archivePath, password)
  }

  function onExtractLog(callback: (line: string) => void) {
    return EventsOn('extract-log', callback)
  }

  return {
    selectFile,
    selectDirectory,
    extract,
    extractWithPassword,
    getPasswordCandidates,
    savePassword,
    onExtractLog,
  }
}
