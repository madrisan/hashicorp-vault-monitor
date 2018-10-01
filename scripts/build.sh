#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

# Set build tags
BUILD_TAGS="${BUILD_TAGS:-"hashicorp-vault-monitor"}"

# Get the git commit
GIT_COMMIT="$(git rev-parse HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd netbsd solaris}
XC_OSARCH=${XC_OSARCH:-"linux/amd64 linux/arm linux/arm64 darwin/amd64 windows/amd64 freebsd/amd64 freebsd/arm openbsd/amd64 openbsd/arm netbsd/amd64 solaris/amd64"}

# If its dev mode, only build for ourself
if [ "${VAULT_MONITOR_DEV_BUILD}x" != "x" ] && [ "${XC_OSARCH}x" == "x" ]; then
    XC_OS=$(go env GOOS)
    XC_ARCH=$(go env GOARCH)
    XC_OSARCH=$(go env GOOS)/$(go env GOARCH)
elif [ "${XC_OSARCH}x" != "x" ]; then
    IFS='/' read -ra SPLITXC <<< "${XC_OSARCH}"
        DEV_PLATFORM="./pkg/${SPLITXC[0]}_${SPLITXC[1]}"
fi

GOPATH=${GOPATH:-$(go env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Build!
echo "==> Building..."
gox \
    -osarch="${XC_OSARCH}" \
    -gcflags "${GCFLAGS}" \
    -ldflags "${LD_FLAGS}-X github.com/madrisan/hashicorp-vault-monitor/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}'" \
    -output "pkg/{{.OS}}_{{.Arch}}/hashicorp-vault-monitor" \
    -tags="${BUILD_TAGS}" \
    .

# Move all the compiled things to the $GOPATH/bin
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Copy our OS/Arch to the bin/ directory
DEV_PLATFORM=${DEV_PLATFORM:-"./pkg/$(go env GOOS)_$(go env GOARCH)"}
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f >/dev/null 2>&1); do
    cp ${F} bin/
    cp ${F} ${MAIN_GOPATH}/bin/
done

if [ "${VAULT_MONITOR_DEV_BUILD}x" = "x" ]; then
    # Zip and copy to the dist dir
    echo "==> Packaging..."
    for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
        OSARCH=$(basename ${PLATFORM})
        echo "--> ${OSARCH}"

        pushd $PLATFORM >/dev/null 2>&1
        zip ../${OSARCH}.zip ./*
        popd >/dev/null 2>&1
    done
fi

# Done!
echo
echo "==> Results:"
ls -hl bin/
