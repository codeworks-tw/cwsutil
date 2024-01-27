/*
 * File: func_test.go
 * Created Date: Friday, January 26th 2024, 11:54:08 am
 *
 * Last Modified: Fri Jan 26 2024
 * Modified By: Howard Ling-Hao Kung
 */

package cwsutil

import (
	"cws"
	"fmt"
	"os"
	"testing"
)

func TestUtil(t *testing.T) {
	fmt.Println("\n================ Testing functions ================")

	cws.InitBasicLocalizationData()
	fmt.Println("Test Location Message 200: ", cws.GetLocalizationMessage("200"))

	fmt.Println("\nTest general encryption...")
	m := map[string]any{
		"a": "b",
	}
	fmt.Println("Before encryption: ", m)

	os.Setenv("CRYPTO_KEY_HEX", "4d8f1e227bc2115d1008e98965abd753a420dd0d27a2ee66c284606981867ee0")
	os.Setenv("CRYPTO_IV_HEX", "fdbfcd1c11e7ec1a2d7073e0f45b39c4")
	s, err := cws.EncryptMap(m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After encryption: ", s)

	d, err := cws.DecryptToMap(s)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("After decryption: ", d)

	fmt.Println("================ Testing functions passed ================")
}
