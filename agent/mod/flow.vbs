FLOW_USAGE_INTERVAL = 300

Set objWMIService = GetObject("winmgmts:\\.\root\CIMV2") 
Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_NetworkAdapter WHERE PhysicalAdapter=True",,48) 
Dim nics()
count = 0
For Each objItem in colItems
	ReDim Preserve nics(count+1)
	NicName = Replace(objItem.Name, "#", "_")
	NicName = Replace(NicName, "(", "[")
	NicName = Replace(NicName, ")", "]")
	NicName = Replace(NicName, "/", "_")
	nics(count) = NicName
	count = count + 1
Next

eth0_rx_bytes_1 = 0
eth0_rx_pkt_1 = 0
eth0_tx_bytes_1 = 0
eth0_tx_pkt_1 = 0
eth1_rx_bytes_1 = 0
eth1_rx_pkt_1 = 0
eth1_tx_bytes_1 = 0
eth1_tx_pkt_1 = 0
eth0_rx_bytes_2 = 0
eth0_rx_pkt_2 = 0
eth0_tx_bytes_2 = 0
eth0_tx_pkt_2 = 0
eth1_rx_bytes_2 = 0
eth1_rx_pkt_2 = 0
eth1_tx_bytes_2 = 0
eth1_tx_pkt_2 = 0

Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface",,48)
NicIndex = 0
For Each objItem in colItems
	For Each nic in nics
		if objItem.Name = nic then
		    if NicIndex = 0 then
				eth0_rx_bytes_1 = objItem.BytesReceivedPersec
				eth0_tx_bytes_1 = objItem.BytesSentPersec
				eth0_rx_pkt_1 = objItem.PacketsReceivedPersec
				eth0_tx_pkt_1 = objItem.PacketsSentPersec
			end if
			if NicIndex = 1 then
				eth0_rx_bytes_1 = objItem.BytesReceivedPersec
				eth0_tx_bytes_1 = objItem.BytesSentPersec
				eth0_rx_pkt_1 = objItem.PacketsReceivedPersec
				eth0_tx_pkt_1 = objItem.PacketsSentPersec
			end if
			NicIndex = NicIndex + 1
		End If
	Next
Next
Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfRawData_Tcpip_TCPv4",,48)
For Each objItem in colItems 
    passive_opens_1 = objItem.ConnectionsPassive
Next

wscript.Sleep 1000 * FLOW_USAGE_INTERVAL

Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface",,48)
NicIndex = 0
For Each objItem in colItems
	For Each nic in nics
		if objItem.Name = nic then
		    if NicIndex = 0 then
				eth0_rx_bytes_2 = objItem.BytesReceivedPersec
				eth0_tx_bytes_2 = objItem.BytesSentPersec
				eth0_rx_pkt_2 = objItem.PacketsReceivedPersec
				eth0_tx_pkt_2 = objItem.PacketsSentPersec
			end if
			if NicIndex = 1 then
				eth0_rx_bytes_2 = objItem.BytesReceivedPersec
				eth0_tx_bytes_2 = objItem.BytesSentPersec
				eth0_rx_pkt_2 = objItem.PacketsReceivedPersec
				eth0_tx_pkt_2 = objItem.PacketsSentPersec
			end if
			NicIndex = NicIndex + 1
		End If
	Next
Next
Set colItems = objWMIService.ExecQuery("SELECT * FROM Win32_PerfRawData_Tcpip_TCPv4",,48)
For Each objItem in colItems 
    passive_opens_2 = objItem.ConnectionsPassive
    tcp_cur_estab = objItem.ConnectionsEstablished
Next

eth0_rx_bits_rate = ((eth0_rx_bytes_2 - eth0_rx_bytes_1) * 8) / FLOW_USAGE_INTERVAL
eth0_rx_pkt_rate = (eth0_rx_pkt_2 - eth0_rx_pkt_1) / FLOW_USAGE_INTERVAL
eth0_tx_bits_rate = ((eth0_tx_bytes_2 - eth0_tx_bytes_1) * 8) / FLOW_USAGE_INTERVAL
eth0_tx_pkt_rate = (eth0_tx_pkt_2 - eth0_tx_pkt_1) / FLOW_USAGE_INTERVAL

eth1_rx_bits_rate = ((eth1_rx_bytes_2 - eth1_rx_bytes_1) * 8) / FLOW_USAGE_INTERVAL
eth1_rx_pkt_rate = (eth1_rx_pkt_2 - eth1_rx_pkt_1) / FLOW_USAGE_INTERVAL
eth1_tx_bits_rate = ((eth1_tx_bytes_2 - eth1_tx_bytes_1) * 8) / FLOW_USAGE_INTERVAL
eth1_tx_pkt_rate = (eth1_tx_pkt_2 - eth1_tx_pkt_1) / FLOW_USAGE_INTERVAL

tcp_passive_opens_rate = (passive_opens_2 - passive_opens_1) / FLOW_USAGE_INTERVAL

wscript.echo Round(eth0_rx_bits_rate) & chr(9) & _
            Round(eth0_rx_pkt_rate) & chr(9) & _
            Round(eth0_tx_bits_rate) & chr(9) & _
            Round(eth0_tx_pkt_rate) & chr(9) & _
            Round(eth1_rx_bits_rate) & chr(9) & _
            Round(eth1_rx_pkt_rate) & chr(9) & _
            Round(eth1_tx_bits_rate) & chr(9) & _
            Round(eth1_tx_pkt_rate) & chr(9) & _
            Round(tcp_passive_opens_rate) & chr(9) & _
            tcp_cur_estab
