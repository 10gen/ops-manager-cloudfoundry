FROM ubuntu
RUN apt-get update &&   DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends     build-essential git ruby-dev openssh-client     unzip wget curl rsync vim-nox python     && apt-get autoremove     && apt-get clean     && rm -rf /var/lib/apt/lists/* 
RUN wget https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 &&   mv jq-linux64 /usr/local/bin/jq &&   chmod +x /usr/local/bin/jq

ENV BOSH_CLI_VERSION=5.4.0
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-${BOSH_CLI_VERSION}-linux-amd64 -O bosh &&   mv bosh /usr/local/bin/bosh &&   chmod +x /usr/local/bin/bosh
ENV OM_VERSION=0.47.0 
RUN wget https://github.com/pivotal-cf/om/releases/download/${OM_VERSION}/om-linux &&   mv om-linux /usr/local/bin/om &&   chmod +x /usr/local/bin/om
ENV CREDHUB_CLI_VERSION=2.2.0 
RUN wget https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/${CREDHUB_CLI_VERSION}/credhub-linux-${CREDHUB_CLI_VERSION}.tgz -O credhub.tgz &&   tar xzf ./credhub.tgz &&   mv credhub /usr/local/bin/credhub &&   chmod +x /usr/local/bin/credhub &&   rm -f ./credhub.tgz
ENV CF_CLI_VERSION=6.41.0
RUN wget https://s3-us-west-1.amazonaws.com/cf-cli-releases/releases/v${CF_CLI_VERSION}/cf-cli_${CF_CLI_VERSION}_linux_x86-64.tgz -O cf.tgz &&   tar xzf ./cf.tgz cf &&   mv cf /usr/local/bin/cf &&   chmod +x /usr/local/bin/cf &&   rm -f ./cf.tgz
ENV PIVNET_CLI_VERSION=0.0.55
RUN wget https://github.com/pivotal-cf/pivnet-cli/releases/download/v${PIVNET_CLI_VERSION}/pivnet-linux-amd64-${PIVNET_CLI_VERSION} -O pivnet &&   mv pivnet /usr/local/bin/pivnet &&   chmod +x /usr/local/bin/pivnet
ENV DUMB_INIT_VERSION=1.2.2
RUN wget https://github.com/Yelp/dumb-init/releases/download/v${DUMB_INIT_VERSION}/dumb-init_${DUMB_INIT_VERSION}_amd64 -O dumb-init &&   mv dumb-init /usr/local/bin/dumb-init &&   chmod +x /usr/local/bin/dumb-init
ENV YQ_VERSION=2.2.0
RUN wget https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -O yq &&   mv yq /usr/local/bin/yq &&   chmod +x /usr/local/bin/yq
ENV TILE_GENERATOR_VERSION=12.0.8
RUN wget https://github.com/cf-platform-eng/tile-generator/releases/download/v${TILE_GENERATOR_VERSION}/pcf_linux-64bit -O pcf &&   mv pcf /usr/local/bin/pcf &&   chmod +x /usr/local/bin/pcf 
RUN wget https://github.com/cf-platform-eng/tile-generator/releases/download/v${TILE_GENERATOR_VERSION}/tile_linux-64bit -O tile &&   mv tile /usr/local/bin/tile &&   chmod +x /usr/local/bin/tile
RUN curl -sSL https://sdk.cloud.google.com | bash 
ENV PATH=/root/google-cloud-sdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV GOPATH=/go
ENV PATH=/go/bin:/usr/local/go/bin:/root/google-cloud-sdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV GOLANG_VERSION=1.11.4
RUN wget https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz -P /tmp &&   tar xzf /tmp/go${GOLANG_VERSION}.linux-amd64.tar.gz -C /usr/local &&   mkdir $GOPATH &&   go get -u github.com/onsi/ginkgo/... &&   rm -rf /tmp/* 
ENV TERRAFORM_VERSION=0.11.11
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -P /tmp &&   unzip /tmp/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin &&   chmod +x /usr/local/bin/terraform &&   rm -rf /tmp/* 

ENV BBL_VERSION=7.0.0
RUN wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v${BBL_VERSION}/bbl-v${BBL_VERSION}_linux_x86-64 -O bbl &&   mv bbl /usr/local/bin/bbl &&   chmod +x /usr/local/bin/bbl
