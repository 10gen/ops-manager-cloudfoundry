#!/bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

for module in ops-manager-cloudfoundry/src/{mongodb-config-agent,mongodb-service-adapter,smoke-tests}; do
    pushd $module

    echo Linting "$module"

    # run linters
    golangci-lint run --modules-download-mode vendor --timeout 15m

    # fix whitespace
    find vendor -type f -exec grep -Iq . {} \; -print0 | xargs -0 sed -i 's/\r\n$/\n/'

    # get folder hash before running go mod vendor
    OLDHASHES=$(find vendor -type f -print0 | sort -z | xargs -0 sha1sum)
    OLDSUM=$(echo "$OLDHASHES" | sha1sum)

    go mod vendor

    # fix whitespace again (we don't care if it's wrong)
    find vendor -type f -exec grep -Iq . {} \; -print0 | xargs -0 sed -i 's/\r\n$/\n/'

    # get folder hash after go mod vendor
    NEWHASHES=$(find vendor -type f -print0 | sort -z | xargs -0 sha1sum)
    NEWSUM=$(echo "$NEWHASHES" | sha1sum)

    popd

    if [[ $OLDSUM != $NEWSUM ]]; then
        echo Vendor is out of date!
        LEFT=$(mktemp)
        RIGHT=$(mktemp)
        echo "$OLDHASHES" >$LEFT
        echo "$NEWHASHES" >$RIGHT
        echo Hash diff:
        # diff --context=0 $LEFT $RIGHT | grep -E "^\+ |^\- "
        diff $LEFT $RIGHT
        echo $OLDSUM
        echo $NEWSUM
        # exit 1
    fi
done
