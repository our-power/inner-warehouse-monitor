#!/bin/sh

function get_one_cpu_usage()
{
    CPULOG_1=$1
    CPULOG_2=$2

    SYS_IDLE_1=$(echo $CPULOG_1 | awk -F"," '{print $4}')
    Total_1=$(echo $CPULOG_1 | awk -F"," '{print $1+$2+$3+$4+$5+$6+$7}')

    SYS_IDLE_2=$(echo $CPULOG_2 | awk -F"," '{print $4}')
    Total_2=$(echo $CPULOG_2 | awk -F"," '{print $1+$2+$3+$4+$5+$6+$7}')

    SYS_IDLE=`expr $SYS_IDLE_2 - $SYS_IDLE_1`

    Total=`expr $Total_2 - $Total_1`
    SYS_USAGE=`expr $SYS_IDLE/$Total*100 |bc -l`
    SYS_Rate=`expr 100-$SYS_USAGE |bc -l`

    #Disp_SYS_Rate=`expr "scale=2; $SYS_Rate/1" |bc -l`
    Disp_SYS_Rate=$(awk -v Rate=$SYS_Rate BEGIN'{printf "%.2f", Rate}')
    echo $Disp_SYS_Rate
}

function get_cpu_usage()
{
    ##echo user nice system idle iowait irq softirq
    CPULOG_ALL_1=$(cat /proc/stat | grep 'cpu ' | awk '{print $2","$3","$4","$5","$6","$7","$8" "}')

    sleep 15

    CPULOG_ALL_2=$(cat /proc/stat | grep 'cpu ' | awk '{print $2","$3","$4","$5","$6","$7","$8" "}')

    CPULOG_ALL_1=$(echo "$CPULOG_ALL_1" |  tr -d '\n' | sed 's/[[:space:]]*$//')
    CPULOG_ALL_2=$(echo "$CPULOG_ALL_2" |  tr -d '\n' | sed 's/[[:space:]]*$//')

    ## string to array

    CPULOG_ARRAY_1=($CPULOG_ALL_1)
    CPULOG_ARRAY_2=($CPULOG_ALL_2)

    ARR_LEN=${#CPULOG_ARRAY_1[*]}

    CPU_USAGE=""

    for ((i=0;i<$ARR_LEN;i++))
    do
        SINGLE_CPU_USAGE=`get_one_cpu_usage "${CPULOG_ARRAY_1[$i]}" "${CPULOG_ARRAY_2[$i]}"`

        if [ $i -eq 0 ]; then
            CPU_USAGE="$SINGLE_CPU_USAGE"
        else
            CPU_USAGE="$CPU_USAGE"" ""$SINGLE_CPU_USAGE"
        fi
    done

    echo -e $CPU_USAGE
}

CPU_USAGE=0

for ((i=0;i<20;i++))
do
    CPU_USAGE_TMP=$(get_cpu_usage)
    CPU_USAGE_TMP=($CPU_USAGE_TMP)

    if [[ "${CPU_USAGE_TMP[0]}" > "$CPU_USAGE" ]] ; then
        CPU_USAGE=${CPU_USAGE_TMP[0]}
    fi
done

LOAD_AVG=$(cat /proc/loadavg | awk '{print $2}')
LOAD_AVG=$(awk -v Load=$LOAD_AVG BEGIN'{printf "%d", Load * 100}')

echo -e "$CPU_USAGE""\t""$LOAD_AVG"
