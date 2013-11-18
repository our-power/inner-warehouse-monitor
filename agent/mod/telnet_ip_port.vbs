IP_PORT = "http://10.24.178.60:8000"

Set x=CreateObject("msxml2.serverXMLHTTP") 
x.Open "POST", IP_PORT
x.Send("test") 
If ERr.NuMbEr=0 Or ERr.NuMbEr=-2147012866 Or ERr.NuMbEr=-2147012894 Or ERr.NuMbEr=-2147012744 Or ERr.NuMbEr=-2147467259 Then 
    wsh.Echo "IP:PORT OK" 
End If 