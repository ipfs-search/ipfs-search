# Using the Queue-pinservice BETA

Queue-pinservice is an ipfs-pinservice implementation, you can use the queue-pinservice to make ipfs-search.com pick up your new/changed CIDs automatically.

It is an  MVP, in BETA stage. It is free to use but no guarantees are given about service quality or availability, nor that it will stay completely free forever, and using it is on one's own responsibility. If you are interested in high-volume usage of the service, please contact info@ipfs-search.com

For documentation about using pinning services, see https://docs.ipfs.tech/how-to/work-with-pinning-services/

## Authentication
Note that (for now), authentication has been disabled, because there is no persistent data storage.
Nonetheless, the ipfs client expects an authentication key and won't work without one. You can use anything as an authentication key, except for nothing.

## IPFS Desktop or IPFS Web UI
Add a Custom service as described here: https://docs.ipfs.tech/how-to/work-with-pinning-services/#use-an-existing-pinning-service

- Name: ipfs-search
- API endpoint: `https://api.ipfs-search.com/v1/queue-pinservice/`
- Secret access token: "secretAccessToken" (see authentication above)

## Command line usage
Setting up your IPFS client:
```
ipfs pin remote service add ipfs-search https://api.ipfs-search.com/v1/queue-pinservice/ anyAuthenticationKey
```

Sending a CID to this queue pinning service:
```
ipfs pin remote add --service=ipfs-search --name=war-and-peace.txt bafybeib32tuqzs2wrc52rdt56cz73sqe3qu2deqdudssspnu4gbezmhig4
```

**N.b.** Because the IPFS client immediately after **Add pin** checks for the status of the request using **Get pin object**, this gives a not-implemented-error (code `456`).
This does not mean the call did not come through! There is simply no persistent data to retrieve about the call, and no way to reconstruct this information (at least for now).

## Pinning service API spec implementation

The queue-pinning service complies with https://ipfs.github.io/pinning-services-api-spec/,
but only [Add Pin](https://ipfs.github.io/pinning-services-api-spec/#operation/addPin) has been implemented.
[Replace pin object](bafybeib32tuqzs2wrc52rdt56cz73sqe3qu2deqdudssspnu4gbezmhig4) is routed to the **Add pin** service.

**List pin objects** returns an empty object.

Other calls throw a not-implemented-error with code `456`.

## Source code/repository
https://github.com/ipfs-search/ipfs-search-queue-pinservice
