#!/bin/sh

freshclam --daemon --checks=$NO_OF_CHECKS_FOR_DB_UPDATE &
clamd &

/usr/bin/hawk -address $IPADDR -port $PORT  $INDEXES &

#get All PIDs
pid_list=$(jobs -p)

exit_code=0

function shutdownAll() {
    trap "" SIGINT

    for pid in $pid_list; do
        if ! kill -0 "$pid" 2> /dev/null; then
            wait "$pid"
            exit_code=$?
        fi
    done

    kill "$pid_list" 2> /dev/null
}

# shutdown All
trap shutdownAll SIGINT
wait -n

# exit with code
exit $exit_code