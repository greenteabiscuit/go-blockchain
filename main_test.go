package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBlockchain_LastBlock(t *testing.T) {
	type fields struct {
		Chain               []*Block
		CurrentTransactions []*Transaction
	}
	tests := []struct {
		name   string
		fields fields
		want   *Block
	}{
		{
			name: "only 1 block: success",
			fields: fields{
				Chain: []*Block{
					{
						Index:        1,
						Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
						Transactions: nil,
						Proof:        1,
						PreviousHash: "exampleHash",
					},
				},
				CurrentTransactions: nil,
			},
			want: &Block{
				Index:        1,
				Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
				Transactions: nil,
				Proof:        1,
				PreviousHash: "exampleHash",
			},
		},
		{
			name: "only 1 block: success",
			fields: fields{
				Chain: []*Block{
					{
						Index:     1,
						Timestamp: time.Now().Unix(), // TODO: this test needs fixing
						Transactions: []*Transaction{
							{
								Sender:    "exampleSender",
								Recipient: "exampleRecipient",
								Amount:    100,
							},
						},
						Proof:        1,
						PreviousHash: "exampleHash",
					},
				},
				CurrentTransactions: nil,
			},
			want: &Block{
				Index:     1,
				Timestamp: time.Now().Unix(), // TODO: this test needs fixing
				Transactions: []*Transaction{
					{
						Sender:    "exampleSender",
						Recipient: "exampleRecipient",
						Amount:    100,
					},
				},
				Proof:        1,
				PreviousHash: "exampleHash",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blockchain{
				Chain:               tt.fields.Chain,
				CurrentTransactions: tt.fields.CurrentTransactions,
			}
			if got := b.LastBlock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockchain_NewBlock(t *testing.T) {
	type fields struct {
		Chain               []*Block
		CurrentTransactions []*Transaction
	}
	type args struct {
		proof        int
		previousHash string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Block
	}{
		{
			name: "Creating first block: success",
			fields: fields{
				Chain:               nil,
				CurrentTransactions: nil,
			},
			args: args{
				proof:        1,
				previousHash: "exampleHash",
			},
			want: &Block{
				Index:        1,
				Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
				Transactions: nil,
				Proof:        1,
				PreviousHash: "exampleHash",
			},
		},
		{
			name: "Creating second block: success",
			fields: fields{
				Chain: []*Block{
					{
						Index:        1,
						Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
						Transactions: nil,
						Proof:        1,
						PreviousHash: "exampleHash",
					},
				},
				CurrentTransactions: nil,
			}, args: args{
				proof:        1,
				previousHash: "exampleHash",
			},
			want: &Block{
				Index:        2,
				Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
				Transactions: nil,
				Proof:        1,
				PreviousHash: "exampleHash",
			},
		},
		{
			name: "Creating first block with transactions: success",
			fields: fields{
				Chain: nil,
				CurrentTransactions: []*Transaction{
					{
						Sender:    "exampleSender",
						Recipient: "exampleRecipient",
						Amount:    100,
					},
				},
			},
			args: args{
				proof:        1,
				previousHash: "exampleHash",
			},
			want: &Block{
				Index:     1,
				Timestamp: time.Now().Unix(), // TODO: this test needs fixing
				Transactions: []*Transaction{
					{
						Sender:    "exampleSender",
						Recipient: "exampleRecipient",
						Amount:    100,
					},
				},
				Proof:        1,
				PreviousHash: "exampleHash",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blockchain{
				Chain:               tt.fields.Chain,
				CurrentTransactions: tt.fields.CurrentTransactions,
			}
			if got := b.NewBlock(tt.args.proof, tt.args.previousHash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockchain_NewTransaction(t *testing.T) {
	type fields struct {
		Chain               []*Block
		CurrentTransactions []*Transaction
	}
	type args struct {
		sender    string
		recipient string
		amount    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "adding first transaction: success",
			fields: fields{
				Chain: []*Block{
					{
						Index:        1,
						Timestamp:    time.Now().Unix(), // TODO: this test needs fixing
						Transactions: nil,
						Proof:        1,
						PreviousHash: "exampleHash",
					},
				},
				CurrentTransactions: nil,
			},
			args: args{
				sender:    "exampleSender",
				recipient: "exampleRecipient",
				amount:    100,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blockchain{
				Chain:               tt.fields.Chain,
				CurrentTransactions: tt.fields.CurrentTransactions,
			}
			if got := b.NewTransaction(tt.args.sender, tt.args.recipient, tt.args.amount); got != tt.want {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		block Block
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "success, when block is nil",
			args: args{block: Block{
				Index:        0,
				Timestamp:    0,
				Transactions: nil,
				Proof:        0,
				PreviousHash: "",
			}},
			want: []byte{49, 152, 133, 160, 103, 111, 24, 28, 104, 155, 146, 98, 31, 211, 70, 192, 119, 130, 105, 216, 54, 41, 249, 238, 41, 32, 11, 199, 247, 122, 146, 107},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
