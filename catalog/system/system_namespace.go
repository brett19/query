//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package system

import (
	"github.com/couchbaselabs/query/catalog"
	"github.com/couchbaselabs/query/errors"
)

type namespace struct {
	site      *site
	id        string
	name      string
	keyspaces map[string]catalog.Keyspace
}

func (p *namespace) SiteId() string {
	return p.site.Id()
}

func (p *namespace) Id() string {
	return p.id
}

func (p *namespace) Name() string {
	return p.name
}

func (p *namespace) KeyspaceIds() ([]string, errors.Error) {
	return p.KeyspaceNames()
}

func (p *namespace) KeyspaceNames() ([]string, errors.Error) {
	rv := make([]string, len(p.keyspaces))
	i := 0
	for k, _ := range p.keyspaces {
		rv[i] = k
		i = i + 1
	}
	return rv, nil
}

func (p *namespace) KeyspaceById(id string) (catalog.Keyspace, errors.Error) {
	return p.KeyspaceByName(id)
}

func (p *namespace) KeyspaceByName(name string) (catalog.Keyspace, errors.Error) {
	b, ok := p.keyspaces[name]
	if !ok {
		return nil, errors.NewError(nil, "Keyspace "+name+" not found.")
	}

	return b, nil
}

// newNamespace creates a new namespace.
func newNamespace(s *site) (*namespace, errors.Error) {
	p := new(namespace)
	p.site = s
	p.id = NAMESPACE_ID
	p.name = NAMESPACE_NAME
	p.keyspaces = make(map[string]catalog.Keyspace)

	e := p.loadKeyspaces()
	if e != nil {
		return nil, e
	}
	return p, nil
}

func (p *namespace) loadKeyspaces() (e errors.Error) {

	sb, e := newSitesKeyspace(p)
	if e != nil {
		return e
	}
	p.keyspaces[sb.Name()] = sb

	pb, e := newNamespacesKeyspace(p)
	if e != nil {
		return e
	}
	p.keyspaces[pb.Name()] = pb

	bb, e := newKeyspacesKeyspace(p)
	if e != nil {
		return e
	}
	p.keyspaces[bb.Name()] = bb

	db, e := newDualKeyspace(p)
	if e != nil {
		return e
	}
	p.keyspaces[db.Name()] = db

	ib, e := newIndexesKeyspace(p)
	if e != nil {
		return e
	}
	p.keyspaces[ib.Name()] = ib

	return nil
}
