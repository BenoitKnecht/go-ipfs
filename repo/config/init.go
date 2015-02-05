package config

import "io"

func Init(out io.Writer, nBitsForKeypair int) (*Config, error) {
	ds, err := datastoreConfig()
	if err != nil {
		return nil, err
	}

	identity, err := InitIdentity(out, nBitsForKeypair)
	if err != nil {
		return nil, err
	}

	bootstrapPeers, err := DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}

	conf := &Config{

		// setup the node's default addresses.
		// Note: two swarm listen addrs, one tcp, one utp.
		Addresses: Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/4001",
				// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
			},
			API: "/ip4/127.0.0.1/tcp/5001",
		},

		Bootstrap: BootstrapPeerStrings(bootstrapPeers),
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
