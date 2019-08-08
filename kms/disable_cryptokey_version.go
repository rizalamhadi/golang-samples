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

import (
	"context"
	"fmt"
	"io"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

// [START kms_disable_cryptokey_version]

// disableCryptoKeyVersion disables a specified key version on KMS.
// example keyVersionName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func disableCryptoKeyVersion(w io.Writer, keyVersionName string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.UpdateCryptoKeyVersionRequest{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:  keyVersionName,
			State: kmspb.CryptoKeyVersion_DISABLED,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"state"},
		},
	}
	// Call the API.
	result, err := client.UpdateCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Disabled crypto key version: %s", result)
	return nil
}

// [END kms_disable_cryptokey_version]
