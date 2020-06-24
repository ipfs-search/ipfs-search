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
	    "dynamic": "strict",
	    "dynamic_templates": [
	        {
	            "default_noindex": {
	                "match": "*",
	                "mapping": {
	                    "index": false,
	                    "doc_values": false,
      					"norms": false
	                }
	            }
	        }
	    ],
	    "properties": {
	        "first-seen": {
	            "type": "date",
	            "format": "date_time_no_millis"
	        },
	        "last-seen": {
	            "type": "date",
	            "format": "date_time_no_millis"
	        },
	        "content":  {
	            "type": "text"
	        },
	        "ipfs_tika_version": {
	        	"type": "keyword"
	        },
	        "language": {
	        	"properties": {
	        		"confidence": {
	        			"type": "keyword"
	        		},
	        		"language": {
	        			"type": "keyword"
	        		},
	        		"rawScore": {
	        			"type": "double"
	        		}
	        	}
	        },
	        "metadata":  {
	            "dynamic":  "true",
	            "properties": {
	                "title" : {
	                    "type": "text"
	                },
	                "name": {
	                    "type": "text"
	                },
	                "author": {
	                    "type": "text"
	                },
	                "description": {
	                    "type": "text"
	                },
	                "producer": {
	                    "type": "text"
	                },
	                "publisher": {
	                    "type": "text"
	               	},
	                "isbn": {
	                    "type": "keyword"
	               	},
	                "language": {
	                    "type": "keyword"
	                },
	                "resourceName": {
	                	"type": "keyword"
	                },
	                "keywords": {
	                   "type": "text"
	                },
	                "xmpDM:album": {
	                    "type": "text"
	                },
	                "xmpDM:albumArtist": {
	                    "type": "text"
	                },
	                "xmpDM:artist": {
	                    "type": "text"
	                },
	                "xmpDM:composer": {
	                    "type": "text"
	               	},
	                "Content-Type": {
	                    "type": "keyword"
	                },
	                "X-Parsed-By": {
	                    "type": "keyword"
	                },
	                "created": {
	                    "type": "date",
	                    "format": "date_optional_time"
	               	},
	                "date": {
	                    "type": "date",
	                    "format": "date_optional_time"
	               	},
	                "modified": {
	                    "type": "date",
	                    "format": "date_optional_time"
	               	}
	            }
	        },
	        "urls": {
	            "type": "keyword"
	        },
	        "size": {
	            "type": "long",
	            "ignore_malformed": true
	        },
	        "references":  {
	            "properties": {
	                "name": {
	                    "type": "text"
	                },
	                "hash": {
	                    "type": "keyword"
	                },
	                "parent_hash": {
	                    "type": "keyword"
	                }
	            }
	        }
	    }
    }
}`

const dirMappingJSON = `{
    "_doc": {
	    "dynamic": "strict",
	    "properties": {
	        "first-seen": {
	            "type": "date",
	            "format": "date_time_no_millis"
	        },
	        "last-seen": {
	            "type": "date",
	            "format": "date_time_no_millis"
	        },
	        "links":  {
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
	                   "ignore_malformed": true
	                },
	                "Type": {
	                   "type": "keyword"
	                }
	             }
	        },
	        "size": {
	            "type": "long",
	            "ignore_malformed": true
	        },
	        "references":  {
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
	                    "include_in_all": false,
	                    "norms": false
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
