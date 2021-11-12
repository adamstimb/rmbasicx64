; Script generated by the HM NIS Edit Script Wizard.

; HM NIS Edit Wizard helper defines
!define PRODUCT_NAME "RM BASICx64"
!define PRODUCT_VERSION "0.20"
!define PRODUCT_PUBLISHER "Tim Adams"
!define PRODUCT_WEB_SITE "https://adamstimb.github.io/rmbasicx64site"
!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\rmbasicx64.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"

; MUI 1.67 compatible ------
!include "MUI.nsh"

; MUI Settings
!define MUI_ABORTWARNING
!define MUI_ICON "rmbasicx64-ico-256.ico"
!define MUI_UNICON "rmbasicx64-ico-256.ico"

; Welcome page
!insertmacro MUI_PAGE_WELCOME
; License page
!insertmacro MUI_PAGE_LICENSE "license.txt"
; Installartion directory page
!insertmacro MUI_PAGE_DIRECTORY
; Workspace directory page
Var WorkspaceDir
!define MUI_PAGE_HEADER_SUBTEXT "Choose where to create your workspace directory."
!define MUI_DIRECTORYPAGE_TEXT_TOP "RM BASICx64 will create a folder called 'RMBASICx64 Workspace' in this location if it does not already exist.  Your BASIC programs will be stored here.  To use a differenct folder, click Browse and select another folder. Click Next to continue."
!define MUI_DIRECTORYPAGE_VARIABLE $WorkspaceDir
!insertmacro MUI_PAGE_DIRECTORY
; Instfiles page
!insertmacro MUI_PAGE_INSTFILES
; Finish page
!define MUI_FINISHPAGE_RUN "$INSTDIR\rmbasicx64.exe"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_INSTFILES

; Language files
!insertmacro MUI_LANGUAGE "English"

; MUI end ------

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "rmbasicx64setup.exe"
InstallDir "$PROGRAMFILES\RM BASICx64"
InstallDirRegKey HKLM "${PRODUCT_DIR_REGKEY}" ""
ShowInstDetails show
ShowUnInstDetails show

Section "MainSection" SEC01
  SetOutPath "$INSTDIR"
  SetOverwrite ifnewer
  File "rmbasicx64.exe"
  File "rmbasicx64-ico-256.ico"
  CreateDirectory "$SMPROGRAMS\RM BASICx64"
  CreateShortCut "$SMPROGRAMS\RM BASICx64\RM BASICx64.lnk" "$INSTDIR\rmbasicx64.exe" "" "$INSTDIR\rmbasicx64-ico-256.ico" 0
  CreateShortCut "$DESKTOP\RM BASICx64.lnk" "$INSTDIR\rmbasicx64.exe" "" "$INSTDIR\rmbasicx64-ico-256.ico" 0
SectionEnd

Section -AdditionalIcons
  WriteIniStr "$INSTDIR\${PRODUCT_NAME}.url" "InternetShortcut" "URL" "${PRODUCT_WEB_SITE}"
  CreateShortCut "$SMPROGRAMS\RM BASICx64\Website.lnk" "$INSTDIR\${PRODUCT_NAME}.url"
  CreateShortCut "$SMPROGRAMS\RM BASICx64\Uninstall.lnk" "$INSTDIR\uninst.exe"
SectionEnd

Section -Post
  WriteUninstaller "$INSTDIR\uninst.exe"
  WriteRegStr HKLM "${PRODUCT_DIR_REGKEY}" "" "$INSTDIR\rmbasicx64.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\rmbasicx64.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
SectionEnd


Function un.onUninstSuccess
  HideWindow
  MessageBox MB_ICONINFORMATION|MB_OK "$(^Name) was successfully removed from your computer."
FunctionEnd

Function un.onInit
  MessageBox MB_ICONQUESTION|MB_YESNO|MB_DEFBUTTON2 "Are you sure you want to completely remove $(^Name) and all of its components?" IDYES +2
  Abort
FunctionEnd

Section Uninstall
  Delete "$INSTDIR\${PRODUCT_NAME}.url"
  Delete "$INSTDIR\uninst.exe"
  Delete "$INSTDIR\rmbasicx64.exe"
  Delete "$INSTDIR\rmbasicx64-ico-256.ico"

  Delete "$SMPROGRAMS\RM BASICx64\Uninstall.lnk"
  Delete "$SMPROGRAMS\RM BASICx64\Website.lnk"
  Delete "$DESKTOP\RM BASICx64.lnk"
  Delete "$SMPROGRAMS\RM BASICx64\RM BASICx64.lnk"

  RMDir "$SMPROGRAMS\RM BASICx64"
  RMDir "$INSTDIR"

  DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}"
  DeleteRegKey HKLM "${PRODUCT_DIR_REGKEY}"
  SetAutoClose true
SectionEnd