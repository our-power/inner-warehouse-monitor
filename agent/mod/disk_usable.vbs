Set objWMIService = GetObject("winmgmts:\\.\root\CIMV2") 
Function Usable(disk)
	Set colItems = objWMIService.ExecQuery("SELECT * FROM WIN32_LogicalDisk WHERE DriveType=3 AND DeviceId='" & disk & "'",,48) 
	For Each objItem in colItems 
	    Usable = Round(objItem.FreeSpace/objItem.Size*100, 2)
	    Exit Function
	Next
	Usable = -1
End Function
wscript.echo Usable("C:") & chr(9) & Usable("D:") & chr(9) & Usable("E:")