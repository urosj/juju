// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api

// TODO(ericsnow) Eliminate the dependence on apiserver if possible.

import (
	"strings"

	"github.com/juju/errors"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"
	"gopkg.in/juju/names.v2"
	"gopkg.in/macaroon.v1"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/charmstore"
)

// NewListResourcesArgs returns the arguments for the ListResources endpoint.
func NewListResourcesArgs(services []string) (params.ListResourcesArgs, error) {
	var args params.ListResourcesArgs
	var errs []error
	for _, service := range services {
		if !names.IsValidApplication(service) {
			err := errors.Errorf("invalid application %q", service)
			errs = append(errs, err)
			continue
		}
		args.Entities = append(args.Entities, params.Entity{
			Tag: names.NewApplicationTag(service).String(),
		})
	}
	if err := resolveErrors(errs); err != nil {
		return args, errors.Trace(err)
	}
	return args, nil
}

// NewAddPendingResourcesArgs returns the arguments for the
// AddPendingResources API endpoint.
func NewAddPendingResourcesArgs(applicationID string, chID charmstore.CharmID, csMac *macaroon.Macaroon, resources []charmresource.Resource) (params.AddPendingResourcesArgs, error) {
	var args params.AddPendingResourcesArgs

	if !names.IsValidApplication(applicationID) {
		return args, errors.Errorf("invalid application %q", applicationID)
	}
	tag := names.NewApplicationTag(applicationID).String()

	var apiResources []params.CharmResource
	for _, res := range resources {
		if err := res.Validate(); err != nil {
			return args, errors.Trace(err)
		}
		apiRes := CharmResource2API(res)
		apiResources = append(apiResources, apiRes)
	}
	args.Tag = tag
	args.Resources = apiResources
	if chID.URL != nil {
		args.URL = chID.URL.String()
		args.Channel = string(chID.Channel)
		args.CharmStoreMacaroon = csMac
	}
	return args, nil
}

// XXX
func resolveErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		msgs := make([]string, len(errs))
		for i, err := range errs {
			msgs[i] = err.Error()
		}
		return errors.New(strings.Join(msgs, "\n"))
	}
}
