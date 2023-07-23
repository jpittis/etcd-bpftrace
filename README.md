## Usage

You can build the bpftrace script and execute it like so:

```
$ go run main.go $PATH_TO_ETCD && sudo bpftrace -p $ETCD_PROCESS_ID etcd.bt
```

Then if you run some operations against etcd like:

```
$ etcdctl put foo bar --ignore-lease=true
$ etcdctl get hello --limit=42 --consistency=s
```

The full gRPC request structs will be printed out:

```
Put key=foo value=bar lease=0 prev_kv=0 ignore_value=0 ignore_lease=1
Range key=hello range_end= limit=42 revision=0 sort_order=0 sort_limit=0 serializable=1 keys_only=0 count_only=0 min_mod_revision=0 max_mod_revision=0 min_create_revision=0 max_create_revision=0
```

## Architecture

The script is designed to work on amd64 and arch64.

## Future Work

This is just a demonstration/prototype. In the future, I'm considering:

- Extending the script to support all etcd KV RPCs.
- Figuring out how to compose more complex programs with bpftrace.
  - For example histograms, etc.
- Rewriting in a lower-level language so it's less janky.
