#!/usr/local/bin/dumb-init /bin/bash

#save config.pie from pipeline
make_env_config() {
    file=$1
    if [ -f $file ] ; then
        rm $file
    fi
    echo "$CONFIG" >> config.pie
}


