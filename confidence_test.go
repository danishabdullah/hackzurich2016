package main

import (
    "testing"
    "math"
)

func TestBinomCDF(t *testing.T) {
	t.Log("Testing the binomial cumulative distribution function with some value triples...")
	
	check(t,12,20,.5,0.868412017822)
	check(t,7,20,.5,1-0.868412017822)
	check(t,5,500,.01,0.615962131821)
	check(t,494,500,.99,1-0.615962131821)
	check(t,0,2,.2,0.64)
	check(t,0,200,.2,4.14951556888e-20)
}

func check(t *testing.T, k float64, n float64, p float64, expected float64) {
	const TOL = 1e-9
	actual := binomCDF(k,n,p)
	if math.Abs(actual-expected) > TOL {
		t.Errorf("binomCDF failed for k=%f,n=%f,p=%f. Expected: %.12g, Actual: %.12g",k,n,p,expected,actual)
	}
}