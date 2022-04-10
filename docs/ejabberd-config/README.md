### ejabberd configuration instructions for a full featured I2P federation.

* Install ejabberd:
```
dpkg-reconfigure exim4-config
```
* Edit "smarthost" section of /etc/exim4/conf.d/router/200_exim4-config_primary to be:
```
smarthost:
  debug_print = "R: smarthost for $local_part@$domain"
  driver = manualroute
  domains = ! +local_domains
  transport = tunneled_smtp
  route_list = * 127.0.1.1
  host_find_failed = pass
  same_domain_copy_routing = yes
  no_more

.endif
```
* Create /etc/exim4/conf.d/transport/40_tunneled_smtp to be:
```
tunneled_smtp:
    debug_print = "T: tunneled remote_smtp for $local_part@$domain"
    driver = smtp
    multi_domain = false
    allow_localhost = true
    port = 7658
    hosts = 127.0.0.1
    hosts_override = true
```
* Proceed to building and configuing smtprelay-I2P project.


