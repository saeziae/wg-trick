```
 _       __ ______      ______ ____   ____ ______ __ __
| |     / // ____/     /_  __// __ \ /  _// ____// //_/
| | /| / // / __ ______ / /  / /_/ / / / / /    / ,<
| |/ |/ // /_/ //_____// /  / _, _/_/ / / /___ / /| |
|__/|__/ \____/       /_/  /_/ |_|/___/ \____//_/ |_|
```

**WG-trick** is a tool that helps you configure the routes of WireGuard.

This program does:

- configure the routes of WireGuard
- intend for PC client

This program does **not**:

- configure all from scratch for you
- download key pairs from server
- run in a complicated environment (home use only!)

So at least you need to:

- generate key pairs on your client PC
- configure the peer on the server side

## Server side

### Server Install

```shell
git clone https://github.com/saeziae/wg-trick && cd wg-trick
make
sudo make install
```

the usage of the program is like:

```shell
wg-trick-server -l 127.0.0.1:8964 -c /etc/wireguard/wg0.conf
```

**The configuration needs modification, see next title.**

If you use systemd, here is the example of a daemon

```ini
[Unit]
Description=WG-trick-server Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/wg-trick-server -l 127.0.0.1:8964 -c /etc/wireguard/wg0.conf

Restart=on-failure

[Install]
WantedBy=multi-user.target

```

Lastly, you'll need a proxy like [Caddy](https://github.com/caddyserver/caddy) or Nginx to deal with https, here is a Caddy config example:

```caddyfile
vpn.example.com {
    # Redirect HTTP to HTTPS
    @http {
        protocol http
    }
    redir @http https://{host}{uri} permanent

    # Proxy to wg-trick-server
    reverse_proxy 127.0.0.1:8964

    # Automatic TLS with ACME
    tls {
        issuer acme
    }
}

```

### Server configure file

It just reads a WireGuard conf file, but a bit of additional options, here is the toy example:

```ini
[Interface]
Address = 192.168.1.1/24,10.0.0.2/8
Mask = 192.168.1.0/24 #additional
Endpoint = vpn.example.com #additional
PublicKey = #additional
PrivateKey =
ListenPort =
MTU =

[Peer]
# Wildchicken University
Endpoint = vpn.example.edu:1145
IsGateway = True #additional
PublicKey =
AllowedIPs = 10.0.0.0/8

[Peer]
# PC
PublicKey =
AllowedIPs = 192.168.1.2/32
```

- The `Mask` (not really a mask) indicates the subnet used by the WireGuard server.
- The `PublicKey` under `Interface` is not required by WireGuard but we use it to distribute to client.
- `Endpoint` under `Interface` is used to distribute your server domain (IP) to the client.
- `IsGateway` indicates another server not in our subnet to which we forward the packets targeting its subnet.

That's all!

## Client usage

`wg-trick` is the client script, written in bash.

Install:

```shell
wget https://raw.githubusercontent.com/saeziae/wg-trick/refs/heads/main/wg-trick
chmod +x wg-trick
sudo cp wg-trick /usr/local/bin/
```

Use:

```shell
sudo wg-trick connect vpn.example.com
```

The interface will use the domain name, you can also specify the private key if it is not the one in `/etc/wireguard/privatekey`:

```shell
sudo wg-trick connect vpn.example.com /path/to/private/key
```

Other commands alike to `wg-quick` also work, like turning it off:

```shell
sudo wg-trick down vpn.example.com
```
