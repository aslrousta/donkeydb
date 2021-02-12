# DonkeyDB

DonkeyDB is an impractical disk-backed database intented to study key-value
databases.

## Roadmap

- [x] Disk-backed storage (small memory footprint)
- [x] Support basic `GET` and `SET` commands for string values
- [x] Can be used both as a server and embedded
- [ ] Benchmark (correctness, memory, performance, compare with Redis)
- [ ] Support for `DEL` command
- [ ] Compress and reuse empty pages
- [ ] Support for integer values
