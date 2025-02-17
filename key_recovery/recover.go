package keyrecovery

import (
	"errors"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func Recover(threshold int, shares []Share) (string, string, error) {

	// Check if the number of shares is less than the threshold
	if len(shares) < threshold {
		return "", "", errors.New("number of shares is less than the threshold")
	}

	// Truncate the shares to the threshold value
	// there is no need to use more than t=threshold shares for the interpolation
	shares = shares[0:threshold]

	var points = make([]point, threshold)

	// Transform the shares from string to field element before interpolating
	for i := range points {
		_, errX1 := points[i].xS.SetString("0x" + shares[i].Point)
		_, errX2 := points[i].xV.SetString("0x" + shares[i].Point)
		_, errY1 := points[i].yS.SetString("0x" + shares[i].SpendingEval)
		_, errY2 := points[i].yV.SetString("0x" + shares[i].ViewingEval)

		if errX1 != nil {
			return "", "", errX1
		}
		if errX2 != nil {
			return "", "", errX2
		}
		if errY1 != nil {
			return "", "", errY1
		}
		if errY2 != nil {
			return "", "", errY2
		}
	}

	// Define the spending key and viewing key as field elements for calculations
	var sk SECP256K1_fr.Element
	var vk BN254_fr.Element

	// Initialize the spending key and viewing key as zero
	sk.SetZero()
	vk.SetZero()

	// Recover the spending key and viewing key from the shares
	//  		n             n
	// Result = Σ  y_i   *    Π ( x_j / (x_j - x_i) )
	// 		   i=1           j=1
	// 						j ≠ i

	for i := range points {
		// Current term = y_i
		var lsk SECP256K1_fr.Element
		var lvk BN254_fr.Element
		lsk.Set(&points[i].yS)
		lvk.Set(&points[i].yV)

		// Denominators for the Lagrange basis Li
		denominatorS := SECP256K1_fr.One()
		denominatorV := BN254_fr.One()

		// Apply Li to the current term
		for j := range points {
			if i == j {
				continue
			}

			// current term * xj
			lsk.Mul(&lsk, &points[j].xS)
			lvk.Mul(&lvk, &points[j].xV)

			// xj - xi
			var diffS SECP256K1_fr.Element
			var diffV BN254_fr.Element
			diffS.Sub(&points[j].xS, &points[i].xS)
			diffV.Sub(&points[j].xV, &points[i].xV)

			// denominator *= (xj - xi)
			denominatorS.Mul(&denominatorS, &diffS)
			denominatorV.Mul(&denominatorV, &diffV)
		}

		// Divide by the denominator
		lsk.Div(&lsk, &denominatorS)
		lvk.Div(&lvk, &denominatorV)

		// Add the current term to the spending key and viewing key
		sk.Add(&sk, &lsk)
		vk.Add(&vk, &lvk)
	}

	return sk.Text(16), vk.Text(16), nil
}

type point struct {
	xS SECP256K1_fr.Element
	xV BN254_fr.Element
	yS SECP256K1_fr.Element
	yV BN254_fr.Element
}
