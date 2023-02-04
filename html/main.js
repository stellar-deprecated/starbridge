const requiredSigs = 2;

var validatorUrls = [
    "https://starbridge1.prototypes.kube001.services.stellar-ops.com",
    "https://starbridge2.prototypes.kube001.services.stellar-ops.com",
    "https://starbridge3.prototypes.kube001.services.stellar-ops.com"
];

function genRandomHash() {
    const hex = '0123456789abcdef';
    let output = '';
    for (let i = 0; i < 64; ++i) {
        output += hex.charAt(Math.floor(Math.random() * hex.length));
    }
    return output;
}

document.getElementById("deposit").onclick = function() {
    hash = genRandomHash()
    var form = new FormData();
    form.append("hash", hash);
    form.append("stellar_address", document.getElementById("stellar_address").value);

    var promises = [];

    for (var i = 0; i < validatorUrls.length; i++) {
        promises[i] = axios.post(validatorUrls[i]+"/deposit", form);
    }

    Promise.all(promises).then(results => {
        document.getElementById("eth_hash").textContent = hash;
        document.getElementById("transaction_hash").value = hash;
    });
}

document.getElementById("withdraw").onclick = function() {
    document.getElementById("validators").style.display = "block";

    var promises = [];

    for (var i = 0; i < validatorUrls.length; i++) {
        promises[i] = function(j) {
            return new Promise(async (resolve, reject) => {
                var form = new FormData();
                form.append("transaction_hash", document.getElementById("transaction_hash").value);
                form.append("log_index", "1");
                form.append("tx_expiration_timestamp", Math.ceil(Date.now()/1000)+5*60);

                var status = document.getElementById("status"+j)
                status.textContent = "Sending request..."
                while (true) {
                    const response = await axios.post(
                        validatorUrls[j]+"/ethereum/withdraw",
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
        }(i);
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