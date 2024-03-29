# make sure to run this from the root of the project directory
#
# ./scripts/client_distributed.sh MSG_COUNT HOSTS_PATH

if [ $# -lt 2 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/client_distributed.sh MSG_COUNT HOSTS_PATH"
    exit 2
fi

# splitAndGet str del idx
splitAndGet () {
    IFS=$2
    read -ra ADDR <<< "$1"
    echo ${ADDR[$3]}
}

MSG_COUNT=$1
HOSTS_PATH=$2

for l in $(cat $HOSTS_PATH)
do
    ID=$(splitAndGet $l ',' 0)
    HOSTNAME=$(splitAndGet $l ',' 1)

    PORT=$((4000+$ID))
    curl -d '{"reqCount": '$MSG_COUNT'}' -H "Content-Type: application/json" -X POST http://$HOSTNAME:$PORT/client/launch
    echo "Launched client that will inject $MSG_COUNT messages on $HOSTNAME"
done
