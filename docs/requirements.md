# Starbridge Requirements

Build a transfer bridge that facilitates the transfer of value from Stellar to Ethereum, and from Ethereum to Stellar. When assets are transferred on either chain they should become unusable while they are usable on the other chain.

The bridge has the following properties:
- Value transfer bridge
- Symmetrically bidirectional
    - Assets from Stellar may be transferred and usable on Ethereum.
    - Assets from Ethereum may be transferred and usable on Stellar.
- Throughput and Latency
    - Constrained only by the throughput and latency of the chains.
    - No more than 5s latency introduced by validators.
- Users
    - Transfers are facilitated by a user sending and a user receiving.
    - The sender and receiver are the same user, but not necessarily the same wallet.
    - The success and timing of transfers from one user are unaffected by the transfer of another user.
    - The sender and receiver must not be required to be online at the same time.
    - The receiver must not be required to be online at a specific time.
    - The sender does not require a method to cancel or reverse a transfer.
- Security and Control
    - Decentralized closed participation of N validators.
    - Assets locked in escrow accounts / contracts controlled by m-of-n signer configurations.
    - Assets not vulnerable to fraud where 0 to (n – m) validators are nefarious or compromised.
    - M-of-n signer configuration configurable through validator agreement.
- Fees
    - Network fees paid by the sender or receiver not the bridge.
    - Validators can collect fees from senders.
- Asset Support
    - Limited preselected set of assets.
    - Configurable through validator agreement.
    - Setup of new assets is a manual operation.
    - Assets with any auth required, auth clawback enabled, and auth immutable state supported but behavior may be undefined if assets clawed back, trust line deauthorized, or trust line not authorized.
- Operational
    - Validators can always reconstruct their local state using on-chain data only. This makes it easy to restart after data loss.
