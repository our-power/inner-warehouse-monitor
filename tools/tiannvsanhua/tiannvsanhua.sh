#!/bin/bash

REPO_URL="git@10.24.178.60:inner-warehouse-monitor.git"
TARGET_DIR="./inner-warehouse-monitor"

type git > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "You need Git"
    exit 1
fi

type go > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "You need go compiler!"
    exit 1
fi

function fetch_code()
{
    if [! -d $TARGET_DIR ]; then
        git clone $REPO_URL
    fi

    pushd $TARGET_DIR

    if [ $(git branch | grep "\* multiSqlite" | wc -l) -eq 0 ]; then
        git checkout multiSqlite
    fi

    git pull origin multiSqlite
    popd
}

function compile_all()
{
    pushd $TARGET_DIR

    # compile NSQ
    $ROOT=`pwd`
    NSQ_DIR="$ROOT/nsq"

    if [ ! -d $NSQ_DIR ]; then
        echo "Not Exist NSQ source code!"
        exit 1
    fi

    export GOPATH=$NSQ_DIR

    go install github.com/bitly/nsq/nsqd
    go install github.com/bitly/nsq/nsqadmin
    go install github.com/bitly/nsq/nsqlookupd

    echo "Finish to compile nsq!"

    # compile NSQ client
    CLIENT_DIR="$ROOT/client"
    if [ ! -d "$CLIENT_DIR" ]; then
        echo "Not Exist NSQ client source code!"
        exit 1
    fi
    export GOPATH=$CLIENT_DIR
    go install nsq_client
    echo "Finish to compile nsq client"

    # compile local server webapp
    LOCAL_SERVER_WEBAPP_DIR="$ROOT/webapp/local_server"
    LOCAL_SERVER_WEBAPP_BIN="$LOCAL_SERVER_WEBAPP_DIR/src/big_brother/big_brother"
    if [ ! -d "$LOCAL_SERVER_WEBAPP_DIR" ]; then
        echo "Not Exist local server webapp source code!"
        exit 1
    fi

    export GOPATH=$LOCAL_SERVER_WEBAPP_DIR
    go install big_brother
    mv $LOCAL_SERVER_WEBAPP_DIR/bin/big_brother $LOCAL_SERVER_WEBAPP_BIN
    echo "Finish to compile local server webapp!"

    popd
}

fetch_code
compile_all

# 还写个毛线的代码啊
