/*
 * File: cws_test.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/codeworks-tw/cwsutil/cwsbase"
)

func TestUtil(t *testing.T) {
	fmt.Println("\n================ Testing functions ================")

	cwsbase.InitBasicLocalizationData()
	fmt.Println("Test Location Message 200: ", cwsbase.GetLocalizationMessage("200"))

	fmt.Println("\nTest general encryption...")
	m := map[string]any{
		"a": "b",
	}
	fmt.Println("Before encryption: ", m)

	os.Setenv("CRYPTO_KEY_HEX", "4d8f1e227bc2115d1008e98965abd753a420dd0d27a2ee66c284606981867ee0")
	os.Setenv("CRYPTO_IV_HEX", "fdbfcd1c11e7ec1a2d7073e0f45b39c4")
	s, err := cwsbase.EncryptMap(m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After encryption: ", s)

	d, err := cwsbase.DecryptToMap(s)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After decryption: ", d)

	fmt.Println("================ Testing functions end ================")
}
