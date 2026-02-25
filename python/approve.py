"""Approve builder fee (one-time setup)."""

from client import rpc, sign_hash

MAX_FEE = "1%"

res = rpc("hl_buildApproveBuilderFee", {"maxFeeRate": MAX_FEE})
sig = sign_hash(res["result"]["hash"])

rpc("hl_sendApproval", {
    "nonce": res["result"]["nonce"],
    "signature": sig,
    "maxFeeRate": MAX_FEE,
})

print("Builder fee approved.")
