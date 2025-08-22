package main

import (
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	// Example: Creating a v0 transaction with durable nonce and address lookup tables

	// Initialize RPC client
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	// Define accounts
	feePayer := solana.MustPublicKeyFromBase58("YourFeePayerPublicKeyHere")
	nonceAccount := solana.MustPublicKeyFromBase58("YourNonceAccountHere")
	nonceAuthority := solana.MustPublicKeyFromBase58("YourNonceAuthorityHere")

	// Get nonce value from on-chain nonce account
	nonceAccountInfo, err := rpcClient.GetAccountInfo(nonceAccount)
	if err != nil {
		log.Fatal("Failed to get nonce account info:", err)
	}

	// Parse nonce value (simplified - in practice you'd parse the actual nonce data)
	nonceValue := solana.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	// Define address lookup table
	addressTable := map[solana.PublicKey]solana.PublicKeySlice{
		solana.MustPublicKeyFromBase58("YourLUTAccountHere"): {
			solana.MustPublicKeyFromBase58("Address1InLUT"),
			solana.MustPublicKeyFromBase58("Address2InLUT"),
			solana.MustPublicKeyFromBase58("Address3InLUT"),
		},
	}

	// Create a simple instruction (example: system transfer)
	transferInstruction := system.NewTransferInstruction(
		1000000, // lamports
		feePayer,
		solana.MustPublicKeyFromBase58("RecipientAddressHere"),
	).Build()

	// Method 1: Using NewVersionedTransactionWithLUTAndNonce convenience function
	vtx1, err := solana.NewVersionedTransactionWithLUTAndNonce(
		[]solana.Instruction{transferInstruction},
		nonceValue,
		nonceAccount,
		nonceAuthority,
		addressTable,
	)
	if err != nil {
		log.Fatal("Failed to create v0 transaction:", err)
	}

	fmt.Printf("Created v0 transaction with LUT and nonce: %s\n", vtx1.Message.GetVersion())

	// Method 2: Using the builder pattern
	builder := solana.NewVersionedTransactionBuilder().
		SetFeePayer(feePayer).
		SetAddressTables(addressTable).
		SetDurableNonce(nonceAccount, nonceAuthority).
		SetNonceValue(nonceValue).
		AddInstruction(transferInstruction)

	vtx2, err := builder.Build()
	if err != nil {
		log.Fatal("Failed to build v0 transaction:", err)
	}

	fmt.Printf("Built v0 transaction: %s\n", vtx2.Message.GetVersion())

	// Method 3: Step by step construction
	vtx3, err := solana.NewVersionedTransaction(
		[]solana.Instruction{transferInstruction},
		solana.Hash{}, // will be replaced by nonce value
		solana.TransactionDurableNonce(nonceAccount, nonceAuthority),
		solana.TransactionWithNonceValue(nonceValue),
		solana.TransactionAddressTables(addressTable),
	)
	if err != nil {
		log.Fatal("Failed to create step-by-step v0 transaction:", err)
	}

	fmt.Printf("Step-by-step v0 transaction: %s\n", vtx3.Message.GetVersion())

	// All three methods produce equivalent v0 transactions with:
	// - MessageV0 format
	// - Address lookup tables support
	// - Durable nonce (AdvanceNonceAccount instruction first)
	// - Nonce value as recentBlockhash
}
