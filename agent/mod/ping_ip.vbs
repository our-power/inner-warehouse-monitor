IPs = "127.0.0.1,256.0.0.0,10.24.178.60,1.1.1.1"

Function Ping(strHost)
    On Error Resume Next
    Dim oPing, oRetStatus, bReturn
    Set oPing = GetObject("winmgmts:{impersonationLevel=impersonate}").ExecQuery("select * from Win32_PingStatus where address='" & strHost & "'")
    For Each oRetStatus In oPing
        If IsNull(oRetStatus.StatusCode) Or oRetStatus.StatusCode <> 0 Then
            bReturn = False
            Ping = "Status code is " & oRetStatus.StatusCode
        Else
            bReturn = True
            Ping = "ResponseTime=" & oRetStatus.ResponseTime
        End If
        Set oRetStatus = Nothing
    Next
    Set oPing = Nothing
End Function

Wscript.Echo "1"
IP_list = Split(IPs, ",")
for each IP in IP_list
    Wscript.Echo IP & "," & Ping(IP)
next