; NSIS script to generate rmbasicx64.msi Windows installer
; Adapted from https://nsis.sourceforge.io/Sample_installation_script_for_an_application
 
; -------------------------------
; Start
 
  !define MUI_FILE "rmbasicx64.exe"
  !define MUI_NAME "RM BASICx64"
  !define MUI_BRANDINGTEXT "RM BASICx64"
  CRCCheck On
 
  ; We should test if we must use an absolute path 
  !include "${NSISDIR}\Contrib\Modern UI\System.nsh"
 
  ; Modern UI
  !include "MUI.nsh"
 
;---------------------------------
;General
 
  OutFile "rmbasicx64"
  ShowInstDetails "nevershow"
  ShowUninstDetails "nevershow"
  ;SetCompressor "bzip2"


;--------------------------------
;Folder selection page
 
  InstallDir "$PROGRAMFILES\${MUI_NAME}"
 
 
;--------------------------------
;Modern UI Configuration
 
  !define MUI_WELCOMEPAGE  
  !define MUI_LICENSEPAGE
  !define MUI_DIRECTORYPAGE
  !define MUI_ABORTWARNING
  !define MUI_UNINSTALLER
  !define MUI_UNCONFIRMPAGE
  !define MUI_FINISHPAGE  
 
 
;--------------------------------
;Language
 
  !insertmacro MUI_LANGUAGE "English"
 

;--------------------------------
;Data
 
  LicenseData "license.txt"
 
 
;-------------------------------- 
;Installer Sections     
Section "install" Installation info
 
;Add files
  SetOutPath "$INSTDIR"
 
  File "${MUI_FILE}.exe"
 
;create desktop shortcut
  CreateShortCut "$DESKTOP\${MUI_NAME}.lnk" "$INSTDIR\${MUI_FILE}.exe" ""
 
;create start-menu items
  CreateDirectory "$SMPROGRAMS\${MUI_NAME}"
  CreateShortCut "$SMPROGRAMS\${MUI_NAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe" "" "$INSTDIR\Uninstall.exe" 0
  CreateShortCut "$SMPROGRAMS\${MUI_NAME}\${MUI_NAME}.lnk" "$INSTDIR\${MUI_FILE}.exe" "" "$INSTDIR\${MUI_FILE}.exe" 0
 
;write uninstall information to the registry
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_NAME}" "DisplayName" "${MUI_NAME} (remove only)"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_NAME}" "UninstallString" "$INSTDIR\Uninstall.exe"
 
  WriteUninstaller "$INSTDIR\Uninstall.exe"
 
SectionEnd
 
 
;--------------------------------    
;Uninstaller Section  
Section "Uninstall"
 
;Delete Files 
  RMDir /r "$INSTDIR\*.*"    
 
;Remove the installation directory
  RMDir "$INSTDIR"
 
;Delete Start Menu Shortcuts
  Delete "$DESKTOP\${MUI_NAME}.lnk"
  Delete "$SMPROGRAMS\${MUI_NAME}\*.*"
  RmDir  "$SMPROGRAMS\${MUI_NAME}"
 
;Delete Uninstaller And Unistall Registry Entries
  DeleteRegKey HKEY_LOCAL_MACHINE "SOFTWARE\${MUI_NAME}"
  DeleteRegKey HKEY_LOCAL_MACHINE "SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_NAME}"  
 
SectionEnd
 
 
;--------------------------------    
;MessageBox Section
 
 
;Function that calls a messagebox when installation finished correctly
Function .onInstSuccess
  MessageBox MB_OK "You have successfully installed ${MUI_NAME}. Use the desktop icon to start the program."
FunctionEnd
 
Function un.onUninstSuccess
  MessageBox MB_OK "You have successfully uninstalled ${MUI_NAME}."
FunctionEnd
 
 
;eof