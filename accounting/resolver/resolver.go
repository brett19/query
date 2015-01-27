//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package resolver

import (
	"strings"

	"github.com/couchbaselabs/query/accounting"
	"github.com/couchbaselabs/query/accounting/gometrics"
	"github.com/couchbaselabs/query/accounting/stub"

	"github.com/couchbaselabs/query/errors"
)

func NewAcctstore(uri string) (accounting.AccountingStore, errors.Error) {
	if strings.HasPrefix(uri, "stub:") {
		return accounting_stub.NewAccountingStore(uri)
	}

	if strings.HasPrefix(uri, "gometrics:") {
		return accounting_gm.NewAccountingStore(), nil
	}

	return nil, errors.NewAdminInvalidURL("AccountingStore", uri)
}
