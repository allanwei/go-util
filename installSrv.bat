@ECHO off
ECHO create GroundCast Emulator windows service
set mypath=%CD%
ECHO step 1: create GroundCastemulator start=auto binPath=%mypath%\go_util.exe DisplayName="GroundCast Emulator"
sc create GroundCastemulator start=auto binPath=%mypath%\go_util.exe DisplayName="GroundCast Emulator"
ECHO step 2: start GroundCastemulator
sc start GroundCastemulator