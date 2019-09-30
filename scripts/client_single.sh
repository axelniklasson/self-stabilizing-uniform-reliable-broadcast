# make sure to run this from the root of the project directory
#
# ./scripts/client_single.sh HOST NODE_ID MSG_COUNT

if [ $# -lt 3 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/client_single.sh HOST NODE_ID MSG_COUNT"
    exit 2
fi

HOST=$1
NODE_ID=$2
MSG_COUNT=$3
PORT=$((4000+$NODE_ID))
curl -d '{"reqCount": '$MSG_COUNT'}' -H "Content-Type: application/json" -X POST http://$HOST:$PORT/client/launch
