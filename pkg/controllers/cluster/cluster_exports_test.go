// Copyright 2024-2026 SAP SE or an SAP affiliate company and Greenhouse contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	headscalev1 "github.com/juanfont/headscale/gen/go/headscale/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ExportServiceAccountName        = serviceAccountName
	ExportTailscaleAuthorizationKey = tailscaleAuthorizationKey
)

func ExportSetHeadscaleGRPCClientOnHAR(r *HeadscaleAccessReconciler, c headscalev1.HeadscaleServiceClient) {
	r.headscaleGRPCClient = c
}

func ExportSetRestClientGetterFunc(r *HeadscaleAccessReconciler, f func(restClientGetter genericclioptions.RESTClientGetter, proxy string, headscaleAddress string) (client.Client, error)) {
	r.getHeadscaleClientFromRestClientGetter = f
}
