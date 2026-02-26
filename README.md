# Tamagotchi through SSH!
https://github.com/user-attachments/assets/d7d3fd87-cbe6-49e6-ba79-cc897cd059d6


# Deploying
You can build this project like this
```sh
$ go install
```
This will install it at $GOPATH/bin, which usually is ~/go/bin.

The easiest way to run it on your server is to use a SystemD service unit, you can use the following template to do so (this goes in /etc/systemd/system/mlssh.service):
```txt
[Unit]
Description=tamagotchi through ssh
After=network.target

[Service]
Type=simple
User=<your user>
WorkingDirectory=<a directory where the database will be stored, User should have read+write permissions>
ExecStart=/home/<User>/go/bin/mlssh
Restart=on-failure
Environment="MLSSH_PORT=22"
Environment="MLSSH_HOST=0.0.0.0"

[Install]
WantedBy=multi-user.target
```
Make sure to fill in the gaps, replacing &lt;placeholders&gt;.
By default this runs the server on port 22, to make it easier to connect to it.

This means that you have 2 options:

1. Use a non-default port to run this program (change MLSSH_PORT variable) and keep 22 for good-old ssh
2. Use the standard 22, but change your OpenSSH port to something else like 2222.



If you choose option 1 you need to give the executable permission to bind to privileged ports, like this:
```sh
$ sudo setcap 'cap_net_bind_service=+ep' $GOPATH/bin/mlssh
```

When you made your choice just enable the service like this and voilà:
```sh
$ sudo systemctl enable --now mlssh
```

