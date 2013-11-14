#!/bin/sh

MEM_INFO=$(cat /proc/meminfo)

MEM_TOTAL=$(echo "$MEM_INFO" | grep '^MemTotal' | awk '{print $2}')
MEM_FREE=$(echo "$MEM_INFO" | grep '^MemFree' | awk '{print $2}')
BUFFERS=$(echo "$MEM_INFO" | grep '^Buffers' | awk '{print $2}')
CACHED=$(echo "$MEM_INFO" | grep '^Cached' | awk '{print $2}')
SWAP_TOTAL=$(echo "$MEM_INFO" | grep '^SwapTotal' | awk '{print $2}')
SWAP_FREE=$(echo "$MEM_INFO" | grep '^SwapFree' | awk '{print $2}')

APP_USED=`expr $MEM_TOTAL - $MEM_FREE - $BUFFERS - $CACHED`
APP_USED=`expr "scale=1;$APP_USED/1024" | bc`

MEM_USED=`expr $MEM_TOTAL - $MEM_FREE`

SWAP_USED=`expr $SWAP_TOTAL - $SWAP_FREE`

MEM_PERCENTAGE=`expr "scale=2;$MEM_USED/$MEM_TOTAL*100" | bc`

echo -e "$APP_USED\t$MEM_USED\t$SWAP_USED\t$MEM_PERCENTAGE"

