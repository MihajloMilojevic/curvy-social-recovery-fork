package keyrecovery

import (
	"fmt"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
	"testing"
)

func TestSplitAndRecover(t *testing.T) {
	spendingKeyStr, viewingKeyStr, err := setupRandomKeys()
	if err != nil {
		t.Fatal(err)
	}

	n := 20
	threshold := 14

	t.Logf("Spending key: %s\n", spendingKeyStr)
	t.Logf("Viewing key: %s\n", viewingKeyStr)

	shares, err := Split(threshold, n, spendingKeyStr, viewingKeyStr)
	if err != nil {
		t.Fatalf("Failed to construct the shares: %s\n", err)
	}

	t.Logf("Shares: %v\n", shares)

	t.Run("xi != 0 for all i", func(t *testing.T) {
		for _, share := range shares {
			if share.Point == "0" {
				t.Fatal("One share evaluation point is zero")
			}
		}
	})

	t.Run("xi != xj for i != j", func(t *testing.T) {
		for i, share := range shares {
			for j := i + 1; j < len(shares); j++ {
				if share.Point == shares[j].Point {
					t.Fatal("Two shares have the same evaluation point")
				}
			}
		}
	})

	t.Run("given shares = t", func(t *testing.T) {
		newSpendingKeyStr, newViewingKeyStr, err := Recover(threshold, shares[0:threshold])
		t.Logf("Recovered spending key: %s\n", newSpendingKeyStr)
		t.Logf("Recovered viewing key: %s\n", newViewingKeyStr)
		if err != nil {
			t.Fatalf("Failed to recover the shares: %s", err)
		}
		if newSpendingKeyStr != spendingKeyStr {
			t.Error("Spending keys do not match")
		}
		if newViewingKeyStr != viewingKeyStr {
			t.Error("Viewing keys do not match")
		}
	})

	t.Run("given shares = t + 1 > t", func(t *testing.T) {
		newSpendingKeyStr, newViewingKeyStr, err := Recover(threshold, shares[0:threshold+1])
		t.Logf("Recovered spending key: %s\n", newSpendingKeyStr)
		t.Logf("Recovered viewing key: %s\n", newViewingKeyStr)
		if err != nil {
			t.Fatalf("Failed to recover the shares: %s\n", err)
		}
		if newSpendingKeyStr != spendingKeyStr {
			t.Error("Spending keys do not match")
		}
		if newViewingKeyStr != viewingKeyStr {
			t.Error("Viewing keys do not match")
		}
	})

	t.Run("given shares = t - 1 < t", func(t *testing.T) {
		_, _, err := Recover(threshold, shares[0:threshold-1])
		if err == nil {
			t.Error("Keys should not have been reconstructed")
		}
	})
}

func BenchmarkSplit(b *testing.B) {
	spendingKeyStr, viewingKeyStr, err := setupRandomKeys()
	if err != nil {
		b.Fatal(err)
	}

	n := 20
	threshold := 14

	b.ResetTimer()
	for range b.N {
		Split(threshold, n, spendingKeyStr, viewingKeyStr)
	}
}

func BenchmarkRecover(b *testing.B) {
	spendingKeyStr, viewingKeyStr, err := setupRandomKeys()
	if err != nil {
		b.Fatal(err)
	}

	n := 20
	threshold := 14

	shares, err := Split(threshold, n, spendingKeyStr, viewingKeyStr)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for range b.N {
		Recover(threshold, shares)
	}
}

func setupRandomKeys() (string, string, error) {
	var spendingKey SECP256K1_fr.Element
	var viewingKey BN254_fr.Element

	_, err := spendingKey.SetRandom()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random spending key: %v\n", err)
	}
	_, err = viewingKey.SetRandom()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random viewing key: %v\n", err)
	}

	spendingKeyStr := spendingKey.Text(16)
	viewingKeyStr := viewingKey.Text(16)

	return spendingKeyStr, viewingKeyStr, nil
}
