package oci

import (
	"io"
)

type (
	// ClientOption allows specifying various settings configurable by the user for overriding the defaults
	// used when creating a new default client
	ClientOption func(*Client)
)

// ClientOptDebug returns a function that sets the debug setting on client options set
func ClientOptDebug(debug bool) ClientOption {
	return func(client *Client) {
		client.debug = debug
	}
}

// ClientOptWriter returns a function that sets the writer setting on client options set
func ClientOptWriter(out io.Writer) ClientOption {
	return func(client *Client) {
		client.out = out
	}
}

// ClientOptResolver returns a function that sets the resolver setting on client options set
func ClientOptResolver(resolver *Resolver) ClientOption {
	return func(client *Client) {
		client.resolver = resolver
	}
}

// ClientOptAuthorizer returns a function that sets the authorizer setting on client options set
func ClientOptAuthorizer(authorizer *Authorizer) ClientOption {
	return func(client *Client) {
		client.authorizer = authorizer
	}
}

// ClientOptRootPath returns a function that sets the rootpath setting on client options set
func ClientOptRootPath(rootPath string) ClientOption {
	return func(client *Client) {
		client.rootPath = rootPath
	}
}

// ClientOptCache returns a function that sets the cache setting on a client options set
func ClientOptCache(cache *Cache) ClientOption {
	return func(client *Client) {
		client.cache = cache
	}
}
