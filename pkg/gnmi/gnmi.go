package gnmi

import (
	"context"
	"strings"
	"sync"
	"time"

	log "k8s.io/klog"

	// "github.com/google/gnxi/utils"

	"github.com/google/gnxi/utils/xpath"

	// "github.com/openconfig/ygot/ygot"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	defaultUsername = "admin"
	defaultPassword = "admin"
)

// passCred is an username/password implementation of credentials.Credentials.
type passCred struct {
	username string
	password string
	secure   bool
}

// GetRequestMetadata returns the current request metadata, including
// username and password in this case.
// This implements the required interface fuction of credentials.Credentials.
func (pc *passCred) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"username": pc.username,
		"password": pc.password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
// This implements the required interface fuction of credentials.Credentials.
func (pc *passCred) RequireTransportSecurity() bool {
	return pc.secure
}

// newPassCred returns a newly created passCred as credentials.Credentials.
func newPassCred(username, password string, secure bool) credentials.PerRPCCredentials {
	return &passCred{
		username: username,
		password: password,
		secure:   secure,
	}
}

// Duration just wraps time.Duration
type Duration struct {
	Duration time.Duration
}

// TelemetryGNMI plugin instance
type TelemetryGNMI struct {
	Addresses     []string
	Subscriptions []Subscription

	// Optional subscription configuration
	Encoding    string
	Origin      string
	Prefix      string
	Target      string
	UpdatesOnly bool

	Username string
	Password string

	// Redial
	Redial Duration

	// GRPC TLS settings
	EnableTLS bool
	// internaltls.ClientConfig

	// Internal state
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Subscription for a GNMI client
type Subscription struct {
	Name   string
	Origin string
	Path   string

	// Subscription mode and interval
	SubscriptionMode string
	SampleInterval   Duration

	// Duplicate suppression
	SuppressRedundant bool
	HeartbeatInterval Duration
}

// Client is a holder for a gNMI client
type Client struct {
	Target    string
	Username  string
	Password  string
	EnableTLS bool
	Redial    Duration

	authorizedUser credentials.PerRPCCredentials
}

// SetTransaction represents a gNMI set transaction
type SetTransaction struct {
	Path  *gpb.Path
	Value *gpb.TypedValue
}

// New is a constructor function for a gNMI client
func (c *Client) New() {
	if c.Target == "" {
		c.Target = "unix:///opt/srlinux/var/run/sr_gnmi_server"
	}
	if c.Username == "" {
		c.Username = "admin"
	}
	if c.Password == "" {
		c.Password = "admin"
	}
	c.authorizedUser = newPassCred(c.Username, c.Password, c.EnableTLS)
}

// ClientInit initializes a client
//func ClientInit() {
//
//}

// Get does a gNMI get on a path
func Get(path string) {
	c := Client{}
	c.New()
	// Build options
	dialOpts := []grpc.DialOption{}
	dialOpts = append(dialOpts, grpc.WithInsecure())
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(c.authorizedUser))
	conn, err := grpc.Dial(c.Target, dialOpts...)
	if err != nil {
		log.Fatalf("Dialing to %q failed: %v", c.Target, err)
	}
	defer conn.Close()

	cli := gpb.NewGNMIClient(conn)

	ctx := context.Background()

	encoding, ok := gpb.Encoding_value["JSON_IETF"]
	if !ok {
		var gnmiEncodingList []string
		for _, name := range gpb.Encoding_name {
			gnmiEncodingList = append(gnmiEncodingList, name)
		}
		log.Fatalf("Supported encodings: %s", strings.Join(gnmiEncodingList, ", "))
	}

	var gNMIPaths []*gpb.Path
	gNMIPath, err := xpath.ToGNMIPath(path)
	if err != nil {
		log.Fatalf("error in parsing xpath %q to gnmi path", path)
	}
	gNMIPaths = append(gNMIPaths, gNMIPath)
	getRequest := &gpb.GetRequest{
		Encoding: gpb.Encoding(encoding),
		Path:     gNMIPaths,
	}
	callOpts := []grpc.CallOption{}
	callOpts = append(callOpts, grpc.WaitForReady(true))
	resp, err := cli.Get(ctx, getRequest, callOpts...)
	//log.Printf("Got response: %s", proto.MarshalTextString(Resp))

	log.Infof("Get result: %v", resp)
	if err != nil {
		log.Infof("Set failed: %v", err)
	}

}

// Set does a gNMI Set, given a path and value in gpb TypedValue format
func Set(path string, val *gpb.TypedValue) (*gpb.SetResponse, error) {
	var err error
	c := Client{}
	c.New()
	t := SetTransaction{}
	// Build options
	t.Path, err = xpath.ToGNMIPath(path)
	t.Value = val
	if err != nil {
		log.Fatalf("error in parsing xpath %q to gnmi path", path)
	}
	dialOpts := []grpc.DialOption{}
	dialOpts = append(dialOpts, grpc.WithInsecure())
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(c.authorizedUser))
	conn, err := grpc.Dial(c.Target, dialOpts...)
	if err != nil {
		log.Fatalf("Dialing to %q failed: %v", c.Target, err)
	}
	defer conn.Close()

	cli := gpb.NewGNMIClient(conn)

	ctx := context.Background()

	var updateList []*gpb.Update
	updateList = append(updateList, &gpb.Update{Path: t.Path, Val: t.Value})
	setRequest := &gpb.SetRequest{
		Update: updateList,
	}
	callOpts := []grpc.CallOption{}
	callOpts = append(callOpts, grpc.WaitForReady(true))
	log.Infof("Running gNMI SET...")
	output, err := cli.Set(ctx, setRequest, callOpts...)
	//log.Printf("Got response: %s", proto.MarshalTextString(Resp))
	log.Infof("gNMI SET run...")

	if err != nil {
		log.Infof("Set failed: %v", err)
	}
	return output, err
}

// Delete does a gNMI Set delete on the provided path
func Delete(path string) (*gpb.SetResponse, error) {
	var err error
	var pathList []*gpb.Path
	c := Client{}
	c.New()
	t := SetTransaction{}

	log.Infof("Delete called with path: %s", path)
	// Build options
	t.Path, err = xpath.ToGNMIPath(path)
	if err != nil {
		log.Fatalf("error in parsing xpath %q to gnmi path", path)
	}
	pathList = append(pathList, t.Path)

	log.Info("Creating dial options...")
	dialOpts := []grpc.DialOption{}
	dialOpts = append(dialOpts, grpc.WithInsecure())
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(c.authorizedUser))
	log.Info("Dialing...")
	conn, err := grpc.Dial(c.Target, dialOpts...)
	if err != nil {
		log.Fatalf("Dialing to %q failed: %v", c.Target, err)
	}
	defer conn.Close()

	cli := gpb.NewGNMIClient(conn)

	ctx := context.Background()

	setRequest := &gpb.SetRequest{
		Delete: pathList,
	}
	log.Info("Building call options...")
	callOpts := []grpc.CallOption{}
	callOpts = append(callOpts, grpc.WaitForReady(true))
	log.Infof("Running gNMI SET DELETE...")
	output, err := cli.Set(ctx, setRequest, callOpts...)
	//log.Printf("Got response: %s", proto.MarshalTextString(Resp))
	log.Infof("gNMI SET DELETE run...")

	if err != nil {
		log.Infof("Set failed: %v", err)
	}
	return output, err
}
