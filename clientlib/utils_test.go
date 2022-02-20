package clientlib

import "testing"

func TestCheckLimits(t *testing.T) {
	if CheckLimits(zeroBytes(10), 0, 100) == "" || CheckLimits(zeroBytes(10), (100*2*24*60*60), 100) == "" {
		t.Fatal("CheckLimits should fail with abnormal data")
	}
	if CheckLimits(zeroBytes(10), 0, 100) == "" || CheckLimits(zeroBytes(10), (100*2*24*60*60), 100) == "" {
		t.Fatal("CheckLimits should fail with abnormal ttl")
	}
	if CheckLimits(zeroBytes(10), 1, 100) != "" || CheckLimits(zeroBytes(10), (2*24*60*60), 100) != "" || CheckLimits(zeroBytes(10), 500, 100) != "" {
		t.Fatal("CheckLimits should not fail with normal ttl")
	}
	if CheckLimits(zeroBytes(10), 100, 0) == "" || CheckLimits(zeroBytes(10), 100, 100*1000000) == "" {
		t.Fatal("CheckLimits should fail with abnormal hits")
	}
	if CheckLimits(zeroBytes(10), 100, 1) != "" || CheckLimits(zeroBytes(10), 100, 1000000) != "" || CheckLimits(zeroBytes(10), 100, 500) != "" {
		t.Fatal("CheckLimits should not fail with normal hits")
	}
	if CheckLimits(zeroBytes(0), 100, 100) == "" || CheckLimits(zeroBytes(5*1024*1024+1), 100, 100) == "" {
		t.Fatal("CheckLimits should fail with abnormal data")
	}
	if CheckLimits(zeroBytes(1), 100, 100) != "" || CheckLimits(zeroBytes(5*1024*1024), 100, 100) != "" || CheckLimits(zeroBytes(500), 100, 100) != "" {
		t.Fatal("CheckLimits should not fail with normal data")
	}
}

func zeroBytes(n uint) []byte {
	outBytes := []byte{}
	for i := uint(0); i < n; i++ {
		outBytes = append(outBytes, byte(0))
	}
	return outBytes
}
