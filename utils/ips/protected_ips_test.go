package ips

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidAddressesProtocol(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value string
		valid bool
	}{
		{valid: true, value: "test.domain.test:443"},
		{valid: true, value: "test.test:443"},
		{valid: true, value: "g.g:443"},
		{valid: true, value: "super.long.example.of.sub.domains.to.test.with.and.see.if.it.works.g:443"},
		{valid: false, value: "g..g:443"},
		{valid: false, value: "test..domain.test:443"},
		{valid: false, value: "testtest:443"},
		{valid: false, value: "fe80::215:6dff:fec4:b31d"}, // local ipv6
		{valid: false, value: "http://test.domain.test:443"},
		{valid: false, value: ":"},
		{valid: false, value: "http://::"},
		{valid: false, value: "::"},
		{valid: false, value: "http://localhost:1211"},
		{valid: false, value: "http://127.0.0.1:1211"},
		{valid: false, value: "127.0.0.1:3123"},
		{valid: false, value: "127.0.0.126:3123"},
		{valid: false, value: "http://127.0.0.126:3123"},
		{valid: false, value: "localhost:1211"},
		{valid: false, value: "0.0.0.0:1211"},
		{valid: false, value: "http://0.0.0.0:1211"},
		{valid: false, value: "http://0.0.0.0:1211"},
	}

	for _, tt := range tests {
		t.Run("Test Name: "+tt.value, func(t *testing.T) {
			value := IsValidNetworkAddressProtocol(tt.value)
			require.Equal(t, tt.valid, value)

			value = IsValidNetworkAddressConsensus(tt.value)
			require.Equal(t, tt.valid, value)
		})
	}
}