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
    git clone https://github.com/10gqn/ops-manager-cloudfoundry

    ```
2. Install GO, CF CLI, OM

3. CD out of the repository root. 

4. Create folder artefacts and version
    - 10gen
    |-> artefacts
    |-> ops-manager-cloudfoundry
    |-> version


5. place file "number" into version folder and add version number 1.x.x


6. Execute script on bash ops-manager-cloudfoundry/ci/tasks/build-tile/run.sh 


7. check that artefacts contain tile build with correct timestamp

## Deploy Tile

1. CD out of the repository root.

2. Make sure PCF Ops Manager has no MongoDB Tile deployed with the same or higher version and that BOSH has release has no conflictin error
    To check BOSH login to PCF Ops Manager box; run "bosh ... releases" 
    Make sure that list of saved releases has no conflicting mongodb versions
    to Delete "bosh ..... delete-release mongodb/1.X.X"

3. Make shure config file has correct configuration:
    - Check for MongoDB Ops Manager URL, USer and API Key

4. Run CI script to deploy tile ops-manager-cloudfoundry/ci/tasks/deploy-tile/run.sh 1.X.X '<pcf-opsManager-password>


