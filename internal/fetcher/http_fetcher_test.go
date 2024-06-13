package fetcher

import (
	"context"
	"testing"
)

func TestHttp_Fetch(t *testing.T) {
	type args struct {
		ipfsID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success get ipfs data",
			args: args{
				ipfsID: "bafkreicakci5bdl6pf2ujvxlm7fadu64lgf56phhy5bujyouwqdacjlc5m",
			},
			want: `{"address":"0x17D717EB3DD20a202dce9e8e396A444db1Af1D28","sig":"0x364876b7405e436d3694872d14a54f8f5bd3c33222a1ade68d0491ef5232ce04692dada19d37db54370a28bad2a6c2dc9a6f2fd4da8115e51f6705d031393ae11c","hash":"0xd93d012742b7ff815839d2ed54eae04f87bdf3b270014010b4b9f2145f5054fb","data":{"domain":{"name":"snapshot","version":"0.1.4"},"types":{"CancelProposal":[{"name":"from","type":"address"},{"name":"space","type":"string"},{"name":"timestamp","type":"uint64"},{"name":"proposal","type":"bytes32"}]},"message":{"space":"dfkvote.eth","proposal":"0xb04b92e01daca4fb42cd4bf681e8a6bac97963aa0bae226051a6bdabc89185b8","from":"0x17D717EB3DD20a202dce9e8e396A444db1Af1D28","timestamp":1715839116}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Http{}
			got, err := f.Fetch(context.TODO(), tt.args.ipfsID)
			if err != nil {
				t.Errorf("Fetch() error = %v", err)

				return
			}

			if string(got) != tt.want {
				t.Errorf("Fetch() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
