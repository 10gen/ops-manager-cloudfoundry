package adapter

import (
	"net/http"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
)

func NewPCGCClient(uri, username, password string, acceptedStatusCodes ...int) opsmanager.Client {
	acceptedStatusCodes = append([]int{http.StatusOK, http.StatusCreated}, acceptedStatusCodes...)
	return opsmanager.NewClient(
		opsmanager.WithResolver(httpclient.NewURLResolverWithPrefix(uri, "api/public/v1.0")),
		opsmanager.WithHTTPClient(httpclient.NewClient(
			httpclient.WithDigestAuthentication(username, password),
			httpclient.WithAcceptedStatusCodes(acceptedStatusCodes),
		)),
	)
}
