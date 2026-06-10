!include "MUI2.nsh"

Name "PocketUnzip"
OutFile "PocketUnzip-Setup.exe"
InstallDir "$PROGRAMFILES\PocketUnzip"

Page directory
Page instfiles

Section "Install"
  SetOutPath "$INSTDIR"
  File "build/bin/PocketUnzip.exe"
  File "build/bin/7z.exe"
  File "build/bin/7z.dll"

  CreateDirectory "$SMPROGRAMS\PocketUnzip"
  CreateShortCut "$SMPROGRAMS\PocketUnzip\PocketUnzip.lnk" "$INSTDIR\PocketUnzip.exe"
  CreateShortCut "$DESKTOP\PocketUnzip.lnk" "$INSTDIR\PocketUnzip.exe"

  WriteUninstaller "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\PocketUnzip.exe"
  Delete "$INSTDIR\7z.exe"
  Delete "$INSTDIR\7z.dll"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"

  Delete "$SMPROGRAMS\PocketUnzip\PocketUnzip.lnk"
  RMDir "$SMPROGRAMS\PocketUnzip"
  Delete "$DESKTOP\PocketUnzip.lnk"
SectionEnd
