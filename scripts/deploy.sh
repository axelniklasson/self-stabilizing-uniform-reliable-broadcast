# make sure to run this from the root of the project directory
#
# ./scripts/deploy.sh NODE_COUNT GIT_BRANCH

if [ $# -lt 2 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/deploy.sh NODE_COUNT GIT_BRANCH [args]"
    exit 2
fi

plcli --skip-healthcheck --node-count $1 --git-branch $2 --shuffle-nodes --app-path \$HOME/go/src/github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast --prometheus-sd-path ./heimdall/prometheus/sd.json --node-exporter $3 deploy https://github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast.git
