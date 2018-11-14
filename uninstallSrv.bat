@ECHO off
ECHO Uninstall GroundCast Emulator windows service
ECHO step 1: stop GroundCastemulator
sc stop GroundCastemulator
ECHO step 2: delete GroundCastemulator
sc delete GroundCastemulator