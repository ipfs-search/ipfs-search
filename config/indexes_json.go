package config

// Write JSON here to keep correspondence with browser-based settings
const indexSettingsJSON = `{
    "index": {
        "refresh_interval": "15m",
        "mapping": {
            "total_fields": {
                "limit": "8192"
            }
        },
        "queries": {
            "cache": {
                "enabled": "true"
            }
        },
		"number_of_shards" : "5",
		"number_of_replicas": "1"
    }
}`

const fileMappingJSON = `{
    "_doc": {
	    "dynamic": "false",
	    "dynamic_templates": [
	        {
	            "default_noindex": {
	                "match": "*",
	                "mapping": {
	                    "index": "no",
	                    "doc_values": false
	                }
	            }
	        }
	    ],
	    "properties": {
	        "first-seen": {
	            "type": "date",
	            "format": "strict_date_optional_time||epoch_millis",
	            "index": true,
	            "doc_values": true
	        },
	        "last-seen": {
	            "type": "date",
	            "format": "strict_date_optional_time||epoch_millis",
	            "index": true,
	            "doc_values": true
	        },
	        "content":  {
	            "type": "text",
	            "index": true
	        },
	        "metadata":  {
	            "type":     "object",
	            "dynamic":  true,
	            "properties": {
	                "title" : {
	                    "type": "text",
	                    "index": true
	                },
	                "name": {
	                    "type": "text",
	                    "index": true
	                },
	                "author": {
	                    "type": "text",
	                    "index": true
	                },
	                "description": {
	                    "type": "text",
	                    "index": true
	                },
	                "producer": {
	                    "type": "text",
	                    "index": true
	                },
	                "publisher": {
	                    "type": "text",
	                    "index": true
	               	},
	                "isbn": {
	                    "type": "keyword",
	                    "index": true
	               	},
	                "language": {
	                    "type": "keyword",
	                    "index": true
	                },
	                "keywords": {
	                   "type": "text",
	                   "index": true
	               },
	                "xmpDM:album": {
	                    "type": "text",
	                    "index": true
	                },
	                "xmpDM:albumArtist": {
	                    "type": "text",
	                    "index": true
	                },
	                "xmpDM:artist": {
	                    "type": "text",
	                    "index": true
	                },
	                "xmpDM:composer": {
	                    "type": "text",
	                    "index": true
	               	},
	                "Content-Type": {
	                    "type": "keyword",
	                    "index": true,
	                    "doc_values": true
	                },
	                "X-Parsed-By": {
	                    "type": "keyword",
	                    "index": true
	                },
	                "date": {
	                    "type": "date",
	                    "format": "strict_date_optional_time||epoch_millis",
	                    "index": true
	               	},
	                "modified": {
	                    "type": "date",
	                    "format": "strict_date_optional_time||epoch_millis",
	                    "index": true
	               	}
	            }
	        },
	        "urls": {
	            "type": "keyword",
	            "index": true
	        },
	        "size": {
	            "type": "long",
	            "ignore_malformed": true,
	            "index": true
	        },
	        "references":  {
	            "type":     "object",
	            "dynamic":  true,
	            "properties": {
	                "name": {
	                    "type": "text",
	                    "index": true
	                },
	                "hash": {
	                    "type": "keyword",
	                    "index": true
	                },
	                "parent_hash": {
	                    "type": "keyword",
	                    "index": true
	                }
	            }
	        }
	    }
    }
}`

const dirMappingJSON = `{
    "_doc": {
	    "dynamic": "strict",
	    "dynamic_templates": [
	        {
	            "default_noindex": {
	                "match": "*",
	                "mapping": {
	                    "index": "no",
	                    "doc_values": false
	                }
	            }
	        }
	    ],
	    "properties": {
	        "first-seen": {
	            "type": "date",
	            "format": "strict_date_optional_time||epoch_millis",
	            "index": true,
	            "doc_values": true
	        },
	        "last-seen": {
	            "type": "date",
	            "format": "strict_date_optional_time||epoch_millis",
	            "index": true,
	            "doc_values": true
	        },
	        "links":  {
	            "type":     "object",
	            "dynamic":  true,
	            "properties": {
	                "Hash": {
	                    "type": "keyword",
	                    "index": true
	                },
	                "Name": {
	                    "type": "text"
	                },
	                "Size": {
	                   "type": "long",
	                   "doc_values": true,
	                   "ignore_malformed": true
	                },
	                "Type": {
	                   "type": "keyword",
	                   "index": true,
	                   "doc_values": true
	                }
	             }
	        },
	        "size": {
	            "type": "long",
	            "ignore_malformed": true,
	            "index": true,
	            "doc_values": true
	        },
	        "references":  {
	            "type":     "object",
	            "dynamic":  true,
	            "properties": {
	                "name": {
	                    "type": "text",
	                    "index": true
	                },
	                "hash": {
	                    "type": "keyword",
	                    "index": true
	                },
	                "parent_hash": {
	                    "type": "keyword",
	                    "index": true
	                }
	            }
	        }
	    }
	}
}`

const invalidMappingJSON = `{
    "_doc": {
	    "dynamic_templates": [
	        {
	            "default_noindex": {
	                "match": "*",
	                "mapping": {
	                    "index": "no",
	                    "doc_values": false,
	                    "include_in_all": false
	                }
	            }
	        }
	    ],
	    "properties": {
	       "error": {
	          "type": "text",
	          "index": false
	       }
	    }
    }
}`
