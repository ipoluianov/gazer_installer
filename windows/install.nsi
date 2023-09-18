!define PRODUCT_NAME "GazerNode"
!define PRODUCT_VERSION "2.4.15"
!define GAZER_NODE_EXE "gazer_node.exe"
!define GAZER_CLIENT_EXE "gazer_client.exe"
!define ARP "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"

!include "MUI2.nsh"
!include "FileFunc.nsh"

!define MUI_HEADERIMAGE
;!define MUI_HEADERIMAGE_BITMAP "installer\installer-icon.bmp"

;!define MUI_COMPONENTSPAGE_SMALLDESC ;No value
;!define MUI_UI "myUI.exe" ;Value
;!define MUI_INSTFILESPAGE_COLORS "FFFFFF 000000" ;Two colors

Name "${PRODUCT_NAME}"
Caption "Install ${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "${PRODUCT_NAME}_${PRODUCT_VERSION}.exe"
ShowInstDetails show

InstallDir "$PROGRAMFILES64\gazernode"

;Page directory
;Page instfiles

SetCompressor /SOLID lzma

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_LICENSE "license.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

    # These indented statements modify settings for MUI_PAGE_FINISH
	!define MUI_FINISHPAGE_NOAUTOCLOSE
	!define MUI_FINISHPAGE_RUN
	!define MUI_FINISHPAGE_RUN_CHECKED
	!define MUI_FINISHPAGE_RUN_TEXT "Start Gazer Node"
	!define MUI_FINISHPAGE_RUN_FUNCTION "LaunchLink"
!insertmacro MUI_PAGE_FINISH
  
!insertmacro MUI_LANGUAGE "English"
  
Section
	SetOutPath "$PROGRAMFILES64\gazernode"

	DetailPrint "Stopping service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-stop" SW_HIDE
	DetailPrint "Uninstalling service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-uninstall" SW_HIDE
	
	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -stop'
	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -uninstall'

	File /a /r "bin\"
	
	WriteUninstaller $INSTDIR\uninstaller.exe
	
	WriteRegStr HKLM "${ARP}" "DisplayName" "${PRODUCT_NAME}"
	WriteRegStr HKLM "${ARP}" "UninstallString" "$\"$INSTDIR\uninstaller.exe$\""
	WriteRegStr HKLM "${ARP}" "Publisher" "Poluianov Ivan"
	WriteRegStr HKLM "${ARP}" "DisplayVersion" "${PRODUCT_VERSION}"
	WriteRegStr HKLM "${ARP}" "DisplayIcon" "$INSTDIR\${GAZER_CLIENT_EXE}"
	
	
	${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
	IntFmt $0 "0x%08X" $0
	WriteRegDWORD HKLM "${ARP}" "EstimatedSize" "$0"
	
	
	CreateShortCut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\${GAZER_CLIENT_EXE}" ""
	CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
	CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\uninstaller.lnk" "$INSTDIR\uninstaller.exe" "" "$INSTDIR\uninstaller.exe" 0
	CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\${GAZER_NODE_EXE}" "" "$INSTDIR\${GAZER_NODE_EXE}" 0

	DetailPrint "Installing service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-install" SW_HIDE
	
	DetailPrint "Starting service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-start" SW_HIDE
	
	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -install'
	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -start'
SectionEnd

Section Uninstall
	DetailPrint "Stopping service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-stop" SW_HIDE
	DetailPrint "Uninstalling service ..."
	ExecShellWait "" "$INSTDIR\${GAZER_NODE_EXE}" "-uninstall" SW_HIDE

	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -stop'
	# ExecWait '"$INSTDIR\${GAZER_NODE_EXE}" -uninstall'

	Delete $INSTDIR\uninstaller.exe
	Delete $INSTDIR\${GAZER_NODE_EXE}
	
	Delete "$DESKTOP\${PRODUCT_NAME}.lnk"
	Delete "$SMPROGRAMS\${PRODUCT_NAME}\*.*"
	RMDir /r "$SMPROGRAMS\${PRODUCT_NAME}"
	
	DeleteRegKey HKLM "${ARP}"
	
	RMDir /r $INSTDIR
SectionEnd

Function LaunchLink
  ExecShell "" "$INSTDIR\${GAZER_CLIENT_EXE}"
FunctionEnd
