!include "MUI2.nsh"

Name "PocketZip"
OutFile "PocketZip-Setup.exe"
InstallDir "$PROGRAMFILES\PocketZip"

Page directory
Page instfiles

Section "Install"
  SetOutPath "$INSTDIR"
  File "build/bin/PocketZip.exe"
  File "build/bin/7z.exe"
  File "build/bin/7z.dll"

  CreateDirectory "$SMPROGRAMS\PocketZip"
  CreateShortCut "$SMPROGRAMS\PocketZip\PocketZip.lnk" "$INSTDIR\PocketZip.exe"
  CreateShortCut "$DESKTOP\PocketZip.lnk" "$INSTDIR\PocketZip.exe"

  WriteUninstaller "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\PocketZip.exe"
  Delete "$INSTDIR\7z.exe"
  Delete "$INSTDIR\7z.dll"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"

  Delete "$SMPROGRAMS\PocketZip\PocketZip.lnk"
  RMDir "$SMPROGRAMS\PocketZip"
  Delete "$DESKTOP\PocketZip.lnk"
SectionEnd
