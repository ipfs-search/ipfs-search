package ipfs

import (
	"context"
	"fmt"
	"github.com/dankinder/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

type LsTestSuite struct {
	suite.Suite

	ctx  context.Context
	ipfs *IPFS

	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	responseHeader http.Header
}

func (s *LsTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	s.responseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}

	cfg := DefaultConfig()
	cfg.IPFSAPIURL = s.mockAPIServer.URL()

	s.ipfs = New(cfg, http.DefaultClient, instr.New())
}

func (s *LsTestSuite) TearDownTest() {
	s.mockAPIServer.Close()
}

func (s *LsTestSuite) TestLsEmpty() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq", // file/ls: expected protobuf dag node
		},
	}

	rURL := fmt.Sprintf("/api/v0/ls?arg=%%2Fipfs%%2F%s&resolve-type=false&size=false&stream=true", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(``),
		}).
		Once()

	resultChan := make(chan<- *t.AnnotatedResource)
	err := s.ipfs.Ls(s.ctx, r, resultChan)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *LsTestSuite) TestLsHAMTDirectory() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp", // Wikipedia
		},
	}

	rURL := fmt.Sprintf("/api/v0/ls?arg=%%2Fipfs%%2F%s&resolve-type=false&size=false&stream=true", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`
				{"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Back_of_the_moon.html","Hash":"bafkreidnsi74hf7n2dtidxnqjdyr6lxidnsikdgwxktd7m3duwkuwl2u5u","Size":5169,"Type":2,"Target":""}]}]}
				{"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Munchh..html","Hash":"bafkreice7raasrty3makrm3gyg7sjqimdhhx6pdezh2noh3jlzwmvdcooy","Size":4986,"Type":2,"Target":""}]}]}
			`),
		}).
		Once()

	resultChan := make(chan *t.AnnotatedResource, 2)
	err := s.ipfs.Ls(s.ctx, r, resultChan)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	lsRes := <-resultChan
	s.Equal(lsRes, &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "bafkreidnsi74hf7n2dtidxnqjdyr6lxidnsikdgwxktd7m3duwkuwl2u5u",
		},
		Reference: t.Reference{
			Parent: r.Resource,
			Name:   "Back_of_the_moon.html",
		},
		Stat: t.Stat{
			// Note how, despite us requesting not to resolve, sharded directories
			// will return a type, as the HAMTDirectory will contain this information.
			Type: t.FileType,
			Size: 5169,
		},
	})

	lsRes = <-resultChan
	s.Equal(lsRes, &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "bafkreice7raasrty3makrm3gyg7sjqimdhhx6pdezh2noh3jlzwmvdcooy",
		},
		Reference: t.Reference{
			Parent: r.Resource,
			Name:   "Munchh..html",
		},
		Stat: t.Stat{
			Type: t.FileType,
			Size: 4986,
		},
	})
}

func (s *LsTestSuite) TestNormalDirectory() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv", // IPFS Hello world
		},
	}

	rURL := fmt.Sprintf("/api/v0/ls?arg=%%2Fipfs%%2F%s&resolve-type=false&size=false&stream=true", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`
				{"Objects":[{"Hash":"/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv","Links":[{"Name":"about","Hash":"QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V","Size":0,"Type":0,"Target":""}]}]}
				{"Objects":[{"Hash":"/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv","Links":[{"Name":"contact","Hash":"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y","Size":0,"Type":0,"Target":""}]}]}
			`),
		}).
		Once()

	resultChan := make(chan *t.AnnotatedResource, 2)
	err := s.ipfs.Ls(s.ctx, r, resultChan)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	lsRes := <-resultChan
	s.Equal(lsRes, &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V",
		},
		Reference: t.Reference{
			Parent: r.Resource,
			Name:   "about",
		},
		Stat: t.Stat{
			// Note how, due to the parameters resolve-type and size being false, child
			// nodes are not resolved and hence the type is undefined.
			Type: t.UndefinedType,
			Size: 0,
		},
	})

	lsRes = <-resultChan
	s.Equal(lsRes, &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		},
		Reference: t.Reference{
			Parent: r.Resource,
			Name:   "contact",
		},
		Stat: t.Stat{
			Type: t.UndefinedType,
			Size: 0,
		},
	})
}

func (s *LsTestSuite) TestLsInvalid() {
	errors := []string{
		"proto: required field \"Type\" not set",             // Example: QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8
		"proto: unixfs_pb.Data: illegal tag 0 (wire type 0)", // Example: QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj
		"unexpected EOF",                 // Example: QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq
		"unrecognized object type: 144",  // Example: z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW
		"not unixfs node (proto or raw)", // Example: z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE
	}

	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8", // file/ls: proto: required field "Type" not set
		},
	}

	rURL := fmt.Sprintf("/api/v0/ls?arg=%%2Fipfs%%2F%s&resolve-type=false&size=false&stream=true", r.ID)

	for _, errStr := range errors {
		msgStruct := &struct {
			Message string
			Code    int
			Type    string
		}{
			errStr, 0, "error",
		}

		s.mockAPIHandler.
			On("Handle", "POST", rURL, mock.Anything).
			Return(httpmock.Response{
				Header: s.responseHeader,
				Status: 500,
				Body:   httpmock.ToJSON(msgStruct),
			}).
			Once()

		resultChan := make(chan<- *t.AnnotatedResource)
		err := s.ipfs.Ls(s.ctx, r, resultChan)

		s.Error(err)
		s.mockAPIHandler.AssertExpectations(s.T())

		s.True(s.ipfs.IsInvalidResourceErr(err))
	}
}

func TestLsTestSuite(t *testing.T) {
	suite.Run(t, new(LsTestSuite))
}

// bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq -> file/ls: expected protobuf dag node
// $ curl -X POST "http://127.0.0.1:5001/api/v0/ls?arg=/ipfs/bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
// 100   102    0   102    0     0  91479      0 --:--:-- --:--:-- --:--:--   99k
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq",
//       "Links": []
//     }
//   ]
// }

// Wikipedia
// $ curl -v -N -X POST "http://127.0.0.1:5001/api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp"
// *   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 200 OK
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < X-Chunked-Output: 1
// < Date: Sun, 06 Dec 2020 21:15:42 GMT
// < Transfer-Encoding: chunked
// <
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Back_of_the_moon.html","Hash":"bafkreidnsi74hf7n2dtidxnqjdyr6lxidnsikdgwxktd7m3duwkuwl2u5u","Size":5169,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Munchh..html","Hash":"bafkreice7raasrty3makrm3gyg7sjqimdhhx6pdezh2noh3jlzwmvdcooy","Size":4986,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Canning_Town_DLR_station.html","Hash":"bafkreigm7ods7ezplrhkwddelmy4sf34qaduljkvtn3blo4g5te7g5ddhe","Size":78818,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Mangho_Pir.html","Hash":"bafkreifctzlb3r3uahnlj5qcuv5a6oc5oy7yxjvlr7jhd6rucpslnderyq","Size":6018,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Prashnopanishad.html","Hash":"bafkreic3bkjyy6nqc7gv4efp53rhrefguhfndhgga5u36ke6hxo3s7bwwq","Size":135859,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"SYN3_(gene).html","Hash":"bafkreid2vibdl3wrw2r64i4vys5lfvixuntlgxiedljzhx5wl3zweq5zjq","Size":43011,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Anarchist_Federation_(Britain_\u0026_Ireland).html","Hash":"bafkreicnewuujeyrav5xjvyrlicmprahrz25hnfv3eiopkz2phrhwypegm","Size":74081,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Hypolamprus_melilialis.html","Hash":"bafkreie2huyjxtg3cc4vpmiub4kmk7cv2wxeyb4c3bzgc2n7w545xp2gnq","Size":6696,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Jermaine_Brown_(Caymanian_footballer).html","Hash":"bafkreiay5c75rxlddzl76wnl3vv5jhspkgpxn4oiu42yfa7yeu5wyzte2i","Size":7518,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"Michael_Tsiselsky.html","Hash":"bafkreiewoqdyekzqv24uiqbedrmk6rpalokv5rfbbqdhmqa4fpqa4fihza","Size":17731,"Type":2,"Target":""}]}]}
// {"Objects":[{"Hash":"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp","Links":[{"Name":"President_of_Rice_University.html","Hash":"bafkreibjv5doenpjoleuq6fz4ppq7bmbsiyt6rrhiu2nymmp7nusopi7zi","Size":153247,"Type":2,"Target":""}]}]}
// ...

// QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8 -> file/ls: proto: required field "Type" not set
// $ curl -v -N -X POST "http://127.0.0.1:5001/api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8"
// *   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8 HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 21:17:27 GMT
// < Transfer-Encoding: chunked
// <
// {"Message":"proto: required field \"Type\" not set","Code":0,"Type":"error"}

// Multiblock file
// $ curl -v -N -X POST "http://127.0.0.1:5001/api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/ls?resolve-type=false&size=false&stream=true&arg=/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 200 OK
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < X-Chunked-Output: 1
// < Date: Sun, 06 Dec 2020 21:18:19 GMT
// < Transfer-Encoding: chunked
// <
// { [768 bytes data]
// 100  3348    0  3348    0     0   678k      0 --:--:-- --:--:-- --:--:--  817k
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmddrRa6PVSnPTyMRBsPTpqnWTvc8n8kfqdt2iVGx5gv3m",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmRCYNxaJKaXEQEZYbzANjB9uCsiVYrDuY6TNqWtQQDamq",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmaNYRP83ARdjELmoQWTLoJ31vxn3zmBxK3d7vR6gAfiLZ",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmVNz2znpRpwPFafwbb6TJCN7FWrxv6eprQeJLnkA8sDqh",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmdVS3CcQMfJt8XvPcswNWgnQ2mHyk9wjvFSqrYxjFW83u",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmRP87qouYU5AewinSGvxog8d5zEuYonLtM9cTcSL8Rdr6",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmRd5xGHUtY2mNcrb8uSG8VrhcynChqJU9z3oYx2oLcdmJ",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmaEAM7XWkY9P8A4nmBK4qFDGxnEqadmevfkCpoodpjnna",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmQFWxRHjCKjPPZ92GP56cdtw7kCGwP2bNp6poHfmCuh6t",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmciHZfp2BT5yJyEW7U9w5LGPcsxpbPR5aVr4fcTVFdW97",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmWwXc53v3xS1Z8BHmKco8dTnNqQyveU5GdsF7aoRaAKka",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmYemNQVTQseKv4U5EpSbVSJHXtdt6HpHWrNFJtZmYgB4m",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmcoUbH5cDdkjbMXWWq5nw64UJbJQY7NFaGdXzoF4ptGV3",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmRQY4bwaot6wVtJiFHeK4VPYK1Z29BkDfLP9E3nLbTgVn",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmasVpGLWXGDs4aHha7DMLW3JU1b6mNF2VnHozjHiePNTq",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmQ4FawzQY28kUpRNdrCxp78D2jtf8q6n41Se6wvouBUhD",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
// {
//   "Objects": [
//     {
//       "Hash": "/ipfs/Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//       "Links": [
//         {
//           "Name": "",
//           "Hash": "QmNbPfbbvh3qvBjZtWSzxCrx2o6Vc6iV4xvyBf29eqzZZz",
//           "Size": 0,
//           "Type": 0,
//           "Target": ""
//         }
//       ]
//     }
//   ]
// }
