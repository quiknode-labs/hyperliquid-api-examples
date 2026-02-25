"""Place a resting order (GTC, 3% below mid), then cancel it."""

import json
from client import rpc, sign_hash, get_mid

COIN = "BTC"

mid = get_mid(COIN)
if mid == 0:
    print(f"Could not fetch {COIN} mid price")
    raise SystemExit(1)

sz = round(11.0 / mid, 5)
rest_px = int(mid * 0.97)

print(f"{COIN} mid: ${mid:,.2f}")
print(f"Placing resting BUY {sz} @ {rest_px} (GTC, 3% below mid)\n")

order_action = {
    "type": "order",
    "orders": [{
        "a": COIN,
        "b": True,
        "p": str(rest_px),
        "s": str(sz),
        "r": False,
        "t": {"limit": {"tif": "Gtc"}},
    }],
    "grouping": "na",
}

res = rpc("hl_buildOrder", {"action": order_action})
sig = sign_hash(res["result"]["hash"])

result = rpc("hl_sendOrder", {
    "action": res["result"].get("action", order_action),
    "nonce": res["result"]["nonce"],
    "signature": sig,
})

exchange = result["result"]["exchangeResponse"]
statuses = exchange.get("response", {}).get("data", {}).get("statuses", [])

oid = None
for s in statuses:
    if isinstance(s, dict) and "resting" in s:
        oid = s["resting"].get("oid")
        break

if oid is None:
    print("Could not extract OID from resting order")
    print(json.dumps(exchange, indent=2))
    raise SystemExit(1)

print(f"Order resting (OID: {oid})")
print("Cancelling...\n")

cancel_action = {
    "type": "cancel",
    "cancels": [{"a": COIN, "o": oid}],
}

res = rpc("hl_buildCancel", {"action": cancel_action})
sig = sign_hash(res["result"]["hash"])

result = rpc("hl_sendCancel", {
    "action": cancel_action,
    "nonce": res["result"]["nonce"],
    "signature": sig,
})

print(json.dumps(result["result"]["exchangeResponse"], indent=2))
print("\nOrder cancelled.")
