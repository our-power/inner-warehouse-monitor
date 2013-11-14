#!/bin/sh

DISK_USAGE=$(df -l)

#ROOT_USAGE=$(df -l | grep '/$' | awk '{printf "%.2f", ($4 * 1.0)/$2}')
#DATA_USAGE=$(df -l | grep '/data$' | awk '{printf "%.2f", ($4 * 1.0)/$2}')

ROOT_USABLE=-1.00

ROOT_TOTAL=$(echo "$DISK_USAGE" | grep '/$' | awk '{print $2}')
if [ "$ROOT_TOTAL" ]; then
	ROOT_AVAL=$(echo "$DISK_USAGE" | grep '/$' | awk '{print $4}')
	ROOT_USABLE=$(expr "$ROOT_AVAL/$ROOT_TOTAL*100" | bc -l)
	ROOT_USABLE=$(awk -v USABLE=$ROOT_USABLE BEGIN'{printf "%.2f", USABLE}')
fi

DATA_USABLE=-1.00
DATA1_USABLE=-1.00

DATA_TOTAL=$(echo "$DISK_USAGE" | grep '/data$' | awk '{print $2}')
if [ "$DATA_TOTAL" ]; then
	DATA_AVAL=$(echo "$DISK_USAGE" | grep '/data$' | awk '{print $4}')
	DATA_USABLE=$(expr "$DATA_AVAL/$DATA_TOTAL*100" | bc -l)
	DATA_USABLE=$(awk -v USABLE=$DATA_USABLE BEGIN'{printf "%.2f", USABLE}')
fi

DATA1_TOTAL=$(echo "$DISK_USAGE" | grep '/data1$' | awk '{print $2}')
if [ "$DATA1_TOTAL" ]; then
	DATA1_AVAL=$(echo "$DISK_USAGE" | grep '/data1$' | awk '{print $4}')
	DATA1_USABLE=$(expr "$DATA1_AVAL/$DATA1_TOTAL*100" | bc -l)
	DATA1_USABLE=$(awk -v USABLE=$DATA1_USABLE BEGIN'{printf "%.2f", USABLE}')
fi


echo -e "$ROOT_USABLE""\t""$DATA_USABLE""\t""$DATA1_USABLE"
