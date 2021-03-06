package config

import (
	"encoding/base64"
	"fmt"
	"io"

	ci "github.com/jbenet/go-ipfs/p2p/crypto"
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	errors "github.com/jbenet/go-ipfs/util/debugerror"
)

func Init(out io.Writer, nBitsForKeypair int, addresses Addresses) (*Config, error) {
	ds, err := datastoreConfig()
	if err != nil {
		return nil, err
	}

	identity, err := identityConfig(out, nBitsForKeypair)
	if err != nil {
		return nil, err
	}

	bootstrapPeers, err := DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}

	snr, err := initSNRConfig()
	if err != nil {
		return nil, err
	}

	if addresses.Swarm == nil {
		addresses.Swarm = []string{
			"/ip4/0.0.0.0/tcp/4001",
			// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
		}
	}
	if addresses.API == "" {
		addresses.API = "/ip4/127.0.0.1/tcp/5001"
	}
	if addresses.Gateway == "" {
		addresses.Gateway = "/ip4/127.0.0.1/tcp/8080"
	}

	conf := &Config{

		// setup the node's default addresses.
		// Note: two swarm listen addrs, one tcp, one utp.
		Addresses: addresses,

		Bootstrap: BootstrapPeerStrings(bootstrapPeers),
		SupernodeRouting:       *snr,
		Datastore: *ds,
		Identity:  identity,
		Log: Log{
			MaxSizeMB:  250,
			MaxBackups: 1,
		},

		// setup the node mount points.
		Mounts: Mounts{
			IPFS: "/ipfs",
			IPNS: "/ipns",
		},

		// tracking ipfs version used to generate the init folder and adding
		// update checker default setting.
		Version: VersionDefaultValue(),

		Gateway: Gateway{
			RootRedirect: "",
			Writable:     false,
		},
	}

	return conf, nil
}

func datastoreConfig() (*Datastore, error) {
	dspath, err := DataStorePath("")
	if err != nil {
		return nil, err
	}
	return &Datastore{
		Path: dspath,
		Type: "leveldb",
	}, nil
}

// identityConfig initializes a new identity.
func identityConfig(out io.Writer, nbits int) (Identity, error) {
	// TODO guard higher up
	ident := Identity{}
	if nbits < 1024 {
		return ident, errors.New("Bitsize less than 1024 is considered unsafe.")
	}

	fmt.Fprintf(out, "generating %v-bit RSA keypair...", nbits)
	sk, pk, err := ci.GenerateKeyPair(ci.RSA, nbits)
	if err != nil {
		return ident, err
	}
	fmt.Fprintf(out, "done\n")

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	fmt.Fprintf(out, "peer identity: %s\n", ident.PeerID)
	return ident, nil
}
