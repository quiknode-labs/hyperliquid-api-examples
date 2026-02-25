"""Place a perp limit order (BTC, IOC, ~$10 notional)."""

import json
from client import rpc, sign_hash, get_mid

COIN = "BTC"

mid = get_mid(COIN)
if mid == 0:
    print(f"Could not fetch {COIN} mid price")
    raise SystemExit(1)

sz = round(11.0 / mid, 5)
buy_px = int(mid * 1.03)

print(f"{COIN} mid: ${mid:,.2f}")
print(f"BUY {sz} @ {buy_px} (IOC, ~${sz * mid:.2f} notional)")

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
