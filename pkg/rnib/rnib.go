// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rnib

import (
	"context"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	toposdk "github.com/onosproject/onos-ric-sdk-go/pkg/topo"
)

// TopoClient R-NIB client interface
type TopoClient interface {
	WatchE2Connections(ctx context.Context, ch chan topoapi.Event) error
	GetCells(ctx context.Context, nodeID topoapi.ID) ([]*topoapi.E2Cell, error)
	GetE2NodeAspects(ctx context.Context, nodeID topoapi.ID) (*topoapi.E2Node, error)
	E2NodeIDs(ctx context.Context) ([]topoapi.ID, error)
}

// NewClient creates a new topo SDK client
func NewClient() (Client, error) {
	sdkClient, err := toposdk.NewClient()
	if err != nil {
		return Client{}, err
	}
	cl := Client{
		client: sdkClient,
	}

	return cl, nil

}

// Client topo SDK client
type Client struct {
	client toposdk.Client
}

// E2NodeIDs lists all of connected E2 nodes
func (c *Client) E2NodeIDs(ctx context.Context) ([]topoapi.ID, error) {
	objects, err := c.client.List(ctx, toposdk.WithListFilters(getControlRelationFilter()))
	if err != nil {
		return nil, err
	}

	e2NodeIDs := make([]topoapi.ID, len(objects))
	for _, object := range objects {
		relation := object.Obj.(*topoapi.Object_Relation)
		e2NodeID := relation.Relation.TgtEntityID
		e2NodeIDs = append(e2NodeIDs, e2NodeID)
	}

	return e2NodeIDs, nil
}

// GetE2NodeAspects gets E2 node aspects
func (c *Client) GetE2NodeAspects(ctx context.Context, nodeID topoapi.ID) (*topoapi.E2Node, error) {
	object, err := c.client.Get(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	e2Node := &topoapi.E2Node{}
	err = object.GetAspect(e2Node)

	return e2Node, err

}

// GetCells get list of cells for each E2 node
func (c *Client) GetCells(ctx context.Context, nodeID topoapi.ID) ([]*topoapi.E2Cell, error) {
	filter := &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: string(nodeID),
			RelationKind: topoapi.CONTAINS,
			TargetKind:   ""}}

	objects, err := c.client.List(ctx, toposdk.WithListFilters(filter))
	if err != nil {
		return nil, err
	}
	var cells []*topoapi.E2Cell
	for _, obj := range objects {
		targetEntity := obj.GetEntity()
		if targetEntity.GetKindID() == topoapi.E2CELL {
			cellObject := &topoapi.E2Cell{}
			err = obj.GetAspect(cellObject)
			if err == nil {
				cells = append(cells, cellObject)
			}
		}
	}

	if len(cells) == 0 {
		return nil, errors.New(errors.NotFound, "there is no cell to subscribe for e2 node %s", nodeID)
	}

	return cells, nil
}

func getControlRelationFilter() *topoapi.Filters {
	controlRelationFilter := &topoapi.Filters{
		KindFilter: &topoapi.Filter{
			Filter: &topoapi.Filter_Equal_{
				Equal_: &topoapi.EqualFilter{
					Value: topoapi.CONTROLS,
				},
			},
		},
	}
	return controlRelationFilter
}

// WatchE2Connections watch e2 node connection changes
func (c *Client) WatchE2Connections(ctx context.Context, ch chan topoapi.Event) error {
	err := c.client.Watch(ctx, ch, toposdk.WithWatchFilters(getControlRelationFilter()))
	if err != nil {
		return err
	}
	return nil
}

var _ TopoClient = &Client{}
