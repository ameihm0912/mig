#!/bin/bash
# Julien Vehent - 2014
# watch the MIG source code directory and rebuild the components
# when a file is saved.

if [[ -z $1 || ! -d "$1" ]]; then
    cat <<EOF
usage: $0 ./src/mig

$0 watches the content of the MIG source directory and
rebuilds components when a source code file is modified.
It must be ran from the directory of the MIG Makefile.
EOF
    exit 1
fi

echo "starting inotify listener on $1"
# feed the inotify events into a while loop that creates
# the variables 'date' 'time' 'dir' 'file' and 'event'
inotifywait -mr --timefmt '%d/%m/%y %H:%M' --format '%T %w %f %e' \
-e modify $1 \
| while read date time dir file event
do
    if [[ "$file" =~ \.go$ && "$dir" =~ src\/mig ]]; then
        dontskip=true
    else
        #echo skipping $date $time $event $dir $file
        continue
    fi

    #echo NEW EVENT: $event $dir $file

    if [[ "$dir" =~ src\/mig\/$ ]]; then
        make mig-agent && \
        make mig-action-generator && \
        make mig-action-verifier && \
        make mig-api && \
        make mig-scheduler

    elif [[ "$dir" =~ agent && "$file" != "configuration.go" ]] ; then
        make mig-agent

    elif [[ "$dir" =~ api ]] ; then
        make mig-api

    elif [[ "$dir" =~ generator ]] ; then
        make mig-action-generator

    elif [[ "$dir" =~ verifier ]] ; then
        make mig-action-verifier

    elif [[ "$dir" =~ pgp ]] ; then
        make mig-agent && \
        make mig-action-generator && \
        make mig-action-verifier && \
        make mig-api

    elif [[ "$dir" =~ scheduler ]] ; then
        make mig-scheduler

    fi
done
