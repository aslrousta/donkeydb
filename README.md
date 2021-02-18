# DonkeyDB

DonkeyDB is an impractical disk-backed database intented to study key-value
databases.

## Roadmap

- [x] Disk-backed storage (small memory footprint)
- [x] Support basic `GET` and `SET` commands for string values
- [x] Can be used both as a server and embedded
- [x] Support for `DEL` command
- [x] Reuse empty pages
- [ ] Support for integer values
- [ ] Compress pages
- [ ] Benchmark (correctness, memory, performance, compare with Redis)
