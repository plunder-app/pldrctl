# pldrctl

Plunderctl - The control utility for a plunder server

The `pldrctl` utility is designed to interact with the **alpha** API exposed from a [plunder](github.com/plunder-app/plunder) server, and provide a CLI/UX that allows an end user to quickly and simply manage the deployment of servers and Operating Systems.


## Building

If you wish to build the code yourself then this can be done simply by running:

```
go get github.com/plunder-app/pldrctl
```

Alternatively clone the repository and either `go build` or `make install`, note that using the makefile will ensure that the current git commit and version number are returned by `pldrctl version`.

**NOTE** Releases will be created soon.

## Usage!

Due to a severe lack of originality the usage of `pldrctl` mirrors a LOT of the typical interactions that an end user may have with `kubectl`, this is both lazyness on my part and homage to the good design of the latter tool.

### Set up Plunder

**Quick Overview** - More steps will need to be followed for a completely working configuration

Once Plunder is installed the following things need/can be done:

1. Create a config `plunder config server -o yaml > config.yaml`
2. Edit the configuration so that `DHCP/TFTP/HTTP` are enabled (or use the `-d <NIC>` for a different interface on the above command)
3. Edit the configuration so that a `preseed` config exists pointing to a kernel/initrd etc..
4. Generate an example deployment `plunder config deployment -o yaml > deployment.yaml`
5. Modify the global configuration to use sensible settings for your environment
6. Generate the required apiServer configs `plunder config apiserver server/client`
7. Start `plunder` with a `reboot` configuration for unknown servers `plunder server --config ./config.yaml --deployment ./deployment.json --defaultBoot reboot`

(Step 6 is important as the client cert is needed for `pldrctl`)

### Set up pldrctl

The client configuration will either need to be in your working directory or pointed at with the `-p <path to config>` flag.

The below example will copy the configuraiton from my deployment server `deploy01` and test connectivity with a `get config`

```
$ scp deploy01:plunderclient.yaml .
plunderclient.yaml                            100% 2829     2.6MB/s   00:00    
$ pldrctl get config
Adapter:      ens192
Enable DHCP:  true
              DHCP Start Address:    192.168.1.2
              DHCP Server Address:   192.168.1.1
              DHCP Gateway Address:  192.168.1.1
              DHCP DNS Address:      192.168.1.1
              DHCP Lease Pool Size:  20
Enable TFTP:  true
              TFTP Server Address:  192.168.1.1
Enable HTTP:  true
              DHCP Server Address:  192.168.1.1
              PXE File Name:        undionly.kpxe
$
```

### Using pldrctl to view the infrastructure

If the Plunder has been running for a while it will have started to see some MAC addresses that it's not aware of, these addresses aren't given an IP address lease and therefore are defined as *unleased*. These nodes can be viewed as shown below:

```
$ pldrctl get unleased
Mac Address        Time Seen                 Time since  Hardware Vendor
00:50:56:9b:6e:fc  Sat Aug 10 00:11:32 2019  39h54m11s   VMware, Inc.
00:50:56:9b:3a:3d  Sat Aug 10 00:13:47 2019  39h51m56s   VMware, Inc.
00:50:56:9b:d1:e4  Sat Aug 10 00:20:26 2019  39h45m17s   VMware, Inc.
00:50:56:9b:de:2d  Sat Aug 10 10:27:59 2019  29h37m44s   VMware, Inc.
00:0e:ab:11:22:33  Sun Aug 11 16:05:38 2019  5s          Cray Inc
00:0f:4b:fe:41:54  Sun Aug 11 16:05:34 2019  9s          Oracle Corporation
ec:9a:74:9b:6e:fc  Sun Aug 11 16:05:43 2019  0s          Hewlett Packard
```

**Note** I don't actually have a Cray, i fudged a mac address on a VM :-( 

### "Applying" with pldrctl

We can look at an existing deployment, and export it's configuration with `pldrctl` very easily..

```
$ pldrctl get deployments
Mac Address        Deploymemt  Hostname  IP Address
00:50:56:a3:64:a2  preseed     etcd01    192.168.1.3
00:50:56:a3:4c:da  preseed     etcd02    192.168.1.4
00:50:56:a3:e1:da  preseed     etcd03    192.168.1.5
00:50:56:a3:7e:ee  preseed     master01  192.168.1.6
00:50:56:a3:91:ee  vsphere     master02  192.168.1.7
$ pldrctl describe deployment 00:50:56:a3:64:a2 -o yaml
{...}
$ pldrctl describe deployment 00:50:56:a3:64:a2 -o yaml > newConfigForCray.yaml
```

We can modify the newConfigForCray.yaml so that it only has configuration we're interested in, plunder will automatically inherit `config:` and look up a `bootConfigName:`.

Example config, with `mac`/`address` and `hostname` modified:

```
definition: deployment
resource:
  bootConfigName: preseed
  config:
    address: 192.168.1.100
    hostname: newNode01
  mac: 00:0e:ab:11:22:33
```

With out new configuration created we can apply it!

```
# Either directly
$ pldrctl apply -f ./test.yaml
# Or through a pipe
$ cat test.yaml | pldrctl appy -
```

We will see our new node now listed in the deployments:

```
$ pldrctl get deployments | grep newNode
00:0e:ab:11:22:33  preseed     newNode01  192.168.1.100
```

### "Creating" with pldrctl

We can also create deployments from `pldrctl` without creating any configuration files. The create option will allow creating a simple deployment with a few flags:

```
$ pldrctl create -t deployment \
     -c preseed                \
     -m 00:0f:4b:fe:41:54      \
     -a 192.168.1.110          \
     -n newNode02
```

We will see our new node now listed in the deployments:

```
$ pldrctl get deployments | grep newNode
00:0e:ab:11:22:33  preseed     newNode01  192.168.1.100
00:0f:4b:fe:41:54  preseed     newNode02  192.168.1.110
```

### Further awesome-ness with pldrctl

We can dig deeper into whats happening a deployment by `describing` the boot process:

```
$ pldrctl describe boot 00:0e:ab:11:22:33
Config:     
            Deployment Type:  preseed
            Kernel:           ubuntu/install/netboot/ubuntu-installer/amd64/linux
            Initrd:           ubuntu/install/netboot/ubuntu-installer/amd64/initrd.gz
            cmdline:          
            Adapter:          
            Server Name:      newNode01
            IP Address:       192.168.1.100
Phase One:  
            Action:  DHCP Request -> TFTP Boot -> iPXE boot with config file
            Config:  http://deploy01/00-0e-ab-11-22-33.ipxe
Phase Two:  
            Action:  OS Bootstraps with config
            Config:  http://deploy01/00-0e-ab-11-22-33.cfg
```

The URLs should be accessible in order to debug the configurations that `plunder` has created.

To view the various `bootConfigName` options we can get them:

```
$ pldrctl get boot
Config Name  Kernel Path                                          Initrd Path                                              Command Line
vsphere      vsphere/mboot.c32                                    vsphere/boot.cfg                                         
preseed      ubuntu/install/netboot/ubuntu-installer/amd64/linux  ubuntu/install/netboot/ubuntu-installer/amd64/initrd.gz  
```

## NEXT STEPS

- Documentation
- Tidy a lot of bad bad bad code
- Ability to upload SSH keys to remote plunder
- Pass through commands
- Get availability (ping) of hosts
- Get uptime of hosts