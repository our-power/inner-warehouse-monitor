IP_PORTs = "10.24.178.60:8000,10.24.178.60:8009,10.24.178.60:8080"

Wscript.Echo "2"
list = Split(IP_PORTs, ",")
for each IP_PORT in list
    Wscript.Echo IP_PORT & "," & TEST_IP_PORT(IP_PORT)
next

Function TEST_IP_PORT(IP_PORT)
    On Error Resume Next
    Set x=CreateObject("msxml2.serverXMLHTTP")
    x.Open "POST", "http://" & IP_PORT
    x.Send("test")
    If ERr.NuMbEr=0 Or ERr.NuMbEr=-2147012866 Or ERr.NuMbEr=-2147012894 Or ERr.NuMbEr=-2147012744 Or ERr.NuMbEr=-2147467259 Then
        TEST_IP_PORT = "OK"
    End If
End Function