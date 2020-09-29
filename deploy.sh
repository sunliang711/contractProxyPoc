#!/bin/bash
rpath="$(readlink ${BASH_SOURCE})"
if [ -z "$rpath" ];then
    rpath=${BASH_SOURCE}
fi
this="$(cd $(dirname $rpath) && pwd)"
cd "$this"
export PATH=$PATH:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

user="${SUDO_USER:-$(whoami)}"
home="$(eval echo ~$user)"

red=$(tput setaf 1)
green=$(tput setaf 2)
yellow=$(tput setaf 3)
blue=$(tput setaf 4)
cyan=$(tput setaf 5)
bold=$(tput bold)
reset=$(tput sgr0)
function _runAsRoot(){
    verbose=0
    while getopts ":v" opt;do
        case "$opt" in
            v)
                verbose=1
                ;;
            \?)
                echo "Unknown option: \"$OPTARG\""
                exit 1
                ;;
        esac
    done
    shift $((OPTIND-1))
    cmd="$@"
    if [ -z "$cmd" ];then
        echo "${red}Need cmd${reset}"
        exit 1
    fi

    if [ "$verbose" -eq 1 ];then
        echo "run cmd:\"${red}$cmd${reset}\" as root."
    fi

    if (($EUID==0));then
        sh -c "$cmd"
    else
        if ! command -v sudo >/dev/null 2>&1;then
            echo "Need sudo cmd"
            exit 1
        fi
        sudo sh -c "$cmd"
    fi
}
editor=vi
if command -v vim >/dev/null 2>&1;then
    editor=vim
fi
if command -v nvim >/dev/null 2>&1;then
    editor=nvim
fi
###############################################################################
# write your code below (just define function[s])
# function is hidden when begin with '_'
###############################################################################
# TODO
init(){
    cd ${this}/node/

    if [ -d datadirs ];then
        echo "clear datadir..."
        /bin/rm -rf datadirs
    fi

    echo "init genesis..."
    ./00-init.sh config

    echo "coppy key file..."
    cp ../cmd/exportUTC/UTC--2020-09-28T05-17-16.010505420Z--77c3e6dedbf1157026c7d8d4b66d6b19004b7b40 datadirs/node1-datadir/keystore

    cd ${this}
}

startNode(){
    cd ${this}/node
    ./02-start.sh config
}

console(){
    cd ${this}/node
    ./03-console.sh config
}

log(){
    cd ${this}/node
    ./04-logfile.sh config
}

# contract(){
#     cd ${this}/cmd/createContract
#     echo "create contract anc..."
#     go run main.go  -b anc.bytecode --rpc http://localhost:8546 --sk 1979d4ce44d5fa7181310b0f8108a701c9bae2f793e86d0216d549eda714b002
#     echo "create contract Proxy..."
#     go run main.go  -b Proxy.bytecode --rpc http://localhost:8546 --sk 1979d4ce44d5fa7181310b0f8108a701c9bae2f793e86d0216d549eda714b002
# }
_build(){
    local name=${1:?'missing foldername'}
    cd ${this}/cmd/${name}
    echo "build ${name}..."
    go build -o ${name} main.go

    cd ${this}
}

buildall(){
    _build createContract
    _build encodeAbi
    _build exportUTC
    _build readContract
    _build writeContract
}

rpcURL="http://localhost:8546"
create(){
    #Address:77c3e6dedbf1157026c7d8d4b66d6b19004b7b40 
    #PrivateKey:1979d4ce44d5fa7181310b0f8108a701c9bae2f793e86d0216d549eda714b002 

    sk1=$(perl -lne 'print $1 if /PrivateKey:(\w+)/' cmd/exportUTC/account)
    addr1=$(perl -lne 'print $1 if /Address:(\w+)/' cmd/exportUTC/account)

    sk2=$(perl -lne 'print $1 if /PrivateKey:(\w+)/' cmd/exportUTC/account2)
    addr2=$(perl -lne 'print $1 if /Address:(\w+)/' cmd/exportUTC/account2)

cat<<-EOF
	rpcURL: ${rpcURL}
	addr1: ${addr1}
	sk1: ${sk1}

	addr2: ${addr2}
	sk2: ${sk2}

EOF
    local name=anc
    # create contract
    echo "create contract '${name}' output to : contracts/${name}.output"
    ${this}/cmd/createContract/createContract -b contracts/${name}.bytecode --rpc ${rpcURL} --sk ${sk1} -o contracts/${name}.output

    name=Proxy
    echo "create contract '${name}' output to : contracts/${name}.output"
    ${this}/cmd/createContract/createContract -b contracts/${name}.bytecode --rpc ${rpcURL} --sk ${sk1} -o contracts/${name}.output

}

read(){
    # read contract
    local contractName=${1:?'missing contract name'}
    local contractAddress=${2:?'missing contract address'}
    local fromaddr=${3:?'missing from addr'}
    local methodName=${4:?'missing method name'}
    local args=${5}
    if [ -n "${args}" ];then
        local parg="--args ${args}"
    fi
    ${this}/cmd/readContract/readContract --rpc ${rpcURL} --abi contracts/${contractName}.abi --addr ${contractAddress} --fromaddr ${fromaddr} --method=${methodName} ${parg}
}


write(){
    # write contract
    local contractName=${1:?'missing contract name'}
    local contractAddress=${2:?'missing contract address'}
    local sk=${3:?'missing sk'}
    local methodName=${4:?'missing method name'}
    local args=${5}
    if [ -n "${args}" ];then
        local parg="--args ${args}"
    fi
    ${this}/cmd/writeContract/writeContract --rpc ${rpcURL} --abi contracts/${contractName}.abi --addr ${contractAddress} --sk ${sk} --method ${methodName} ${parg}
}


em(){
    $editor $0
}

###############################################################################
# write your code above
###############################################################################
function _help(){
    cat<<EOF2
Usage: $(basename $0) ${bold}CMD${reset}

${bold}CMD${reset}:
EOF2
    # perl -lne 'print "\t$1" if /^\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE})
    # perl -lne 'print "\t$2" if /^\s*(function)?\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE}) | grep -v '^\t_'
    perl -lne 'print "\t$2" if /^\s*(function)?\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE}) | perl -lne "print if /^\t[^_]/"
}

function _loadENV(){
    if [ -z "$INIT_HTTP_PROXY" ];then
        echo "INIT_HTTP_PROXY is empty"
        echo -n "Enter http proxy: (if you need) "
        read INIT_HTTP_PROXY
    fi
    if [ -n "$INIT_HTTP_PROXY" ];then
        echo "set http proxy to $INIT_HTTP_PROXY"
        export http_proxy=$INIT_HTTP_PROXY
        export https_proxy=$INIT_HTTP_PROXY
        export HTTP_PROXY=$INIT_HTTP_PROXY
        export HTTPS_PROXY=$INIT_HTTP_PROXY
        git config --global http.proxy $INIT_HTTP_PROXY
        git config --global https.proxy $INIT_HTTP_PROXY
    else
        echo "No use http proxy"
    fi
}

function _unloadENV(){
    if [ -n "$https_proxy" ];then
        unset http_proxy
        unset https_proxy
        unset HTTP_PROXY
        unset HTTPS_PROXY
        git config --global --unset-all http.proxy
        git config --global --unset-all https.proxy
    fi
}


case "$1" in
     ""|-h|--help|help)
        _help
        ;;
    *)
        "$@"
        ;;
esac
