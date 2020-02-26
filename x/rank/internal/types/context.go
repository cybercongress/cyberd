package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/link"
	"sort"
)

type CalculationContext struct {
	CidsCount  int64
	LinksCount int64

	inLinks  map[link.CidNumber]link.CidLinks
	outLinks map[link.CidNumber]link.CidLinks

	stakes map[AccNumber]uint64

	FullTree bool

	DampingFactor float64
	Tolerance 	  float64
}

func NewCalcContext(
	ctx sdk.Context, linkIndex LinkIndexedKeeper, numberKeeper CidNumberKeeper,
	stakeKeeper StakeKeeper, fullTree bool, dampingFactor float64, tolerance float64) *CalculationContext {

	return &CalculationContext{
		CidsCount:  int64(numberKeeper.GetCidsCount(ctx)),
		LinksCount: int64(linkIndex.GetLinksCount(ctx)),

		inLinks:  linkIndex.GetInLinks(),
		outLinks: linkIndex.GetOutLinks(),

		stakes: stakeKeeper.GetTotalStakes(),

		FullTree: fullTree,

		DampingFactor: dampingFactor,
		Tolerance: tolerance,
	}
}

func (c *CalculationContext) GetInLinks() map[link.CidNumber]link.CidLinks {
	return c.inLinks
}

func (c *CalculationContext) GetOutLinks() map[link.CidNumber]link.CidLinks {
	return c.outLinks
}

func (c *CalculationContext) GetCidsCount() int64 {
	return c.CidsCount
}

func (c *CalculationContext) GetStakes() map[AccNumber]uint64 {
	return c.stakes
}

func (c *CalculationContext) GetTolerance() float64 {
	return c.Tolerance
}

func (c *CalculationContext) GetDampingFactor() float64 {
	return c.DampingFactor
}

func (с *CalculationContext) GetSortedInLinks(cid link.CidNumber) (link.CidLinks, []link.CidNumber, bool) {
	links := с.inLinks[cid]

	if len(links) == 0 {
		return nil, nil, false
	}

	numbers := make([]link.CidNumber, 0, len(links))
	for num := range links {
		numbers = append(numbers, num)
	}

	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })

	return links, numbers, true
}
