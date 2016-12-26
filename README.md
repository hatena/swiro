# swiro

swiro is a switching route tool for AWS to realize VIP (Virtual IP) with Routing-Based High Availability pattern.

This pattern is possible to perform failover (switching of the connection destination) of the EC2 redundant across the subnet (AZ).


## Usage

* Switching routes

```
$ swiro switch -r rtb-xxxxxx -v 10.0.0.1 -I i-xxxxxx
```


## Example

In most cases you can switch the routing with the Route Table ID as follows:

```
$ swiro switch -r rtb-xxxxxx -v 10.0.0.1 -I i-xxxxxx
Switch the route below setting:
============================================
Route Table: route_table_name (rtb-xxxxxx)
Virtual IP:  10.0.0.1 -------- Src:  src_instance_name (i-yyyyyy)
                      \\
                       ======> Dest: i-xxxxxx
============================================
Are you sure? (y/n) [y]: y
Success!!
```

You can also switch by specifying Route Table Name instead of Route Table ID.

```
$ swiro switch -r route_table_name -v 10.0.0.1 -I instance_name
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/taku-k/swiro
```

## Contribution

1. Fork ([https://github.com/taku-k/swiro/fork](https://github.com/taku-k/swiro/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[taku-k](https://github.com/taku-k)
