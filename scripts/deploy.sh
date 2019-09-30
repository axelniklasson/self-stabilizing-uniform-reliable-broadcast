# make sure to run this from the root of the project directory
#
# ./scripts/deploy.sh NODE_COUNT GIT_BRANCH [blacklisted_nodes] [shuffle]

if [ $# -lt 2 ]; then
    echo 1>&2 "$0: not enough arguments, run as ./scripts/deploy.sh NODE_COUNT GIT_BRANCH [blacklisted_nodes] [shuffle]"
    exit 2
fi

if [ $4 == "shuffle" ]; then
    plcli --node-count $1 --git-branch $2 --shuffle-nodes --blacklist $3 --app-path \$HOME/go/src/github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast --prometheus-sd-path ./heimdall/prometheus/sd.json --node-exporter deploy https://github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast.git
else
    plcli --node-count $1 --git-branch $2 --blacklist $3 --app-path \$HOME/go/src/github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast --prometheus-sd-path ./heimdall/prometheus/sd.json --node-exporter deploy https://github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast.git
fi
