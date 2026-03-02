"""Place a market order — no price needed, worker computes it automatically."""

import json
from client import exchange, sign_hash

COIN = "BTC"
SIZE = "0.00011"

print(f"Market BUY {SIZE} {COIN}\n")

# Optional: custom slippage (default 3%, range 0.1%-10%)
# res = exchange({"action": {...}, "slippage": 0.05})  # 5% slippage

res = exchange({
    "action": {
        "type": "order",
        "orders": [{
            "asset": COIN,
            "side": "buy",
            "size": SIZE,
            "tif": "market",
        }],
    },
})

computed_price = res["action"]["orders"][0]["p"]
print(f"Computed price (mid + slippage, default 3%): {computed_price}")
print(f"Builder fee: {res['builderFee']}")

sig = sign_hash(res["hash"])

result = exchange({
    "action": res["action"],
    "nonce": res["nonce"],
    "signature": sig,
})

print(json.dumps(result["exchangeResponse"], indent=2))
