package policy_test

import (
	"capact.io/capact/internal/ptr"
	"capact.io/capact/pkg/engine/k8s/policy"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
)

func fixPolicyWithTypeRef() policy.Policy {
	return policy.Policy{
		Rules: policy.RulesList{
			{
				Interface: types.ManifestRefWithOptRevision{
					Path: "cap.*",
				},
				OneOf: []policy.Rule{
					{
						ImplementationConstraints: policy.ImplementationConstraints{},
						Inject: &policy.InjectData{
							RequiredTypeInstances: []policy.RequiredTypeInstanceToInject{
								{
									RequiredTypeInstanceReference: policy.RequiredTypeInstanceReference{
										ID: "id",
									},
									TypeRef: &types.ManifestRef{
										Path:     "cap.type.sample",
										Revision: "0.1.0",
									},
								},
								{
									RequiredTypeInstanceReference: policy.RequiredTypeInstanceReference{
										ID:          "id2",
										Description: ptr.String("ID 2"),
									},
									TypeRef: &types.ManifestRef{
										Path:     "cap.type.sample2",
										Revision: "0.2.0",
									},
								},
							},
							AdditionalTypeInstances: []policy.AdditionalTypeInstanceToInject{
								{
									AdditionalTypeInstanceReference: policy.AdditionalTypeInstanceReference{
										ID:   "id3",
										Name: "name",
									},
									TypeRef: &types.ManifestRef{
										Path:     "cap.type.sample3",
										Revision: "0.3.0",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
