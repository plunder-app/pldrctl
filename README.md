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

### Remote Execution with pldrctl

`pldrctl` can automate the deployment of applications or configuration of deployed servers in two different methods:

#### Scripted automation with parlay

The below file is the contents of kube_reset.yaml and is configured to remove pacakges and tidy up etc..

```
definition: parlay
resource:
  deployments:
  - actions:
    - command: kubeadm reset -f
      commandSudo: root
      name: Reset Kubernetes
      type: command
      ignoreFail: true
    - command: dpkg -r kubeadm kubelet kubectl cri-tools kubernetes-cni
      commandSudo: root
      name: Remove packages
      type: command
    - command: rm -rf /opt/cni/bin
      commandSudo: root
      name: Remove any remaining cni directories 
      type: command
    hosts:
    - 192.168.1.129
    name: Reset any Kubernetes configuration (and remove packages)
```

We can apply this to plunder with `pldrctl` in the normal manner:

`pldrctl apply -f ./kube_reset.yaml` 

Plunder will examine the resource definition and execute the below actions on the hosts listed in the host array and follow the actions in order. 

#### Viewing and Managing logs of the remote commands

Plunder will keep all of the remote logs in-memory to be retrieved at a later date, we can view and delete these logs through the `pldrctl` utility.

Using the `pldrctl get logs <address>` command we can pull the logs from the plunder server, below we can see the output from the above command. 

```
pldrctl get logs  192.168.1.129 
Logs:
Started: Tue Aug 20 08:47:42 2019	Task Name: Reset Kubernetes
Output:
[reset] Reading configuration from the cluster...
[reset] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -oyaml'
[preflight] Running pre-flight checks
[reset] Removing info for node "master" from the ConfigMap "kubeadm-config" in the "kube-system" Namespace
[reset] Stopping the kubelet service
W0820 15:45:36.585713    2141 reset.go:158] [reset] failed to remove etcd member: error syncing endpoints with etc: etcdclient: no available endpoints
.Please manually remove this etcd member using etcdctl
[reset] unmounting mounted directories in "/var/lib/kubelet"
[reset] Deleting contents of stateful directories: [/var/lib/etcd /var/lib/kubelet /etc/cni/net.d /var/lib/dockershim /var/run/kubernetes]
[reset] Deleting contents of config directories: [/etc/kubernetes/manifests /etc/kubernetes/pki]
[reset] Deleting files: [/etc/kubernetes/admin.conf /etc/kubernetes/kubelet.conf /etc/kubernetes/bootstrap-kubelet.conf /etc/kubernetes/controller-manager.conf /etc/kubernetes/scheduler.conf]

The reset process does not reset or clean up iptables rules or IPVS tables.
If you wish to reset iptables, you must do so manually.
For example:
iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -X

If your cluster was setup to utilize IPVS, run ipvsadm --clear (or similar)
to reset your system's IPVS tables.
Started: Tue Aug 20 08:47:43 2019	Task Name: Remove packages
Output:
(Reading database ... 54421 files and directories currently installed.)
Removing kubeadm (1.14.0-00) ...
Removing kubelet (1.14.0-00) ...
Removing kubectl (1.14.0-00) ...
Removing cri-tools (1.12.0-00) ...
Removing kubernetes-cni (0.7.5-00) ...
dpkg: warning: while removing kubernetes-cni, directory '/opt/cni/bin' not empty so not removed
Started: Tue Aug 20 08:47:43 2019	Task Name: Remove any remaining cni directories
Output:
[No Output]
Task Status: Completed
```

We can also watch running execution with the following command:

`kubectl get logs -w <timeout> <address>` 

This command will keep refreshing the logs until the **Task Status:** changes from **Running**. 

#### Single automation and deleting old logs

On the same host we can delete all of the logs from `plunder` with the following command:

`pldrctl delete -t <resource_type> <address>`

Then we can run a single remote command with the following:

`pldrctl exec -a <address> -c <command>`

Below in an example of the same server from above:

```
pldrctl delete -t logs 192.168.1.129

pldrctl exec -a 192.168.1.129 -c "uptime"

pldrctl get logs  192.168.1.129 

Logs:
Started: Tue Aug 20 08:55:48 2019	Task Name: pldrctl command
Output:
 15:53:46 up 16:15,  1 user,  load average: 0.03, 0.04, 0.07
Task Status: Completed
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