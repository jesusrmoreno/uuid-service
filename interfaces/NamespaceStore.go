package intf

// IDStore is the interface that must be satisfied to save IDs in a namespace
type IDStore interface {
	Store(ns, id string) (bool, error)
	Contains(ns, id string) (bool, error)
	All(ns string) ([]string, error)
}
