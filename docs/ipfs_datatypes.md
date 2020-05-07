# IPFS data types
Over the course of time, IPFS has seen a progression in data types. This document services as a reference to aid in the understanding of and interaction with the various ways in which data can be stored.

## CLI examples and HTTP API
Note: the examples in this document represent the output of go-ipfs version 0.23. It is important to note that, at least in theory, all the CLI commands have corresponding HTTP API calls.

Reference: https://docs.ipfs.io/reference/api/http/

## CID
Content identifiers CID in IPFS currently exist in two different versions: CIDv0 and CIDv1. Version 1 is self-descriptive in that the base, codec (format), CID version and hash are contained within the identifier. In version 0 these predefined in the format.

As of the time of writing, CIDv0 is still the default for files and directories being added to IPFS. With the release of go-ipfs 0.5 this should be updated to CIDv1.

Reference: https://github.com/multiformats/cid

### CIDv0
Example file: `QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv`

Base: `base58btc` (implicit)
Codec: `dag-pb` (implicit)
CID: `0` (implicit)
Hash: `sha2-256`

Format: `cidv0 ::= <multihash-content-address>`

### CIDv1
Example file: `bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4`

Base: `base32` (default)
Codec: `raw` (default)
CID: `1`
Hash: `sha2-256`

Format: `<cidv1> ::= <multibase-prefix><cid-version><multicodec-content-type><multihash-content-address>`

Any CIDv0 can be converted to CIDv1.

## IPLD
The superset of data addressable with CID, extends beyond IPFS to cover compatible content-addressable hash-linked datastructures in a unified information space. It is designed to be extensible, so that any compatible hash-linked structure may be addressed within the same Merkle forest.

Currently, this includes specifications to address and interact with the data formats of the following protocols:

* IPFS
* Bitcoin
* Ethereum
* GIT

Reference: https://ipld.io/

## IPFS data types

### Objects
In IPFS, all CIDs refer to objects with the following properties:

* `cid`: the identifier (CIDv0 or CIDv1, as described above).
* `size`: the cumulative size of the object, including the size of the objects linked to.
* `links`: a list of links to other objects in CID format, except for `raw` objects.
* `data`: encoded data according to the IPLD multicodec content type in the CID.

Objects may be explored with the IPLD explorer: https://explore.ipld.io/

### Data types

Within IPLD, the following data formats are currently usage within IPFS:

1. `dag-pb`: Protobuf wrapper around UnixFS protobuf format for files, raw data, HAMT, and directories.
    The `Data` field of the [dag-pb protobuf](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-pb.md) holds protobuf-encoded UnixFS data.
    See below for a dissection of a few exemplary blocks to make this more clear. 
2. `dag-cbor`: JSON-like generic datastructures for objects with links to other objects.
3. `raw`: Unencoded binary data.

The full list of IPLD formats (including various address format specifiers and other stuff) is defined here:
https://github.com/multiformats/multicodec/blob/master/table.csv

#### DAG protobuf (`dag-pb`)
This is a protobuf wrapper which seems to be only used for UnixFSv1-encoded data with links.
A draft specification can be found at:
https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-pb.md

Inside the `Data` field is the legacy IPFS UnixFS format, encoding files and directories according to the following Protobuf format:

https://github.com/ipfs/go-unixfs/blob/master/pb/unixfs.proto

