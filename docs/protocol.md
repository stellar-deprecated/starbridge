# Starbridge Protocol

## Participants
- W1, W2 = Wallets
- A = Asset (Local asset)
- wA = Wrapped asset (Foreign asset)
- Vn = Validators
- T, T1, T2 = Transaction or Message native to the chain it is generated
- B = Bridge account on Stellar, or contract on Ethereum.
- N = Next sequence number for the source.

## Actions

### 1. Action: Sending a local asset that may be used as a foreign asset

#### i. Transfer
- W1 pays A to Bs, annotated with dest=W2.
- W2 requests Vn to sign T.
- Vn observes, produces T, signs T, and broadcasts T.
- W2 observes T with signatures.
- W2 submits T (or calls B with T).
- B mints wA, and pays wA to W2.

### ii. Tx Fails in Anyway
- W1 pays A to Bs, annotated with dest=W2.
- W2 requests Vn to sign T1.
- Vn observes, produces T1, signs T1, and broadcasts T1.
- W2 observes T1 with signatures.
- W2 submits T1 (or calls B with T1).
- T1 fails.
- W2 waits for T1 ledger bounds to expire.
- W2 requests Vn to sign T2.
- Vn observes, produces T2, signs T2, and broadcasts T2.
- W2 observes T2 with signatures.
- W2 submits T2 (or calls B with T2).
- B mints wA, and pays wA to W2.

### iii. Reversal/Refund
- W1 pays A to Bs, annotated with dest=W2.
- W1 requests Vn to sign T.
- Vn observes, produces T, signs T, and broadcasts T.
- W1 observes T with signatures.
- W1 submits T (or calls B with T).
- B unlocks A, and pays A to W1.

## 2. Action: Returning a foreign asset that may be used as a local asset

### i. Transfer
- W1 pays wA to Bs, annotated with dest=W2.
- W2 requests Vn to sign T.
- Vn observes, produces T, signs T, and broadcasts T.
- W2 observes T with signatures.
- W2 submits T (or calls B with T).
- B unlocks A, and pays A to W2.

### ii. Tx Fails in Anyway
- W1 pays wA to Bs, annotated with dest=W2.
- W2 requests Vn to sign T1.
- Vn observes, produces T1, signs T1, and broadcasts T1.
- W2 observes T1 with signatures.
- W2 submits T1 (or calls B with T1).
- T1 fails.
- W2 waits for T1 ledger bounds to expire.
- W2 requests Vn to sign T2.
- Vn observes, produces T2, signs T2, and broadcasts T2.
- W2 observes T2 with signatures.
- W2 submits T2 (or calls B with T2).
- B unlocks A, and pays A to W2.

### iii. Reversal/Refund
- W1 pays wA to Bs, annotated with dest=W2.
- W1 requests Vn to sign T.
- Vn observes, produces T, signs T, and broadcasts T.
- W1 observes T with signatures.
- W1 submits T (or calls B with T).
- B mints wA, and pays wA to W1.
