// SPDX-License-Identifier: Apache-2.0

package cli

import mfxsdk "github.com/cloustone/pandas/sdk/go"

// Keep SDK handle in global var
var sdk mfxsdk.SDK

// SetSDK sets mainflux SDK instance.
func SetSDK(s mfxsdk.SDK) {
	sdk = s
}
