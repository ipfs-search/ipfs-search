# Snapshots
ipfs-search makes daily [OpenSearch snapshots](https://opensearch.org/docs/latest/opensearch/rest-api/snapshots/index/) of the indexed data.

Our most current index snapshots are available at https://ipfs-search-snapshots-v9.s3.eu-central-1.amazonaws.com/ over S3/HTTPS and can be loaded directly into an OpenSearch or OpenSearch cluster with sufficient disk space.

Our current production indexes are:
* `ipfs_files_v9`
* `ipfs_directories_v9`
* `ipfs_invalids_v8`
* `ipfs_partials_v9`

We highly recommend users to only restore these; other indexes might not be complete or up to date (although you're welcome to play with them!). As of the time of writing (April 4, 2022) these indexes together take up about 22 TB.

## Restoring
Our snapshots can be configured as a [read-only URL snapshot repository](https://www.elastic.co/guide/en/opensearch/reference/current/snapshots-read-only-repository.html) into an OpenSearch cluster. In order to do so, configure the following URL as the repository: https://ipfs-search-snapshots-v9.s3.eu-central-1.amazonaws.com/

### Steps
1.  Ensure that an OpenSearch 7+/OpenSearch cluster with sufficient disk space is available at `localhost:9200`.
2.  Add our repository URL to the `repositories.url.allowed_urls` setting in `opensearch.yml`:
    ```yaml
    allowed_urls: ["https://ipfs-search-snapshots-v9.s3.eu-central-1.amazonaws.com/*"]
    ```
3.  Restart your cluster for the config changes to take affect.
4.  Configure our snapshot repo as a read-only URL repository:
    ```sh
    curl -X PUT "localhost:9200/_snapshot/ipfs_search?pretty" -H 'Content-Type: application/json' -d'
    {
      "type": "url",
      "settings": {
        "url": "https://ipfs-search-snapshots-v9.s3.eu-central-1.amazonaws.com/"
      }
    }
    '
    ```
5.  List available snapshots:
    ```sh
    curl -X GET "localhost:9200/_snapshot/ipfs_search/_all?pretty"
    ```

    Reference: https://opensearch.org/docs/latest/opensearch/rest-api/snapshots/get-snapshot/

6.  Pick a *succesful* snapshot (substitute `<snapshot_id>` from the available snapshots above) from the list and start restoring it:
    ```sh
    curl -X POST "localhost:9200/_snapshot/ipfs_search/<snapshot_id>/_restore?pretty"
    ```
    **WARNING**: This initiates a large transfer and will take a considerable amount of time! Make sure you have a fast & reliable connection!

    Reference: https://opensearch.org/docs/latest/opensearch/rest-api/snapshots/restore-snapshot/

7.  Track progress of running snapshot restore task:
    ```sh
    curl "localhost:9200/_recovery/"
    ```

Once recovered, you should have *all* of our data available. As long as you don't make updates, future restores should be incremental and, hence, a lot faster.

## License
The ipfs-search.com index snapshots are available under the Open Database License, which can be found on: http://opendatacommons.org/licenses/odbl/1.0/.

In short, you are free:

* To share: To copy, distribute and use the database.
* To create: To produce works from the database.
* To adapt: To modify, transform and build upon the database.

As long as you:

* Attribute: You must attribute any public use of the database, or works produced from the database, in the manner specified in the ODbL. For any use or redistribution of the database, or works produced from it, you must make clear to others the license of the database and keep intact any notices on the original database.
* Share-Alike: If you publicly use any adapted version of this database, or works produced from an adapted database, you must also offer that adapted database under the ODbL.
* Keep open: If you redistribute the database, or an adapted version of it, then you may use technological measures that restrict the work (such as DRM) as long as you also redistribute a version without such measures.
