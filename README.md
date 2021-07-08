go-blockstore-tfs-uring
=======================

Why the long name?  Here's each chunk:

- It's in golang.
- It's a "blockstore" (it's for IPLD/IPFS -- storing chunks of content-addressed data).
- "TFS" stands for "trust the filesystem" (we're gonna assume you have a filesystem that does B+ reasonably well, and lean on that, rather than introducing other indexes).
- "uring" is a reference to the linux kernel `io_uring` API, which we're gonna use for maximum performance.