Directories may be self-contained as `Directory` or may be sharded in Protobuf [HAMT](https://en.wikipedia.org/wiki/Hash_array_mapped_trie) format as `HAMTShard`. Files may be stored as `File` (binary data encoded in self-descriptive UnixFS protobuf) or as `Raw` (unencoded binary data).

With regards to HAMT, it is important to note that a generic (content-agnostic) HAMT is being developed on top of `dag-cbor`, to be use with UnixFS v2 (see below) and other formats.

The current implementation reference implementation lives in `go-unixfs`.

Source code: https://github.com/ipfs/go-unixfs
Godoc: https://godoc.org/github.com/ipfs/go-unixfs

#### DAG CBOR (`dag-cbor`)
This is a  [CBOR](https://cbor.io/)-based (extendable, binary JSON-compatible format) encoding for linked data. It's essentially JSON, encoded as a binary, with the addition of a tagged format for links (encoded as CID) to other objects.

Hence, CBOR may be converted to JSON, in which case the format tag for links will be lost.

#### UnixFS v2
It is important to know that a new version of the specification for UnixFS, the encoding for files and directories, is in active development. It will feature an IPLD-compliant format for encoding files, likely as an extendible subformat of DAG-CBOR.

Reference: https://docs.ipfs.io/guides/concepts/unixfs/

### Files and directories

#### Files
Files added to IPFS may be encoded in several different ways:

1. Raw binary data
2. Chunked UnixFS protobuf
3. Trickle-DAG UnixFS protobuf
4. Self-contained UnixFS protobuf files

It is important to note that, for any objects which may be represented as files, the `Size` field in the `ipfs files stat` output is the only way to acquire the size of the original file.

Moreover, the `Type` field in the `ipfs files stat` output can be used to exclude directories but cannot guarantee that the object stored is in fact a file, as it may well be that the object contains raw binary data which might not be considered a file.

Moreover, there is no reliable way to discriminate chunks of files although the `Size` proves to be a strong indicator.

##### 1. Raw binary data
Files are added as raw binary data to the object, without any encoding. The multiformat is `raw`. When files are added to a UnixFS HAMT directory, they are typically added as `raw`.

Examples:
``` $ ipfs add --cid-version 1 README.md
added bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4 README.md
```
Note that currently, it seems that CIDv1 seems to imply that files are added as `raw`. In this case, `--raw-leaves` has the same result.

``` $ ipfs --encoding=json files stat /ipfs/bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4 | jq
{
  "Hash": "bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4",
  "Size": 6060,
  "CumulativeSize": 6060,
  "Blocks": 0,
  "Type": "file"
}
```

``` $ ipfs --encoding=json object stat /ipfs/bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4 | jq
{
  "Hash": "bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4",
  "NumLinks": 0,
  "BlockSize": 0,
  "LinksSize": 0,
  "DataSize": 6060,
  "CumulativeSize": 6060
}
```

##### 2. Chunked UnixFS protobuf
When added files are more than a certain size, they are automatically chunked. The default chunker creates equal chunks of 262144 bytes.

A chunked file will have the following fields set:

* `type`: `"file"`
* `data`: `undefined`
* `blockSizes`: `[262144, 262144, ..., <LastChunkSize>]`
* `links`: `[[0, <Chunk1CID>], [1, <Chunk2CID>], ...]`
* `size`: The cumulative size of the object, including the metadata for encoding the chunks.

The chunks themselves are encoded as regular files, hence there is no reliable way to discriminate between chunked and regular files. It is only possible to identify the 'header' of chunked files although a `Size` of exactly `262144` is a strong indicator that we are in fact dealing with a chunk and not a full file - although this includes files which naturally happen to have this size.

Examples:
```
$ ipfs add 19022011252.jpg
added QmcsE5YiF8NVkdf5SwL4o2N5WDLkrUfhWwL91gv3MDosvx 19022011252.jpg
```

```
$ ipfs --encoding=json files stat /ipfs/QmcsE5YiF8NVkdf5SwL4o2N5WDLkrUfhWwL91gv3MDosvx | jq
{
  "Hash": "QmcsE5YiF8NVkdf5SwL4o2N5WDLkrUfhWwL91gv3MDosvx",
  "Size": 1720471,
  "CumulativeSize": 1720913,
  "Blocks": 7,
  "Type": "file"
}
```

```
$ ipfs --encoding=json object stat /ipfs/QmcsE5YiF8NVkdf5SwL4o2N5WDLkrUfhWwL91gv3MDosvx | jq
{
  "Hash": "QmcsE5YiF8NVkdf5SwL4o2N5WDLkrUfhWwL91gv3MDosvx",
  "NumLinks": 7,
  "BlockSize": 344,
  "LinksSize": 310,
  "DataSize": 34,
  "CumulativeSize": 1720913
}
```

Example for the first chunk:
```
ipfs --encoding=json files stat /ipfs/QmNd4BFWZ5tU5C3AVtL3RTg5r2hVLBKtGtDynK4cM36G3y | jq
{
  "Hash": "QmNd4BFWZ5tU5C3AVtL3RTg5r2hVLBKtGtDynK4cM36G3y",
  "Size": 262144,
  "CumulativeSize": 262158,
  "Blocks": 0,
  "Type": "file"
}
```
(Note the `Size` of `262144` bytes.)

```
ipfs --encoding=json object stat /ipfs/QmNd4BFWZ5tU5C3AVtL3RTg5r2hVLBKtGtDynK4cM36G3y | jq
{
  "Hash": "QmNd4BFWZ5tU5C3AVtL3RTg5r2hVLBKtGtDynK4cM36G3y",
  "NumLinks": 0,
  "BlockSize": 262158,
  "LinksSize": 4,
  "DataSize": 262154,
  "CumulativeSize": 262158
}
```

##### 3. Trickle-DAG UnixFS protobuf
Trickle-DAG is another way of chunking files in IPFS. It is unclear at this time how to discriminate between the normal (balanced) DAG and the trickle-DAG formats, or as to how or whether they need a different way to access.

> Trickle-dag is optimized for reading data in sequence, while the merkle-dag improves random access time. It might make sense to use trickle-dag for long videos, but in my experience it’s not a massive difference.

Reference: https://discuss.ipfs.io/t/what-is-the-difference-between-trickle-dag-and-merkle-dag/265/3

Examples (same file as in chunked UnixFS above):
```
$ ipfs add -t 19022011252.jpg
added QmWPzAVU681LaAt2GrveH4JEMTdXKJr7LjeZRpwN6C7zVj 19022011252.jpg
```

```
$ ipfs --encoding=json files stat /ipfs/QmWPzAVU681LaAt2GrveH4JEMTdXKJr7LjeZRpwN6C7zVj | jq
{
  "Hash": "QmWPzAVU681LaAt2GrveH4JEMTdXKJr7LjeZRpwN6C7zVj",
  "Size": 1720471,
  "CumulativeSize": 1720913,
  "Blocks": 7,
  "Type": "file"
}
```

```
$ ipfs --encoding=json object stat /ipfs/QmWPzAVU681LaAt2GrveH4JEMTdXKJr7LjeZRpwN6C7zVj | jq
{
  "Hash": "QmWPzAVU681LaAt2GrveH4JEMTdXKJr7LjeZRpwN6C7zVj",
  "NumLinks": 7,
  "BlockSize": 344,
  "LinksSize": 310,
  "DataSize": 34,
  "CumulativeSize": 1720913
}
```

Example for the first chunk:
```
$ ipfs --encoding=json files stat /ipfs/QmQ8VoFUWDKJGCBvRx6enp5CSijCWVCFDqjZ9gaCJAwr3P | jq
{
  "Hash": "QmQ8VoFUWDKJGCBvRx6enp5CSijCWVCFDqjZ9gaCJAwr3P",
  "Size": 262144,
  "CumulativeSize": 262158,
  "Blocks": 0,
  "Type": "file"
}
```

```
$ ipfs --encoding=json object stat /ipfs/QmQ8VoFUWDKJGCBvRx6enp5CSijCWVCFDqjZ9gaCJAwr3P | jq
{
  "Hash": "QmQ8VoFUWDKJGCBvRx6enp5CSijCWVCFDqjZ9gaCJAwr3P",
  "NumLinks": 0,
  "BlockSize": 262158,
  "LinksSize": 4,
  "DataSize": 262154,
  "CumulativeSize": 262158
}
```

##### 4. Self-contained UnixFS protobuf files
Like `raw` data, these are self-contained files, where the actual file data sits inside a protobuf container describing it as a file.

Examples:
```
$ ipfs add README.md
added QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv README.md
```

```
$ ipfs --encoding=json files stat /ipfs/QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv | jq
{
  "Hash": "QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv",
  "Size": 6060,
  "CumulativeSize": 6071,
  "Blocks": 0,
  "Type": "file"
}
```

```
$ ipfs --encoding=json object stat /ipfs/QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv | jq
{
  "Hash": "QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv",
  "NumLinks": 0,
  "BlockSize": 6071,
  "LinksSize": 3,
  "DataSize": 6068,
  "CumulativeSize": 6071
}
```

We can examine the `dag-pb(UnixFSv1(data))` wrapping through decoding the raw block data:
```
$ ipfs block get QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB | xxd
  00000000: 0acb 0808 0212 c308 4865 6c6c 6f20 616e  ........Hello an
  ...
  00000440: 7269 7479 2d6e 6f74 6573 0a18 c308       rity-notes....
```

We know from the CID that this is `protobuf`, so we can use [the dag-pb spec](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-pb.md) to decode it:
* `0acb08` is the interesting part: `0a` is `1010`, which tells us this is field number `1`, type `010`, which is an embedded message or bunch of bytes (the latter in our case).
    `cb08` is the length as `varint`, so that's `1100 1011 0000 1000` -> `100 1011 000 1000` -> reverse to `000 1000 100 1011` -> `1099` in decimal.
    The referenced spec tells us that field `1` is the `Data` field.

Due to historical reasons, the `Link` field is encoded before the `Data` field, but we have no links as this is a self-contained file.
We can parse the bytes following this as UnixFS protobuf.
We can also use `ipfs object data` to output just that, the `Data` field of an Object.    

```
$ ipfs object data QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB | xxd
  00000000: 0802 12c3 0848 656c 6c6f 2061 6e64 2057  .....Hello and W
  ...
  00000440: 792d 6e6f 7465 730a 18c3 08              y-notes....
```

Using the UnixFSv1 protobuf schema from above, this decodes as follows:
- `0802` -> `08` is `1000`, which tells us this is field number `1`, type `000`, which is a `varint` in this case.
    `02` is the value in `varint` encoding, which is just `2` in decimal.
    This is the `Type` field in the [unixfs protobuf spec](https://github.com/ipfs/go-unixfs/blob/master/pb/unixfs.proto).
    The type is `file`.
- `12c308` -> `12` is `1 0010`, which tells us this is field number `2`, type `010`, which is a `bytes` in this case.
    `c308` is the length, as `varint`: `1100 0011 0000 1000` -> `100 0011 000 1000` -> reverse to `000 1000 100 0011` -> `1091` in decimal.
- some data (which is actually `1091` bytes long, nice!)
- `18c3 08` -> `18` is `1 1000`, which tells us this is field number `3`, type `000`, which is a `varint` in this case.
    `c308` is, again `1091` in decimal, which is good, because this is a file and that's what you'd expect.
    This is the `filesize` field.

#### Directories
Directories added to IPFS may be encoded in several different ways:

1. UnixFS protobuf directory
2. UnixFS protobuf directory with raw leaves
3. UnixFS protobuf HAMT directory

##### `ls` commands and output
One relevant note here is that IPFS has a confusing plethora of `ls` commands and, while they are aware of the problem, resolving it (e.g. through deprecation) does not currently seem a prioerity.

Here's an overview, with the official description from the command's help text:

* `ipfs ls`: List directory contents for Unix filesystem objects.
* `ipfs files ls`: List directories in the local mutable namespace.
* `ipfs file ls`: List directory contents for Unix filesystem objects.

As `ipfs file` is currently not even listed anymore in `ipfs --help`, it does seem that `ipfs ls` is what we should go forward with. However, it is important to note that `ipfs file ls` and `ipfs ls` are *not* the same and that both the arguments as well as the output differ.

In these examples, we will be talking to the HTTP API equivalent to `ipfs ls` in the CLI as the latter currently seems to be unable to return JSON output.

Reference: https://github.com/ipfs/go-ipfs/issues/7050

##### Directory sizes
It seems that the only reliable way to determine the size of a directory is the `CumulativeSize` output of the `files stat` and `object stat` command, which represents the size of the directory and all its contents, including metadata (e.g. links to HAMT shards).

##### 1. UnixFS protobuf directory
This is the default, where files get added wrapped in protobuf containers.

Examples:
```
$ ipfs add -w README.md
added QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv README.md
added QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF
```

``` $ ipfs --encoding=json files stat /ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF | jq
{
  "Hash": "QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF",
  "Size": 0,
  "CumulativeSize": 6127,
  "Blocks": 1,
  "Type": "directory"
}
```

``` $ ipfs --encoding=json object stat /ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF | jq
{
  "Hash": "QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF",
  "NumLinks": 1,
  "BlockSize": 56,
  "LinksSize": 54,
  "DataSize": 2,
  "CumulativeSize": 6127
}
```

```
curl "http://localhost:5001/api/v0/ls?arg=/ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF",
      "Links": [
        {
          "Name": "README.md",
          "Hash": "QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv",
          "Size": 6060,
          "Type": 2,
          "Target": ""
        }
      ]
    }
  ]
}
```

Faster, streaming unordered results, without resolving file types and sizes, so that child objects do not need to be requested:
```
$ curl -s "http://localhost:5001/api/v0/ls?arg=/ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF&size=false&resolve-type=false&stream=true" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmQy8FTQCdJGtNHg92pc4B6F4cjnuSdrFx9gdRmkrE1rVF",
      "Links": [
        {
          "Name": "README.md",
          "Hash": "QmWyDJmrr6cRwEpTF2VGhWDi4uytrDHT8S5BptVdkbhjpv",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
```

We can again decode the protobuf encoding to understand what's going on behind the scenes:
```
$ ipfs block get QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv | xxd
  00000000: 122e 0a22 1220 a52c 3602 030c b912 edfe  ...". .,6.......
  00000010: 4de9 7002 fdad f9d4 5666 c3be 122a 2efb  M.p.....Vf...*..
  00000020: 5db9 3c1d 5fa6 1205 6162 6f75 7418 980d  ].<._...about...
  00000030: 1230 0a22 1220 929a 303c 39da 8a0b 67c0  .0.". ..0<9...g.
  ...
  00000150: 6375 7269 7479 2d6e 6f74 6573 1895 090a  curity-notes....
  00000160: 0208 01                                  ...
```

The block data decodes like this:
- `122e` - `12` is `0001 0010`, so field `2`, type embedded message.
    `2e` is the size: `0010 1110`, which is `46` in decimal.
    Fields:
    - `0a22` - `0a` is `0000 1010`, field `1`, type `bytes`, this is the `Hash` field.
        `22`, varint, `34` bytes. 
    - `1220 a52c 3602 030c b912 edfe 4de9 7002 fdad f9d4 5666 c3be 122a 2efb 5db9 3c1d 5fa6` binary hash.
        This is 34 bytes because it's a multihash, including one byte for the hash type (`sha2-256`) and the length (32 bytes)
    - `1205` - `12` is `0001 0010`, field `2`, type `string`, this is the `Name` field.
        `05`, varint, `5` bytes.
    - `6162 6f75 74` binary name, `about`.
    - `18 980d` - `18` is `0001 1000`, field `3`, type `varint`.
        `980d` is `1001 1000 0000 1101` -> `001 1000 000 1101` -> `000 1101 001 1000` -> `1688` in decimal.
        That's also what `ipfs object links` reports. Good.
- More links.
    They each start with `12`, indicating field `2`, embedded message.
- `0a02` at the end is the `Data` field, which is two bytes long.
    The actual data can be read out with `object data` again.

```
$ ipfs object data QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv | xxd
  00000000: 0801                                     ..
```

It reads "I am a directory".

The interesting thing to note here is that the UnixFS protobuf marks this block as a directory, but does not actually contain links to the files contained in the directory.
Those links are stored in the surrounding `dag-pb` Object.

##### 2. UnixFS protobuf with raw leaves

Examples:
```
$ ipfs add -w --raw-leaves README.md
added bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4 README.md
added QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm
```

``` $ ipfs --encoding=json files stat /ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm | jq
{
  "Hash": "QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm",
  "Size": 0,
  "CumulativeSize": 6118,
  "Blocks": 1,
  "Type": "directory"
}
```

``` $ ipfs --encoding=json object stat /ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm | jq
{
  "Hash": "QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm",
  "NumLinks": 1,
  "BlockSize": 58,
  "LinksSize": 56,
  "DataSize": 2,
  "CumulativeSize": 6118
}
```

```
curl -s "http://localhost:5001/api/v0/ls?arg=/ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm",
      "Links": [
        {
          "Name": "README.md",
          "Hash": "bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4",
          "Size": 6060,
          "Type": 2,
          "Target": ""
        }
      ]
    }
  ]
}
```

Faster, streaming unordered results, without resolving file types and sizes, so that child objects do not need to be requested:
```
$ curl -s "http://localhost:5001/api/v0/ls?arg=/ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm&size=false&resolve-type=false&stream=true" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmZzwcXprWah5w7qFPQ42UdGmokC4buH9ApNTJxmXhjZBm",
      "Links": [
        {
          "Name": "README.md",
          "Hash": "bafkreihqmkkhyq35uwiis5ed5mtudmv5abzdzzgop2urwp44uxutczahv4",
          "Size": 6060,
          "Type": 2,
          "Target": ""
        }
      ]
    }
  ]
}
```
Note that our request not to resolve the types and sizes for protobuf objects with raw leaves seem to be ignored.

From a protobuf perspective, raw-leaves directories behave identically to "normal" protobuf directories.
The only difference is that the CIDs referenced in the `Links` field of the `dag-pb` shell are `raw`.

##### 3. UnixFS protobuf HAMT directory
Protobuf HAMT directories are a special kind of directories which are accessed as a [HAMT](https://en.wikipedia.org/wiki/Hash_array_mapped_trie). This allows for fast indexing of very large directories, as protobuf directories require scanning all the blocks (chunks) to find a particular link.

The current HAMT implementation seems to have been done with the specific use case of distributing an uncensorable mirror of Wikipedia over IPFS. For the examples here, we will be using the hash `QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco` which is an IPFS-distributed snapshot of Wikipedia.

Note that there currently seems no reliable way to discriminate between HAMT and normal directories, even though the protobuf type is indeed different (HAMTShard instead of Directory).

References:
* https://blog.ipfs.io/24-uncensorable-wikipedia/

Examples:

``` $ ipfs --encoding=json files stat /ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco | jq
{
  "Hash": "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
  "Size": 0,
  "CumulativeSize": 658038834798,
  "Blocks": 5,
  "Type": "directory"
}
```

``` $ ipfs --encoding=json object stat /ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco | jq
{
  "Hash": "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
  "NumLinks": 5,
  "BlockSize": 283,
  "LinksSize": 253,
  "DataSize": 30,
  "CumulativeSize": 658038834798
}
```

```
curl -s "http://localhost:5001/api/v0/ls?arg=/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "-",
          "Hash": "QmPQVLHXAcDLvdf6ni24YWgGwitVTwtpiFaKkMfzZKquUB",
          "Size": 0,
          "Type": 1,
          "Target": ""
        },
        {
          "Name": "I",
          "Hash": "QmNYBYYjwtCYXdA2KC68rX8RqXBu9ajBM6Gi8CBrXjvk1j",
          "Size": 0,
          "Type": 1,
          "Target": ""
        },
        {
          "Name": "M",
          "Hash": "QmaeP3RagknCH4gozhE6VfCzTZRU7U2tgEEfs8QMoexEeG",
          "Size": 0,
          "Type": 1,
          "Target": ""
        },
        {
          "Name": "index.html",
          "Hash": "QmdgiZFqdzUGa7vAFycnA5Xv2mbcdHSsPQHsMyhpuzm9xb",
          "Size": 154,
          "Type": 2,
          "Target": ""
        },
        {
          "Name": "wiki",
          "Hash": "QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp",
          "Size": 0,
          "Type": 1,
          "Target": ""
        }
      ]
    }
  ]
}
```

Faster, streaming unordered results, without resolving file types and sizes, so that child objects do not need to be requested:
```
$ curl -s "http://localhost:5001/api/v0/ls?arg=/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco&size=false&resolve-type=false&stream=true" | jq
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "-",
          "Hash": "QmPQVLHXAcDLvdf6ni24YWgGwitVTwtpiFaKkMfzZKquUB",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "I",
          "Hash": "QmNYBYYjwtCYXdA2KC68rX8RqXBu9ajBM6Gi8CBrXjvk1j",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "wiki",
          "Hash": "QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "M",
          "Hash": "QmaeP3RagknCH4gozhE6VfCzTZRU7U2tgEEfs8QMoexEeG",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
{
  "Objects": [
    {
      "Hash": "/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
      "Links": [
        {
          "Name": "index.html",
          "Hash": "QmdgiZFqdzUGa7vAFycnA5Xv2mbcdHSsPQHsMyhpuzm9xb",
          "Size": 0,
          "Type": 0,
          "Target": ""
        }
      ]
    }
  ]
}
```
