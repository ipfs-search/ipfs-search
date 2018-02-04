/* jshint node: true, esnext: true */

const elasticsearch = require('elasticsearch');
const http = require('http');
const url = require('url');
const htmlEncode = require('js-htmlencode');
const downsize = require('downsize');

const server_port = 9615;
const result_description_length = 250;

var client = new elasticsearch.Client({
  host: 'localhost:9200',
  log: 'info'
});

function query(q, page, page_size) {
  var body = {
      "query": {
          "query_string": {
              "query": q,
              "default_operator": "AND"
          }
      },
      "highlight": {
          "order" : "score",
          "require_field_match": false,
          "encoder": "html",
          "fields": {
              "*": {
                  "number_of_fragments" : 1,
                  "fragment_size" : result_description_length
              }
          }
      },
      "_source": [
        "metadata.title", "metadata.name", "metadata.description",
        "references", "size", "last-seen", "first-seen"
      ]
  };

  return client.search({
    index: 'ipfs',
    body: body,
    size: page_size,
    from: page*page_size
  });
}

function error_response(response, code, error) {
  console.trace(code+": "+error.message);

  response.writeHead(code, {"Content-Type": "application/json"});
  response.write(JSON.stringify({
    "error": error
  }));
  response.end();
}

function get_title(result) {
  // Get title from result

  // Highlights take preference
  var hl = result.highlight;

  if (hl) {
    const highlight_priority = [
      "metadata.title", "references.name"
    ];

    // Return the first one from the priority list
    for (var i=0; i<highlight_priority.length; i++) {
      if (hl[highlight_priority[i]]) {
        return hl[highlight_priority[i]][0];
      }
    }
  }

  // Try metadata
  var src = result._source;
  var titles = [];

  if ("metadata" in src) {
    const metadata_priority = [
      "title", "name"
    ];

    metadata_priority.forEach(function (item) {
      if (src.metadata[item]) {
        titles.push(src.metadata[item][0]);
      }
    });
  }

  // Try references
  src.references.forEach(function (item) {
    if (item.name) {
      titles.push(item.name[0]);
    }
  });

  // Pick longest title
  if (titles.length > 0) {
    titles.sort(function (a, b) { return b.length - a.length; });

    return htmlEncode.htmlEncode(titles[0]);
  } else {
    // Fallback to id
    return htmlEncode.htmlEncode(result._id);
  }
}

function get_description(result) {
  // Use highlights, if available
  if (result.highlight) {
    if (result.highlight.content) {
      return result.highlight.content[0];
    }

    if (result.highlight["links.Name"]) {
      // Reference name matching
      return "Links to &ldquo;"+result.highlight["links.Name"][0]+"&rdquo;";
    }

    if (result.highlight["links.Hash"]) {
      // Reference name matching
      return "Links to &ldquo;"+result.highlight["links.Hash"][0]+"&rdquo;";
    }
  }

  var metadata = result._source.metadata;
  if (metadata) {
    // Description, if available
    if (metadata.description) {
      return htmlEncode.htmlEncode(
        downsize(metadata.description[0], {
          "characters": result_description_length, "append": "..."
        })
      );
    }

  }

  // Default to nothing
  return null;
}

function transform_results(results) {
  var hits = [];

  results.hits.forEach(function (item) {
    hits.push({
      "hash": item._id,
      "title": get_title(item),
      "description": get_description(item),
      "type": item._type,
      "size": item._source.size,
      "first-seen": item._source['first-seen'],
      "last-seen": item._source['last-seen']
    });
  });

  // Overwrite existing list of hits
  results.hits = hits;
}

console.info("Starting server on http://localhost:"+server_port+"/");

http.createServer(function(request, response) {
  const page_size = 15;
  var parsed_url;

  try {
    try {
      parsed_url = url.parse(request.url, true);
    } catch(err) {
      error_response(response, 400, err.message);
    }

    if (parsed_url.pathname === "/search") {
      if (!"q" in parsed_url.query) {
        error_response(response, 422, "query argument missing");
      }

      var page = 0;
      const max_page = 100;

      if ("page" in parsed_url.query) {
        page = parseInt(parsed_url.query.page, 10);

        // For performance reasons, don't allow paging too far down
        if (page > 100) {
          error_response(422, "paging not allowed beyond 100");
        }
      }

      query(parsed_url.query.q, page, page_size).then(function (body) {
        console.info(request.url + " 200: Returning " + body.hits.hits.length + " results");

        body.hits.page_size = page_size;
        body.hits.page_count = Math.ceil(body.hits.total/page_size);

        transform_results(body.hits);

        response.writeHead(200, {"Content-Type": "application/json"});
        response.write(JSON.stringify(body.hits, null, 4));
        response.end();
      }, function (error) {
        throw(error);
      });

    } else {
        error_response(response, 404, "file not found");
    }

  } catch(err) {
    // Catch generic errors
    error_response(response, 500, err.message);
  } finally {

  }
}).listen(server_port);

