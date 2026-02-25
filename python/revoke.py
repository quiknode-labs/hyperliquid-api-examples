"""Revoke builder fee approval."""

from client import rpc, sign_hash

res = rpc("hl_buildRevokeBuilderFee", {})
sig = sign_hash(res["result"]["hash"])

rpc("hl_sendRevocation", {
    "nonce": res["result"]["nonce"],
    "signature": sig,
    "maxFeeRate": "0%",
})

print("Builder fee revoked.")
