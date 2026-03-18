//go:build integ
// +build integ

/*
 * Copyright The Kmesh Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kmesh

import (
	"testing"

	"istio.io/istio/pkg/test/echo/common/scheme"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/check"
)

// TestBoundaryTraffic verifies traffic between applications inside and outside the Kmesh mesh.
func TestBoundaryTraffic(t *testing.T) {
	framework.NewTest(t).
		Run(func(t framework.TestContext) {
			// Case 1: In-Mesh Client -> Out-of-Mesh Server
			// Verify that a meshed application can successfully access an unmeshed application.
			t.NewSubTest("In-Mesh to Out-of-Mesh").Run(func(t framework.TestContext) {
				for _, src := range apps.EnrolledToKmesh {
					for _, dst := range apps.Unmeshed {
						t.NewSubTestf("from %v to %v", src.Config().Service, dst.Config().Service).Run(func(t framework.TestContext) {
							src.CallOrFail(t, echo.CallOptions{
								To:     dst,
								Port:   echo.Port{Name: "http"},
								Scheme: scheme.HTTP,
								Count:  1,
								Check:  check.OK(),
							})
						})
					}
				}
			})

			// Case 2: Out-of-Mesh Client -> In-Mesh Server
			// Verify that an unmeshed application can successfully access a meshed application.
			t.NewSubTest("Out-of-Mesh to In-Mesh").Run(func(t framework.TestContext) {
				for _, src := range apps.Unmeshed {
					for _, dst := range apps.EnrolledToKmesh {
						t.NewSubTestf("from %v to %v", src.Config().Service, dst.Config().Service).Run(func(t framework.TestContext) {
							src.CallOrFail(t, echo.CallOptions{
								To:     dst,
								Port:   echo.Port{Name: "http"},
								Scheme: scheme.HTTP,
								Count:  1,
								Check:  check.OK(),
							})
						})
					}
				}
			})
		})
}
