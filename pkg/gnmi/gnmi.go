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
	defaultUsername     = "admin"
	defaultPassword     = "admin"
	defaultRetries      = 5
	defaultRetryTimeout = time.Duration(3) * time.Second
	defaultEncoding     = "JSON_IETF"
	defaultCallOptions  = []grpc.CallOption{
		grpc.WaitForReady(true),
	}
	defaultDialOptions = []grpc.DialOption{
		grpc.WithInsecure(),
	}
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

	Encoding        string
	DialOptions     []grpc.DialOption
	CallOptions     []grpc.CallOption
	Connection      *grpc.ClientConn
	GNMI            gpb.GNMIClient
	SubscribeClient gpb.GNMI_SubscribeClient

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
		c.Username = defaultUsername
	}
	if c.Password == "" {
		c.Password = defaultPassword
	}
	if c.Encoding == "" {
		c.Encoding = defaultEncoding
	}
	if len(c.DialOptions) == 0 {
		c.DialOptions = defaultDialOptions
	}
	if len(c.CallOptions) == 0 {
		c.CallOptions = defaultCallOptions
	}
	c.authorizedUser = newPassCred(c.Username, c.Password, c.EnableTLS)
	c.DialOptions = append(c.DialOptions, grpc.WithPerRPCCredentials(c.authorizedUser))
}

// Dial dials a gRPC connection and returns a reference to it
func (c *Client) Dial() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(c.Target, c.DialOptions...)
	if err != nil {
		log.Fatalf("Dialing to %q failed: %v", c.Target, err)
	}
	return conn, err
}

// InitGNMI initializes a gNMI client over a gRPC connection
func (c *Client) InitGNMI(conn *grpc.ClientConn) *gpb.GNMIClient {
	client := gpb.NewGNMIClient(conn)
	return &client
}

func CheckEncoding(encodingString string) int32 {
	encoding, ok := gpb.Encoding_value[encodingString]
	if !ok {
		var gnmiEncodingList []string
		for _, name := range gpb.Encoding_name {
			gnmiEncodingList = append(gnmiEncodingList, name)
		}
		log.Fatalf("Supported encodings: %s", strings.Join(gnmiEncodingList, ", "))
	}
	return encoding
}

// Get does a gNMI get on a path, returning the response
func Get(path string) (*gpb.GetResponse, error) {
	c := Client{}
	c.New()
	encoding := CheckEncoding(c.Encoding)
	c.Connection, _ = c.Dial()
	defer c.Connection.Close()

	ctx := context.Background()

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
	c.GNMI = *c.InitGNMI(c.Connection)
	resp, err := c.GNMI.Get(ctx, getRequest, c.CallOptions...)
	if err != nil {
		log.Infof("Get failed: %v", err)
	}
	//log.Printf("Got response: %s", proto.MarshalTextString(Resp))

	log.Infof("Get result: %v", resp)
	return resp, err
}

// Set does a gNMI Set, which can be used for both updates and deletes
func Set(ctx context.Context, path string, c *Client, setRequest *gpb.SetRequest) (*gpb.SetResponse, error) {
	var err error
	var output *gpb.SetResponse

	// RETRYLOOP:
	for j := 0; j < defaultRetries; j++ {
		output, err = c.GNMI.Set(ctx, setRequest, c.CallOptions...)
		if err != nil {
			log.Warningf("Set request failed for path: %s (attempt: %d): %v", path, j, err)
			time.Sleep(defaultRetryTimeout)
			continue
		}
		// break RETRYLOOP
		break
	}

	if err != nil {
		log.Errorf("Set failed: %v", err)
	} else {
		//log.Printf("Got response: %s", proto.MarshalTextString(Resp))
		log.Infof("gNMI SET run successfully, path: %s", path)
	}
	return output, err
}

// Subscribe subscribes to a path
func Subscribe(path string) (*Client, error) {
	var err error
	c := Client{}
	c.New()
	c.Connection, _ = c.Dial()

	c.GNMI = *c.InitGNMI(c.Connection)
	ctx := context.Background()

	encoding := CheckEncoding(defaultEncoding)

	xPath, err := xpath.ToGNMIPath(path)
	if err != nil {
		log.Fatalf("error in parsing xpath %q to gnmi path", "/interface[name=*]")
	}

	subsList := []*gpb.Subscription{
		&gpb.Subscription{
			Path: xPath,
		},
	}

	req := &gpb.SubscribeRequest_Subscribe{
		Subscribe: &gpb.SubscriptionList{
			Encoding:     gpb.Encoding(encoding),
			Subscription: subsList,
		},
	}

	gnmiSubscribeRequest := &gpb.SubscribeRequest{
		Request: req,
	}

	c.SubscribeClient, err = c.GNMI.Subscribe(ctx)
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}

	if err := c.SubscribeClient.Send(gnmiSubscribeRequest); err != nil {
		log.Fatalf("Failed to send a subscribe request for interface: %v", err)
	}

	return &c, err
}

// Set does a gNMI Set, given a path and value in gpb TypedValue format
func Update(ctx context.Context, path string, val *gpb.TypedValue) (*gpb.SetResponse, error) {
	var err error
	var output *gpb.SetResponse
	c := Client{}
	c.New()
	t := SetTransaction{}
	// Build options
	t.Path, err = xpath.ToGNMIPath(path)
	t.Value = val
	if err != nil {
		log.Fatalf("error in parsing xpath %q to gnmi path", path)
	}
	c.Connection, _ = c.Dial()
	defer c.Connection.Close()
	c.GNMI = *c.InitGNMI(c.Connection)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var updateList []*gpb.Update
	updateList = append(updateList, &gpb.Update{Path: t.Path, Val: t.Value})
	setRequest := &gpb.SetRequest{
		Update: updateList,
	}
	log.Infof("Running gNMI SET, path: %s", path)
	output, err = Set(ctx, path, &c, setRequest)
	return output, err
}

// Delete does a gNMI Set delete on the provided path
func Delete(ctx context.Context, path string) (*gpb.SetResponse, error) {
	var err error
	var output *gpb.SetResponse
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

	c.Connection, _ = c.Dial()
	defer c.Connection.Close()
	c.GNMI = *c.InitGNMI(c.Connection)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	setRequest := &gpb.SetRequest{
		Delete: pathList,
	}
	log.Infof("Running gNMI SET DELETE, path: %s...", path)
	output, err = Set(ctx, path, &c, setRequest)
	return output, err
}
