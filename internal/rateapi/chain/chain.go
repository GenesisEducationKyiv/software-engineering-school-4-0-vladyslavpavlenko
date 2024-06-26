package chain

import (
	"context"
)

type Fetcher interface {
	Fetch(ctx context.Context, base, target string) (string, error)
}

// Chain interface defines a chain of responsibility for rate fetching.
type Chain interface {
	Fetcher
	SetNext(Chain)
}

// Node is a concrete implementation of the Chain interface.
type Node struct {
	fetcher Fetcher
	next    Chain
}

// NewNode creates a new Node with the given rateapi.Fetcher.
func NewNode(fetcher Fetcher) *Node {
	return &Node{
		fetcher: fetcher,
	}
}

// SetNext sets the next chain in the responsibility chain.
func (n *Node) SetNext(next Chain) {
	n.next = next
}

// Fetch fetches the rate and delegates to the next chain if necessary.
func (n *Node) Fetch(ctx context.Context, base, target string) (string, error) {
	rate, err := n.fetcher.Fetch(ctx, base, target)
	if err != nil {
		next := n.next
		if next == nil {
			return "", err
		}

		return next.Fetch(ctx, base, target)
	}

	return rate, nil
}
