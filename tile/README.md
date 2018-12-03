# How to build MongoDB On Demand for PCF using tile generator

## Prerequisites

### Install Dependencies

The following should be installed on your local machine
- [bosh-cli](https://bosh.io/docs/cli-v2.html)
- [tile-generator](https://github.com/emiloserdov/tile-generator)
- [yq](https://github.com/mikefarah/yq)

## Build Tile

1. Check out PCF MongoDB On Demand tile generator repo:

    ```bash
    git clone https://github.com/Altoros/ops-manager-cloudfoundry.git
    ```

2. Download the following releases into [tile/resources](https://github.com/Altoros/ops-manager-cloudfoundry/tree/master/tile/resources) folder.

    - [Pivotal Cloud Foundry On Demand Service Broker Release](https://s3.amazonaws.com/mongodb-tile-ci/on-demand-service-broker-0.22.0-ubuntu-trusty-3586.36.tgz)
    - [Pivotal Cloud Foundry MongoDB Helpers Release](https://s3.amazonaws.com/mongodb-tile-ci/pcf-mongodb-helpers-0.0.1.tgz)
    - [Pivotal Cloud Foundry Syslog Migration Release](https://s3.amazonaws.com/mongodb-tile-ci/syslog-migration-11.1.1-ubuntu-trusty-3586.36.tgz)
    - [Pivotal Cloud Foundry BOSH Process Manager Release](https://s3.amazonaws.com/mongodb-tile-ci/bpm-release-0.12.2-ubuntu-trusty-3586.36.tgz)

3. Update build version:

   Notice, that it should be in X.X.X format (for example: 3.2.1) 
    ```bash
    export VERSION_NUMBER=
    ```

4. Create release tarball and place it into [tile/resources](https://github.com/Altoros/ops-manager-cloudfoundry/tree/master/tile/resources) folder.

    ```bash
    cd ops-manager-cloudfoundry
    tarball_path="$PWD/tile/resources/mongodb-${VERSION_NUMBER}.tgz"
    bosh -n create-release --sha2 --tarball="$tarball_path" --version="${VERSION_NUMBER}"
    ```

5. Edit `tile.yml` file and check path and versions for ops-manager-cloudfoundry.
   Ensure that tile file configured with version which was specified in step 4.

    ```bash
    cd ops-manager-cloudfoundry/tile
    yq w -i tile.yml packages.[4].path "$(ls resources/mongodb-*.tgz)"
    yq w -i tile.yml packages.[4].jobs[0].properties.service_deployment.releases[0].version "${VERSION_NUMBER}"
    yq w -i tile.yml runtime_configs[0].runtime_config.releases[0].version "${VERSION_NUMBER}"
    ```

6. Build your tile, after build is finished you can find product file in `tile/product` subdirectory.

    ```bash
    tile build "${VERSION_NUMBER}"
    ```
