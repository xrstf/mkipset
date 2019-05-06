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

## Usage

Create a configuration file by copying the `config.yaml.dist` and adjust to your
liking. Note that the set name can be at most 20 characters long, even though
ipset allows 31 characters. This is to allow `mkipset` to generate unique names
during set swaps.

After creating the config file, define a blacklist file. You can use Text/JSON/YAML
and its type will be determined by the file extension (.txt, .json or .yml/.yaml).
Look at the same blacklist files in the `examples/` directory for more information.

Set up a cronjob to run `mkipset` and make sure that your `PATH` is set up
properly so that it can find the `ipset` binary:

    PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

    # update ipsets
    * * * * * mkipset -config /etc/mkipset/web-config.yaml /etc/mkipset/web-config.txt

And finally create as many iptables rules as you like, referencing your set. In
this example, the set is named `blacklist-web` because it contains IPs that we want
to deny HTTP/HTTPS access to.

    iptables -A INPUT -p tcp --dport 80  -m set --match-set blacklist-web src -j DROP
    iptables -A INPUT -p tcp --dport 443 -m set --match-set blacklist-web src -j DROP

## License

MIT
