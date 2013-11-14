Set objWMIService = GetObject("winmgmts:\\.\root\CIMV2") 

Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PhysicalMemory",,48) 
PhysicalMemory = 0
For Each objItem in colItems 
    PhysicalMemory = PhysicalMemory + objItem.Capacity
Next

Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfFormattedData_PerfOS_Memory",,48) 
Available = 0
For Each objItem in colItems 
    Available = Available + objItem.AvailableBytes
Next

UsedBytes = PhysicalMemory - Available

if PhysicalMemory <> 0 then 
	Usage = Round(UsedBytes/PhysicalMemory*100, 2)
else
	Usage = -1
end if
 
AvailableMB = Round(UsedBytes/1024/1024)
AvailableKB = Round(UsedBytes/1024)

Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PageFileUsage",,48) 
PageFileUsage = 0
For Each objItem in colItems 
    PageFileUsage = PageFileUsage + objItem.CurrentUsage*1024
Next

wscript.echo AvailableMB & chr(9) & AvailableKB & chr(9) & PageFileUsage & chr(9) & Usage