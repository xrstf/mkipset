# mkipset

This repository contains an application that can be used to read IP addresses and
subnets from a text/JSON/YAML file and convert the entries into entries in an
[ipset](http://ipset.netfilter.org/) set. Its primary function is to allow webmasters
without root permissions to easily define a blacklist of IPs and have the server
administrator use that list to automatically manage the firewall. Therefore the
goal is to make it near impossible to fuck up the firewall by inexperienced
webmasters and to separate system permissions.

In its current state, `mkipset` should be run as a regular cronjob (the more frequent,
the better). A daemon mode that permanently watches blacklist files may be added
in the future.

Entries in the blacklists can have a start and/or ending date, but these do **not**
map to ipset's native timing functionality. Instead they only serve as a hint
to `mkipset`. Likewise, it is assumed that `mkipset` runs as a cronjob, so there
is currently no support for saving/restoring ipset sets on server restores.
