IP = "127.0.0.1"

Function Ping(strHost)
    Dim oPing, oRetStatus, bReturn
    Set oPing = GetObject("winmgmts:{impersonationLevel=impersonate}").ExecQuery("select * from Win32_PingStatus where address='" & strHost & "'")
    For Each oRetStatus In oPing
        If IsNull(oRetStatus.StatusCode) Or oRetStatus.StatusCode <> 0 Then
            bReturn = False
            WScript.Echo "Status code is " & oRetStatus.StatusCode
        Else
            bReturn = True
            Wscript.Echo "ResponseTime=" & oRetStatus.ResponseTime
        End If
        Set oRetStatus = Nothing
    Next
    Set oPing = Nothing
    Ping = bReturn
End Function

Wscript.Echo IP
Ping(IP)