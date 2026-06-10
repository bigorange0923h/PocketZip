import {
  SelectFile,
  SelectFiles,
  SelectDirectory,
  Extract,
  ExtractWithPassword,
  ExtractWithStrategy,
  BatchExtract,
  GetPasswordCandidates,
  SavePassword,
  GetHistory,
  OpenDirectory,
  TestArchive,
  PreviewArchive,
  GetPasswordRecords,
  DeletePasswordRecord,
  GetPasswordStats,
  GetTheme,
  SetTheme,
  GetExtractStrategies,
  SaveExtractStrategy
} from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export function useApp() {
  async function selectFile(): Promise<string> {
    return await SelectFile()
  }

  async function selectFiles(): Promise<string[]> {
    return await SelectFiles()
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

  async function extractWithStrategy(archivePath: string, strategyName: string): Promise<void> {
    return await ExtractWithStrategy(archivePath, strategyName)
  }

  async function batchExtract(archivePaths: string[], outputDir: string): Promise<any[]> {
    return await BatchExtract(archivePaths, outputDir)
  }

  async function getPasswordCandidates(archivePath: string): Promise<string[]> {
    return await GetPasswordCandidates(archivePath)
  }

  async function savePassword(archivePath: string, password: string): Promise<void> {
    return await SavePassword(archivePath, password)
  }

  async function getHistory(limit: number): Promise<any[]> {
    return await GetHistory(limit)
  }

  async function openDirectory(dirPath: string): Promise<void> {
    return await OpenDirectory(dirPath)
  }

  async function testArchive(archivePath: string): Promise<boolean> {
    return await TestArchive(archivePath)
  }

  async function previewArchive(archivePath: string): Promise<any[]> {
    return await PreviewArchive(archivePath)
  }

  async function getPasswordRecords(): Promise<any[]> {
    return await GetPasswordRecords()
  }

  async function deletePasswordRecord(id: number): Promise<void> {
    return await DeletePasswordRecord(id)
  }

  async function getPasswordStats(): Promise<any> {
    return await GetPasswordStats()
  }

  async function getTheme(): Promise<string> {
    return await GetTheme()
  }

  async function setTheme(theme: string): Promise<void> {
    return await SetTheme(theme)
  }

  async function getExtractStrategies(): Promise<any[]> {
    return await GetExtractStrategies()
  }

  async function saveExtractStrategy(strategy: any): Promise<void> {
    return await SaveExtractStrategy(strategy)
  }

  function onExtractLog(callback: (line: string) => void) {
    return EventsOn('extract-log', callback)
  }

  return {
    selectFile,
    selectFiles,
    selectDirectory,
    extract,
    extractWithPassword,
    extractWithStrategy,
    batchExtract,
    getPasswordCandidates,
    savePassword,
    getHistory,
    openDirectory,
    testArchive,
    previewArchive,
    getPasswordRecords,
    deletePasswordRecord,
    getPasswordStats,
    getTheme,
    setTheme,
    getExtractStrategies,
    saveExtractStrategy,
    onExtractLog,
  }
}
