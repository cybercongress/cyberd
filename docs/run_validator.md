
# Join Cyberd testnet as a Validator

**Note** The current active testnet is `euler-6` (substitute <testnet_chain_id> with that value, do not forget to remove the `<` and the `>` symbols).

## Prepare your server

First of all, you should set up a server.
Your node should be running (online) constantly. This means that you will need a reliable server to keep it running.
You may also consider using any cloud service with a dedicated GPU, like Hetzner (or use a local machine). Whatever you'll choose, for better stability and consistency we recommend to use a dedicated server for each separate validator node.

Cyberd is based on Cosmos-SDK and is written in Go.
It should work on any platform which can compile and run programs in Go.
However, we strongly recommend running the validator node on a Linux-based server.

The rank calculations are done via GPU computations.
They are easy to parallelize. This is why we recommended using a GPU.

Recommended requirements:

```js
CPU: 6 cores
RAM: 32 GB
SSD: 256 GB
Connection: 100Mb, Fiber, Stable and low-latency connection
GPU: Nvidia GeForce(or Tesla/Titan/Quadro) with CUDA-cores; at least 6gb of memory*
Software: Docker, Ubuntu 16.04/18.04 LTS
```

*Cyberd runs well on consumer-grade cards like Geforce GTX 1070, but expect load growth. We advise to use Error Correction compatible cards from Tesla or Quadro families.*

But, of course, the hardware is your own choice and technically it might be possible to run the chain on "even - 1 CUDA core GPU", but, you should be aware of stability and a decline in calculation speed.

## Validator setup

### Third-party software

