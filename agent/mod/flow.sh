#!/bin/sh

PACKET_LOG_ALL_1=$(cat /proc/net/dev | grep eth | awk -F":" '{if (NR<3) print $2}' | awk '{print $1","$2","$9","$10" "}')

TCP_LOG_1=$(cat /proc/net/snmp | grep "Tcp:" | awk -F":" '{if (NR==2) print $2}')

sleep 300

PACKET_LOG_ALL_2=$(cat /proc/net/dev | grep eth | awk -F":" '{if (NR<3) print $2}' | awk '{print $1","$2","$9","$10" "}')

TCP_LOG_2=$(cat /proc/net/snmp | grep "Tcp:" | awk -F":" '{if (NR==2) print $2}')

PACKET_LOG_ALL_1=$(echo "$PACKET_LOG_ALL_1" |  tr -d '\n' | sed 's/[[:space:]]*$//')
PACKET_LOG_ALL_2=$(echo "$PACKET_LOG_ALL_2" |  tr -d '\n' | sed 's/[[:space:]]*$//')

#echo $PACKET_LOG_ALL_1
#echo $PACKET_LOG_ALL_2

PACKET_LOG_ARRAY_1=($PACKET_LOG_ALL_1)
PACKET_LOG_ARRAY_2=($PACKET_LOG_ALL_2)

ARR_LEN=${#PACKET_LOG_ARRAY_1[*]}

#for ((i=0;i<$ARR_LEN;i++))
#do
#    echo "${PACKET_LOG_ARRAY_1[$i]}"
#done

function get_eth_info()
{
    PKT_LOG_1=$1
    PKT_LOG_2=$2

    RX_BYTE_1=$(echo $PKT_LOG_1 | awk -F"," '{print $1}')
    RX_PKT_1=$(echo $PKT_LOG_1 | awk -F"," '{print $2}')
    TX_BYTE_1=$(echo $PKT_LOG_1 | awk -F"," '{print $3}')
    TX_PKT_1=$(echo $PKT_LOG_1 | awk -F"," '{print $4}')

    RX_BYTE_2=$(echo $PKT_LOG_2 | awk -F"," '{print $1}')
    RX_PKT_2=$(echo $PKT_LOG_2 | awk -F"," '{print $2}')
    TX_BYTE_2=$(echo $PKT_LOG_2 | awk -F"," '{print $3}')
    TX_PKT_2=$(echo $PKT_LOG_2 | awk -F"," '{print $4}')

    RX_BYTE_RATE=$(awk -v v1=$RX_BYTE_1 -v v2=$RX_BYTE_2 BEGIN'{printf "%d", (v2 - v1) * 8 / 300}')
    RX_PKT_RATE=$(awk -v v1=$RX_PKT_1 -v v2=$RX_PKT_2 BEGIN'{printf "%d", (v2 - v1) / 300}')

    TX_BYTE_RATE=$(awk -v v1=$TX_BYTE_1 -v v2=$TX_BYTE_2 BEGIN'{printf "%d", (v2 - v1) * 8 / 300}')
    TX_PKT_RATE=$(awk -v v1=$TX_PKT_1 -v v2=$TX_PKT_2 BEGIN'{printf "%d", (v2 - v1) / 300}')

    echo -e "$RX_BYTE_RATE""\t""$RX_PKT_RATE""\t""$TX_BYTE_RATE""\t""$TX_PKT_RATE"
}

ETH_USAGE=""

for ((i=0;i<$ARR_LEN;i++))
do
    SINGLE_ETH_USAGE=`get_eth_info "${PACKET_LOG_ARRAY_1[$i]}" "${PACKET_LOG_ARRAY_2[$i]}"`

    if [ $i -eq 0 ]; then
        ETH_USAGE="$SINGLE_ETH_USAGE"
    else
        ETH_USAGE="$ETH_USAGE""\t""$SINGLE_ETH_USAGE"
    fi
done

if [ $ARR_LEN -lt 2 ]; then
    ETH_USAGE="$ETH_USAGE""\t0\t0\t0\t0"
fi

PASSIVE_OPENS_1=$(echo $TCP_LOG_1 | awk '{print $6}')
PASSIVE_OPENS_2=$(echo $TCP_LOG_2 | awk '{print $6}')
PASSIVE_OPENS_INC=$(expr $PASSIVE_OPENS_2 - $PASSIVE_OPENS_1)
PASSIVE_OPENS_RATE=$(awk -v Passive=$PASSIVE_OPENS_INC BEGIN'{printf "%d", Passive / 300}')

CUR_TCP_ESTB=$(echo $TCP_LOG_2 | awk '{print $9}')

echo -e "$ETH_USAGE""\t""$PASSIVE_OPENS_RATE""\t""$CUR_TCP_ESTB"
