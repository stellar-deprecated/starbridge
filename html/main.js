const requiredSigs = 1;

var validatorUrls = [
    "https://starbridge1.prototypes.kube001.services.stellar-ops.com",
    "https://starbridge2.prototypes.kube001.services.stellar-ops.com",
    "https://starbridge3.prototypes.kube001.services.stellar-ops.com"
];

document.getElementById("deposit").onclick = function() {
    var form = new FormData();
    form.append("stellar_address", document.getElementById("stellar_address").value);

    axios.post(validatorUrls[0]+"/deposit", form).then(response => {
        document.getElementById("eth_hash").textContent = response.data
    });
}

document.getElementById("withdraw").onclick = function() {
    document.getElementById("validators").style.display = "block";

    var promises = [];

    for (var i = 0; i < validatorUrls.length; i++) {
        var j = i;
        
        promises[j] = new Promise(async (resolve, reject) => {
            var form = new FormData();
            form.append("transaction_hash", document.getElementById("transaction_hash").value);
            form.append("tx_expiration_timestamp", Math.ceil(Date.now()/1000)+5*60);

            var status = document.getElementById("status"+j)
            status.textContent = "Sending request..."
            while (true) {
                const response = await axios.post(
                    validatorUrls[j]+"/stellar/get_inverse_transaction/ethereum",
                    form,
                    {validateStatus: null}
                );
                switch (response.status) {
                    case 202: // Accepted
                        status.textContent = "Signing request sent, waiting..."
                        break;
                    case 404: // NotFound
                        status.textContent = "Tx not found"
                        resolve(response);
                        return;
                    case 200: // OK
                        status.textContent = "Signature received!"
                        resolve(response);
                        return
                    default: // Error
                        status.textContent = "Unknown error"
                        resolve(response);
                        return
                }

                await sleep(1000);
            }
        });
    }

    Promise.all(promises).then(results => {
        var sigs = 0;
        var tx;
        for (var i = 0; i < results.length; i++) {
            if (sigs == requiredSigs) break;

            // No signature
            if (results[i].data == "") {
                continue;
            }

            localTx = StellarSdk.TransactionBuilder.fromXDR(results[i].data,  "Test SDF Network ; September 2015")

            sigs++
            if (!tx) {
                tx = localTx
                continue
            }

            tx._signatures.push(localTx._signatures[0]);
        }

        if (sigs != requiredSigs) {
            // Missing signatures
            return
        }

        document.getElementById("signed").style.display = "block";
        document.getElementById("signed_tx").textContent = tx.toXDR()
        document.getElementById("signed_link").href = "https://laboratory.stellar.org/#txsigner?xdr="+encodeURIComponent(tx.toXDR())+"&network=test";
    });
};

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}