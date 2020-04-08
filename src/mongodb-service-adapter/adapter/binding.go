package adapter

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	mgo "gopkg.in/mgo.v2"
)

type Binder struct {
	Logger *log.Logger
}

func (b *Binder) logf(msg string, v ...interface{}) {
	if b.Logger != nil {
		b.Logger.Printf(msg, v...)
	}
}

var (
	adminDB        = "admin"
	defaultDB      = "default"
	serverKeyPath  = "/var/vcap/jobs/mongodb_service_adapter/config/server.key"
	serverCertPath = "/var/vcap/jobs/mongodb_service_adapter/config/server.crt"
)

func (b Binder) CreateBinding(params serviceadapter.CreateBindingParams) (serviceadapter.Binding, error) {
	// create an admin level user
	username := mkUsername(params.BindingID)
	password := GenerateString(32, true)

	properties := params.Manifest.Properties["mongo_ops"].(map[interface{}]interface{})
	adminPassword := properties["admin_password"].(string)
	URL := properties["url"].(string)
	adminUsername := properties["username"].(string)
	adminAPIKey := properties["admin_api_key"].(string)
	ssl := properties["require_ssl"].(bool)
	groupID := properties["group_id"].(string)

	b.logf("properties: %v", properties)

	servers := make([]string, len(params.DeploymentTopology["mongod_node"]))
	for i, node := range params.DeploymentTopology["mongod_node"] {
		servers[i] = fmt.Sprintf("%s:28000", node)
	}

	plan := properties["plan_id"].(string)
	if plan == PlanShardedCluster {
		routers := properties["routers"].(int)
		configServers := properties["config_servers"].(int)
		replicas := properties["replicas"].(int)

		cluster, err := NodesToCluster(servers, routers, configServers, replicas)
		if err != nil {
			return serviceadapter.Binding{}, errors.Wrap(err, "cannot convert node list to cluster")
		}
		servers = cluster.Routers
	}

	if ssl {
		omClient := NewPCGCClient(URL, adminUsername, adminAPIKey)
		hosts, err := omClient.GetHosts(groupID)
		if err != nil {
			return serviceadapter.Binding{}, errors.Wrap(err, "cannot get hosts")
		}

		servers = ToEndpointList(hosts)
	}

	session, err := GetWithCredentials(servers, adminPassword, ssl)
	if err != nil {
		return serviceadapter.Binding{}, errors.Wrap(err, "cannot get session")
	}
	defer session.Close()

	// add user to admin database with admin privileges
	user := &mgo.User{
		Username: username,
		Password: password,
		Roles: []mgo.Role{
			mgo.RoleUserAdmin,
			mgo.RoleDBAdmin,
			mgo.RoleReadWrite,
		},
		OtherDBRoles: map[string][]mgo.Role{
			defaultDB: {
				mgo.RoleUserAdmin,
				mgo.RoleDBAdmin,
				mgo.RoleReadWrite,
			},
		},
	}

	if err = session.DB(adminDB).UpsertUser(user); err != nil {
		return serviceadapter.Binding{}, errors.Wrap(err, "cannot upsert user")
	}

	connectionOpts := url.Values{}
	connectionOpts.Add("authSource", adminDB)
	if ssl {
		connectionOpts.Add("ssl", "true")
	}
	if plan == PlanReplicaSet {
		connectionOpts.Add("replicaSet", "pcf_repl")
	}

	url := url.URL{
		Scheme:   "mongodb",
		Host:     strings.Join(servers, ","),
		User:     url.UserPassword(username, password),
		Path:     defaultDB,
		RawQuery: connectionOpts.Encode(),
	}

	return serviceadapter.Binding{
		Credentials: map[string]interface{}{
			"username": username,
			"password": password,
			"database": defaultDB,
			"servers":  servers,
			"ssl":      ssl,
			"uri":      url.String(),
		},
	}, nil
}

func (b Binder) DeleteBinding(params serviceadapter.DeleteBindingParams) error {
	// create an admin level user
	username := mkUsername(params.BindingID)
	properties := params.Manifest.Properties["mongo_ops"].(map[interface{}]interface{})
	adminPassword := properties["admin_password"].(string)
	ssl := properties["require_ssl"].(bool)
	URL := properties["url"].(string)
	adminUsername := properties["username"].(string)
	adminAPIKey := properties["admin_api_key"].(string)
	groupID := properties["group_id"].(string)

	servers := make([]string, len(params.DeploymentTopology["mongod_node"]))
	for i, node := range params.DeploymentTopology["mongod_node"] {
		servers[i] = fmt.Sprintf("%s:28000", node)
	}

	if ssl {
		var err error
		omClient := NewPCGCClient(URL, adminUsername, adminAPIKey)
		hosts, err := omClient.GetHosts(groupID)
		if err != nil {
			return errors.Wrap(err, "cannot get hosts")
		}

		servers = ToEndpointList(hosts)
	}

	session, err := GetWithCredentials(servers, adminPassword, ssl)
	if err != nil {
		return errors.Wrap(err, "cannot get session")
	}
	defer session.Close()

	err = session.DB(adminDB).RemoveUser(username)
	if err == mgo.ErrNotFound {
		b.logf("ERROR: could not find user %q in DB %q for removal! It could be that something or someone has already removed it. Deleting the binding anyway...", username, adminDB)
		return nil
	}

	return errors.Wrapf(err, "cannot remove user %q from DB %q", username, adminDB)
}

func GetWithCredentials(addrs []string, adminPassword string, ssl bool) (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Username:  "admin",
		Password:  adminPassword,
		Mechanism: "SCRAM-SHA-1",
		Database:  adminDB,
		FailFast:  true,
	}
	if ssl {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot load keypair from %q, %q", serverCertPath, serverKeyPath)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}

		dialInfo.DialServer = func(addrs *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addrs.String(), tlsConfig)
			return conn, errors.Wrapf(err, "cannot tls.Dial %q", addrs.String())
		}
	}

	result, err := mgo.DialWithInfo(dialInfo)
	return result, errors.Wrapf(err, "cannot dial mongo cluster %q (ssl: %v)", addrs, ssl)
}

func mkUsername(bindingID string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(bindingID))
	return fmt.Sprintf("pcf_%x", md5.Sum([]byte(b64)))
}
