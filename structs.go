// Copyright 2013 Belogik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goes

import (
	"net/http"
	"net/url"
)

// Client represents a connection to elasticsearch
type Client struct {
	// The host to connect to
	Host string

	// The port to use
	Port string

	// Client is the http client used to make requests, allowing settings things
	// such as timeouts etc
	Client *http.Client
}

// Request holds a single request to elasticsearch
type Request struct {
	// Which connection will be used
	Conn *Client

	// A search query
	Query interface{}

	// Which index to search into
	IndexList []string

	// Which type to search into
	TypeList []string

	// HTTP Method to user (GET, POST ...)
	Method string

	// Which api keyword (_search, _bulk, etc) to use
	API string

	// Bulk data
	BulkData []byte

	// Request body
	Body []byte

	// A list of extra URL arguments
	ExtraArgs url.Values

	// Used for the id field when indexing a document
	ID string
}

// Response holds an elasticsearch response
type Response struct {
	Acknowledged bool
	Error        string
	Errors       bool
	Status       uint64
	Took         uint64
	TimedOut     bool  `json:"timed_out"`
	Shards       Shard `json:"_shards"`
	Hits         Hits
	Index        string `json:"_index"`
	ID           string `json:"_id"`
	Type         string `json:"_type"`
	Version      int    `json:"_version"`
	Found        bool
	Count        int

	// Used by the _stats API
	All All `json:"_all"`

	// Used by the _bulk API
	Items []map[string]Item `json:"items,omitempty"`

	// Used by the GET API
	Source map[string]interface{} `json:"_source"`
	Fields map[string]interface{} `json:"fields"`

	// Used by the _status API
	Indices map[string]IndexStatus

	// Scroll id for iteration
	ScrollID string `json:"_scroll_id"`

	Aggregations map[string]Aggregation `json:"aggregations,omitempty"`

	Raw map[string]interface{}
}

// Aggregation holds the aggregation portion of an ES response
type Aggregation map[string]interface{}

// Bucket represents a bucket for aggregation
type Bucket map[string]interface{}

// Document holds a document to send to elasticsearch
type Document struct {
	// XXX : interface as we can support nil values
	Index       interface{}
	Type        string
	ID          interface{}
	BulkCommand string
	Fields      interface{}
}

// Item holds an item from the "items" field in a _bulk response
type Item struct {
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Index   string `json:"_index"`
	Version int    `json:"_version"`
	Error   string `json:"error"`
	Status  uint64 `json:"status"`
}

// All represents the "_all" field when calling the _stats API
// This is minimal but this is what I only need
type All struct {
	Indices   map[string]StatIndex   `json:"indices"`
	Primaries map[string]StatPrimary `json:"primaries"`
}

// StatIndex contains stats for a specific index
type StatIndex struct {
	Primaries map[string]StatPrimary `json:"primaries"`
}

// StatPrimary contains stats for a primary index
type StatPrimary struct {
	// primary/docs:
	Count   int
	Deleted int
}

// Shard holds the "shard" struct as returned by elasticsearch
type Shard struct {
	Total      uint64
	Successful uint64
	Failed     uint64
}

// Hit holds a hit returned by a search
type Hit struct {
	Index     string                 `json:"_index"`
	Type      string                 `json:"_type"`
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Highlight map[string]interface{} `json:"highlight"`
	Fields    map[string]interface{} `json:"fields"`
}

// Hits holds the hits structure as returned by elasticsearch
type Hits struct {
	Total uint64
	// max_score may contain the "null" value
	MaxScore interface{} `json:"max_score"`
	Hits     []Hit
}

// SearchError holds errors returned from an ES search
type SearchError struct {
	Msg        string
	StatusCode uint64
}

// IndexStatus holds the status for a given index for the _status command
type IndexStatus struct {
	// XXX : problem, int will be marshaled to a float64 which seems logical
	// XXX : is it better to use strings even for int values or to keep
	// XXX : interfaces and deal with float64 ?
	Index map[string]interface{}

	Translog map[string]uint64
	Docs     map[string]uint64
	Merges   map[string]interface{}
	Refresh  map[string]interface{}
	Flush    map[string]interface{}

	// TODO: add shards support later, we do not need it for the moment
}
