package types

import (
	. "github.com/cybercongress/go-cyber/types"
)

// map of map, where first key is cid, second key is account.String()
// second map is used as set for fast contains check
type Links map[CidNumber]CidLinks
type CidLinks map[CidNumber]map[AccNumber]struct{}

type Cid string // 32 byte string
type CidNumber uint64

func (links Links) Put(from CidNumber, to CidNumber, acc AccNumber) {
	cidLinks := links[from]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users := cidLinks[to]
	if users == nil {
		users = make(map[AccNumber]struct{})
	}
	users[acc] = struct{}{}
	cidLinks[to] = users
	links[from] = cidLinks
}

func (links Links) PutAll(newLinks Links) {
	for from := range newLinks {
		for to := range newLinks[from] {
			for u := range newLinks[from][to] {
				links.Put(from, to, u)
			}
		}
	}
}

func (links Links) Copy() Links {

	linksCopy := make(Links, len(links))

	for from := range links {
		fromLinks := make(CidLinks, len(links[from]))
		for to := range links[from] {
			users := make(map[AccNumber]struct{}, len(links[from][to]))
			for u := range links[from][to] {
				users[u] = struct{}{}
			}
			fromLinks[to] = users
		}
		linksCopy[from] = fromLinks
	}
	return linksCopy
}

func (links Links) IsAnyLinkExist(from CidNumber, to CidNumber) bool {

	toLinks, fromExists := links[from]
	if fromExists {
		linkAccs, toExists := toLinks[to]

		if toExists && len(linkAccs) != 0 {
			return true
		}
	}
	return false
}

func (links Links) IsLinkExist(from CidNumber, to CidNumber, acc AccNumber) bool {

	toLinks, fromExists := links[from]
	if fromExists {
		linkAccs, toExists := toLinks[to]

		if toExists && len(linkAccs) != 0 {
			_, exists := linkAccs[acc]
			return exists
		}
	}
	return false
}
