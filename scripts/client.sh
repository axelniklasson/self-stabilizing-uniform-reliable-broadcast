# make sure to run this from the root of the project directory
#
# ./scripts/client.sh BROADCAST_ENDPOINT NUMBER_OF_MESSAGES

if [ $# -lt 2 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/client.sh BROADCAST_ENDPOINT NUMBER_OF_MESSAGES"
    exit 2
fi

ENDPOINT=$1
MSG_COUNT=$2

for (( i=0; i<=$(($MSG_COUNT-1)); i++ ))
do
    curl -d '{"text": "Message '$i'"}' -H "Content-Type: application/json" -X POST $ENDPOINT
done