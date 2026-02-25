"""Place an IOC order on a HIP-3 market by name (e.g. xyz:SILVER)."""

import json
import requests
from client import rpc, sign_hash

COIN = "xyz:SILVER"
DEX = COIN.split(":")[0]

r = requests.post("https://api.hyperliquid.xyz/info", json={"type": "allMids", "dex": DEX})
mid = float(r.json().get(COIN, 0))
if mid == 0:
    print(f"Could not fetch {COIN} mid price, using fallback")
    mid = 78.0

sz = round(10.0 / mid, 2)
buy_px = round(mid * 1.03, 2)

print(f"{COIN} mid: ${mid:,.2f}")
print(f"BUY {sz} @ {buy_px} (IOC, ~${sz * mid:.2f} notional)\n")

action = {
    "type": "order",
    "orders": [{
        "a": COIN,
        "b": True,
        "p": str(buy_px),
        "s": str(sz),
        "r": False,
        "t": {"limit": {"tif": "Ioc"}},
    }],
    "grouping": "na",
}

res = rpc("hl_buildOrder", {"action": action})
sig = sign_hash(res["result"]["hash"])

result = rpc("hl_sendOrder", {
    "action": res["result"].get("action", action),
    "nonce": res["result"]["nonce"],
    "signature": sig,
})

print(json.dumps(result["result"]["exchangeResponse"], indent=2))
