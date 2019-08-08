// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package kms contains samples for asymmetric keys feature of Cloud Key Management Service
// https://cloud.google.com/kms/
package kms

// [START kms_remove_member_from_keyring_policy]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// removeMemberRingPolicy removes a specified member from an IAM role for the key ring
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func removeMemberRingPolicy(w io.Writer, keyRingName, member string, role iam.RoleName) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Get the KeyRing.
	keyRingObj, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: keyRingName})
	if err != nil {
		return err
	}
	// Get IAM Policy.
	handle := client.KeyRingIAM(keyRingObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return err
	}

	// Remove Member.
	policy.Remove(member, role)
	err = handle.SetPolicy(ctx, policy)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "Removed member from keyring policy.")
	return nil
}

// [END kms_remove_member_from_keyring_policy]
