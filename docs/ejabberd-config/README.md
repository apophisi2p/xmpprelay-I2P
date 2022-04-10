## ejabberd configuration instructions for a full featured I2P federation

### Installing and configuring ejabberd XMPP server on I2P - STAND ALONE
* Modify your I2P tunnel configuration to allow for ejabberd services. See tunnels.conf file provided here.

* Download and install packages appropriate for your distribution from https://www.process-one.net/en/ejabberd/downloads/. These packages install ejabberd in /opt/ejabberd-XX.YY/. Your configuration and Mnesia database are available in /opt/ejabberd/.

* Copy ejabberd.yml file provided here to /opt/ejabberd/conf/ejabberd.yml and Modify it for your I2P domain.

* Generate certificates required by XMPP clients:
```
openssl genrsa -out /opt/ejabberd/conf/xxx.b32.i2p.key 2048
openssl req -new -x509 -key /opt/ejabberd/conf/xxx.b32.i2p.key -out /opt/ejabberd.conf//xxx.b32.i2p.crt -days 3650
chown ejabberd:ejabberd /opt/ejabberd/conf/*.b32.i2p.{key,crt}
chmod 640 /opt/ejabberd/conf/*.b32.i2p.{key,crt}
```

* Restart ejabbered service and create an admin XMPP account:
```
/opt/ejabberd-XX.YY/bin/ejabberdctl register admin example.i2p password
```

### Installing and configuring ejabberd XMPP server on I2P - FULL FEDERATON

* Clone and build xmpprelay-I2P project (specifics, logging and services will be added shortly)
* Run the relay as shown:
```
xmpprelay -l 127.0.0.1:9626 -r 127.0.0.1:4447
```
or
```
nohup xmpprelay -l 127.0.0.1:9626 -r 127.0.0.1:4447 &
``` 

* Install and configure dnsmasq service to "trick" XMPP DNS resolution for I2P domains.
Basically, sets all DNS resolution requests to return 127.0.0.1. (DNS record resolution is a big part of XMPP specs and operations):

	- Modify your /etc/resolv.conf to look similar to (add nameserver 127.0.0.1 to be the FIRST):

		```
		nameserver 127.0.0.1
		nameserver DNS SERVER1
		nameserver DNS SERVER2
		```
	- Modify you /etc/dnsmasq.conf to have the following directive at the END of the file:
		```
		address=/i2p/127.0.0.1
		```
	- Restart dnsmasq service and test to see that .i2p domains resolve to 127.0.0.1

* Your ejabberd XMPP server is FULLY FEDERATED now. All OPTIONS such as group chats, media transfers and others - are fully and transparently federated on I2P.


