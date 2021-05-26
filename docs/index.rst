.. ipfs-search documentation master file, created by
   sphinx-quickstart on Fri Mar 12 16:50:16 2021.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to ipfs-search documentation!
===========================================
.. image:: images/ipfs_search_og_image.png

This is the main documentation repository for ipfs-search.com

Search engine for the `Interplanetary Filesystem <https://ipfs.io>`_. Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using `ipfs-tika <https://github.com/ipfs-search/ipfs-tika>`_, searching is done using ElasticSearch 7, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

The ipfs-search command consists of two components: the crawler and the sniffer. The sniffer extracts hashes from the gossip between nodes. The crawler extracts data from the hashes and indexes them.

**Docs:** Documentation is hosted on here on `Read the Docs <https://ipfs-search.readthedocs.io/en/latest/>`_, based on files contained in the GitHub `docs <https://github.com/ipfs-search/ipfs-search/tree/master/docs>`_ folder. In addition, there's extensive `Go docs <https://pkg.go.dev/github.com/ipfs-search/ipfs-search>`_ for the internal API as well as `SwaggerHub OpenAPI documentation <https://app.swaggerhub.com/apis-docs/ipfs-search/ipfs-search/>`_ for the REST API.

**Contact:** Please find us on our Freenode/`Riot/Matrix <https://riot.im/app/#/room/#ipfs-search:chat.weho.st>`_ channel `#ipfs-search:chat.weho.st <https://matrix.to/#/#ipfs-search:chat.weho.st>`_.

.. figure:: images/ipfs-search-arch-inv.png
   :alt: ipfs-search current architecture
   :align: center
   :width: 800
   
   Current ipfs-search architecture. Credit: `Nina <https://niverel.tymyrddin.space/doku.php?id=en/start>`_. 

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   project_docs
   guides
   data_and_application_architecture
   distributed_search


