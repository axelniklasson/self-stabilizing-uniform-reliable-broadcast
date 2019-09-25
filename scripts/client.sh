# make sure to run this from the root of the project directory
#
# ./scripts/client.sh NODE_ID MSG_COUNT

if [ $# -lt 2 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/client.sh NODE_ID MSG_COUNT"
    exit 2
fi

NODE_ID=$1
MSG_COUNT=$2
PORT=$((4000+$NODE_ID))
curl -d '{"reqCount": '$MSG_COUNT'}' -H "Content-Type: application/json" -X POST http://localhost:$PORT/client/launch