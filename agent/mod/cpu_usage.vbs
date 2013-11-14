CPU_INTERVAL_ONCE = 15
TIMES = 12

maxLoad = 0
maxQueue = 0
cpuString = ""

Sub GetLoad()
	Set objWMIService = GetObject("winmgmts:\\.\root\CIMV2") 
	Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_Processor",,48) 
	total = 0
	count = 0
	LoadPercentage = -1
	For Each objItem in colItems 
	    total = total + objItem.LoadPercentage
	    count = count + 1
	Next
	if count <> 0 then
		LoadPercentage = total/count
	end if

	total = 0
	count = 0
	ProcessorQueueLength = -1
	Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfRawData_PerfOS_System",,48) 
	For Each objItem in colItems 
	    total = total + objItem.ProcessorQueueLength
	    count = count + 1
	Next
	if count <> 0 then
		ProcessorQueueLength = total/count
	end if

	if cpuString <> "" then
	    cpuString = cpuString & "=" & LoadPercentage
	else
	    cpuString = LoadPercentage
	end if

	if maxLoad < LoadPercentage then
		maxLoad = LoadPercentage
	end if
	if maxQueue < ProcessorQueueLength then
		maxQueue = ProcessorQueueLength
	end if
End Sub

for i = 1 to TIMES
	GetLoad
	wscript.sleep CPU_INTERVAL_ONCE * 1000
next

wscript.echo maxLoad & chr(9) & maxQueue & chr(9) & cpuString