To access the GPU, cyberd uses Nvidia drivers version **410+** and the [Nvidia CUDA toolkit](https://developer.nvidia.com/cuda-downloads) should be installed on the hosting system. 

You may skip any sections of the guide if you already have any of the necessary software configured. 

As long as the current implementation of `cyber` is written in [Go](https://golang.org/), you will need to install Go.

### Installing Go

For `euler-6` Cyberd requires at least Go version 1.13+. Install it according to the official [guide](https://golang.org/doc/install):

1. Download the archive:

```bash
wget https://dl.google.com/go/go1.13.9.linux-amd64.tar.gz
```

2. Extract it into `/usr/local`, creating a Go tree in `/usr/local/go`:

```bash
tar -C /usr/local -xzf go1.13.9.linux-amd64.tar.gz
```

3. Add `/usr/local/go/bin` to the PATH environment variable. You can do this by adding this line to your `/etc/profile` (for installation on the whole system) or `$HOME/.profile`:

```bash
export PATH=$PATH:/usr/local/go/bin
```

4. Do `source` for the file with your `$PATH` variable or just log-out/log-in:

```bash
source /etc/profile
```

or

```bash
source $HOME/.profile
```

5. To check your installation run

```bash
go version
```

This will let you know if everything was installed correctly. As an output, you should see the following (version number may vary, of course):

```bash
go version go1.13.9 linux/amd64
```

#### Installing Nvidia drivers 

1. To proceed, first, add the `ppa:graphics-drivers/ppa` repository into your system (you might see some warnings - press `enter`):

```bash
sudo add-apt-repository ppa:graphics-drivers/ppa
```

```bash
sudo apt update
```

2. Install Ubuntu-drivers:

```bash
sudo apt install -y ubuntu-drivers-common
```

3. Next, identify your graphic card model and the recommended drivers:

```bash
ubuntu-drivers devices
```

You should see something similar to this:

```bash
== /sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0 ==
modalias : pci:v000010DEd00001BA1sv00001462sd000011E4bc03sc00i00
vendor   : NVIDIA Corporation
model    : GP104M [GeForce GTX 1070 Mobile]
driver   : nvidia-driver-418 - third-party free
driver   : nvidia-driver-410 - third-party free
driver   : nvidia-driver-430 - third-party free
driver   : nvidia-driver-440 - third-party free recommended
driver   : xserver-xorg-video-nouveau - distro free builtin
```

4. We need the **410+** drivers release. As you can see that v440 is recommended. The command below will install the recommended version of the drivers:

```bash
sudo ubuntu-drivers autoinstall
```

The driver installation takes approximately 10 minutes.

```bash
DKMS: install completed.
Setting up libxdamage1:i386 (1:1.1.4-3) ...
Setting up libxext6:i386 (2:1.3.3-1) ...
Setting up libxfixes3:i386 (1:5.0.3-1) ...
Setting up libnvidia-decode-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up build-essential (12.4ubuntu1) ...
Setting up libnvidia-gl-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-encode-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up nvidia-driver-415 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libxxf86vm1:i386 (1:1.1.4-1) ...
Setting up libglx-mesa0:i386 (18.0.5-0ubuntu0~18.04.1) ...
Setting up libglx0:i386 (1.0.0-2ubuntu2.2) ...
Setting up libgl1:i386 (1.0.0-2ubuntu2.2) ...
Setting up libnvidia-ifr1-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-fbc1-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Processing triggers for libc-bin (2.27-3ubuntu1) ...
Processing triggers for initramfs-tools (0.130ubuntu3.1) ...
update-initramfs: Generating /boot/initrd.img-4.15.0-45-generic
```

5. Reboot the system for the changes to take effect.

6. Check the installed drivers:

```bash
nvidia-smi
```

You should see this:
(Some version/driver numbers might differ. You might also have some processes already running)

```bash
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 440.26       Driver Version: 430.14       CUDA Version: 10.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1070    Off  | 00000000:01:00.0 Off |                  N/A |
| 26%   36C    P5    26W / 180W |      0MiB /  8119MiB |      2%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

### Install CUDA toolkit

Simply run:

```bash
apt install nvidia-cuda-toolkit
```

Any version above 9.1 is OK. To check the version run: `

```bash
nvcc --version
```

The possible output will look like this:

```bash
nvcc: NVIDIA (R) Cuda compiler driver
Copyright (c) 2005-2017 NVIDIA Corporation
Built on Fri_Nov__3_21:07:56_CDT_2017
Cuda compilation tools, release 9.1, V9.1.85
```

### Launching Cyberd fullnode

1. Add environment variables:

```bash
export DAEMON_HOME=$HOME/.cyberd
export DAEMON_NAME=cyberd
```

To make those variables persistent add them to the end of the **`$HOME/.profile`** and log-out/log-in, or do:

```bash
source ~/.profile
```

2. Make a directories tree for storing your daemon:

```bash
mkdir $HOME/.cyberd
mkdir -p $DAEMON_HOME/upgrade_manager
mkdir -p $DAEMON_HOME/upgrade_manager/genesis
mkdir -p $DAEMON_HOME/upgrade_manager/genesis/bin
```

3. Download cosmosd repo and build cosmosd:

```bash
git clone https://github.com/regen-network/cosmosd
cd cosmosd
go build
mv cosmosd $DAEMON_HOME/
chmod +x $DAEMON_HOME/cosmosd
```

4. Clone the go-cyber repo, checkout to the necessary version (`master` by default):

```bash
cd ~
git clone https://github.com/cybercongress/go-cyber
```

5. Build cyber~Rank CUDA kernel:

```bash
cd ~/go-cyber/x/rank/cuda/
make
```

5. Build cyber daemon (as a result you should see `cyberd` and `cyberdcli` files inside of the `go-cyber/build/` folder):

```bash
cd ~/go-cyber
make build
```

6. Copy the binaries to an apropriate locations and add execute permitions to them:

```bash
cp build/cyberd $DAEMON_HOME/upgrade_manager/genesis/bin
cp build/cyberdcli /usr/local/bin/
cp build/cyberd /usr/local/bin/
chmod +x $DAEMON_HOME/upgrade_manager/genesis/bin/cyberd
```

7. Initialize cyber daemon:

```bash
cd $DAEMON_HOME/upgrade_manager/genesis/bin
./cyberd init <your_node_moniker> --home $DAEMON_HOME
```

8. Your folder with cyberd must look like this after initialization:

```bash
root@node:~/.cyberd# tree
.
├── config
│   ├── app.toml
│   ├── config.toml
│   ├── node_key.json
│   └── priv_validator_key.json
├── cosmosd
├── data
│   └── priv_validator_state.json
└── upgrade_manager
    └── genesis
        └── bin
            └── cyberd
```

As a result of this operation, the `data` and `config` folders should appear inside of your *$DAEMON_HOME/* folder.

8. Download `genesis.json` and place into your `.cyberd/config`:

```bash
cd $DAEMON_HOME/config
wget -O genesis.json https://ipfs.io/ipfs/QmZHpLc3H5RMXp3Z4LURNpKgNfXd3NZ8pZLYbjNFPL6T5n
```

9. Setup private peers in `config.toml`. You can find them on our [forum](https://ai.cybercongress.ai/t/euler-6-testnet-faq/65).

### Setup cyberd as a service (Ubuntu)

1. Increase resource limits for [Tendermint](https://tendermint.com):

```bash
ulimit -n 4096
```

2. Make cyberd a system service. This will help you easily start/stop cyberd and run it in the background:

```bash
sudo nano /etc/systemd/system/cyberd.service
```

Paste the following (replace `ubuntu` with your username, or if you running as `root` replce the whole */home/ubuntu/* to `/root/`):

```bash
[Unit]
Description=Cyber Node
After=network-online.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/.cyberd/
ExecStart=/home/ubuntu/.cyberd/cosmosd start --compute-rank-on-gpu=true
Environment=DAEMON_HOME=/home/ubuntu/.cyberd
Environment=DAEMON_NAME=cyberd
Environment=GAIA_HOME=/home/ubuntu/.cyberd
Restart=always
RestartSec=3m
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

If you need to enable search of the node add the flag `--allow-search=true` right after `--compute-rank-on-gpu=true`. If you need to run a rest-server alongside `cyberd` here is a service file for it (do `sudo nano /etc/systemd/system/cyberdcli-rest.service` and paste the following), just make sure you'll replace `ubuntu` to your user name and group:

```bash
[Unit]
Description=Cyberdcli REST Server

[Service]
User=ubuntu
Group=ubuntu
ExecStart=/usr/local/bin/cyberdcli rest-server --laddr tcp://0.0.0.0:1317 --chain-id euler-6
Restart=always
TimeoutSec=120
RestartSec=30

[Install]
WantedBy=multi-user.target
```

There's a possibility to build and run swagger-ui for your node to get a better experience with the rest-server. To get it up you'll have to install `static` library for Go: 

```bash
go get github.com/rakyll/statik
```

Then `cd` to the go-cyber repo and set static file for swagger-ui:

```bash
cd <path_to_go-cyber>/go-cyber/
statik -src=cmd/cyberdcli/swagger-ui -dest=cmd/cyberdcli/lcd -f
```

Rebuild cyberdcli and replace the one in `/usr/local/bin` (no worries, you won't lose your keys, if you already have keys imported):

```bash
make build
cp build/cyberdcli /usr/local/bin/
```

When all of the above steps are completed and cyberdcli-rest service has been started, you should have Swagger-ui available at `http://localhost:1317/swagger-ui/` .

3. Run cyberd:

Reload `systemd` after the creation of the new service:

```bash
systemctl daemon-reload
```

Start the node:

```bash
sudo systemctl start cyberd
```

Check node status:

```bash
sudo systemctl status cyberd
```

Enable service:

```bash
sudo systemctl enable cyberd
```

Check logs:

```bash
journalctl -u cyberd -f --lines 50
```

If you need to stop the node:

```bash
sudo systemctl stop cyberd
```

All commands in this section are also applicable to `cyberdcli-rest.service`.

At this point your cyberd should be running in the backgroud and you should be able to call `cyberdcli` to operate with the client. Try calling `cyberdcli status`. A possible output looks like this:

```bash
{"node_info":{"protocol_version":{"p2p":"6","block":"9","app":"0"},"id":"93b776d3eb3f3ce9d9bda7164bc8af3acacff7b6","listen_addr":"tcp://0.0.0.0:26656","network":"euler-6","version":"0.32.7","channels":"4020212223303800","moniker":"anon","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"686B4E65415D4E56D3B406153C965C0897D0CE27004E9CABF65064B6A0ED4240","latest_app_hash":"0A1F6D260945FD6E926785F07D41049B8060C60A132F5BA49DD54F7B1C5B2522","latest_block_height":"4553","latest_block_time":"2019-11-24T09:49:19.771375108Z","catching_up":false},"validator_info":{"address":"66098853CF3B61C4313DD487BA21EDF8DECACDF0","pub_key":{"type":"tendermint/PubKeyEd25519","value":"uZrCCdZTJoHE1/v+EvhtZufJgA3zAm1bN4uZA3RyvoY="},"voting_power":"0"}}
```

Your node has started to sync. If that didn't happen, check your config.toml file located at `$DAEMON_HOME/config/config.toml` and add at least a couple of addresses to <persistent_peers = ""> and <seeds = "">, some of those you can find on our [forum](https://ai.cybercongress.ai/).

Additional information about the chain is available via an API endpoint at: `localhost:26657` (access via your browser)

E.G. the number of active validators is available at: `localhost:26657/validators`

4. If your node did not launch correctly from the genesis, you need to set the current link to cosmosd for cyber daemon:

```bash
ln -s $DAEMON_HOME/upgrade_manager/genesis current
```

If you joined the testnet **after** a chain upgrade happened, you must point your current link to a new location (with an approptiatly upgraded binary file):

```bash
mkdir $DAEMON_HOME/upgrade_manager/upgrades
cp <path_to_upgraded_cyberd> $DAEMON_HOME/upgrade_manager/upgrades
ln -s $DAEMON_HOME/upgrade_manager/upgrades current
```

## Validator start

After your node has successfully synced, you can run a validator.

### Prepare a staking address

We included 1 million Ethereum addresses, over 10000 Cosmos addresses and all of `euler-4` validators addresses into the genesis file. This means that there's a huge chance that you already have some EUL tokens. Here are 3 ways to check this:

If you already have a cyberd address with EUL and know the seed phrase or your private key, just restore it into your local Keystore:

```bash
cyberdcli keys add <your_key_name> --recover
cyberdcli keys show <your_key_name>
```

If you have an Ethereum address that had ~0.2Eth or more at block 8080808 (on the ETH network), you probably received a gift and may import your Ethereum private key. To check your gift balance, paste your Ethereum address on [cyber.page](https://cyber.page).

> Please do not import high-value Ethereum accounts. This is not safe! cyberd software is new and has not been audited yet.

```bash
cyberdcli keys add private <your_key_name>
cyberdcli keys show <your_key_name>
```

If you want to create a new account, use the command below:
(You should send coins to that address to bound them later during the launch of the validator)

```bash
cyberdcli keys add <your_key_name>
cyberdcli keys show <your_key_name>
```

You could use your Ledger device, with the Cosmos app installed on it to sign and store cyber addresses: [guide here](https://github.com/cybercongress/cyberd/blob/0.1.5/docs/cyberd_Ledger_guide.md).
In most cases use the --ledger flag, with your commands:

```bash
cyberdcli keys add <your_key_name> --ledger
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to this parameter <your_key_name> later when you use the keys to sign transactions.
It will ask you to enter your password twice to encrypt the key.
You will also need to enter your password when you use your key to sign any transaction.

The command returns the address, a public key and a seed phrase, which you can use to
recover your account if you forget your password later.
Keep the seed phrase at a safe place (not in hot storage) in case you have to use it.

The address shown here is your account address. Let’s call this **<your_account_address>**.
It stores your assets.

**Important note**: Starting with v.38 cosmos-SDK uses os-native keyring to store all your keys. We've noticed that in certain cases it does not work well by default (for example if you don't have any GUI installed on your machine). If during the execution of the `cyberdcli keys add` command, you are getting this type of error:

```bash
panic: No such  interface 'org.freedesktop.DBus.Properties' on object at path /

goroutine 1 [running]:
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeInfo(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:479 +0x38c
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeLocalKey(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:465 +0x189
github.com/cosmos/cosmos-sdk/crypto/keys.baseKeybase.CreateAccount(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x146aa00, 0xc000b15630, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keybase_base.go:171 +0x192
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.CreateAccount(...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:107
github.com/cosmos/cosmos-sdk/client/keys.RunAddCmd(0xc000f0b400, 0xc000f125f0, 0x1, 0x1, 0x148dcc0, 0xc000aca550, 0xc000ea75c0, 0xc000ae1c08, 0x5e93b7)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/client/keys/add.go:273 +0xa8b
... etc
```

You will have to use another keyring backend to keep your keys. Here are 2 options: store the files within the cli folder or a `pass` manager.

Setting keyring backend to a **local file**:

Execute:

```bash
cyberdcli config keyring-backend file
```

As a result you migth see following: `configuration saved to /root/.cybercli/config/config.toml`

Execute:

```bash
cyberdcli config --get keyring-backend
```

The result should be the following:

```bash
user@node:~# cyberdcli config --get keyring-backend
file
```

This means that you've set your keyring-backend to a local file. *Note*, in this case, all the keys in your keyring will be encrypted using the same password. If you would like to set up a unique password for each key, you should set a unique `--home` folder for each key. To do that, just use `--home=/<unique_path_to_key_folder>/` with setup keyring backend and at all interactions with keys when using cyberdcli:

```bash
cyberdcli config keyring-backend file --home=/<unique_path_to_key_folder>/
cyberdcli keys add <your_second_key_name> --home=/<unique_path_to_key_folder>/
cyberdcli keys list --home=/<unique_path_to_key_folder>/
```

Set keyring backend to [**pass manager**](https://github.com/cosmos/cosmos-sdk/blob/9cce836c08d14dc6836d07164dd964b2b7226f36/crypto/keyring/doc.go#L30):

Pass utility uses a GPG key to encrypt your keys (but again, it uses the same GPG for all the keys). To install and generate your GPG key you should follow [this guide](https://www.passwordstore.org/) or this very [detailed guide](http://tuxlabs.com/?p=450). When you'll get your `pass` set, configure `cyberdcli` to use it as a keyring backend:

```bash
cyberdcli config keyring-backend pass
```

And verify that all has been set as planned:

```bash
cyberdcli config --get keyring-backend
pass
```

#### Send the create validator transaction

Validators are actors on the network committing to new blocks by submitting their votes.
This refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the nodes public key,
not the public key of the address you have just created.

To get the nodes public key run the following command:

```bash
cyberd tendermint show-validator
```

It will return a bech32 public key. Let’s call it **<your_node_pubkey>**.
The next step is to declare a validator candidate.
The validator candidate is the account which stakes the coins.
So the validator candidate is the account this time.
To declare a validator candidate, run the following command adjusting the staked amount and the other fields:

```bash
cyberdcli tx staking create-validator \
  --amount=10000000eul \
  --min-self-delegation "1000000" \
  --pubkey=<your_node_pubkey> \
  --moniker=<your_node_nickname> \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=euler-6
```

#### Verify that you are validating

```bash
cyberdcli query staking validators --trust-node=true
```

If you see your `<your_node_nickname>` with status `Bonded` and Jailed `false`, everything is good.
You are validating the network.

## Maintenance of the validator

### Jailing

If your validator got slashed, it will get jailed.
If it happens the operator must unjail the validator manually:

```bash
cyberdcli tx slashing unjail --from=<your_key_name> --chain-id euler-5
```
