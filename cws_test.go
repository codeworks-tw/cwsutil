/*
 * File: cws_test.go
 * Created Date: Friday, January 26th 2024, 11:54:08 am
 *
 * Last Modified: Sat Jan 27 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"cwsutil/baseutil"
	"fmt"
	"os"
	"testing"
)

func TestUtil(t *testing.T) {
	fmt.Println("\n================ Testing functions ================")

	baseutil.InitBasicLocalizationData()
	fmt.Println("Test Location Message 200: ", baseutil.GetLocalizationMessage("200"))

	fmt.Println("\nTest general encryption...")
	m := map[string]any{
		"a": "b",
	}
	fmt.Println("Before encryption: ", m)

	os.Setenv("CRYPTO_KEY_HEX", "4d8f1e227bc2115d1008e98965abd753a420dd0d27a2ee66c284606981867ee0")
	os.Setenv("CRYPTO_IV_HEX", "fdbfcd1c11e7ec1a2d7073e0f45b39c4")
	s, err := baseutil.EncryptMap(m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After encryption: ", s)

	d, err := baseutil.DecryptToMap(s)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After decryption: ", d)

	fmt.Println("================ Testing functions end ================")
}
